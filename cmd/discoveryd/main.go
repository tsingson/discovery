package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/discovery/model"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/discovery"
	"github.com/tsingson/discovery/http"
	xhttp "github.com/tsingson/discovery/lib/http"
	xtime "github.com/tsingson/discovery/lib/time"
)

func main() {

	flag.Parse()

	cfg := &conf.Config{
		Env: &conf.Env{
			Region:    "test",
			Zone:      "test",
			DeployEnv: "test",
			Host:      "test_server",
		},
		Nodes: []string{"127.0.0.1:7171"},
		HTTPServer: &conf.ServerConfig{
			Addr: "127.0.0.1:7171",
		},
		HTTPClient: &xhttp.ClientConfig{
			Dial:      xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second * 30),
		},
	}
	_ = cfg.Fix()
	svr, cancel := discovery.New(cfg)
	 	svr.Register(context.Background(), defRegDiscovery(), time.Now().UnixNano(), false)
	http.Init(cfg, svr)

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("discovery get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancel()
			time.Sleep(time.Second)
			log.Info("discovery quit !!!")
			// log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func defRegDiscovery() *model.Instance {
	return &model.Instance{
		AppID:           "infra.discovery",
		Hostname:        "test_server",
		Zone:            "test",
		Env:             "test",
		Status:          1,
		Addrs:           []string{"http://10.0.0.111:7171"},
		LatestTimestamp: time.Now().UnixNano(),
	}
}
