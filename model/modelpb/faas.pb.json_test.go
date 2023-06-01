package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestFaasToModelJSON(t *testing.T) {
	coldstart := true

	testCases := map[string]struct {
		proto    *Faas
		expected *modeljson.FAAS
	}{
		"empty": {
			proto:    &Faas{},
			expected: &modeljson.FAAS{},
		},
		"full": {
			proto: &Faas{
				Id:               "id",
				ColdStart:        &coldstart,
				Execution:        "execution",
				TriggerType:      "triggertype",
				TriggerRequestId: "triggerrequestid",
				Name:             "name",
				Version:          "version",
			},
			expected: &modeljson.FAAS{
				ID:        "id",
				Coldstart: &coldstart,
				Execution: "execution",
				Trigger: modeljson.FAASTrigger{
					Type:      "triggertype",
					RequestID: "triggerrequestid",
				},
				Name:    "name",
				Version: "version",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.FAAS
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
