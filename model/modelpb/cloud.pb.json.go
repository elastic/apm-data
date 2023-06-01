package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (c *Cloud) toModelJSON(out *modeljson.Cloud) {
	*out = modeljson.Cloud{
		AvailabilityZone: c.AvailabilityZone,
		Provider:         c.Provider,
		Region:           c.Region,
		Account: modeljson.CloudAccount{
			ID:   c.AccountId,
			Name: c.AccountName,
		},
		Service: modeljson.CloudService{
			Name: c.ServiceName,
		},
		Project: modeljson.CloudProject{
			ID:   c.ProjectId,
			Name: c.ProjectName,
		},
		Instance: modeljson.CloudInstance{
			ID:   c.InstanceId,
			Name: c.InstanceName,
		},
		Machine: modeljson.CloudMachine{
			Type: c.MachineType,
		},
	}
	if c.Origin != nil {
		out.Origin = modeljson.CloudOrigin{
			Provider: c.Origin.Provider,
			Region:   c.Origin.Region,
			Account: modeljson.CloudAccount{
				ID: c.Origin.AccountId,
			},
			Service: modeljson.CloudService{
				Name: c.Origin.ServiceName,
			},
		}
	}

}
