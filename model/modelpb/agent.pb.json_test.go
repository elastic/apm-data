package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestAgentToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Agent
		expected *modeljson.Agent
	}{
		"empty": {
			proto:    &Agent{},
			expected: &modeljson.Agent{},
		},
		"full": {
			proto: &Agent{
				Name:             "name",
				Version:          "version",
				EphemeralId:      "ephemeralid",
				ActivationMethod: "activationmethod",
			},
			expected: &modeljson.Agent{
				Name:             "name",
				Version:          "version",
				EphemeralID:      "ephemeralid",
				ActivationMethod: "activationmethod",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Agent
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
