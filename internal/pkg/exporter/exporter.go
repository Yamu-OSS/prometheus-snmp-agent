package exporter

import (
	"errors"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/textparse"

	"github.com/Yamu-OSS/snmp-agent/internal/pkg/log"
)

type Exporter struct {
	Client *resty.Client
}

func New(addr string) *Exporter {
	return &Exporter{
		Client: resty.New().SetBaseURL(addr),
	}
}

func (c *Exporter) QueryMetric() ([]byte, error) {
	res, err := c.Client.R().Get("")
	if err != nil {
		return res.Body(), err
	}

	return res.Body(), nil
}

func (c *Exporter) ParseMetric(metric []byte) (map[string][]*Metric, error) {
	metrics := make(map[string][]*Metric)

	metric = append(metric, []byte("# EOF")...)

	p := textparse.NewOpenMetricsParser(metric)

	var res labels.Labels
	for {
		et, err := p.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			continue
		}

		switch et {
		case textparse.EntrySeries:
			m, t, v := p.Series()
			p.Metric(&res)
			if len(res) < 1 {
				log.Error(errors.New("parse fail"), "failed to parse metric!", "line", string(m))
				continue
			}
			if _, ok := metrics[res[0].Value]; !ok {
				metrics[res[0].Value] = make([]*Metric, 0)
			}
			metrics[res[0].Value] = append(metrics[res[0].Value], &Metric{
				Value:     fmt.Sprintf("%f", v),
				TimeStamp: t,
				Labels:    res,
			})
		case textparse.EntryType, textparse.EntryHelp, textparse.EntryUnit, textparse.EntryComment:
		}
	}

	return metrics, nil
}

// GetMetrics query and parse metrics
func (c *Exporter) GetMetrics() (map[string][]*Metric, error) {
	m, err := c.QueryMetric()
	if err != nil {
		return nil, err

	}

	return c.ParseMetric(m)
}
