package modelpb

import (
	"net/netip"

	"github.com/elastic/apm-data/model/internal/modeljson"
)

func (s *Source) toModelJSON(out *modeljson.Source) {
	*out = modeljson.Source{
		Domain: s.Domain,
		Port:   int(s.Port),
	}
	if ip, err := netip.ParseAddr(s.Ip); err == nil {
		out.IP = modeljson.IP(ip)
	}
	if s.Nat != nil {
		if ip, err := netip.ParseAddr(s.Nat.Ip); err == nil {
			out.NAT.IP = modeljson.IP(ip)
		}
	}
}
