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

package metricexport

import (
	"go.opencensus.io/metric/metricdata"
	"testing"
)

func TestProducerSet_AddProducer(t *testing.T) {
	ps := &ProducerSet{}
	ps.AddProducer(&constProducer{[]*metricdata.Metric{
		{Descriptor: metricdata.Descriptor{Name: "m1"}},
	}})
	ps.AddProducer(&constProducer{[]*metricdata.Metric{
		{Descriptor: metricdata.Descriptor{Name: "m2"}},
	}})
	ms := ps.ReadAll()
	if got, want := len(ms), 2; got != want {
		t.Fatalf("len(ms) = %d; want %d", got, want)
	}
}

type constProducer struct {
	ms []*metricdata.Metric
}

func (c *constProducer) Read() []*metricdata.Metric {
	return c.ms
}
