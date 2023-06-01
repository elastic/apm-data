package modelpb

import (
	"net/netip"

	"github.com/elastic/apm-data/model/internal/modeljson"
)

func (c *Client) toModelJSON(out *modeljson.Client) {
	*out = modeljson.Client{
		Domain: c.Domain,
		Port:   int(c.Port),
	}
	if c.Ip != "" {
		if _, err := netip.ParseAddr(c.Ip); err == nil {
			out.IP = c.Ip
		}
	}
}
