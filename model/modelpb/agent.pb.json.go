package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (a *Agent) toModelJSON(out *modeljson.Agent) {
	*out = modeljson.Agent{
		Name:             a.Name,
		Version:          a.Version,
		EphemeralID:      a.EphemeralId,
		ActivationMethod: a.ActivationMethod,
	}
}
