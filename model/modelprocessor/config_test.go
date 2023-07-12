// ELASTICSEARCH CONFIDENTIAL
// __________________
//
//  Copyright Elasticsearch B.V. All rights reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Elasticsearch B.V. and its suppliers, if any.
// The intellectual and technical concepts contained herein
// are proprietary to Elasticsearch B.V. and its suppliers and
// may be covered by U.S. and Foreign Patents, patents in
// process, and are protected by trade secret or copyright
// law.  Dissemination of this information or reproduction of
// this material is strictly forbidden unless prior written
// permission is obtained from Elasticsearch B.V.

package modelprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestConfigOptions(t *testing.T) {
	tp := trace.NewTracerProvider()

	for _, tt := range []struct {
		name           string
		opts           []ConfigOption
		expectedConfig config
	}{
		{
			name: "with no option",
			expectedConfig: config{
				TracerProvider: otel.GetTracerProvider(),
			},
		},
		{
			name: "a custom tracer provider",
			opts: []ConfigOption{
				WithTracerProvider(tp),
			},
			expectedConfig: config{
				TracerProvider: tp,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(tt.opts...)
			assert.Equal(t, tt.expectedConfig, cfg)
		})
	}
}
