package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestNetworkToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Network
		expected *modeljson.Network
	}{
		"empty": {
			proto:    &Network{},
			expected: &modeljson.Network{},
		},
		"connection": {
			proto: &Network{
				Connection: &NetworkConnection{
					Type:    "type",
					Subtype: "subtype",
				},
			},
			expected: &modeljson.Network{
				Connection: modeljson.NetworkConnection{
					Type:    "type",
					Subtype: "subtype",
				},
			},
		},
		"carrier": {
			proto: &Network{
				Carrier: &NetworkCarrier{
					Name: "name",
					Mcc:  "mcc",
					Mnc:  "mnc",
					Icc:  "icc",
				},
			},
			expected: &modeljson.Network{
				Carrier: modeljson.NetworkCarrier{
					Name: "name",
					MCC:  "mcc",
					MNC:  "mnc",
					ICC:  "icc",
				},
			},
		},
		"full": {
			proto: &Network{
				Connection: &NetworkConnection{
					Type:    "type",
					Subtype: "subtype",
				},
				Carrier: &NetworkCarrier{
					Name: "name",
					Mcc:  "mcc",
					Mnc:  "mnc",
					Icc:  "icc",
				},
			},
			expected: &modeljson.Network{
				Connection: modeljson.NetworkConnection{
					Type:    "type",
					Subtype: "subtype",
				},
				Carrier: modeljson.NetworkCarrier{
					Name: "name",
					MCC:  "mcc",
					MNC:  "mnc",
					ICC:  "icc",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Network
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
