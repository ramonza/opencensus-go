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

package metric

import (
	"go.opencensus.io/resource"
	"strings"
	"testing"
	"time"
)

func TestNewLogExporter(t *testing.T) {
	r := NewRegistry()
	le := NewLogExporterWithRegistry(r)
	r.AddProducer(producerFunc(func() []*Metric {
		return []*Metric{
			{
				Descriptor: &Descriptor{
					Unit:        UnitBytes,
					Type:        TypeCumulativeDistribution,
					LabelKeys:   []string{"k1"},
					Description: "Test metric",
					Name:        "m1",
				},
				Resource: &resource.Resource{
					Type: "test",
					Labels: map[string]string{
						"zone": "a1",
					},
				},
				TimeSeries: []*TimeSeries{
					{
						Points: []Point{
							NewInt64Point(time.Time{}, 1),
						},
						LabelValues: []LabelValue{
							NewLabelValue("v1"),
						},
					},
				},
			},
		}
	}))

	ch := make(chan []interface{})
	le.Logger = logToChan(ch)
	le.ReportingPeriod = 100 * time.Millisecond

	go le.Run()

	args := <-ch

	line := args[0].(string)

	if !strings.Contains(line, "m1") {
		t.Errorf("Should include metric name")
	}
	if !strings.Contains(line, "Test metric") {
		t.Errorf("Should include metric description")
	}
	if !strings.Contains(line, "k1") {
		t.Errorf("Should include label keys")
	}
	if !strings.Contains(line, "v1") {
		t.Errorf("Should include label values")
	}
}

type logToChan chan []interface{}

func (ch logToChan) Println(vals ...interface{}) {
	ch <- vals
}

type producerFunc func() []*Metric

func (tp producerFunc) Read() []*Metric {
	return tp()
}
