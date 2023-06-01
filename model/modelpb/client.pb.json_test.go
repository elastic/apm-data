package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestClientToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Client
		expected *modeljson.Client
	}{
		"empty": {
			proto:    &Client{},
			expected: &modeljson.Client{},
		},
		"full": {
			proto: &Client{
				Ip:     "127.0.0.1",
				Domain: "example.com",
				Port:   443,
			},
			expected: &modeljson.Client{
				IP:     "127.0.0.1",
				Domain: "example.com",
				Port:   443,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Client
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
