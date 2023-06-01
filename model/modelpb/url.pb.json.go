package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (u *URL) toModelJSON(out *modeljson.URL) {
	*out = modeljson.URL{
		Original: u.Original,
		Scheme:   u.Scheme,
		Full:     u.Full,
		Domain:   u.Domain,
		Path:     u.Path,
		Query:    u.Query,
		Fragment: u.Fragment,
		Port:     int(u.Port),
	}
}
