package exporter

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
)

func TestParseMetric(t *testing.T) {
	c := New("127.0.0.1")

	int64p := func(x int64) *int64 { return &x }

	input := `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
# UNIT go_gc_duration_seconds seconds
go_gc_duration_seconds{quantile="0"} 4.9351e-05
go_gc_duration_seconds{quantile="0.25"} 7.424100000000001e-05
go_gc_duration_seconds{quantile="0.5",a="b"} 8.3835e-05
# HELP nohelp1 
# HELP help2 escape \ \n \\ \" \x chars
# UNIT nounit 
go_gc_duration_seconds{quantile="1.0",a="b"} 8.3835e-05
go_gc_duration_seconds_count 99
some:aggregate:rate5m{a_b="c"} 1
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 33 123.123
# TYPE hh histogram
hh_bucket{le="+Inf"} 1
# TYPE gh gaugehistogram
gh_bucket{le="+Inf"} 1
# TYPE hhh histogram
hhh_bucket{le="+Inf"} 1 # {id="histogram-bucket-test"} 4
hhh_count 1 # {id="histogram-count-test"} 4
# TYPE foo counter
foo_total 17.0 1520879607.789 # {id="counter-test"} 5`

	input += "\n"

	want := map[string][]*Metric{
		"go_gc_duration_seconds": []*Metric{
			{
				Value: "0.000049",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "go_gc_duration_seconds",
					},
					{
						Name:  "quantile",
						Value: "0",
					},
				},
			},
			{
				Value: "0.000074",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "go_gc_duration_seconds",
					},
					{
						Name:  "quantile",
						Value: "0.25",
					},
				},
			},
			{
				Value: "0.000084",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "go_gc_duration_seconds",
					},
					{
						Name:  "a",
						Value: "b",
					},
					{
						Name:  "quantile",
						Value: "0.5",
					},
				},
			},
			{
				Value: "0.000084",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "go_gc_duration_seconds",
					},
					{
						Name:  "a",
						Value: "b",
					},
					{
						Name:  "quantile",
						Value: "1.0",
					},
				},
			},
		},
		"go_gc_duration_seconds_count": []*Metric{
			{
				Value: "99.000000",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "go_gc_duration_seconds_count",
					},
				},
			},
		},
		"some:aggregate:rate5m": []*Metric{
			{
				Value: "1.000000",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "some:aggregate:rate5m",
					},
					{
						Name:  "a_b",
						Value: "c",
					},
				},
			},
		},
		"go_goroutines": []*Metric{
			{
				Value:     "33.000000",
				TimeStamp: int64p(123123),
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "go_goroutines",
					},
				},
			},
		},
		"hh_bucket": []*Metric{
			{
				Value: "1.000000",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "hh_bucket",
					},
					{
						Name:  "le",
						Value: "+Inf",
					},
				},
			},
		},
		"gh_bucket": []*Metric{
			{
				Value: "1.000000",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "gh_bucket",
					},
					{
						Name:  "le",
						Value: "+Inf",
					},
				},
			},
		},
		"hhh_bucket": []*Metric{
			{
				Value: "1.000000",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "hhh_bucket",
					},
					{
						Name:  "le",
						Value: "+Inf",
					},
				},
			},
		},
		"hhh_count": []*Metric{
			{
				Value: "1.000000",
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "hhh_count",
					},
				},
			},
		},
		"foo_total": []*Metric{
			{
				Value:     "17.000000",
				TimeStamp: int64p(1520879607789),
				Labels: []labels.Label{
					{
						Name:  "__name__",
						Value: "foo_total",
					},
				},
			},
		},
	}

	get, err := c.ParseMetric([]byte(input))
	assert.NoError(t, err)
	assert.Equal(t, want, get)
}
