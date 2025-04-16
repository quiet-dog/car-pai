package manage

import (
	"server/global"
	commonReq "server/model/common/request"
	commonRes "server/model/common/response"
	manageReq "server/model/manage/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CarApi struct{}

// GetCarList
// @Tags      CarApi
// @Summary   获取车牌号列表
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]manage.CarModel,msg=string}
// @Router    /car/getCarList [post]
func (a *CarApi) GetCarList(c *gin.Context) {
	var carPageInfo manageReq.SearchCar
	if err := c.ShouldBindJSON(&carPageInfo); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if list, total, err := carService.GetCarList(c, carPageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("获取失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(commonRes.PageResult{
			Page:     carPageInfo.Page,
			PageSize: carPageInfo.PageSize,
			List:     list,
			Total:    total,
		}, "获取成功", c)
	}
}

// CreateCar
// @Tags      CarApi
// @Summary   获取车牌号列表
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]manage.CarModel,msg=string}
// @Router    /car/createCar [post]
func (a *CarApi) CreateCar(c *gin.Context) {
	var addCar manageReq.AddCar
	if err := c.ShouldBindJSON(&addCar); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if car, err := carService.CreateCar(addCar); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("创建失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(car, "创建成功", c)
	}
}

func (a *CarApi) EditCar(c *gin.Context) {
	var editCar manageReq.EditCar
	if err := c.ShouldBindJSON(&editCar); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if car, err := carService.EditCar(editCar); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error(err.Error(), zap.Error(err))
	} else {
		commonRes.OkWithDetailed(car, "编辑成功", c)
	}
}

func (a *CarApi) DeleteCar(c *gin.Context) {
	var cId commonReq.CId
	if err := c.ShouldBindJSON(&cId); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if err := carService.DeleteCar(cId); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("删除失败", zap.Error(err))
	} else {
		commonRes.OkWithMessage("删除成功", c)
	}
}

func (a *CarApi) GetSelectCar(c *gin.Context) {
	if list, err := carService.GetCarSelect(c); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("获取失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(list, "获取成功", c)
	}
}
