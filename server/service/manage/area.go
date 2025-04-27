package manage

import (
	"encoding/csv"
	"fmt"
	"server/global"
	"server/model/authority"
	commonReq "server/model/common/request"
	mangeModel "server/model/manage"
	mangeReq "server/model/manage/request"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

// 导出区域下的excel
func (as *AreaService) ExportHKAreaExcel(c *gin.Context, id uint) (err error) {
	var areaModel mangeModel.AreaModel
	if err = global.TD27_DB.Where("id = ?", id).First(&areaModel).Error; err != nil {
		return fmt.Errorf("区域不存在")
	}

	var carModel []mangeModel.CarModel
	carIdsQuery := global.TD27_DB.
		Table("car_device").
		Unscoped().
		Where("device_model_id in (?)", global.TD27_DB.Model(&mangeModel.DeviceModel{}).Where("area_id = ?", id).Select("id")).
		Select("car_model_id")
	global.TD27_DB.Where("id in (?)", carIdsQuery).
		Find(&carModel)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置标头
	// 车牌号码 车牌颜色 车牌类型 开始时间 结束时间
	f.SetCellValue("Sheet1", "A1", "车牌号码")
	f.SetCellValue("Sheet1", "B1", "车牌颜色")
	f.SetCellValue("Sheet1", "C1", "车牌类型")
	f.SetCellValue("Sheet1", "D1", "开始时间")
	f.SetCellValue("Sheet1", "E1", "结束时间")
	// 设置单元格格式
	// 便利数据
	for i, v := range carModel {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), v.CarNum)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), mangeModel.TCG205EPlateColor(v.Color))
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), mangeModel.HK_ExcelFormatType(v.CarType))
		fmt.Println(v.StartTime, v.EndTime)
		fmt.Println(time.UnixMilli(v.StartTime).Format(time.DateTime), time.UnixMilli(v.EndTime).Format(time.DateTime))
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), time.UnixMilli(v.StartTime).Format(time.DateTime))
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), time.UnixMilli(v.EndTime).Format(time.DateTime))
	}
	// 设置当前活动的工作表
	f.SetActiveSheet(index)
	// 设置文件名
	fileName := fmt.Sprintf("海康_区域_%s_车牌数据.xlsx", areaModel.Name)
	// 返回给前端
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("File-Name", fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	// c.Header("Content-Length", fmt.Sprintf("%d", len(f.)))
	// 获取文件长度
	if err := f.Write(c.Writer); err != nil {
		fmt.Println(err)
		c.String(500, "导出失败")
	}
	// c.Writer.Flush()

	return
}

func (as *AreaService) ExportDHAreaCsv(c *gin.Context, id uint) (err error) {
	var areaModel mangeModel.AreaModel
	if err = global.TD27_DB.Where("id = ?", id).First(&areaModel).Error; err != nil {
		return fmt.Errorf("区域不存在")
	}

	var carModel []mangeModel.CarModel
	carIdsQuery := global.TD27_DB.
		Table("car_device").
		Unscoped().
		Where("device_model_id in (?)", global.TD27_DB.Model(&mangeModel.DeviceModel{}).Where("area_id = ?", id).Select("id")).
		Select("car_model_id")
	global.TD27_DB.Where("id in (?)", carIdsQuery).
		Find(&carModel)
	// 返回csv文件
	fileName := fmt.Sprintf("大华_区域_%s_车牌数据.xlsx", areaModel.Name)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Cache-Control", "no-cache")
	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"开始时间", "结束时间", "车主姓名", "车牌号"})
	for _, v := range carModel {
		writer.Write([]string{
			time.UnixMilli(v.StartTime).Format("2006/01/02 15:04:05"),
			time.UnixMilli(v.EndTime).Format("2006/01/02 15:04:05"),
			v.Name,
			v.CarNum,
		})
	}
	writer.Flush()

	return
}
