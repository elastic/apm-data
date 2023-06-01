package modelpb

import (
	"testing"
	"time"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestEventToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Event
		expected *modeljson.Event
	}{
		"empty": {
			proto:    &Event{},
			expected: &modeljson.Event{},
		},
		"full": {
			proto: &Event{
				Outcome:  "outcome",
				Action:   "action",
				Dataset:  "dataset",
				Kind:     "kind",
				Category: "category",
				Type:     "type",
				SuccessCount: &SummaryMetric{
					Count: 1,
					Sum:   2,
				},
				Duration: durationpb.New(3 * time.Second),
				Severity: 4,
			},
			expected: &modeljson.Event{
				Outcome:  "outcome",
				Action:   "action",
				Dataset:  "dataset",
				Kind:     "kind",
				Category: "category",
				Type:     "type",
				SuccessCount: modeljson.SummaryMetric{
					Count: 1,
					Sum:   2,
				},
				Duration: int64(3 * time.Second),
				Severity: 4,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Event
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
