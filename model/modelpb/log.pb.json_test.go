package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestLogToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Log
		expected *modeljson.Log
	}{
		"empty": {
			proto:    &Log{},
			expected: &modeljson.Log{},
		},
		"no pointers": {
			proto: &Log{
				Level:  "level",
				Logger: "logger",
			},
			expected: &modeljson.Log{
				Level:  "level",
				Logger: "logger",
			},
		},
		"full": {
			proto: &Log{
				Level:  "level",
				Logger: "logger",
				Origin: &LogOrigin{
					FunctionName: "functionname",
					File: &LogOriginFile{
						Name: "name",
						Line: 1,
					},
				},
			},
			expected: &modeljson.Log{
				Level:  "level",
				Logger: "logger",
				Origin: modeljson.LogOrigin{
					Function: "functionname",
					File: modeljson.LogOriginFile{
						Name: "name",
						Line: 1,
					},
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Log
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
