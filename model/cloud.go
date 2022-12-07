// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package model

// Cloud holds information about the cloud computing environment
// in which a service is running.
type Cloud struct {
	AccountID        string
	AccountName      string
	AvailabilityZone string
	InstanceID       string
	InstanceName     string
	MachineType      string
	ProjectID        string
	ProjectName      string
	Provider         string
	Region           string
	ServiceName      string

	Origin *CloudOrigin
}

type CloudOrigin struct {
	AccountID   string
	Provider    string
	Region      string
	ServiceName string
}

func (c *Cloud) fields() map[string]any {
	var fields mapStr

	var account, instance, machine, project, service mapStr
	account.maybeSetString("id", c.AccountID)
	account.maybeSetString("name", c.AccountName)
	instance.maybeSetString("id", c.InstanceID)
	instance.maybeSetString("name", c.InstanceName)
	machine.maybeSetString("type", c.MachineType)
	project.maybeSetString("id", c.ProjectID)
	project.maybeSetString("name", c.ProjectName)
	service.maybeSetString("name", c.ServiceName)

	fields.maybeSetMapStr("account", map[string]any(account))
	fields.maybeSetString("availability_zone", c.AvailabilityZone)
	fields.maybeSetMapStr("instance", map[string]any(instance))
	fields.maybeSetMapStr("machine", map[string]any(machine))
	fields.maybeSetMapStr("project", map[string]any(project))
	fields.maybeSetMapStr("service", map[string]any(service))
	fields.maybeSetString("provider", c.Provider)
	fields.maybeSetString("region", c.Region)
	if c.Origin != nil {
		var origin mapStr
		origin.maybeSetString("account.id", c.Origin.AccountID)
		origin.maybeSetString("provider", c.Origin.Provider)
		origin.maybeSetString("region", c.Origin.Region)
		origin.maybeSetString("service.name", c.Origin.ServiceName)
		fields.maybeSetMapStr("origin", map[string]any(origin))
	}
	return map[string]any(fields)
}
