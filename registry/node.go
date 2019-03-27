package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/errors"
	"github.com/tsingson/discovery/lib/http"
	"github.com/tsingson/discovery/model"

	log "github.com/tsingson/zaplogger"
)

const (
	_registerURL = "/discovery/register"
	_cancelURL   = "/discovery/cancel"
	_renewURL    = "/discovery/renew"
	_setURL      = "/discovery/set"
)

// Node represents a peer node to which information should be shared from this node.
//
// This struct handles replicating all update operations like 'Register,Renew,Cancel,Expiration and Status Changes'
// to the <Discovery Server> node it represents.
type Node struct {
	c *conf.Config

	// url
	client       *http.Client
	pRegisterURL string
	registerURL  string
	cancelURL    string
	renewURL     string
	setURL       string

	addr      string
	status    model.NodeStatus
	zone      string
	otherZone bool
}

// newNode return a node.
func newNode(c *conf.Config, addr string) (n *Node) {
	n = &Node{
		c: c,
		// url
		client:      http.NewClient(c.HTTPClient),
		registerURL: fmt.Sprintf("http://%s%s", addr, _registerURL),
		cancelURL:   fmt.Sprintf("http://%s%s", addr, _cancelURL),
		renewURL:    fmt.Sprintf("http://%s%s", addr, _renewURL),
		setURL:      fmt.Sprintf("http://%s%s", addr, _setURL),

		addr:   addr,
		status: model.NodeStatusLost,
	}
	return
}

// Register send the registration information of Instance receiving by this node to the peer node represented.
func (n *Node) Register(c context.Context, i *model.Instance) (err error) {
	err = n.call(c, model.Register, i, n.registerURL, nil)
	if err != nil {
		log.Warningf("node be called(%s) register instance(%v) error(%v)", n.registerURL, i, err)
	}
	return
}

// Cancel send the cancellation information of Instance receiving by this node to the peer node represented.
func (n *Node) Cancel(c context.Context, i *model.Instance) (err error) {
	err = n.call(c, model.Cancel, i, n.cancelURL, nil)
	if err != nil {
		log.Warningf("node be called(%s) instance(%v) already canceled", n.cancelURL, i)
	}
	return
}

// Renew send the heartbeat information of Instance receiving by this node to the peer node represented.
// If the instance does not exist the node, the instance registration information is sent again to the peer node.
func (n *Node) Renew(c context.Context, i *model.Instance) (err error) {
	var res *model.Instance
	err = n.call(c, model.Renew, i, n.renewURL, &res)
	if err == errors.ServerErr {
		log.Warningf("node be called(%s) instance(%v) error(%v)", n.renewURL, i, err)
		n.status = model.NodeStatusLost
		return
	}
	n.status = model.NodeStatusUP
	if err == errors.NothingFound {
		log.Warningf("node be called(%s) instance(%v) error(%v)", n.renewURL, i, err)
		err = n.call(c, model.Register, i, n.registerURL, nil)
		return
	}
	// NOTE: register response instance whitch in conflict with peer node
	if err == errors.Conflict && res != nil {
		err = n.call(c, model.Register, res, n.pRegisterURL, nil)
	}
	return
}

func (n *Node) call(c context.Context, action model.Action, i *model.Instance, uri string, data interface{}) (err error) {
	params := url.Values{}
	params.Set("zone", i.Zone)
	params.Set("env", i.Env)
	params.Set("appid", i.AppID)
	params.Set("hostname", i.Hostname)
	if n.otherZone {
		params.Set("replication", "false")
	} else {
		params.Set("replication", "true")
	}
	switch action {
	case model.Register:
		params.Set("addrs", strings.Join(i.Addrs, ","))
		params.Set("status", strconv.FormatUint(uint64(i.Status), 10))
		params.Set("version", i.Version)
		meta, _ := json.Marshal(i.Metadata)
		params.Set("metadata", string(meta))
		params.Set("reg_timestamp", strconv.FormatInt(i.RegTimestamp, 10))
		params.Set("dirty_timestamp", strconv.FormatInt(i.DirtyTimestamp, 10))
		params.Set("latest_timestamp", strconv.FormatInt(i.LatestTimestamp, 10))
	case model.Renew:
		params.Set("dirty_timestamp", strconv.FormatInt(i.DirtyTimestamp, 10))
	case model.Cancel:
		params.Set("latest_timestamp", strconv.FormatInt(i.LatestTimestamp, 10))
	}
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err = n.client.Post(c, uri, "", params, &res); err != nil {
		log.Errorf("node be called(%s) instance(%v) error(%v)", uri, i, err)
		return
	}
	if res.Code != 0 {
		log.Errorf("node be called(%s) instance(%v) response code(%v)", uri, i, res.Code)
		if err = errors.Int(res.Code); err == errors.Conflict {
			_ = json.Unmarshal([]byte(res.Data), data)
		}
	}
	return
}
