package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestDeviceToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Device
		expected *modeljson.Device
	}{
		"empty": {
			proto:    &Device{},
			expected: &modeljson.Device{},
		},
		"no pointers": {
			proto: &Device{
				Id:           "id",
				Manufacturer: "manufacturer",
			},
			expected: &modeljson.Device{
				ID:           "id",
				Manufacturer: "manufacturer",
			},
		},
		"full": {
			proto: &Device{
				Id: "id",
				Model: &DeviceModel{
					Name:       "name",
					Identifier: "identifier",
				},
				Manufacturer: "manufacturer",
			},
			expected: &modeljson.Device{
				ID: "id",
				Model: modeljson.DeviceModel{
					Name:       "name",
					Identifier: "identifier",
				},
				Manufacturer: "manufacturer",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Device
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
