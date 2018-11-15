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
	"fmt"
	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
	"go.opencensus.io/trace"
	"log"
	"sync"
	"time"
)

// PushExporter exports metrics at regular intervals.
// Use Init to initialize a new PushExporter metric exporter.
// Do not modify fields of PushExporter after calling Run.
type PushExporter struct {
	push      PushFunc
	quitChan  chan struct{}
	flushChan chan *sync.WaitGroup
	wg        sync.WaitGroup

	// ProducerSet is the set of producers that will be exported.
	// If unset, the default producer set is used.
	ProducerSet *metricexport.ProducerSet

	// ReportingPeriod is the interval between exports. A new export will
	// be started every ReportingPeriod, even if the previous one is not
	// yet completed.
	// ReportingPeriod should be set higher than Timeout.
	ReportingPeriod time.Duration

	// OnError may be provided to customize the logging of errors returned
	// from push calls. By default, the error message is logged to the
	// standard error stream.
	OnError func(err error)

	// Timeout that will be applied to the context for each push call.
	Timeout time.Duration

	// NewContext may be set to override the default context that will be
	// used for exporting.
	// By default, this is set to context.Background.
	NewContext func() context.Context
}

// PushFunc is a function that exports metrics, for example to a monitoring
// backend.
type PushFunc func(context.Context, []*metricdata.Metric) error

const (
	defaultReportingPeriod = 10 * time.Second
	defaultTimeout         = 5 * time.Second
)

// NewPushExporter creates a new PushExporter metrics exporter that reads
// metrics from the given registry and pushes them to the given PushFunc.
func NewPushExporter(pushFunc PushFunc) *PushExporter {
	pe := &PushExporter{
		push:            pushFunc,
		quitChan:        make(chan struct{}),
		flushChan:       make(chan *sync.WaitGroup),
		ProducerSet:     metricexport.DefaultProducerSet(),
		ReportingPeriod: defaultReportingPeriod,
		OnError: func(err error) {
			log.Printf("Error exporting metrics: %s", err)
		},
		Timeout:    defaultTimeout,
		NewContext: context.Background,
	}
	pe.wg.Add(1)
	return pe
}

// Run exports metrics periodically.
// Run should only be called once, and returns when Stop is called.
func (p *PushExporter) Run() {
	ticker := time.NewTicker(p.ReportingPeriod)
	defer func() {
		ticker.Stop()
		p.wg.Done()
	}()
	for {
		select {
		case <-ticker.C:
			go p.Export()
		case flushed := <-p.flushChan:
			p.Export()
			flushed.Done()
		case <-p.quitChan:
			p.Export()
			return
		}
	}
}

// Export reads all metrics from the registry and exports them.
// Most users will rely on Run, which calls Export in a loop until stopped.
func (p *PushExporter) Export() {
	ms := p.ProducerSet.ReadAll()
	if len(ms) == 0 {
		return
	}

	ctx, done := context.WithTimeout(p.NewContext(), p.Timeout)
	defer done()

	// Create a Span that is never sampled to avoid sending many uninteresting
	// traces.
	ctx, span := trace.StartSpan(
		ctx,
		"go.opencensus.io/metric/exporter.Export",
		trace.WithSampler(trace.ProbabilitySampler(0.0)),
	)
	defer span.End()

	defer func() {
		e := recover()
		if e != nil {
			p.OnError(fmt.Errorf("PushFunc panic: %s", e))
		}
	}()

	err := p.push(ctx, ms)
	if err != nil {
		p.OnError(err)
	}
}

// Flush causes this push exporter to immediately read and push the current
// values of all metrics.
func (p *PushExporter) Flush() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	p.flushChan <- wg
	wg.Wait()
}

// Stop causes Run to return after exporting metrics one last time.
// Only call stop after Run has been called.
// Stop may only be called once.
func (p *PushExporter) Stop() {
	close(p.quitChan)
	p.wg.Wait()
}
