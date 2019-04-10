package registry

import (
	"encoding/json"
	"sync"

	"github.com/sanity-io/litter"
	"golang.org/x/xerrors"

	"github.com/tsingson/discovery/model"

	log "github.com/tsingson/zaplogger"
)

// Scheduler info.
type scheduler struct {
	schedulers map[string]*model.Scheduler
	mutex      sync.RWMutex
	r          *Registry
}

func newScheduler(r *Registry) *scheduler {
	return &scheduler{
		schedulers: make(map[string]*model.Scheduler),
		r:          r,
	}
}

// Load load scheduler info.
func (s *scheduler) Load(conf []byte) {
	schs := make([]*model.Scheduler, 0)
	err := json.Unmarshal(conf, &schs)
	if err != nil {
		log.Errorf("load scheduler  info  err %v", err)
	}
	for _, sch := range schs {
		s.schedulers[appsKey(sch.AppID, sch.Env)] = sch
	}
}

// Load load scheduler info.
func (s *scheduler) Build(schs map[string]*model.Scheduler) {
	litter.Dump(schs)

	if len(schs) == 0 {
		log.Errorf("schemuler is nil: %v", xerrors.New("schemuler is nil "))
		return
	}
	for _, sch := range schs {
		s.schedulers[appsKey(sch.AppID, sch.Env)] = sch
	}
	litter.Dump(s.schedulers)
}

// TODO:dynamic reload scheduler config.
// func (s *scheduler)Reolad(){
//
//}

// Get get scheduler info.
func (s *scheduler) Get(appid, env string) *model.Scheduler {
	s.mutex.RLock()
	sch := s.schedulers[appsKey(appid, env)]
	s.mutex.RUnlock()
	return sch
}
