package dahua

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type DaHuaGateway struct {
	daHuaMap sync.Map
}

func NewDaHuaGateway() *DaHuaGateway {
	return &DaHuaGateway{}
}

func NewClient(conf Config) *Client {
	c := &Client{Config: conf}
	c.clinet = resty.New().
		SetBaseURL(conf.Host + ":" + conf.Port).
		// SetDigestAuth(conf.Username, conf.Password).
		SetDebug(true).
		SetTimeout(5 * time.Second)
	return c
}

func (d *DaHuaGateway) AddClient(key string, conf Config) {
	d.daHuaMap.Store(key, NewClient(conf))
}

func (d *DaHuaGateway) DeleteClient(key string) {
	d.daHuaMap.Delete(key)
}

func (d *DaHuaGateway) ToolgateInfo(key string, q TollgateInfoReq) (err error) {
	client, ok := d.daHuaMap.Load(key)
	if !ok {
		return fmt.Errorf("client not found")
	}
	err = client.(*Client).ToolgateInfo(q)
	return
}

func (d *DaHuaGateway) DeviceInfo(key string) (err error) {
	client, ok := d.daHuaMap.Load(key)
	if !ok {
		return fmt.Errorf("client not found")
	}
	err = client.(*Client).DeviceInfo()
	return
}

func (d *DaHuaGateway) KeepAlive(key string) (err error) {
	client, ok := d.daHuaMap.Load(key)
	if !ok {
		return fmt.Errorf("client not found")
	}
	err = client.(*Client).Keeplive()
	return
}
