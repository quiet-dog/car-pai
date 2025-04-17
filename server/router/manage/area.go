package manage

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type AreaRouter struct{}

func (a AreaRouter) InitAreaRouter(Router *gin.RouterGroup) {
	areaRouter := Router.Group("area")
	areaApi := api.ApiGroupApp.Manage.AreaApi
	{
		areaRouter.POST("getAreaList", areaApi.GetAreaList)
		areaRouter.POST("createArea", areaApi.CreateArea)
		areaRouter.POST("editArea", areaApi.EditArea)
		areaRouter.POST("deleteArea", areaApi.DeleteArea)
		areaRouter.GET("exportAreaExcel", areaApi.ExportAreaExcel)
	}
}
