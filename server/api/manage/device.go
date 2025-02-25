package manage

import (
	"server/global"
	commonReq "server/model/common/request"
	commonRes "server/model/common/response"
	manageReq "server/model/manage/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceApi struct{}

func (a *DeviceApi) GetDeviceList(c *gin.Context) {
	var devicePageInfo manageReq.SearchDevice
	if err := c.ShouldBindJSON(&devicePageInfo); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if list, total, err := deviceService.GetDeviceList(devicePageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("获取失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(commonRes.PageResult{
			Page:     devicePageInfo.Page,
			PageSize: devicePageInfo.PageSize,
			List:     list,
			Total:    total,
		}, "获取成功", c)
	}
}

// CreateDevice
// @Tags      DeviceApi
// @Summary   获取区域列表
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]manage.DeviceModel,msg=string}
// @Router    /device/createDevice [post]
func (a *DeviceApi) CreateDevice(c *gin.Context) {
	var addDevice manageReq.AddDevice
	if err := c.ShouldBindJSON(&addDevice); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if device, err := deviceService.CreateDevice(addDevice); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("创建失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(device, "创建成功", c)
	}
}

func (a *DeviceApi) EditDevice(c *gin.Context) {
	var editDevice manageReq.EditDevice
	if err := c.ShouldBindJSON(&editDevice); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if device, err := deviceService.EditDevice(editDevice); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("编辑失败", zap.Error(err))
	} else {
		commonRes.OkWithDetailed(device, "编辑成功", c)
	}
}

func (a *DeviceApi) DeleteDevice(c *gin.Context) {
	var cId commonReq.CId
	if err := c.ShouldBindJSON(&cId); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if err := deviceService.DeleteDevice(cId); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		global.TD27_LOG.Error("删除失败", zap.Error(err))
	} else {
		commonRes.OkWithMessage("删除成功", c)
	}
}
