package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestServiceToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Service
		expected *modeljson.Service
	}{
		"empty": {
			proto:    &Service{},
			expected: &modeljson.Service{},
		},
		"no pointers": {
			proto: &Service{
				Name:        "name",
				Version:     "version",
				Environment: "environment",
			},
			expected: &modeljson.Service{
				Name:        "name",
				Version:     "version",
				Environment: "environment",
			},
		},
		"full": {
			proto: &Service{
				Origin: &ServiceOrigin{
					Id:      "origin_id",
					Name:    "origin_name",
					Version: "origin_version",
				},
				Target: &ServiceTarget{
					Name: "target_name",
					Type: "target_type",
				},
				Language: &Language{
					Name:    "language_name",
					Version: "language_version",
				},
				Runtime: &Runtime{
					Name:    "runtime_name",
					Version: "runtime_version",
				},
				Framework: &Framework{
					Name:    "framework_name",
					Version: "framework_version",
				},
				Name:        "name",
				Version:     "version",
				Environment: "environment",
				Node: &ServiceNode{
					Name: "node_name",
				},
			},
			expected: &modeljson.Service{
				Origin: &modeljson.ServiceOrigin{
					ID:      "origin_id",
					Name:    "origin_name",
					Version: "origin_version",
				},
				Target: &modeljson.ServiceTarget{
					Name: "target_name",
					Type: "target_type",
				},
				Language: &modeljson.Language{
					Name:    "language_name",
					Version: "language_version",
				},
				Runtime: &modeljson.Runtime{
					Name:    "runtime_name",
					Version: "runtime_version",
				},
				Framework: &modeljson.Framework{
					Name:    "framework_name",
					Version: "framework_version",
				},
				Name:        "name",
				Version:     "version",
				Environment: "environment",
				Node: &modeljson.ServiceNode{
					Name: "node_name",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Service
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
