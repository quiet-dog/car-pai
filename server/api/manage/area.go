package manage

import (
	"server/global"
	commonReq "server/model/common/request"
	commonRes "server/model/common/response"
	manageReq "server/model/manage/request"
	"strconv"

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

	if list, total, err := areaService.GetAreaList(c, areaPageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
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
		commonRes.FailWithMessage(err.Error(), c)
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
		commonRes.FailWithMessage(err.Error(), c)
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
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("删除失败", zap.Error(err))
	} else {
		commonRes.OkWithMessage("删除成功", c)
	}
}

func (a *AreaApi) ExportAreaExcel(c *gin.Context) {
	idStr := c.Query("id")
	target := c.Query("target")
	if idStr == "" {
		commonRes.FailWithMessage("id不能为空", c)
		return
	}
	// 转换为uint
	id, err := strconv.Atoi(idStr)
	if err != nil {
		commonRes.FailWithMessage("id格式错误", c)
		return
	}
	if target == "hik" {
		areaService.ExportHKAreaExcel(c, uint(id))
	}
	if target == "dh" {
		areaService.ExportDHAreaCsv(c, uint(id))
	}
}
