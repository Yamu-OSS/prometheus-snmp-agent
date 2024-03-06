package snmpd

import (
	"context"
	"fmt"
	"time"

	"github.com/posteo/go-agentx"
	"github.com/posteo/go-agentx/pdu"
	"github.com/posteo/go-agentx/value"
	prometheusModel "github.com/prometheus/common/model"

	"github.com/Yamu-OSS/snmp-agent/internal/config"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/agent"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/client"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/exporter"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/log"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/prometheus"
	"github.com/Yamu-OSS/snmp-agent/internal/service"
)

var (
	ErrNotRegistory = "session is not registered"
)

type SNMPD struct {
	conf        *config.Config
	agentClient *agentx.Client
	session     *agentx.Session
	handler     *agent.CommonHandler
	clients     *client.Manager
}

func New(conf *config.Config, clients *client.Manager) *SNMPD {
	s := &SNMPD{
		conf:    conf,
		handler: &agent.CommonHandler{},
		clients: clients,
	}

	s.Init(conf)

	return s
}

func (s *SNMPD) Reload() error {
	s.conf.Lock()
	defer s.conf.Unlock()

	log.Info("reload session")

	err := s.session.Close()
	if err != nil {
		log.Error(err, "failed to close session")
		return err
	}

	s.Init(s.conf)

	return nil
}

func (s *SNMPD) Init(conf *config.Config) {
	for {
		err := s.newSession(conf)
		if err != nil {
			log.Error(err, "failed to new session, try to reconnect")
			time.Sleep(time.Duration(conf.Snmpd.ReconnectInterval) * time.Second)
			continue
		}

		break
	}

	for {
		if err := s.Register(s.Scan()); err != nil {
			log.Error(err, "failed to register")
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	log.Info("init success")
}

func (s *SNMPD) newSession(conf *config.Config) error {
	if s.agentClient != nil {
		err := s.agentClient.Close()
		if err != nil {
			return err
		}
	}

	client, err := agentx.Dial("tcp", conf.Snmpd.Address)
	if err != nil {
		return err
	}
	client.Timeout = time.Duration(conf.Snmpd.Timeout) * time.Minute
	client.ReconnectInterval = time.Duration(conf.Snmpd.ReconnectInterval) * time.Second

	session, err := client.Session()
	if err != nil {
		return err
	}

	s.session = session
	s.agentClient = client

	return nil
}

// Scan scan exporter to get new data
func (s *SNMPD) Scan() *agent.CommonHandler {
	handler := &agent.CommonHandler{}

	agentItems := make(map[string]*agent.Item)
	exporterClients, prometheusClients := s.clients.Get()

	for _, v := range s.conf.Exporters {
		switch v.Origin {
		case "exporter":
			s.GetExporterData(v, agentItems, exporterClients[v.Server])
		case "prometheus":
			s.GetPromethuesData(v, agentItems, prometheusClients[v.Server])
		default:
			log.Info("Unable to identify origin")
		}
	}

	for oid, item := range agentItems {
		handler.Add(oid, item.T, item.V)
	}

	handler.Sort()

	return handler
}

// GetExporterData get data from exporter
func (s *SNMPD) GetExporterData(
	v *config.Exporter,
	agentItems map[string]*agent.Item,
	client *exporter.Exporter,
) {
	svc := service.New(s.conf, client)
	allMetrics, err := svc.Get()
	if err != nil {
		log.Error(err, "failed to Get metrics")
		return
	}

	for _, item := range v.Items {
		metrics, err := svc.Filter(allMetrics, item)
		if err != nil {
			log.Error(err, "failed to Get metrics")
			continue
		}

		if len(metrics) == 0 {
			continue
		}

		switch item.DataType {
		case SubTypeTable:
			// set metric value and label value into metricMap
			// it is useful for entry to get all table value by label name
			metricMap := make(map[string][]string)
			for i, metric := range metrics {
				if i == 0 {
					metricMap["value"] = []string{metric.Value}

					for _, label := range metric.Labels {
						metricMap[label.Name] = []string{label.Value}
					}
					continue
				}

				metricMap["value"] = append(metricMap["value"], metric.Value)

				for _, label := range metric.Labels {
					metricMap[label.Name] = append(metricMap[label.Name], label.Value)
				}
			}

			MakeTableValue(item, agentItems, metricMap)

		case SubTypeList:
			log.Info("not support SubTypeList")

		default:
			f := valueType[item.ValueType]

			if metrics[0] == nil || f == nil {
				log.Debug("Value or ValueType is nil")
				break
			}

			var t pdu.VariableType
			var v any

			if item.ValueLabel == "" {
				t, v = f(metrics[0].Value)
			} else {
				t, v = f(metrics[0].Labels.Map()[item.ValueLabel])
			}
			agentItems[item.OID] = &agent.Item{
				T: t,
				V: v,
			}
		}
	}
}

// GetPromethuesData get data from ddi-exporter
func (s *SNMPD) GetPromethuesData(
	v *config.Exporter,
	agentItems map[string]*agent.Item,
	client *prometheus.Client,
) {
	ctx := context.Background()
	for _, item := range v.Items {

		res, err := client.Query(ctx, item.Query, time.Now())
		if err != nil {
			log.Error(err, "failed to GetPromethuesData")
			continue
		}

		v, ok := res.(prometheusModel.Vector)
		if !ok {
			continue
		}

		if len(v) == 0 {
			continue
		}

		switch item.DataType {
		case SubTypeTable:
			// set metric value and label value into metricMap
			// it is useful for entry to get all table value by label name
			metricMap := make(map[string][]string)
			for i, vv := range v {
				if i == 0 {
					metricMap["value"] = []string{vv.Value.String()}

					for labelName, labelValue := range vv.Metric {
						metricMap[string(labelName)] = []string{string(labelValue)}
					}
					continue
				}

				metricMap["value"] = append(metricMap["value"], vv.Value.String())

				for labelName, labelValue := range vv.Metric {
					metricMap[string(labelName)] = append(metricMap[string(labelName)], string(labelValue))
				}
			}
			MakeTableValue(item, agentItems, metricMap)

		case SubTypeList:
			for _, list := range item.List {
				filter := make(map[string]string)
				for i := 0; i < len(list.Labels); i += 2 {
					filter[list.Labels[i]] = list.Labels[i+1]
				}

			Filter:
				for _, vv := range v {
					for key, value := range filter {
						if vv.Metric[prometheusModel.LabelName(key)] != prometheusModel.LabelValue(value) {
							continue Filter
						}
					}

					f := valueType[list.ValueType]
					t, v := f(vv.Value.String())
					agentItems[list.OID] = &agent.Item{
						T: t,
						V: v,
					}

					break
				}
			}

		default:
			f := valueType[item.ValueType]

			if v[0] == nil || f == nil {
				log.Debug("Value or ValueType is nil")
				break
			}

			t, v := f(v[0].Value.String())
			agentItems[item.OID] = &agent.Item{
				T: t,
				V: v,
			}
		}
	}
}

func MakeTableValue(item *config.Item, agentItems map[string]*agent.Item, metricMap map[string][]string) {
	for _, entry := range item.Table {
		if entry.OID == "" {
			log.Info("entry must set oid", "entry", entry.Label)
			continue
		}

		for j, leafValue := range metricMap[entry.Label] {
			leafOID := fmt.Sprintf("%s.%v", entry.OID, j+1)
			f := valueType[entry.ValueType]
			t, v := f(leafValue)
			agentItems[leafOID] = &agent.Item{
				T: t,
				V: v,
			}
		}
	}
}

// Unregister unregister from snmpd
func (s *SNMPD) Unregister() error {
	log.Debug("try to Unregister")
	return s.session.Unregister(127, value.MustParseOID(s.conf.Snmpd.BaseOID))
}

// Register register to snmpd
func (s *SNMPD) Register(handler *agent.CommonHandler) error {
	log.Debug("try to Register")
	s.session.Handler = handler
	s.handler = handler
	return s.session.Register(127, value.MustParseOID(s.conf.Snmpd.BaseOID))
}

func (s *SNMPD) AsyncScan() {
	timer := time.NewTicker(time.Second)

	f := func() {
		s.conf.Lock()
		defer s.conf.Unlock()

		timer.Reset(time.Duration(s.conf.Snmpd.ScanInterval) * time.Second)

		newHandler := s.Scan()

		s.handler.UpdateHander(newHandler)
	}

	go func() {
		for range timer.C {
			f()
		}
	}()
}
