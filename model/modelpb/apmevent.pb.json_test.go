package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestAPMEventToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *APMEvent
		expected *modeljson.Document
	}{
		"empty": {
			proto:    &APMEvent{},
			expected: &modeljson.Document{},
		},
		"full": {
			proto: &APMEvent{
				NumericLabels: map[string]*NumericLabelValue{
					"foo": {
						Values: []float64{1, 2, 3},
						Value:  1,
						Global: true,
					},
				},
				Labels: map[string]*LabelValue{
					"bar": {
						Value:  "a",
						Values: []string{"a", "b", "c"},
						Global: true,
					},
				},
				Message: "message",
			},
			expected: &modeljson.Document{
				NumericLabels: map[string]modeljson.NumericLabel{
					"foo": {
						Values: []float64{1, 2, 3},
						Value:  1,
					},
				},
				Labels: map[string]modeljson.Label{
					"bar": {
						Value:  "a",
						Values: []string{"a", "b", "c"},
					},
				},
				Message: "message",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Document
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out,
				cmpopts.IgnoreFields(APMEvent{}, "Timestamp"),
				cmpopts.IgnoreFields(modeljson.Document{}, "Timestamp"),
			)
			require.Empty(t, diff)
		})
	}
}
