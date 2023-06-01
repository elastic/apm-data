package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestDestinationToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Destination
		expected *modeljson.Destination
	}{
		"empty": {
			proto:    &Destination{},
			expected: &modeljson.Destination{},
		},
		"full": {
			proto: &Destination{
				Address: "127.0.0.1",
				Port:    443,
			},
			expected: &modeljson.Destination{
				Address: "127.0.0.1",
				IP:      "127.0.0.1",
				Port:    443,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Destination
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
