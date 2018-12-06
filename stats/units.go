// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stats

import (
	"go.opencensus.io/metric/metricdata"
)

// Units are encoded according to the case-sensitive abbreviations from the
// Unified Code for Units of Measure: http://unitsofmeasure.org/ucum.html
const (
	UnitDimensionless = string(metricdata.UnitDimensionless)
	UnitBytes         = string(metricdata.UnitBytes)
	UnitMilliseconds  = string(metricdata.UnitMilliseconds)
	UnitNone          = UnitDimensionless // Deprecated: Use metric.UnitDimesionless.
)
