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

package pushexporter

import (
	"context"
	"strings"
	"testing"
	"time"

	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
)

func TestPushExporter_Run(t *testing.T) {
	exported := make(chan bool, 1)
	pe := NewPushExporter(func(ctx context.Context, ms []*metricdata.Metric) error {
		_, ok := ctx.Deadline()
		if !ok {
			t.Fatal("Expected a deadline")
		}
		select {
		case exported <- true:
		default:
		}
		return nil
	})
	pe.ProducerSet = &metricexport.ProducerSet{}
	pe.ProducerSet.AddProducer(&constProducer{&metricdata.Metric{}})
	pe.ReportingPeriod = 100 * time.Millisecond

	go pe.Run()
	defer pe.Stop()

	select {
	case _ = <-exported:
	case <-time.After(1 * time.Second):
		t.Fatal("PushFunc should have been called")
	}
}

func TestPushExporter_Run_panic(t *testing.T) {
	errs := make(chan error, 1)
	pe := NewPushExporter(func(ctx context.Context, ms []*metricdata.Metric) error {
		panic("test")
	})
	pe.ProducerSet = &metricexport.ProducerSet{}
	pe.ProducerSet.AddProducer(&constProducer{&metricdata.Metric{}})
	pe.ReportingPeriod = 100 * time.Millisecond
	pe.OnError = func(err error) {
		errs <- err
	}

	go pe.Run()
	defer pe.Stop()

	select {
	case err := <-errs:
		if !strings.Contains(err.Error(), "test") {
			t.Error("Should contain the panic arg")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("OnError should be called")
	}
}

func TestPushExporter_Stop(t *testing.T) {
	exported := make(chan bool, 1)
	pe := NewPushExporter(func(ctx context.Context, ms []*metricdata.Metric) error {
		select {
		case exported <- true:
		default:
			t.Fatal("Export should only be called once")
		}
		return nil
	})
	pe.ProducerSet = &metricexport.ProducerSet{}
	pe.ProducerSet.AddProducer(&constProducer{&metricdata.Metric{}})
	pe.ReportingPeriod = time.Hour // prevent timer-based push

	go pe.Run()
	pe.Stop()

	select {
	case _ = <-exported:
	default:
		t.Fatal("PushFunc should have been called before Stop returns")
	}
}

type constProducer []*metricdata.Metric

func (cp constProducer) Read() []*metricdata.Metric {
	return cp
}
