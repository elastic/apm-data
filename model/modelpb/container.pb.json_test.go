package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestContainerToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Container
		expected *modeljson.Container
	}{
		"empty": {
			proto:    &Container{},
			expected: &modeljson.Container{},
		},
		"full": {
			proto: &Container{
				Id:        "id",
				Name:      "name",
				Runtime:   "runtime",
				ImageName: "imagename",
				ImageTag:  "imagetag",
			},
			expected: &modeljson.Container{
				ID:      "id",
				Name:    "name",
				Runtime: "runtime",
				Image: modeljson.ContainerImage{
					Name: "imagename",
					Tag:  "imagetag",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Container
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
