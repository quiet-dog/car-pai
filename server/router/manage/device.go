package manage

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type DeviceRouter struct{}

func (a DeviceRouter) InitDeviceRouter(Router *gin.RouterGroup) {
	deviceRouter := Router.Group("device")
	deviceApi := api.ApiGroupApp.Manage.DeviceApi
	{
		deviceRouter.POST("getDeviceList", deviceApi.GetDeviceList)
		deviceRouter.POST("createDevice", deviceApi.CreateDevice)
		deviceRouter.POST("editDevice", deviceApi.EditDevice)
		deviceRouter.POST("deleteDevice", deviceApi.DeleteDevice)
	}
}
