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

// ErrorWithResult is an error that contains additional rejected and accepted
// counts. It can be used by ProcessBatch to signal partial failure.
// This should be used sparingly as it will be changed when we move away from
// ProcessBatch and process events individually in the future.
type ErrorWithResult struct {
	error
	Result
}

type Result struct {
	Rejected int64
	Accepted int64
}

func parseResultFromError(err error, totalCount int64) (Result, error) {
	if err == nil {
		return Result{Accepted: totalCount}, nil
	}
	if errorWithCounts, ok := err.(ErrorWithResult); ok {
		return errorWithCounts.Result, errorWithCounts.error
	}
	return Result{Rejected: totalCount}, err
}
