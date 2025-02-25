package manage

import (
	"fmt"
	"server/global"
	"server/model/authority"
	commonReq "server/model/common/request"
	mangeModel "server/model/manage"
	mangeReq "server/model/manage/request"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AreaService struct{}

// 创建区域
func (as *AreaService) CreateArea(req mangeReq.AddArea) (areaModel mangeModel.AreaModel, err error) {
	var users []*authority.UserModel
	global.TD27_DB.Where("id in (?)", req.UserIDs).Find(&users)
	areaModel = mangeModel.AreaModel{
		Name:   req.Name,
		Remark: req.Remark,
		Users:  users,
	}
	err = global.TD27_DB.Create(&areaModel).Error
	return
}

// 更新区域
func (as *AreaService) EditArea(req mangeReq.EditArea) (areaModel mangeModel.AreaModel, err error) {
	if err = global.TD27_DB.Where("id = ?", req.ID).First(&areaModel).Error; err != nil {
		return areaModel, fmt.Errorf("区域不存在")
	}
	var users []*authority.UserModel
	global.TD27_DB.Where("id in (?)", req.UserIDs).Find(&users)

	areaModel.Name = req.Name
	areaModel.Remark = req.Remark
	areaModel.Users = users
	err = global.TD27_DB.Save(&areaModel).Error
	global.TD27_DB.Model(&areaModel).Association("Users").Replace(users)
	return
}

// 删除区域
func (as *AreaService) DeleteArea(req commonReq.CId) (err error) {
	var areaModel mangeModel.AreaModel
	global.TD27_DB.Where("id = ?", req.ID).First(&areaModel)

	err = global.TD27_DB.Delete(&areaModel).Error
	return
}

// 获取区域列表
func (as *AreaService) GetAreaList(c *gin.Context, req mangeReq.SearchArea) (list []mangeModel.AreaModel, total int64, err error) {
	userInfo, err := utils.GetUserInfo(c)
	if err != nil {
		return
	}

	db := global.TD27_DB.Model(&mangeModel.AreaModel{})

	if userInfo.RoleId != 1 {
		db = db.Where("id in (?)", global.TD27_DB.Table("user_area").Select("area_model_id").Where("user_model_id = ?", userInfo.ID))
	}

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+req.Remark+"%")
	}

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	err = db.Count(&total).Limit(limit).Offset(offset).Find(&list).Error
	return
}
