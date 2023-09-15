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
	"encoding/json"
	"strconv"

	"github.com/elastic/apm-data/model/modelpb"
)

// GlobalLabelsFrom populates the Labels and NumericLabels from global labels
// in the metadata object.
func GlobalLabelsFrom(from map[string]any, to *modelpb.APMEvent) {
	to.NumericLabels = make(modelpb.NumericLabels)
	to.Labels = make(modelpb.Labels)
	MergeLabels(from, to)
	for k, v := range to.Labels {
		v.Global = true
		to.Labels[k] = v
	}
	for k, v := range to.NumericLabels {
		v.Global = true
		to.NumericLabels[k] = v
	}
	if len(to.Labels) == 0 {
		to.Labels = nil
	}
	if len(to.NumericLabels) == 0 {
		to.NumericLabels = nil
	}
}

// MergeLabels merges eventLabels into the APMEvent. This is used for
// combining event-specific labels onto (metadata) global labels.
//
// If eventLabels is non-nil, it is first cloned.
func MergeLabels(eventLabels map[string]any, to *modelpb.APMEvent) {
	count := 0
	numericCount := 0
	for _, v := range eventLabels {
		switch v.(type) {
		case string, bool:
			count++
		case float64, json.Number:
			numericCount++
		}
	}

	if count == 0 && numericCount == 0 {
		return
	}

	if to.NumericLabels == nil && numericCount > 0 {
		to.NumericLabels = make(modelpb.NumericLabels, numericCount)
	}
	if to.Labels == nil && count > 0 {
		to.Labels = make(modelpb.Labels, count)
	}
	for k, v := range eventLabels {
		switch v := v.(type) {
		case string:
			modelpb.Labels(to.Labels).Set(k, v)
		case bool:
			modelpb.Labels(to.Labels).Set(k, strconv.FormatBool(v))
		case float64:
			modelpb.NumericLabels(to.NumericLabels).Set(k, v)
		case json.Number:
			if floatVal, err := v.Float64(); err == nil {
				modelpb.NumericLabels(to.NumericLabels).Set(k, floatVal)
			}
		}
	}
}
