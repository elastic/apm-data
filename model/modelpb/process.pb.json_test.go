package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestProcessToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Process
		expected *modeljson.Process
	}{
		"empty": {
			proto:    &Process{},
			expected: &modeljson.Process{},
		},
		"no pointers": {
			proto: &Process{
				Title:       "title",
				CommandLine: "commandline",
				Executable:  "executable",
				Pid:         3,
			},
			expected: &modeljson.Process{
				Title:       "title",
				CommandLine: "commandline",
				Executable:  "executable",
				Pid:         3,
			},
		},
		"full": {
			proto: &Process{
				Ppid: 1,
				Thread: &ProcessThread{
					Name: "name",
					Id:   2,
				},
				Title:       "title",
				CommandLine: "commandline",
				Executable:  "executable",
				Argv:        []string{"argv"},
				Pid:         3,
			},
			expected: &modeljson.Process{
				Title:       "title",
				CommandLine: "commandline",
				Executable:  "executable",
				Args:        []string{"argv"},
				Thread: modeljson.ProcessThread{
					Name: "name",
					ID:   2,
				},
				Parent: modeljson.ProcessParent{
					Pid: 1,
				},
				Pid: 3,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Process
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
