package manage

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type CarLogRouter struct{}

func (a CarLogRouter) InitCarLogRouter(Router *gin.RouterGroup) {
	carLogRouter := Router.Group("carLog")
	carLogApi := api.ApiGroupApp.Manage.CarLogApi
	{
		carLogRouter.POST("getCarLogList", carLogApi.GetCarLogList)
	}
}
