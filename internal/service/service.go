package service

import (
	"github.com/prometheus/prometheus/model/labels"
	"golang.org/x/exp/slices"

	"github.com/Yamu-OSS/snmp-agent/internal/config"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/exporter"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/log"
)

type Service struct {
	conf     *config.Config
	exporter *exporter.Exporter
}

func New(
	conf *config.Config,
	exporter *exporter.Exporter,
) *Service {
	return &Service{
		conf:     conf,
		exporter: exporter,
	}
}

func (s *Service) Get() (map[string][]*exporter.Metric, error) {
	metrics, err := s.exporter.GetMetrics()
	if err != nil {
		log.Error(err, "failed to GetMetrics")
		return nil, err
	}

	return metrics, nil
}

func (s *Service) Filter(metrics map[string][]*exporter.Metric, item *config.Item) ([]*exporter.Metric, error) {
	labs := append([]string{labels.MetricName, item.Name}, item.Labels...)
	curLabels := labels.FromStrings(labs...)

	res := make([]*exporter.Metric, 0)
	for _, metric := range metrics[item.Name] {
		if Contains(metric.Labels, curLabels) {
			res = append(res, metric)
		}
	}

	return res, nil
}

// Contains s contains v
func Contains(s, v labels.Labels) bool {
	for _, vv := range v {
		if !slices.Contains[labels.Label](s, vv) {
			return false
		}
	}
	return true
}
