package manage

import (
	"server/global"
	mangeModel "server/model/manage"
	mangeReq "server/model/manage/request"

	"github.com/gin-gonic/gin"
)

type CarLogService struct{}

func (l *CarLogService) GetCarLogList(c *gin.Context, req mangeReq.SearchCarLog) (list []mangeModel.CarLogModel, total int64, err error) {
	db := global.TD27_DB.Model(&mangeModel.CarLogModel{})

	if req.CarNum != "" {
		db = db.Where("car_num = ?", req.CarNum)
	}

	if req.Color != "" {
		db = db.Where("color = ?", req.Color)
	}

	if req.AreaId != 0 {
		areaDeviceID := global.TD27_DB.Model(&mangeModel.DeviceModel{}).Where("area_id = ?", req.AreaId).Select("id")
		db = db.Where("device_id in (?)", areaDeviceID)
	}

	// 按照创建时间倒序

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	err = db.Preload("Device.Area").Order("created_at desc").Count(&total).Limit(limit).Offset(offset).Find(&list).Error

	return
}
