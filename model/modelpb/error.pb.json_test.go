package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestErrorToModelJSON(t *testing.T) {
	handled := true

	attrs, attrsMap := randomStruct(t)

	testCases := map[string]struct {
		proto    *Error
		expected *modeljson.Error
	}{
		"empty": {
			proto:    &Error{},
			expected: &modeljson.Error{},
		},
		"full": {
			proto: &Error{
				Exception: &Exception{
					Message:    "ex_message",
					Module:     "ex_module",
					Code:       "ex_code",
					Attributes: attrs,
					Type:       "ex_type",
					Handled:    &handled,
					Cause: []*Exception{
						{
							Message: "ex1_message",
							Module:  "ex1_module",
							Code:    "ex1_code",
							Type:    "ex_type",
						},
					},
				},
				Log: &ErrorLog{
					Message:      "log_message",
					Level:        "log_level",
					ParamMessage: "log_parammessage",
					LoggerName:   "log_loggername",
				},
				Id:          "id",
				GroupingKey: "groupingkey",
				Culprit:     "culprit",
				StackTrace:  "stacktrace",
				Message:     "message",
				Type:        "type",
			},
			expected: &modeljson.Error{
				Exception: &modeljson.Exception{
					Message:    "ex_message",
					Module:     "ex_module",
					Code:       "ex_code",
					Attributes: attrsMap,
					Type:       "ex_type",
					Handled:    &handled,
					Cause: []modeljson.Exception{
						{
							Message: "ex1_message",
							Module:  "ex1_module",
							Code:    "ex1_code",
							Type:    "ex_type",
						},
					},
				},
				Log: &modeljson.ErrorLog{
					Message:      "log_message",
					Level:        "log_level",
					ParamMessage: "log_parammessage",
					LoggerName:   "log_loggername",
				},
				ID:          "id",
				GroupingKey: "groupingkey",
				Culprit:     "culprit",
				StackTrace:  "stacktrace",
				Message:     "message",
				Type:        "type",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Error
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
