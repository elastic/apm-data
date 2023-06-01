package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (o *OS) toModelJSON(out *modeljson.OS) {
	*out = modeljson.OS{
		Name:     o.Name,
		Version:  o.Version,
		Platform: o.Platform,
		Full:     o.Full,
		Type:     o.Type,
	}
}
