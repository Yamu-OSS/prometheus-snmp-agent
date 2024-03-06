package config

import (
	"sync"

	"github.com/BurntSushi/toml"

	"github.com/Yamu-OSS/snmp-agent/internal/pkg/log"
)

const (
	OriginExporter   = "exporter"
	OriginPrometheus = "prometheus"
)

type Config struct {
	mutex sync.Mutex

	Snmpd struct {
		Address           string `toml:"address"`
		Timeout           int    `toml:"timeout"`
		ReconnectInterval int    `toml:"reconnect_interval"`
		BaseOID           string `toml:"base_oid"`
		ScanInterval      int    `toml:"scan_interval"`
	}
	Exporters []*Exporter
}
type Exporter struct {
	Origin string
	Server string
	Items  []*Item
}

type Item struct {
	OID        string   `toml:"oid"`
	Name       string   `toml:"name"`
	Labels     []string `toml:"labels"`
	ValueLabel string   `toml:"value_label"`
	Query      string   `toml:"query"`
	ValueType  string   `toml:"value_type"`
	DataType   string   `toml:"data_type"`
	Table      []*Table `toml:"table"`
	List       []*List  `toml:"list"`
}

type Table struct {
	OID       string `toml:"oid"`
	Label     string `toml:"label"`
	ValueType string `toml:"value_type"`
}

type List struct {
	OID       string   `toml:"oid"`
	Labels    []string `toml:"labels"`
	ValueType string   `toml:"value_type"`
}

func Init() *Config {
	return &Config{}
}

func (c *Config) LoadFromToml(filename string) error {
	_, err := toml.DecodeFile(filename, c)
	return err
}

func (c *Config) Reload(filename string) error {
	c.Lock()
	defer c.Unlock()

	log.Info("reload config")

	if err := c.LoadFromToml(filename); err != nil {
		log.Error(err, "failed to reload config")
		return err
	}
	return nil
}

func (c *Config) Lock() {
	c.mutex.Lock()
}

func (c *Config) Unlock() {
	c.mutex.Unlock()
}
