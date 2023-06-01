package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (k *Kubernetes) toModelJSON(out *modeljson.Kubernetes) {
	*out = modeljson.Kubernetes{
		Namespace: k.Namespace,
		Node: modeljson.KubernetesNode{
			Name: k.NodeName,
		},
		Pod: modeljson.KubernetesPod{
			Name: k.PodName,
			UID:  k.PodUid,
		},
	}
}
