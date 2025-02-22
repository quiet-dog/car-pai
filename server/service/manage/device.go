package manage

import (
	"fmt"
	"server/global"
	commonReq "server/model/common/request"
	mangeModel "server/model/manage"
	mangeReq "server/model/manage/request"
)

type DeviceService struct{}

// 创建设备
func (as *DeviceService) CreateDevice(req mangeReq.AddDevice) (deviceModel mangeModel.DeviceModel, err error) {

	global.TD27_DB.Where("host = ?", req.Host).First(&deviceModel)
	if deviceModel.ID != 0 {
		return deviceModel, fmt.Errorf("设备已存在")
	}

	deviceModel = mangeModel.DeviceModel{
		Host:        req.Host,
		Port:        req.Port,
		HikUsername: req.HikUsername,
		HikPassword: req.HikPassword,
		DhUsername:  req.DhUsername,
		DhPassword:  req.DhPassword,
		Remark:      req.Remark,
		AreaId:      req.AreaId,
		Type:        req.Type,
	}
	err = global.TD27_DB.Create(&deviceModel).Error
	return
}

// 更新设备
func (as *DeviceService) EditDevice(req mangeReq.EditDevice) (deviceModel mangeModel.DeviceModel, err error) {
	if err = global.TD27_DB.Where("id = ?", req.ID).First(&deviceModel).Error; err != nil {
		return deviceModel, fmt.Errorf("设备不存在")
	}

	deviceModel.Host = req.Host
	deviceModel.Port = req.Port
	deviceModel.HikUsername = req.HikUsername
	deviceModel.HikPassword = req.HikPassword
	deviceModel.DhUsername = req.DhUsername
	deviceModel.DhPassword = req.DhPassword
	deviceModel.Remark = req.Remark
	deviceModel.AreaId = req.AreaId
	deviceModel.Type = req.Type
	err = global.TD27_DB.Save(&deviceModel).Error
	return
}

// 删除设备
func (as *DeviceService) DeleteDevice(req commonReq.CId) (err error) {
	var deviceModel mangeModel.DeviceModel
	global.TD27_DB.Where("id = ?", req.ID).First(&deviceModel)

	err = global.TD27_DB.Delete(&deviceModel).Error
	return
}

// 获取设备列表
func (as *DeviceService) GetDeviceList(req mangeReq.SearchDevice) (list []mangeModel.DeviceModel, total int64, err error) {
	db := global.TD27_DB.Model(&mangeModel.DeviceModel{})

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	if req.AreaId != 0 {
		db = db.Where("area_id = ?", req.AreaId)
	}

	err = db.Count(&total).Limit(limit).Offset(offset).Find(&list).Error
	return
}
