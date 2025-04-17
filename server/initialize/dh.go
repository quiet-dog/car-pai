package initialize

import (
	"server/global"
	"server/model/manage"
	"server/pkg/dh"
)

func InitDh() {
	global.DhGateway = dh.NewGateway()
	var dhInfo = []*manage.DeviceModel{}
	global.TD27_DB.Where("type = ?", manage.DH).Find(&dhInfo)
	for _, v := range dhInfo {
		global.DhGateway.AddDevice(v.ID, dh.Config{
			Host:     v.Host,
			Port:     v.Port,
			Username: v.DhUsername,
			Password: v.DhPassword,
		})
	}
}
