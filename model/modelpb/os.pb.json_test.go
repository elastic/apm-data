package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestOSToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *OS
		expected *modeljson.OS
	}{
		"empty": {
			proto:    &OS{},
			expected: &modeljson.OS{},
		},
		"full": {
			proto: &OS{
				Name:     "name",
				Version:  "version",
				Platform: "platform",
				Full:     "full",
				Type:     "type",
			},
			expected: &modeljson.OS{
				Name:     "name",
				Version:  "version",
				Platform: "platform",
				Full:     "full",
				Type:     "type",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.OS
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
