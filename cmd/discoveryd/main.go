package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/discovery/lib/file"
	"github.com/tsingson/discovery/model"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/discovery"
	"github.com/tsingson/discovery/http"
)

func main() {

	var cfg = conf.Conf
	path, _ := file.GetCurrentExecDir()
	path = "/Users/qinshen/git/linksmart/bin"
	configToml := path + "/discoveryd-config.toml"

	if _, err := toml.DecodeFile(configToml, &cfg); err != nil {
		log.Info("done")
	}

	// litter.Dump(cfg)

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
		Region:          "sh1",
		AppID:           "infra.discovery",
		Hostname:        "test-host",
		Zone:            "sh1",
		Env:             "dev",
		Status:          1,
		Addrs:           []string{"http://127.0.0.1:7171"},
		LatestTimestamp: time.Now().UnixNano(),
	}
}

