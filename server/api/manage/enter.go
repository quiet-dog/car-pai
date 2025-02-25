package manage

import "server/service"

type ApiGroup struct {
	AreaApi
	DeviceApi
	CarApi
}

var (
	areaService   = service.ServiceGroupApp.Manage.AreaService
	deviceService = service.ServiceGroupApp.Manage.DeviceService
	carService    = service.ServiceGroupApp.Manage.CarService
)
