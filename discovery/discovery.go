package discovery

import (
	"context"
	"sync/atomic"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/lib/http"
	"github.com/tsingson/discovery/registry"
	// 	log "github.com/tsingson/zaplogger"
)

// Discovery discovery.
type Discovery struct {
	c        *conf.Config
	client   *http.Client
	registry *registry.Registry
	nodes    atomic.Value
}

// New get a discovery.
func New(cfg *conf.Config) (d *Discovery, cancel context.CancelFunc) {
	d = &Discovery{
		c:        cfg,
		client:   http.NewClient(cfg.HTTPClient),
		registry: registry.NewRegistry(cfg),
	}
	d.nodes.Store(registry.NewNodes(cfg))
	d.syncUp()
	cancel = d.regSelf()
	go d.nodesproc()
	return
}
