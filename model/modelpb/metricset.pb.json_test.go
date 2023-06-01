package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestMetricsetToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Metricset
		expected *modeljson.Metricset
	}{
		"empty": {
			proto:    &Metricset{},
			expected: &modeljson.Metricset{},
		},
		"full": {
			proto: &Metricset{
				Name:     "name",
				Interval: "interval",
				Samples: []*MetricsetSample{
					{
						Type: MetricType_METRIC_TYPE_COUNTER,
						Name: "name",
						Unit: "unit",
						Histogram: &Histogram{
							Values: []float64{1},
							Counts: []int64{2},
						},
						Summary: &SummaryMetric{
							Count: 3,
							Sum:   4,
						},
						Value: 5,
					},
				},
				DocCount: 1,
			},
			expected: &modeljson.Metricset{
				Name:     "name",
				Interval: "interval",
				Samples: []modeljson.MetricsetSample{
					{
						Type: "METRIC_TYPE_COUNTER",
						Name: "name",
						Unit: "unit",
						Histogram: modeljson.Histogram{
							Values: []float64{1},
							Counts: []int64{2},
						},
						Summary: modeljson.SummaryMetric{
							Count: 3,
							Sum:   4,
						},
						Value: 5,
					},
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Metricset
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
