package manage

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type CarRouter struct{}

func (a CarRouter) InitCarRouter(Router *gin.RouterGroup) {
	carRouter := Router.Group("car")
	carApi := api.ApiGroupApp.Manage.CarApi
	{
		carRouter.POST("getCarList", carApi.GetCarList)
		carRouter.POST("createCar", carApi.CreateCar)
		carRouter.POST("editCar", carApi.EditCar)
		carRouter.POST("deleteCar", carApi.DeleteCar)
		carRouter.GET("getSelectCar", carApi.GetSelectCar)
	}
}
