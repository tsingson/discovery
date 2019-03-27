package discovery

import (
	"context"
	"sync/atomic"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/lib/http"
	"github.com/tsingson/discovery/registry"
)

// Discovery discovery.
type Discovery struct {
	c        *conf.Config
	client   *http.Client
	registry *registry.Registry
	nodes    atomic.Value
}

// New get a discovery.
func New(c *conf.Config) (d *Discovery, cancel context.CancelFunc) {
	d = &Discovery{
		c:        c,
		client:   http.NewClient(c.HTTPClient),
		registry: registry.NewRegistry(c),
	}
	d.nodes.Store(registry.NewNodes(c))
	d.syncUp()
	cancel = d.regSelf()
	go d.nodesproc()
	return
}
