package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (n *Network) toModelJSON(out *modeljson.Network) {
	*out = modeljson.Network{}
	if n.Connection != nil {
		out.Connection = modeljson.NetworkConnection{
			Type:    n.Connection.Type,
			Subtype: n.Connection.Subtype,
		}
	}
	if n.Carrier != nil {
		out.Carrier = modeljson.NetworkCarrier{
			Name: n.Carrier.Name,
			MCC:  n.Carrier.Mcc,
			MNC:  n.Carrier.Mnc,
			ICC:  n.Carrier.Icc,
		}
	}
}
