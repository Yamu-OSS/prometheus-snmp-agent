package exporter

import "github.com/prometheus/prometheus/model/labels"

type Metric struct {
	Value     string
	TimeStamp *int64
	Labels    labels.Labels
}
