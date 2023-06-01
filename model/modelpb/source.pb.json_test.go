package modelpb

import (
	"net/netip"
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestSourceToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Source
		expected *modeljson.Source
	}{
		"empty": {
			proto:    &Source{},
			expected: &modeljson.Source{},
		},
		"full": {
			proto: &Source{
				Ip: "127.0.0.1",
				Nat: &NAT{
					Ip: "127.0.0.2",
				},
				Domain: "domain",
				Port:   443,
			},
			expected: &modeljson.Source{
				IP: modeljson.IP(netip.MustParseAddr("127.0.0.1")),
				NAT: modeljson.NAT{
					IP: modeljson.IP(netip.MustParseAddr("127.0.0.2")),
				},
				Domain: "domain",
				Port:   443,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Source
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out,
				cmp.Comparer(func(a modeljson.IP, b modeljson.IP) bool {
					return netip.Addr(a).Compare(netip.Addr(b)) == 0
				}),
			)
			require.Empty(t, diff)
		})
	}
}
