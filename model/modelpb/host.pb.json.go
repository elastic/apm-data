package modelpb

import (
	"net/netip"

	"github.com/elastic/apm-data/model/internal/modeljson"
)

func (h *Host) toModelJSON(out *modeljson.Host) {
	*out = modeljson.Host{
		Hostname:     h.Hostname,
		Name:         h.Name,
		ID:           h.Id,
		Architecture: h.Architecture,
		Type:         h.Type,
	}

	if n := len(h.Ip); n != 0 {
		ips := make([]string, 0, n)
		for _, ip := range h.Ip {
			if _, err := netip.ParseAddr(ip); err == nil {
				ips = append(ips, ip)
			}
		}
		out.IP = ips
	}

	if h.Os != nil {
		var os modeljson.OS
		h.Os.toModelJSON(&os)
		out.OS = &os
	}
}
