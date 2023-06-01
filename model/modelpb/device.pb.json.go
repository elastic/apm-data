package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (d *Device) toModelJSON(out *modeljson.Device) {
	*out = modeljson.Device{
		ID:           d.Id,
		Manufacturer: d.Manufacturer,
	}
	if d.Model != nil {
		out.Model = modeljson.DeviceModel{
			Name:       d.Model.Name,
			Identifier: d.Model.Identifier,
		}
	}
}
