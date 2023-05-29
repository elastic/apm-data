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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessTransform(t *testing.T) {
	processTitle := "node"
	commandLine := "node run.js"
	executablePath := "/usr/bin/node"

	ppid := uint32(456)
	tests := []struct {
		Process Process
		Output  any
	}{
		{
			Process: Process{},
			Output:  nil,
		},
		{
			Process: Process{
				Pid:         123,
				Ppid:        ppid,
				Title:       processTitle,
				Argv:        []string{"node", "server.js"},
				CommandLine: commandLine,
				Executable:  executablePath,
				Thread: ProcessThread{
					ID:   1,
					Name: "testThread",
				},
			},
			Output: map[string]any{
				"pid": 123.0,
				"parent": map[string]any{
					"pid": 456.0,
				},
				"title":        processTitle,
				"args":         []any{"node", "server.js"},
				"command_line": commandLine,
				"executable":   executablePath,
				"thread": map[string]any{
					"id":   1.0,
					"name": "testThread",
				},
			},
		},
	}

	for _, test := range tests {
		output := transformAPMEvent(APMEvent{Process: test.Process})
		assert.Equal(t, test.Output, output["process"])
	}
}
