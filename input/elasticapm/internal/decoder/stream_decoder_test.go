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

package decoder

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNDStreamReader(t *testing.T) {
	lines := []string{
		`{"key": "value1"}`,
		`{"key": "value2", "too": "long"}`,
		`{invalid-json}`,
		`{"key": "value3"}`,
	}
	expected := []struct {
		errPattern string
		out        map[string]interface{}
		isEOF      bool
		latestLine string
	}{
		{
			out:        map[string]interface{}{"key": "value1"},
			latestLine: `{"key": "value1"}`,
		},
		{
			out:        nil,
			errPattern: "Line exceeded permitted length",
			latestLine: `{"key": "value2", "t`,
		},
		{
			out:        map[string]interface{}{},
			errPattern: "data read error",
			latestLine: `{invalid-json}`,
		},
		{
			out:        map[string]interface{}{"key": "value3"},
			latestLine: `{"key": "value3"}`,
			errPattern: "EOF",
			isEOF:      true,
		},
	}
	buf := bytes.NewBufferString(strings.Join(lines, "\n"))
	n := NewNDJSONStreamDecoder(buf, 20)

	for idx, test := range expected {
		t.Run(fmt.Sprintf("%v", idx), func(t *testing.T) {
			var out map[string]interface{}
			err := n.Decode(&out)
			assert.Equal(t, test.out, out, "Failed at idx %v", idx)
			if test.errPattern == "" {
				assert.Nil(t, err)
			} else {
				require.NotNil(t, err, "Failed at idx %v", idx)
				assert.Contains(t, err.Error(), test.errPattern, "Failed at idx %v", idx)
			}
			assert.Equal(t, test.isEOF, n.IsEOF())
			if test.latestLine == "" {
				assert.Nil(t, n.LatestLine(), "Failed at idx %v", idx)
			} else {
				assert.Equal(t, []byte(test.latestLine), n.LatestLine(), "Failed at idx %v", idx)
			}
		})
	}
}

func TestNDStreamReaderReadAhead(t *testing.T) {
	lines := []string{
		`{"key":"value1"}`,
		`{"a": "b"}`,
	}
	buf := bytes.NewBufferString(strings.Join(lines, "\n"))
	n := NewNDJSONStreamDecoder(buf, 100)

	// Decode reads the next line if it hasn't been buffered already
	var out map[string]interface{}
	require.NoError(t, n.Decode(&out))
	assert.Equal(t, map[string]interface{}{"key": "value1"}, out)
	// ReadAhead buffers the next line, to be consumed by the next call to `Decode`
	var readAheadOut, decodeOut map[string]interface{}
	b, errAhead := n.ReadAhead()
	require.NoError(t, json.Unmarshal(b, &readAheadOut))
	assert.Equal(t, map[string]interface{}{"a": "b"}, readAheadOut)
	errDecode := n.Decode(&decodeOut)
	assert.Equal(t, readAheadOut, decodeOut)
	// ReadAhead and Decode return an error for EOF
	for _, err := range []error{errAhead, errDecode} {
		if assert.Error(t, err) {
			assert.Equal(t, io.EOF, err)
		}
	}
}

func TestNDStreamDecoderTrailing(t *testing.T) {
	dec := NewNDJSONStreamDecoder(strings.NewReader(`{"key":"value1"}junk`+"\n"), 100)

	var out struct{ Key string }
	require.NoError(t, dec.Decode(&out))
	assert.Equal(t, "value1", out.Key)

	dec.Reset(strings.NewReader(`{"key":"value2"}junk` + "\n"))
	err := dec.Decode(&out)
	require.NoError(t, err)
	assert.Equal(t, "value2", out.Key)
}

func BenchmarkNDStreamDecoder(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		r := strings.NewReader("")
		dec := NewNDJSONStreamDecoder(r, 100)
		for pb.Next() {
			dec.Reset(strings.NewReader("\"string\"\n"))
			var out interface{}
			if err := dec.Decode(&out); err != nil {
				b.Fatal(err)
			}
			if out != "string" {
				b.Fatalf("expected \"string\", got %v", out)
			}
		}
	})
}
