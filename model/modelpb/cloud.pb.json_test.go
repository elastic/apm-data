package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestCloudToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Cloud
		expected *modeljson.Cloud
	}{
		"empty": {
			proto:    &Cloud{},
			expected: &modeljson.Cloud{},
		},
		"no pointers": {
			proto: &Cloud{
				AccountId:        "accountid",
				AccountName:      "accountname",
				AvailabilityZone: "availabilityzone",
				InstanceId:       "instanceid",
				InstanceName:     "instancename",
				MachineType:      "machinetype",
				ProjectId:        "projectid",
				ProjectName:      "projectname",
				Provider:         "provider",
				Region:           "region",
				ServiceName:      "servicename",
			},
			expected: &modeljson.Cloud{
				AvailabilityZone: "availabilityzone",
				Provider:         "provider",
				Region:           "region",
				Account: modeljson.CloudAccount{
					ID:   "accountid",
					Name: "accountname",
				},
				Instance: modeljson.CloudInstance{
					ID:   "instanceid",
					Name: "instancename",
				},
				Machine: modeljson.CloudMachine{
					Type: "machinetype",
				},
				Project: modeljson.CloudProject{
					ID:   "projectid",
					Name: "projectname",
				},
				Service: modeljson.CloudService{
					Name: "servicename",
				},
			},
		},
		"full": {
			proto: &Cloud{
				Origin: &CloudOrigin{
					AccountId:   "origin_accountid",
					Provider:    "origin_provider",
					Region:      "origin_region",
					ServiceName: "origin_servicename",
				},
				AccountId:        "accountid",
				AccountName:      "accountname",
				AvailabilityZone: "availabilityzone",
				InstanceId:       "instanceid",
				InstanceName:     "instancename",
				MachineType:      "machinetype",
				ProjectId:        "projectid",
				ProjectName:      "projectname",
				Provider:         "provider",
				Region:           "region",
				ServiceName:      "servicename",
			},
			expected: &modeljson.Cloud{
				AvailabilityZone: "availabilityzone",
				Provider:         "provider",
				Region:           "region",
				Origin: modeljson.CloudOrigin{
					Account: modeljson.CloudAccount{
						ID:   "origin_accountid",
						// TODO this will always be empty
						//Name: "origin_accountname",
					},
					Provider: "origin_provider",
					Region:   "origin_region",
					Service: modeljson.CloudService{
						Name: "origin_servicename",
					},
				},

				Account: modeljson.CloudAccount{
					ID:   "accountid",
					Name: "accountname",
				},
				Instance: modeljson.CloudInstance{
					ID:   "instanceid",
					Name: "instancename",
				},
				Machine: modeljson.CloudMachine{
					Type: "machinetype",
				},
				Project: modeljson.CloudProject{
					ID:   "projectid",
					Name: "projectname",
				},
				Service: modeljson.CloudService{
					Name: "servicename",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Cloud
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
