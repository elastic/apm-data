package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestObserverToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Observer
		expected *modeljson.Observer
	}{
		"empty": {
			proto:    &Observer{},
			expected: &modeljson.Observer{},
		},
		"full": {
			proto: &Observer{
				Hostname: "hostname",
				Name:     "name",
				Type:     "type",
				Version:  "version",
			},
			expected: &modeljson.Observer{
				Hostname: "hostname",
				Name:     "name",
				Type:     "type",
				Version:  "version",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Observer
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
