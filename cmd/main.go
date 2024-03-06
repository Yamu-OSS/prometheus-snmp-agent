package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yamu-OSS/snmp-agent/internal/config"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/client"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/log"
	"github.com/Yamu-OSS/snmp-agent/internal/pkg/snmpd"
)

var (
	buildTime, gitCommitID, buildTag string
	version                          = flag.Bool("V", false, "show version")
	configFile                       = flag.String("c", "snmp-agent.toml", "specify the path of config file")
	pprof                            = flag.Bool("pprof", false, "use pprof")
)

func main() {
	log.InitFlags()
	flag.Parse()

	if *version {
		fmt.Println("build tag: ", buildTag)
		fmt.Println("build date: ", buildTime)
		fmt.Println("git commit: ", gitCommitID)
		return
	}

	conf := config.Init()
	err := conf.LoadFromToml(*configFile)
	if err != nil {
		log.Error(err, "failed to load config file")
		return
	}

	if *pprof {
		go func() {
			log.Info("start pprof")
			_ = http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	clients := client.NewManager()
	clients.Set(conf)
	snmpd := snmpd.New(conf, clients)
	snmpd.AsyncScan()

	var (
		hup       = make(chan os.Signal, 1)
		interrupt = make(chan os.Signal, 1)
	)
	signal.Notify(hup, syscall.SIGHUP)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	fmt.Println("snmp-agent starting")

	for {
		select {
		case <-hup:
			err := conf.Reload(*configFile)
			if err != nil {
				log.Info("snmp-agent reload fail")
				break
			}

			clients.Set(conf)

			err = snmpd.Reload()
			if err != nil {
				log.Info("snmp-agent reload fail")
				break
			}
			log.Info("snmp-agent reload success")
		case <-interrupt:
			log.Info("snmp-agent stopping")
			return
		}
	}
}
