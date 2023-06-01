package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestHostToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Host
		expected *modeljson.Host
	}{
		"empty": {
			proto:    &Host{},
			expected: &modeljson.Host{},
		},
		"no pointers": {
			proto: &Host{
				Hostname:     "hostname",
				Name:         "name",
				Id:           "id",
				Architecture: "architecture",
				Type:         "type",
			},
			expected: &modeljson.Host{
				Hostname:     "hostname",
				Name:         "name",
				ID:           "id",
				Architecture: "architecture",
				Type:         "type",
			},
		},
		"full": {
			proto: &Host{
				Hostname:     "hostname",
				Name:         "name",
				Id:           "id",
				Architecture: "architecture",
				Type:         "type",
				Ip:           []string{"127.0.0.1"},
			},
			expected: &modeljson.Host{
				Hostname:     "hostname",
				Name:         "name",
				ID:           "id",
				Architecture: "architecture",
				Type:         "type",
				IP:           []string{"127.0.0.1"},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Host
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
