package registry

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/model"

	"golang.org/x/sync/errgroup"
)

// Nodes is helper to manage lifecycle of a collection of Nodes.
type Nodes struct {
	nodes    []*Node
	zones    map[string][]*Node
	selfAddr string
}

// NewNodes new nodes and return.
func NewNodes(cfg *conf.Config) *Nodes {
	nodes := make([]*Node, 0, len(cfg.Nodes))
	for _, addr := range cfg.Nodes {
		n := newNode(cfg, addr)
		n.zone = cfg.Env.Zone
		n.pRegisterURL = fmt.Sprintf("http://%s%s", cfg.HTTPServer.Addr, _registerURL)
		nodes = append(nodes, n)
	}
	zones := make(map[string][]*Node)
	for name, addrs := range cfg.Zones {
		var znodes []*Node
		for _, addr := range addrs {
			n := newNode(cfg, addr)
			n.otherZone = true
			n.zone = name
			n.pRegisterURL = fmt.Sprintf("http://%s%s", cfg.HTTPServer.Addr, _registerURL)
			znodes = append(znodes, n)
		}
		zones[name] = znodes
	}
	return &Nodes{
		nodes:    nodes,
		zones:    zones,
		selfAddr: cfg.HTTPServer.Addr,
	}
}

// Replicate replicate information to all nodes except for this node.
func (ns *Nodes) Replicate(c context.Context, action model.Action, i *model.Instance, otherZone bool) (err error) {
	if len(ns.nodes) == 0 {
		return
	}
	eg, c := errgroup.WithContext(c)
	for _, n := range ns.nodes {
		if !ns.Myself(n.addr) {
			ns.action(c, eg, action, n, i)
		}
	}
	if !otherZone {
		for _, zns := range ns.zones {
			if n := len(zns); n > 0 {
				ns.action(c, eg, action, zns[rand.Intn(n)], i)
			}
		}
	}
	err = eg.Wait()
	return
}

func (ns *Nodes) action(c context.Context, eg *errgroup.Group, action model.Action, n *Node, i *model.Instance) {
	switch action {
	case model.Register:
		eg.Go(func() error {
			_ = n.Register(c, i)
			return nil
		})
	case model.Renew:
		eg.Go(func() error {
			_ = n.Renew(c, i)
			return nil
		})
	case model.Cancel:
		eg.Go(func() error {
			_ = n.Cancel(c, i)
			return nil
		})
	}
}

// Nodes returns nodes of local zone.
func (ns *Nodes) Nodes() (nsi []*model.Node) {
	nsi = make([]*model.Node, 0, len(ns.nodes))
	for _, nd := range ns.nodes {
		if nd.otherZone {
			continue
		}
		node := &model.Node{
			Addr:   nd.addr,
			Status: nd.status,
			Zone:   nd.zone,
		}
		nsi = append(nsi, node)
	}
	return
}

// AllNodes returns nodes contain other zone nodes.
func (ns *Nodes) AllNodes() (nsi []*model.Node) {
	nsi = make([]*model.Node, 0, len(ns.nodes))
	for _, nd := range ns.nodes {
		node := &model.Node{
			Addr:   nd.addr,
			Status: nd.status,
			Zone:   nd.zone,
		}
		nsi = append(nsi, node)
	}
	for _, zns := range ns.zones {
		if n := len(zns); n > 0 {
			nd := zns[rand.Intn(n)]
			node := &model.Node{
				Addr:   nd.addr,
				Status: nd.status,
				Zone:   nd.zone,
			}
			nsi = append(nsi, node)
		}
	}
	return
}

// Myself returns whether or not myself.
func (ns *Nodes) Myself(addr string) bool {
	return ns.selfAddr == addr
}

// UP marks status of myself node up.
func (ns *Nodes) UP() {
	for _, nd := range ns.nodes {
		if ns.Myself(nd.addr) {
			nd.status = model.NodeStatusUP
		}
	}
}
