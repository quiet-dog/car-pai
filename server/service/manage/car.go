package manage

import (
	"fmt"
	"server/global"
	commonReq "server/model/common/request"
	mangeModel "server/model/manage"
	mangeReq "server/model/manage/request"
	mangeRes "server/model/manage/response"

	"github.com/gin-gonic/gin"
)

type CarService struct{}

// 创建区域
func (as *CarService) CreateCar(req mangeReq.AddCar) (carModel mangeModel.CarModel, err error) {
	if err = global.TD27_DB.Where("car_num = ? and car_type = ? and color = ?", req.CarNum, req.CarType, req.Color).First(&carModel).Error; err == nil {
		return carModel, fmt.Errorf("车牌号已存在")
	}

	carModel = mangeModel.CarModel{
		Name:      req.Name,
		Remark:    req.Remark,
		Phone:     req.Phone,
		CarNum:    req.CarNum,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		// Devices:   devices,
		CarType:   req.CarType,
		ListType:  req.ListType,
		Color:     req.Color,
		CardNo:    req.CardNo,
		DeviceIds: req.DeviceIDs,
	}

	err = global.TD27_DB.Omit("Devices").Omit("DeviceIds").Create(&carModel).Error

	// areaCarQuery := global.TD27_DB.Table("car_device").Where("car_model_id = ?", carModel.ID).Select("device_model_id")
	// var deviceModel []*mangeModel.DeviceModel
	// global.TD27_DB.Where("id in (?)", areaCarQuery).Find(&deviceModel)

	// for _, v := range deviceModel {
	// 	if v.Type == mangeModel.HIK {
	// 		if v.Model == mangeModel.HIK_DS_TCG225 || v.Model == mangeModel.HIK_DS_2CD9125_KS || v.Model == mangeModel.HIK_DS_TCG2A5_E {
	// 			data := hk_gateway.SetVCLDataReq{}
	// 			vclDataList := hk_gateway.VCLDataList{}
	// 			singlieVCLData := hk_gateway.SingleVCLData{}
	// 			singlieVCLData.RunNum = "0"

	// 			singlieVCLData.ListType = req.ListType
	// 			singlieVCLData.PlateNum = req.CarNum
	// 			singlieVCLData.PlateColor = req.Color
	// 			singlieVCLData.PlateType = req.CarType
	// 			singlieVCLData.CardNo = req.CardNo

	// 			// startTime
	// 			singlieVCLData.StartTime = time.UnixMilli(req.StartTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
	// 			singlieVCLData.EndTime = time.UnixMilli(req.EndTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
	// 			vclDataList.SingleVCLData = append(vclDataList.SingleVCLData, singlieVCLData)
	// 			data.VCLDataList = vclDataList
	// 			err = global.HikGateway.SetVCLData(v.Host, data)
	// 			if err != nil {
	// 				return carModel, err
	// 			}
	// 		}

	// 		if v.Model == mangeModel.HIK_DS_TCG2A5_B {
	// 			data := hk_gateway.SetVCLDataReq{}
	// 			vclDataList := hk_gateway.VCLDataList{}
	// 			singlieVCLData := hk_gateway.SingleVCLData{}
	// 			singlieVCLData.RunNum = "0"
	// 			singlieVCLData.ListType = req.ListType
	// 			singlieVCLData.PlateNum = req.CarNum
	// 			singlieVCLData.PlateColor = req.Color
	// 			singlieVCLData.PlateType = req.CarType
	// 			singlieVCLData.CardNo = req.CardNo
	// 			// startTime
	// 			t := time.UnixMilli(req.StartTime).UTC()
	// 			startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	// 			singlieVCLData.StartTime = startOfDay.Format(time.RFC3339)
	// 			// singlieVCLData.StartTime = time.UnixMilli(c.StartTime).UTC().Format(time.RFC3339)
	// 			t = time.UnixMilli(req.EndTime).UTC()
	// 			endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999000000, time.UTC)
	// 			singlieVCLData.EndTime = endOfDay.Format(time.RFC3339)

	// 			// 将singlieVCLData.EndTime的时间转为最后的时间23:59:59
	// 			vclDataList.SingleVCLData = append(vclDataList.SingleVCLData, singlieVCLData)
	// 			data.VCLDataList = vclDataList
	// 			err = global.HikGateway.TCG2A5EVCLGetCond(v.Host, data)
	// 			if err != nil {
	// 				return carModel, err
	// 			}
	// 		}

	// 		if v.Model == mangeModel.HIK_DS_TCG205_E {
	// 			data := hk_gateway.TCG225EVCLGetCondReq{}
	// 			lic := hk_gateway.LicensePlateInfo{}
	// 			lic.PlateColor = mangeModel.TCG205EPlateColor(req.Color)
	// 			lic.PlateType = mangeModel.TCG205EPlateType(req.CarType)
	// 			lic.ListType = mangeModel.TCG205EListType(req.ListType)
	// 			lic.LicensePlate = req.CarNum
	// 			lic.CardNo = req.CardNo
	// 			lic.ID = "1"
	// 			lic.Operation = "new"
	// 			lic.CreateTime = time.UnixMilli(req.StartTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
	// 			lic.EffectiveStartDate = time.UnixMilli(req.EndTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
	// 			lic.EffectiveTime = time.UnixMilli(req.EndTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
	// 			data.LicensePlateInfo = append(data.LicensePlateInfo, lic)
	// 			err = global.HikGateway.TCG225EVCLGetCond(v.Host, data)
	// 			if err != nil {
	// 				return carModel, err
	// 			}
	// 		}

	// 	}
	// }
	// fmt.Println("carModel err", err)
	return
}

// 更新区域
func (as *CarService) EditCar(req mangeReq.EditCar) (carModel mangeModel.CarModel, err error) {
	if err = global.TD27_DB.Where("id = ?", req.ID).First(&carModel).Error; err != nil {
		return carModel, fmt.Errorf("车牌号不存在")
	}

	var devices []*mangeModel.DeviceModel
	if err = global.TD27_DB.Where("id in (?)", req.DeviceIDs).Find(&devices).Error; err != nil {
		return carModel, fmt.Errorf("区域不存在")
	}

	var deviceIds []uint
	global.TD27_DB.Table("car_device").Where("car_model_id = ? and device_model_id not in (?)", req.ID, req.DeviceIDs).Select("device_model_id").Find(&deviceIds)

	carModel.Name = req.Name
	carModel.Remark = req.Remark
	carModel.CarNum = req.CarNum
	carModel.StartTime = req.StartTime
	carModel.EndTime = req.EndTime
	carModel.Phone = req.Phone
	// carModel.Devices = devices
	carModel.CarType = req.CarType
	carModel.ListType = req.ListType
	carModel.Color = req.Color
	carModel.CardNo = req.CardNo
	carModel.DeviceIds = req.DeviceIDs
	// global.TD27_DB.Model(&carModel).Association("Devices").Replace(devices)
	err = global.TD27_DB.Omit("Devices").Omit("DeviceIds").Save(&carModel).Error
	if len(deviceIds) > 0 {
		err = carModel.DeleteByDevice(deviceIds)
	}
	return
}

// 删除区域
func (as *CarService) DeleteCar(req commonReq.CId) (err error) {
	var carModel mangeModel.CarModel
	global.TD27_DB.Where("id = ?", req.ID).First(&carModel)

	err = global.TD27_DB.Delete(&carModel).Error
	return
}

// 获取区域列表
func (as *CarService) GetCarList(c *gin.Context, req mangeReq.SearchCar) (list []mangeModel.CarModel, total int64, err error) {
	// userInfo, err := utils.GetUserInfo(c)
	// if err != nil {
	// 	return
	// }

	db := global.TD27_DB.Model(&mangeModel.CarModel{}).Preload("Devices")

	// if userInfo.RoleId != 1 {
	// 	areaQuery := global.TD27_DB.Table("user_area").Select("area_model_id").Where("user_model_id = ?", userInfo.ID)
	// 	deviceQuery := global.TD27_DB.Model(mangeModel.DeviceModel{}).Select("id").Where("area_id in (?)", areaQuery)
	// 	carQuery := global.TD27_DB.Table("car_device").Select("car_model_id").Where("device_model_id in (?)", deviceQuery)
	// 	db = db.Where("id in (?)", carQuery)
	// }

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.CarNum != "" {
		db = db.Where("car_num LIKE ?", "%"+req.CarNum+"%")
	}

	if req.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}

	if req.AreaID != 0 {
		areaQuery := global.TD27_DB.Table("device_models").Select("id").Where("area_id = ?", req.AreaID)
		deviceQuery := global.TD27_DB.Table("car_device").Select("car_model_id").Where("device_model_id in (?)", areaQuery)
		db = db.Where("id in (?)", deviceQuery)
	}

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	err = db.Count(&total).Limit(limit).Offset(offset).Find(&list).Error
	return
}

func (as *CarService) GetCarSelect(c *gin.Context) (list []*mangeRes.Select, err error) {
	var areas []*mangeModel.AreaModel
	var devices []*mangeModel.DeviceModel
	if err = global.TD27_DB.Find(&areas).Error; err != nil {
		return
	}
	if err = global.TD27_DB.Find(&devices).Error; err != nil {
		return
	}

	for _, area := range areas {
		list = append(list, &mangeRes.Select{
			Label:    area.Name,
			Value:    area.ID,
			Type:     "area",
			Children: []*mangeRes.Select{},
		})
	}

	for _, device := range devices {
		for _, area := range list {
			if area.Value == device.AreaId {
				area.Children = append(area.Children, &mangeRes.Select{
					Label: device.Host,
					Value: device.ID,
					Type:  "device",
				})
			}
		}
	}

	for _, area := range list {
		if len(area.Children) == 0 {
			area.Disabled = true
		}
	}

	return
}
