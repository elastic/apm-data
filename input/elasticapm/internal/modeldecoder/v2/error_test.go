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

package v2

import (
	"net/http"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/elastic/apm-data/input/elasticapm/internal/decoder"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/modeldecodertest"
	"github.com/elastic/apm-data/model/modelpb"
)

func TestResetErrorOnRelease(t *testing.T) {
	inp := `{"error":{"id":"tr-a"}}`
	root := fetchErrorRoot()
	require.NoError(t, decoder.NewJSONDecoder(strings.NewReader(inp)).Decode(root))
	require.True(t, root.IsSet())
	releaseErrorRoot(root)
	assert.False(t, root.IsSet())
}

func TestDecodeNestedError(t *testing.T) {
	t.Run("decode", func(t *testing.T) {
		now := time.Now().UTC()
		defaultVal := modeldecodertest.DefaultValues()
		_, eventBase := initializedInputMetadata(defaultVal)
		eventBase.Timestamp = timestamppb.New(now)
		input := modeldecoder.Input{Base: eventBase}
		str := `{"error":{"id":"a-b-c","timestamp":1599996822281000,"log":{"message":"abc"}}}`
		dec := decoder.NewJSONDecoder(strings.NewReader(str))
		var batch modelpb.Batch
		require.NoError(t, DecodeNestedError(dec, &input, &batch))
		require.Len(t, batch, 1)
		require.NotNil(t, batch[0].Error)
		assert.Equal(t, time.Unix(1599996822, 281000000).UTC(), batch[0].Timestamp.AsTime())
		assert.Empty(t, cmp.Diff(&modelpb.Error{
			Id:  "a-b-c",
			Log: &modelpb.ErrorLog{Message: "abc"},
		}, batch[0].Error, protocmp.Transform()))

		str = `{"error":{"id":"a-b-c","log":{"message":"abc"},"context":{"experimental":"exp"}}}`
		dec = decoder.NewJSONDecoder(strings.NewReader(str))
		batch = modelpb.Batch{}
		require.NoError(t, DecodeNestedError(dec, &input, &batch))
		// if no timestamp is provided, leave base event time unmodified
		assert.Equal(t, now, batch[0].Timestamp.AsTime())

		err := DecodeNestedError(decoder.NewJSONDecoder(strings.NewReader(`malformed`)), &input, &batch)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "decode")
	})

	t.Run("validate", func(t *testing.T) {
		var out modelpb.Batch
		err := DecodeNestedError(decoder.NewJSONDecoder(strings.NewReader(`{}`)), &modeldecoder.Input{}, &out)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation")
	})
}

func TestDecodeMapToErrorModel(t *testing.T) {
	gatewayIP := netip.MustParseAddr("192.168.0.1")
	randomIP := netip.MustParseAddr("71.0.54.1")

	exceptions := func(key string) bool { return false }
	t.Run("metadata-overwrite", func(t *testing.T) {
		// overwrite defined metadata with event metadata values
		var input errorEvent
		_, out := initializedInputMetadata(modeldecodertest.DefaultValues())
		otherVal := modeldecodertest.NonDefaultValues()
		modeldecodertest.SetStructValues(&input, otherVal)
		mapToErrorModel(&input, out)
		input.Reset()

		// ensure event Metadata are updated where expected
		userAgent := strings.Join(otherVal.HTTPHeader.Values("User-Agent"), ", ")
		assert.Equal(t, userAgent, out.UserAgent.Original)

		// do not overwrite client.ip if already set in metadata
		ip := modeldecodertest.DefaultValues().IP
		assert.Equal(t, ip.String(), out.GetClient().GetIp())
		assert.Equal(t, modelpb.Labels{
			"init0": {Global: true, Value: "init"}, "init1": {Global: true, Value: "init"}, "init2": {Global: true, Value: "init"},
			"overwritten0": {Value: "overwritten"}, "overwritten1": {Value: "overwritten"},
		}, modelpb.Labels(out.Labels))
		// service and user values should be set
		modeldecodertest.AssertStructValues(t, &out.Service, exceptions, otherVal)
		modeldecodertest.AssertStructValues(t, &out.User, exceptions, otherVal)
	})

	t.Run("client-ip-header", func(t *testing.T) {
		var input errorEvent
		var out modelpb.APMEvent
		input.Context.Request.Headers.Set(http.Header{})
		input.Context.Request.Headers.Val.Add("x-real-ip", gatewayIP.String())
		input.Context.Request.Socket.RemoteAddress.Set(randomIP.String())
		mapToErrorModel(&input, &out)
		assert.Equal(t, gatewayIP.String(), out.GetClient().GetIp())
	})

	t.Run("client-ip-socket", func(t *testing.T) {
		var input errorEvent
		var out modelpb.APMEvent
		input.Context.Request.Socket.RemoteAddress.Set(randomIP.String())
		mapToErrorModel(&input, &out)
		assert.Equal(t, randomIP.String(), out.GetClient().GetIp())
	})

	t.Run("error-values", func(t *testing.T) {
		exceptions := func(key string) bool {
			for _, s := range []string{
				// GroupingKey is set by a model processor
				"grouping_key",
				// StackTrace is only set by processor/otel
				"stack_trace",
				// stacktrace original and sourcemap values are set when sourcemapping is applied
				"exception.stacktrace.original",
				"exception.stacktrace.sourcemap",
				"log.stacktrace.original",
				"log.stacktrace.sourcemap",
				// ExcludeFromGrouping is set when processing the event
				"exception.stacktrace.exclude_from_grouping",
				"log.stacktrace.exclude_from_grouping",
				// Message and Type are only set for ECS compatible log event type
				"message",
				"type",
			} {
				if strings.HasPrefix(key, s) {
					return true
				}
			}
			return false
		}
		var input errorEvent
		var out1, out2 modelpb.APMEvent
		defaultVal := modeldecodertest.DefaultValues()
		modeldecodertest.SetStructValues(&input, defaultVal)
		mapToErrorModel(&input, &out1)
		input.Reset()
		modeldecodertest.AssertStructValues(t, out1.Error, exceptions, defaultVal)

		// reuse input model for different event
		// ensure memory is not shared by reusing input model
		otherVal := modeldecodertest.NonDefaultValues()
		modeldecodertest.SetStructValues(&input, otherVal)
		mapToErrorModel(&input, &out2)
		modeldecodertest.AssertStructValues(t, out2.Error, exceptions, otherVal)
		modeldecodertest.AssertStructValues(t, out1.Error, exceptions, defaultVal)
	})

	t.Run("http-headers", func(t *testing.T) {
		var input errorEvent
		input.Context.Request.Headers.Set(http.Header{"a": []string{"b"}, "c": []string{"d", "e"}})
		input.Context.Response.Headers.Set(http.Header{"f": []string{"g"}})
		var out modelpb.APMEvent
		mapToErrorModel(&input, &out)
		assert.Equal(t, map[string]any{"a": []any{"b"}, "c": []any{"d", "e"}}, out.Http.Request.Headers.AsMap())
		assert.Equal(t, map[string]any{"f": []any{"g"}}, out.Http.Response.Headers.AsMap())
	})

	t.Run("page.URL", func(t *testing.T) {
		var input errorEvent
		input.Context.Page.URL.Set("https://my.site.test:9201")
		var out modelpb.APMEvent
		mapToErrorModel(&input, &out)
		assert.Equal(t, "https://my.site.test:9201", out.Url.Full)
	})

	t.Run("page.referer", func(t *testing.T) {
		var input errorEvent
		input.Context.Page.Referer.Set("https://my.site.test:9201")
		var out modelpb.APMEvent
		mapToErrorModel(&input, &out)
		assert.Equal(t, "https://my.site.test:9201", out.Http.Request.Referrer)
	})

	t.Run("exception-code", func(t *testing.T) {
		var input errorEvent
		var out modelpb.APMEvent
		input.Exception.Code.Set(123.456)
		mapToErrorModel(&input, &out)
		assert.Equal(t, "123", out.Error.Exception.Code)
	})

	t.Run("transaction-name", func(t *testing.T) {
		var input errorEvent
		var out modelpb.APMEvent
		input.Transaction.Name.Set("My Transaction")
		mapToErrorModel(&input, &out)
		assert.Equal(t, "My Transaction", out.Transaction.Name)
	})
	t.Run("transaction-id-empty-transaction", func(t *testing.T) {
		var input errorEvent
		var out modelpb.APMEvent
		input.TransactionID.Set("12341231")
		mapToErrorModel(&input, &out)
		assert.Equal(t, "12341231", out.Transaction.Id)
	})
}
