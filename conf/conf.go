package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/discovery/lib/http"
	"github.com/tsingson/discovery/model"
)

var (
	confPath      string
	schedulerPath string
	region        string
	zone          string
	deployEnv     string
	hostname      string
	// Conf conf
	Conf = &Config{}
)

// func init() {
// 	var err error
// 	if hostname, err = os.Hostname(); err != nil || hostname == "" {
// 		hostname = os.Getenv("HOSTNAME")
// 	}
// 	flag.StringVar(&confPath, "conf", "discovery-example.toml", "config path")
// 	flag.StringVar(&region, "region", os.Getenv("REGION"), "avaliable region. or use REGION env variable, value: sh etc.")
// 	flag.StringVar(&zone, "zone", os.Getenv("ZONE"), "avaliable zone. or use ZONE env variable, value: sh001/sh002 etc.")
// 	flag.StringVar(&deployEnv, "deploy.env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
// 	flag.StringVar(&hostname, "hostname", hostname, "machine hostname")
// 	flag.StringVar(&schedulerPath, "scheduler", "scheduler.json", "scheduler info")
// }

// Config config.
type Config struct {
	Nodes      []string
	Zones      map[string][]string
	HTTPServer *ServerConfig
	HTTPClient *http.ClientConfig
	Env        *Env
	// Scheduler  []byte
	Schedulers map[string]*model.Scheduler
}

// Fix fix env config.
func (c *Config) Fix() (err error) {
	if c.Env == nil {
		c.Env = new(Env)
	}
	if c.Env.Region == "" {
		c.Env.Region = region
	}
	if c.Env.Zone == "" {
		c.Env.Zone = zone
	}
	if c.Env.Host == "" {
		c.Env.Host = hostname
	}
	if c.Env.DeployEnv == "" {
		c.Env.DeployEnv = deployEnv
	}
	return
}

// Env is disocvery env.
type Env struct {
	Region    string
	Zone      string
	Host      string
	DeployEnv string
}

// ServerConfig Http Servers conf.
type ServerConfig struct {
	Addr string
}

// Init init conf
func Init() (err error) {
	if _, err = toml.DecodeFile(confPath, &Conf); err != nil {
		return
	}
	if schedulerPath != "" {
		// Conf.Scheduler, _ = ioutil.ReadFile(schedulerPath)
	}
	return Conf.Fix()
}

// ConfigWalther watch configuration file change or not
func ConfigWalther(fp string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(fp)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
