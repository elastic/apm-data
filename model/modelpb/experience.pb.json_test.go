package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestUserExperienceToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *UserExperience
		expected *modeljson.UserExperience
	}{
		"empty": {
			proto:    &UserExperience{},
			expected: &modeljson.UserExperience{},
		},
		"no pointers": {
			proto: &UserExperience{
				CumulativeLayoutShift: 1,
				FirstInputDelay:       2,
				TotalBlockingTime:     3,
				LongTask: &LongtaskMetrics{
					Count: 4,
					Sum:   5,
					Max:   6,
				},
			},
			expected: &modeljson.UserExperience{
				CumulativeLayoutShift: 1,
				FirstInputDelay:       2,
				TotalBlockingTime:     3,
				Longtask: modeljson.LongtaskMetrics{
					Count: 4,
					Sum:   5,
					Max:   6,
				},
			},
		},
		"full": {
			proto: &UserExperience{
				CumulativeLayoutShift: 1,
				FirstInputDelay:       2,
				TotalBlockingTime:     3,
				LongTask: &LongtaskMetrics{
					Count: 4,
					Sum:   5,
					Max:   6,
				},
			},
			expected: &modeljson.UserExperience{
				CumulativeLayoutShift: 1,
				FirstInputDelay:       2,
				TotalBlockingTime:     3,
				Longtask: modeljson.LongtaskMetrics{
					Count: 4,
					Sum:   5,
					Max:   6,
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.UserExperience
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
