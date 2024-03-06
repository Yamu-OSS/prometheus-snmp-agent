package client

import (
	"sync"

	"github.com/Yamu-OSS/snmp-agent/internal/config"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/exporter"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/prometheus"
)

type Manager struct {
	exporterClients   map[string]*exporter.Exporter
	prometheusClients map[string]*prometheus.Client

	lock sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		exporterClients:   make(map[string]*exporter.Exporter),
		prometheusClients: make(map[string]*prometheus.Client),
	}
}

func (c *Manager) Set(conf *config.Config) {
	c.lock.Lock()
	defer c.lock.Unlock()

	newExporterClients := make(map[string]struct{})
	newPrometheusClients := make(map[string]struct{})
	for _, exp := range conf.Exporters {
		switch exp.Origin {
		case config.OriginExporter:
			newExporterClients[exp.Server] = struct{}{}
		case config.OriginPrometheus:
			newPrometheusClients[exp.Server] = struct{}{}
		}
	}

	for key := range c.exporterClients {
		if _, ok := newExporterClients[key]; !ok {
			delete(c.exporterClients, key)
		}
	}

	for key := range c.prometheusClients {
		if _, ok := newPrometheusClients[key]; !ok {
			delete(c.prometheusClients, key)
		}
	}

	for _, exp := range conf.Exporters {
		switch exp.Origin {
		case config.OriginExporter:
			if _, ok := c.exporterClients[exp.Server]; !ok {
				c.exporterClients[exp.Server] = exporter.New(exp.Server)
			}
		case config.OriginPrometheus:
			if _, ok := c.prometheusClients[exp.Server]; !ok {
				c.prometheusClients[exp.Server] = prometheus.New(exp.Server)
			}
		}
	}
}

func (c *Manager) Get() (map[string]*exporter.Exporter, map[string]*prometheus.Client) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.exporterClients, c.prometheusClients
}
