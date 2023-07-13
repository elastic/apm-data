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

package modelprocessor_test

import (
	"testing"

	"github.com/elastic/apm-data/model/modelpb"
	"github.com/elastic/apm-data/model/modelprocessor"
)

func TestSetSpanIdTransaction(t *testing.T) {

	processor := modelprocessor.SetSpanIdTransaction{}

	// span id already present
	transactionWithSpanId := modelpb.APMEvent{
		Transaction: &modelpb.Transaction{
			Id: "123",
		},
		Span: &modelpb.Span{
			Id: "123",
		},
	}
	testProcessBatch(t, processor, &transactionWithSpanId, &transactionWithSpanId)

	transactionWithoutSpanId := modelpb.APMEvent{
		Transaction: &modelpb.Transaction{
			Id: "123",
		},
	}
	// span id copied from transaction id
	testProcessBatch(t, processor, &transactionWithoutSpanId, &transactionWithSpanId)

}
