package main

import (
	"fmt"
	"time"

	"github.com/sanity-io/litter"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/lib/file"
	xhttp "github.com/tsingson/discovery/lib/http"
	xtime "github.com/tsingson/discovery/lib/time"
	"github.com/tsingson/discovery/registry"
)

func main() {

	// var cfg = conf.Conf
	// path, _ := file.GetCurrentPath()
	// configToml := path + "/discoveryd-config.toml"
	//
	// if _, err := toml.DecodeFile(configToml, &cfg); err != nil {
	// 	log.Info("done")
	// }
	// litter.Dump(cfg)
	defaultConfig()

}
func defaultConfig() {

	cfg := &conf.Config{
		Env: &conf.Env{
			Region:    "china",
			Zone:      "gd",
			Host:      "test",
			DeployEnv: "dev",
		},
		Nodes: []string{"127.0.0.1:7171"},
		HTTPServer: &conf.ServerConfig{
			Addr: "127.0.0.1:7171",
		},
		HTTPClient: &xhttp.ClientConfig{
			Dial:      xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second * 120),
		},
	}
	json := []byte(`[{"app_id":"test.service","env":"dev","zones":[{"src":"gd","dst":{"sz01":3}}],"remark":"test"}]`)
	r := registry.NewRegistry(cfg)
	s := registry.NewScheduler(r)
	s.Load(json)
	// litter.Dump(s.Schedulers)

	fmt.Println("--------------------------------------------------->")
	path, _ := file.GetCurrentPath()
	// tomlFile := path + "/schedulers.toml"
	// file.SaveToml(s.Schedulers, tomlFile)

	// if _, err := toml.DecodeFile(tomlFile, &s.Schedulers); err != nil {
	// 	fmt.Println(err)
	//
	// }

	cfg.Schedulers = s.Schedulers
	configToml := path + "/discoveryd-config.toml"

	// fmt.Println("--------------------------------------------------->")
	// s.LoadToml(configToml)
	// litter.Dump(s.Schedulers)
	fmt.Println("--------------------------------------------------->")
	fmt.Println("config file >   ", configToml)
	litter.Dump(cfg)

	file.SaveToml(cfg, configToml)
	// {
	// 	c := make(chan os.Signal, 1)
	// 	log.Info("done")
	// 	conf.ConfigWalther(configToml)
	//
	// 	<-c
	// }

}
