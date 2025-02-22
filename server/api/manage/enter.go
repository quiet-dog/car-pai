package manage

import "server/service"

type ApiGroup struct {
	AreaApi
	DeviceApi
}

var (
	areaService   = service.ServiceGroupApp.Manage.AreaService
	deviceService = service.ServiceGroupApp.Manage.DeviceService
)
