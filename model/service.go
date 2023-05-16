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

// Service bundles together information related to the monitored service and
// the agent used for monitoring
type Service struct {
	Origin      *ServiceOrigin
	Target      *ServiceTarget
	Language    Language
	Runtime     Runtime
	Framework   Framework
	Name        string
	Version     string
	Environment string
	Node        ServiceNode
}

// ServiceOrigin holds information about the service that originated a
// transaction.
type ServiceOrigin struct {
	ID      string
	Name    string
	Version string
}

// ServiceTarget holds information about the target service in case of
// an outgoing event w.r.t. the instrumented service
type ServiceTarget struct {
	Name string
	Type string
}

// Language has an optional version and name
type Language struct {
	Name    string
	Version string
}

// Runtime has an optional version and name
type Runtime struct {
	Name    string
	Version string
}

// Framework has an optional version and name
type Framework struct {
	Name    string
	Version string
}

type ServiceNode struct {
	Name string
}
