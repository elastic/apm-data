package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (o *Observer) toModelJSON(out *modeljson.Observer) {
	*out = modeljson.Observer{
		Hostname: o.Hostname,
		Name:     o.Name,
		Type:     o.Type,
		Version:  o.Version,
	}
}
