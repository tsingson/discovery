package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/sanity-io/litter"
	log "github.com/tsingson/zaplogger"
	"gopkg.in/yaml.v2"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/discovery"
	"github.com/tsingson/discovery/http"
)

func main() {
	// var err error
	cfg = conf.Default()

	runtime.MemProfileRate = 0
	runtime.GOMAXPROCS(128)
	// stopSignal := make(chan struct{})
	/**


	var cntxt = &daemon.Context{
		PidFileName: "pid-discoveryd",
		PidFilePerm: 0644,
		LogFileName: logPath + "/discoveryd-daemon.log",
		LogFilePerm: 0640,
		WorkDir:     path,
		Umask:       027,
		// 	Args:        []string{"aaa-demo"},
	}

	var d, err = cntxt.Reborn()
	if err != nil {
		log.Fatal("cat's reborn ", zap.Error(err))
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	*/
	log.Info("trying to start daemon")

	svr, cancel := discovery.New(cfg)

	http.Init(cfg, svr)

	// 	runtime.Goexit()
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
	// <- stopSignal
}

func loadYaml(fh string) (cfg *conf.Config, err error) {
	fmt.Println("--------------------------->", fh)
	filename, _ := filepath.Abs(fh)
	var yamlFile []byte
	yamlFile, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("------------- file not exists")
		return
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	return

}

func writeYaml(data interface{}, fh string) error {
	fmt.Println("--------------------------->", fh)
	s, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println("marshall error ")
		return err
	}
	litter.Dump(string(s))

	return ioutil.WriteFile(fh, s, 0644)

}
