package view

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.opencensus.io/exemplar"
	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

func TestWorker_Read(t *testing.T) {
	m1 := stats.Float64("m1", "", stats.UnitBytes)
	k1, _ := tag.NewKey("k1")
	v1 := &View{
		Measure:     m1,
		Name:        "v1",
		TagKeys:     []tag.Key{k1},
		Aggregation: Distribution(1, 5, 10),
		Description: "test view v1",
	}
	v2 := &View{
		Measure:     m1,
		Name:        "v2",
		Aggregation: Sum(),
		Description: "test view v2",
	}
	Register(v1, v2)

	ctx, _ := tag.New(context.Background(), tag.Upsert(k1, "k1v1"))
	stats.Record(ctx, m1.M(2.5))
	stats.Record(ctx, m1.M(15.0))

	ms := defaultWorker.Read()

	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Descriptor.Name < ms[j].Descriptor.Name
	})

	want := []*metricdata.Metric{
		{
			Descriptor: metricdata.Descriptor{
				Name:        "v1",
				Description: "test view v1",
				Unit:        metricdata.UnitBytes,
				Type:        metricdata.TypeCumulativeDistribution,
				LabelKeys:   []string{"k1"},
			},
			TimeSeries: []*metricdata.TimeSeries{
				{
					LabelValues: []metricdata.LabelValue{metricdata.NewLabelValue("k1v1")},
					Points: []metricdata.Point{
						{
							Value: &metricdata.Distribution{
								Count:                 2,
								Sum:                   17.5,
								SumOfSquaredDeviation: 78.125,
								BucketOptions: &metricdata.BucketOptions{
									Bounds: []float64{1, 5, 10},
								},
								Buckets: []metricdata.Bucket{
									{},
									{
										Count: 1,
										Exemplar: &exemplar.Exemplar{
											Value:       2.5,
											Attachments: exemplar.Attachments{"tag:k1": "k1v1"},
										},
									},
									{},
									{
										Count: 1,
										Exemplar: &exemplar.Exemplar{
											Value:       15.0,
											Attachments: exemplar.Attachments{"tag:k1": "k1v1"},
										},
									},
								},
							},
						},
					},
				},
			},
			Resource: nil,
		},
		{
			Descriptor: metricdata.Descriptor{
				Name:        "v2",
				Description: "test view v2",
				Unit:        metricdata.UnitBytes,
				Type:        metricdata.TypeCumulativeFloat64,
			},
			TimeSeries: []*metricdata.TimeSeries{
				{Points: []metricdata.Point{{Value: 17.5}}},
			},
			Resource: nil,
		},
	}

	if diff := cmp.Diff(ms, want, cmp.Comparer(func(t1, t2 time.Time) bool {
		return true
	})); diff != "" {
		t.Fatalf("unexpected results -got, +want: %s", diff)
	}
}
