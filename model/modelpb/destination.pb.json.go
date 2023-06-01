package modelpb

import (
	"net/netip"

	"github.com/elastic/apm-data/model/internal/modeljson"
)

func (d *Destination) toModelJSON(out *modeljson.Destination) {
	*out = modeljson.Destination{
		Address: d.Address,
		Port:    int(d.Port),
	}
	if d.Address != "" {
		if _, err := netip.ParseAddr(d.Address); err == nil {
			out.IP = d.Address
		}
	}
}
