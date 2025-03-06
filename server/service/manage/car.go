package manage

import (
	"fmt"
	"server/global"
	commonReq "server/model/common/request"
	mangeModel "server/model/manage"
	mangeReq "server/model/manage/request"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type CarService struct{}

// 创建区域
func (as *CarService) CreateCar(req mangeReq.AddCar) (carModel mangeModel.CarModel, err error) {
	if err = global.TD27_DB.Where("car_num = ?", req.CarNum).First(&carModel).Error; err == nil {
		return carModel, fmt.Errorf("车牌号已存在")
	}

	var areas []*mangeModel.AreaModel
	if err = global.TD27_DB.Where("id in (?)", req.AreaIDs).Find(&areas).Error; err != nil {
		return carModel, fmt.Errorf("区域不存在")
	}

	carModel = mangeModel.CarModel{
		Name:      req.Name,
		Remark:    req.Remark,
		Phone:     req.Phone,
		CarNum:    req.CarNum,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Areas:     areas,
		CarType:   req.CarType,
		ListType:  req.ListType,
		Color:     req.Color,
		CardNo:    req.CarType,
	}
	err = global.TD27_DB.Create(&carModel).Error
	return
}

// 更新区域
func (as *CarService) EditCar(req mangeReq.EditCar) (carModel mangeModel.CarModel, err error) {
	if err = global.TD27_DB.Where("id = ?", req.ID).First(&carModel).Error; err != nil {
		return carModel, fmt.Errorf("车牌号不存在")
	}

	var areas []*mangeModel.AreaModel
	if err = global.TD27_DB.Where("id in (?)", req.AreaIDs).Find(&areas).Error; err != nil {
		return carModel, fmt.Errorf("区域不存在")
	}

	carModel.Name = req.Name
	carModel.Remark = req.Remark
	carModel.CarNum = req.CarNum
	carModel.StartTime = req.StartTime
	carModel.EndTime = req.EndTime
	carModel.Phone = req.Phone
	carModel.Areas = areas
	carModel.CarType = req.CarType
	carModel.ListType = req.ListType
	carModel.Color = req.Color
	carModel.CardNo = req.CardNo
	global.TD27_DB.Model(&carModel).Association("Areas").Replace(areas)
	err = global.TD27_DB.Save(&carModel).Error
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
	userInfo, err := utils.GetUserInfo(c)
	if err != nil {
		return
	}

	db := global.TD27_DB.Model(&mangeModel.CarModel{}).Preload("Areas")

	if userInfo.RoleId != 1 {
		areaQuery := global.TD27_DB.Table("user_area").Select("area_model_id").Where("user_model_id = ?", userInfo.ID)
		carQuery := global.TD27_DB.Table("car_area").Select("car_model_id").Where("area_model_id in (?)", areaQuery)
		db = db.Where("id in (?)", carQuery)
	}

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
		areaQuery := global.TD27_DB.Table("car_area").Select("car_model_id").Where("area_model_id = ?", req.AreaID)
		db = db.Where("id in (?)", areaQuery)
	}

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	err = db.Count(&total).Limit(limit).Offset(offset).Find(&list).Error
	return
}
