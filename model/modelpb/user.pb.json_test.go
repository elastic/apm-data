package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestUserToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *User
		expected *modeljson.User
	}{
		"empty": {
			proto:    &User{},
			expected: &modeljson.User{},
		},
		"full": {
			proto: &User{
				Domain: "domain",
				Id:     "id",
				Email:  "email",
				Name:   "name",
			},
			expected: &modeljson.User{
				Domain: "domain",
				ID:     "id",
				Email:  "email",
				Name:   "name",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.User
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
