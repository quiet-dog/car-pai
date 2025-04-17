package initialize

import (
	"server/global"
	"server/model/manage"
	"server/pkg/hk_gateway"
	"strconv"
)

func InitHikGateway() {
	global.HikGateway = hk_gateway.NewHikGateway()

	var hikInfo = []*manage.DeviceModel{}
	global.TD27_DB.Where("type = ?", manage.HIK).Find(&hikInfo)
	for _, v := range hikInfo {
		port, err := strconv.Atoi(v.Port)
		if err != nil {
			continue
		}
		global.HikGateway.RegisterHikGateway(hk_gateway.HikConfig{
			Ip:       v.Host,
			Port:     port,
			Username: v.HikUsername,
			Password: v.HikPassword,
		})
	}
}
