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

package codec

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/elastic/apm-data/model/modelpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func randomKv(t testing.TB) ([]*modelpb.KeyValue, map[string]any) {
	m := map[string]any{
		t.Name() + ".key." + randString(): t.Name() + ".value." + randString(),
	}

	kv := []*modelpb.KeyValue{}
	for k, v := range m {
		value, _ := structpb.NewValue(v)
		kv = append(kv, &modelpb.KeyValue{
			Key:   k,
			Value: value,
		})
	}

	return kv, m
}

func randomKvPb(t testing.TB) []*modelpb.KeyValue {
	k, _ := randomKv(t)
	return k
}

func randomHTTPHeaders(t testing.TB) []*modelpb.HTTPHeader {
	return []*modelpb.HTTPHeader{
		&modelpb.HTTPHeader{
			Key:   t.Name() + ".key." + randString(),
			Value: []string{t.Name() + ".value." + randString()},
		},
	}
}

func uintPtr(i uint32) *uint32 {
	return &i
}

func uint64Ptr(i uint64) *uint64 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString() string {
	size := 5
	var sb strings.Builder
	sb.Grow(size)
	for i := 0; i < size; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}
