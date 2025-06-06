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

package otlp

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/elastic/apm-data/input"
)

func semAcquire(ctx context.Context, tracer trace.Tracer, sem input.Semaphore, i int64) error {
	ctx, sp := tracer.Start(ctx, "Semaphore.Acquire")
	defer sp.End()

	return sem.Acquire(ctx, i)
}
