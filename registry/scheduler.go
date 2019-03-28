package registry

import (
	"encoding/json"
	"sync"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/model"

	log "github.com/tsingson/zaplogger"
)

// Scheduler info.
type Scheduler struct {
	Schedulers map[string]*model.Scheduler
	mutex      sync.RWMutex
	r          *Registry
}

// NewScheduler  build a new Scheduler
func NewScheduler(r *Registry) *Scheduler {
	return &Scheduler{
		Schedulers: make(map[string]*model.Scheduler),
		r:          r,
	}
}

// Load load Scheduler info.
func (s *Scheduler) Load(conf []byte) {
	schs := make([]*model.Scheduler, 0)
	err := json.Unmarshal(conf, &schs)
	if err != nil {
		log.Errorf("load Scheduler  info  err %v", err)
	}
	for _, sch := range schs {
		s.Schedulers[appsKey(sch.AppID, sch.Env)] = sch
	}
}

func (s *Scheduler) LoadConfig(cfg *conf.Config) error {
	if len(cfg.Schedulers) == 0 {
		return xerrors.New("no schedulers")
	}
	for _, sch := range cfg.Schedulers {
		s.Schedulers[appsKey(sch.AppID, sch.Env)] = sch
	}
	return nil
}

// LoadToml load Scheduler info.
func (s *Scheduler) LoadToml(tomlFile string) error {
	schs := make([]*model.Scheduler, 0)
	if _, err := toml.DecodeFile(tomlFile, &schs); err != nil {
		return err

	}
	for _, sch := range schs {
		s.Schedulers[appsKey(sch.AppID, sch.Env)] = sch
	}
	return nil
}

// TODO:dynamic reload Scheduler config.
// func (s *Scheduler)Reolad(){
//
// }

// Get get Scheduler info.
func (s *Scheduler) Get(appid, env string) *model.Scheduler {
	s.mutex.RLock()
	sch := s.Schedulers[appsKey(appid, env)]
	s.mutex.RUnlock()
	return sch
}
