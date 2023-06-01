package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestKubernetesToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Kubernetes
		expected *modeljson.Kubernetes
	}{
		"empty": {
			proto:    &Kubernetes{},
			expected: &modeljson.Kubernetes{},
		},
		"full": {
			proto: &Kubernetes{
				Namespace: "namespace",
				NodeName:  "nodename",
				PodName:   "podname",
				PodUid:    "poduid",
			},
			expected: &modeljson.Kubernetes{
				Namespace: "namespace",
				Node: modeljson.KubernetesNode{
					Name: "nodename",
				},
				Pod: modeljson.KubernetesPod{
					Name: "podname",
					UID:  "poduid",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Kubernetes
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
