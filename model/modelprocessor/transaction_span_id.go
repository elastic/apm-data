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

package modelprocessor

import (
	"context"

	"github.com/elastic/apm-data/model/modelpb"
)

type SetSpanIdTransaction struct{}

// ProcessBatch copies transaction.id to span.id
func (SetSpanIdTransaction) ProcessBatch(ctx context.Context, b *modelpb.Batch) error {
	for i := range *b {
		setSpanIdTransaction((*b)[i])
	}
	return nil
}

func setSpanIdTransaction(event *modelpb.APMEvent) {
	if event.Transaction != nil && event.GetSpan().GetId() == "" {

		if event.Span == nil {
			event.Span = &modelpb.Span{}
		}
		event.Span.Id = event.Transaction.Id

	}
}
