package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (c *Container) toModelJSON(out *modeljson.Container) {
	*out = modeljson.Container{
		ID:      c.Id,
		Name:    c.Name,
		Runtime: c.Runtime,
		Image: modeljson.ContainerImage{
			Name: c.ImageName,
			Tag:  c.ImageTag,
		},
	}
}
