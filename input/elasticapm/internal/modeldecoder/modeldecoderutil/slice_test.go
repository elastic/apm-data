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

package modeldecoderutil

import (
	"testing"

	"github.com/elastic/apm-data/model/modelpb"
	"github.com/stretchr/testify/assert"
)

func TestReslice(t *testing.T) {
	originalSize := 10
	s := make([]*modelpb.APMEvent, originalSize)
	for i := range s {
		s[i] = &modelpb.APMEvent{}
	}

	downsize := 4
	s = Reslice(s, downsize, nil)
	for _, e := range s {
		assert.NotNil(t, e)
	}
	assert.Equal(t, downsize, len(s))
	assert.Equal(t, originalSize, cap(s))

	upsize := 20
	s = Reslice(s, upsize, func() *modelpb.APMEvent { return &modelpb.APMEvent{} })
	assert.Equal(t, upsize, len(s))
	for _, e := range s {
		assert.NotNil(t, e)
	}
	assert.Equal(t, upsize, cap(s))

}
