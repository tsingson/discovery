package discovery

import (
	"context"
	"sync/atomic"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/lib/xhttp"
	"github.com/tsingson/discovery/registry"
)

// Discovery discovery.
type Discovery struct {
	c        *conf.Config
	client   *xhttp.Client
	registry *registry.Registry
	nodes    atomic.Value
}

// New get a discovery.
func New(c *conf.Config) (d *Discovery, cancelFunc context.CancelFunc) {
	d = &Discovery{
		c:        c,
		client:   xhttp.NewClient(c.HTTPClient),
		registry: registry.NewRegistry(c),
	}
	d.nodes.Store(registry.NewNodes(c))
	d.syncUp()
	cancelFunc = d.regSelf()
	go d.nodesproc()
	return
}
