package dh

import "sync"

func NewGateway() *Gateway {
	return &Gateway{
		Device: sync.Map{},
	}
}

func (g *Gateway) AddDevice(id uint, cfg Config) {
	c := New(cfg)
	g.Device.Store(id, c)
}

func (g *Gateway) GetDevice(id uint) (c *Client, ok bool) {
	if v, ok := g.Device.Load(id); ok {
		c = v.(*Client)
	}
	return
}

func (g *Gateway) DeleteDevice(id uint) {
	g.Device.Delete(id)
}

func (g *Gateway) Insert(id uint, car Car) (err error) {
	g.Device.Range(func(key, value any) bool {
		if key.(uint) == id {
			c := value.(*Client)
			err = c.Insert(car)
			return false
		}
		return true
	})
	return
}

func (g *Gateway) Update(id uint, car Car) (err error) {
	g.Device.Range(func(key, value any) bool {
		if key.(uint) == id {
			c := value.(*Client)
			err = c.Update(car)
			return false
		}
		return true
	})
	return
}

func (g *Gateway) Delete(id uint, car DeleteCar) (err error) {
	g.Device.Range(func(key, value any) bool {
		if key.(uint) == id {
			c := value.(*Client)
			err = c.Delete(car)
			return false
		}
		return true
	})
	return
}

func (g *Gateway) GetCar(id uint, req GetCarReq) (result *GetCarRes, err error) {
	g.Device.Range(func(key, value any) bool {
		if key.(uint) == id {
			c := value.(*Client)
			result, err = c.GetCar(req)
			return false
		}
		return true
	})
	return
}
