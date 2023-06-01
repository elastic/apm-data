package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestMessageToModelJSON(t *testing.T) {
	millis := int64(1)

	testCases := map[string]struct {
		proto    *Message
		expected *modeljson.Message
	}{
		"empty": {
			proto:    &Message{},
			expected: &modeljson.Message{},
		},
		"full": {
			proto: &Message{
				Body: "body",
				Headers: map[string]*HTTPHeaderValue{
					"foo": {
						Values: []string{"bar"},
					},
				},
				AgeMillis:  &millis,
				QueueName:  "queuename",
				RoutingKey: "routingkey",
			},
			expected: &modeljson.Message{
				Body: "body",
				Headers: map[string][]string{
					"foo": {"bar"},
				},
				Age: modeljson.MessageAge{
					Millis: &millis,
				},
				Queue: modeljson.MessageQueue{
					Name: "queuename",
				},
				RoutingKey: "routingkey",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Message
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
