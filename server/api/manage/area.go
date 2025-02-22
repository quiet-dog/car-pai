package manage

import (
	"server/global"
	commonReq "server/model/common/request"
	commonRes "server/model/common/response"
	manageReq "server/model/manage/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AreaApi struct{}

// GetAreaList
// @Tags      AreaApi
// @Summary   获取区域列表
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]manage.AreaModel,msg=string}
// @Router    /area/getAreaList [post]
func (a *AreaApi) GetAreaList(c *gin.Context) {
	var areaPageInfo manageReq.SearchArea
	if err := c.ShouldBindJSON(&areaPageInfo); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if list, total, err := areaService.GetAreaList(areaPageInfo); err != nil {
		commonRes.FailWithMessage("获取失败", c)
		global.TD27_LOG.Error("获取失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(commonRes.PageResult{
			Page:     areaPageInfo.Page,
			PageSize: areaPageInfo.PageSize,
			List:     list,
			Total:    total,
		}, "获取成功", c)
	}
}

// CreateArea
// @Tags      AreaApi
// @Summary   获取区域列表
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]manage.AreaModel,msg=string}
// @Router    /area/createArea [post]
func (a *AreaApi) CreateArea(c *gin.Context) {
	var addArea manageReq.AddArea
	if err := c.ShouldBindJSON(&addArea); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if area, err := areaService.CreateArea(addArea); err != nil {
		commonRes.FailWithMessage("创建失败", c)
		global.TD27_LOG.Error("创建失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(area, "创建成功", c)
	}
}

func (a *AreaApi) EditArea(c *gin.Context) {
	var editArea manageReq.EditArea
	if err := c.ShouldBindJSON(&editArea); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if area, err := areaService.EditArea(editArea); err != nil {
		commonRes.FailWithMessage("编辑失败", c)
		global.TD27_LOG.Error("编辑失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(area, "编辑成功", c)
	}
}

func (a *AreaApi) DeleteArea(c *gin.Context) {
	var cId commonReq.CId
	if err := c.ShouldBindJSON(&cId); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if err := areaService.DeleteArea(cId); err != nil {
		commonRes.FailWithMessage("删除失败", c)
		global.TD27_LOG.Error("删除失败", zap.Error(err))
	} else {
		commonRes.OkWithMessage("删除成功", c)
	}
}
