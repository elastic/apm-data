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

package modelpbtest

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/elastic/apm-data/model/modelpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func RandomStruct(t testing.TB) (*structpb.Struct, map[string]any) {
	m := map[string]any{
		t.Name() + ".key." + randString(): t.Name() + ".value." + randString(),
	}

	s, err := structpb.NewStruct(m)
	require.NoError(t, err)

	return s, m
}

func RandomStructPb(t testing.TB) *structpb.Struct {
	s, _ := RandomStruct(t)
	return s
}

func RandomHTTPHeaders(t testing.TB) []*modelpb.HTTPHeader {
	return []*modelpb.HTTPHeader{
		&modelpb.HTTPHeader{
			Key:   t.Name() + ".key." + randString(),
			Value: []string{t.Name() + ".value." + randString()},
		},
	}
}

func UintPtr(i uint32) *uint32 {
	return &i
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func BoolPtr(b bool) *bool {
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
