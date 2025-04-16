package manage

import (
	"server/global"
	"server/pkg/hk_gateway"
	"time"

	"gorm.io/gorm"
)

type CarModel struct {
	global.TD27_MODEL
	CarNum    string `json:"carNum" gorm:"not null;comment:车牌号" binding:"required"` // 车牌号
	Color     string `json:"color" gorm:"comment:车辆颜色"`                             // 车辆颜色
	CarType   string `json:"carType" gorm:"comment:车辆类型"`                           // 车辆类型
	ListType  string `json:"listType" gorm:"comment:名单类型"`                          // 名单类型 0白 1 黑   blockList黑 allowList白
	Name      string `json:"name" gorm:"comment:车主姓名"`                              // 车主姓名
	Phone     string `json:"phone" gorm:"comment:车主电话"`                             // 车主电话
	CardNo    string `json:"cardNo" gorm:"comment:卡号"`                              // 卡号
	StartTime int64  `json:"startTime" gorm:"not null;comment:开始时间"`                // 开始时间
	EndTime   int64  `json:"endTime" gorm:"not null;comment:结束时间"`                  // 结束时间
	// DeviceID  uint           `json:"deviceId" gorm:"not null;comment:设备ID" binding:"required"` // 设备ID
	// Device    *DeviceModel   `json:"device"`                                                   // 设备
	Remark    string         `json:"remark" gorm:"comment:备注"` // 备注
	Areas     []*AreaModel   `json:"areas" gorm:"many2many:car_area;"`
	Devices   []*DeviceModel `json:"devices" gorm:"many2many:car_device;"`
	DeviceIds []uint         `json:"deviceIds" gorm:"-"`
}

/** carType DS-TCG225 hik: 0 标准民用用车与军车
1 02式民用车牌
2 武警车
3 警车
4 民用车双行尾牌
5 使馆车牌
6 农用车牌
7 摩托车牌
8 新能源车牌

color DS-TCG225 hik: 0 蓝色
1 黄色
2 白色
3 黑色
4 绿色
5 其他
*/

/**
hik DS-TCG205-E
92TypeCivil 92式民用车牌

*/

func (c *CarModel) AfterCreate(tx *gorm.DB) (err error) {

	// 重新加载关联的 Devices
	var deviceModels []*DeviceModel
	tx.Where("id in (?)", c.DeviceIds).Find(&deviceModels)
	// 关联设备
	for _, v := range deviceModels {
		tx.Exec("delete from car_device where car_model_id = ? and device_model_id = ?", c.ID, v.ID)
		tx.Exec("insert into car_device (car_model_id, device_model_id) values (?, ?)", c.ID, v.ID)
	}

	err = c.set(tx)
	if err != nil {
		// 删除
		tx.Delete(c)
	}
	return
}

func (c *CarModel) AfterSave(tx *gorm.DB) (err error) {
	var deviceModels []*DeviceModel
	tx.Where("id in (?)", c.DeviceIds).Find(&deviceModels)
	// 关联设备
	tx.Exec("delete from car_device where car_model_id = ?", c.ID)
	for _, v := range deviceModels {
		tx.Exec("insert into car_device (car_model_id, device_model_id) values (?, ?)", c.ID, v.ID)
	}
	err = c.set(tx)

	return
}

func (c *CarModel) BeforeDelete(tx *gorm.DB) (err error) {
	deviceCarQuery := tx.Table("car_device").Where("car_model_id = ?", c.ID).Select("device_model_id")
	var deviceModel []*DeviceModel
	tx.Where("id in (?)", deviceCarQuery).Find(&deviceModel)
	for _, v := range deviceModel {

		if v.Model == HIK_DS_TCG225 ||
			v.Model == HIK_DS_TCG2A5_B ||
			v.Model == HIK_DS_2CD9125_KS ||
			v.Model == HIK_DS_TCG2A5_E {
			data := hk_gateway.VCLDelCondReq{}
			data.CardNo = c.CardNo
			data.PlateColor = c.Color
			data.PlateNum = c.CarNum
			data.PlateType = c.CarType
			data.DelVCLCond = 1
			err = global.HikGateway.VCLDelCond(v.Host, data)
			continue
		}

		if v.Model == HIK_DS_TCG205_E {
			data := hk_gateway.TCG2A5EVCLDelCondReq{}
			data.DeleteAllEnabled = false
			data.LicensePlate = []string{c.CarNum}
			err = global.HikGateway.TCG205EVCLDelCond(v.Host, data)
		}

	}
	return
}

func TCG205EPlateColor(c string) string {
	switch c {
	case "0":
		{
			return "蓝色"
		}
	case "1":
		{
			return "黄色"
		}
	case "2":
		{
			return "白色"
		}
	case "3":
		{
			return "黑色"
		}
	case "4":
		{
			return "绿色"
		}
	default:
		{
			return "其他"
		}
	}
}

func TCG205EPlateType(c string) string {
	switch c {
	case "0":
		{
			return "92TypeCivil"
		}
	case "1":
		{
			return "stdCivilAndMilitay"
		}
	case "2":
		{
			return "armedPolice"
		}
	case "3":
		{
			return "arm"
		}
	case "4":
		{
			return "civilTwoLine"
		}
	case "5":
		{
			return "embassy"
		}
	case "6":
		{
			return "farmVehicle"
		}
	case "7":
		{
			return "motorola"
		}
	case "8":
		{
			return "newEnergy"
		}
	default:
		{
			return "unknown"
		}
	}
}

func TCG205EListType(c string) string {
	if c == "0" {
		return "allowList"
	}
	return "blockList"

}

func (c *CarModel) set(tx *gorm.DB) (err error) {

	areaCarQuery := tx.Table("car_device").Where("car_model_id = ?", c.ID).Select("device_model_id")
	var deviceModel []*DeviceModel
	tx.Where("id in (?)", areaCarQuery).Find(&deviceModel)

	for _, v := range deviceModel {
		if v.Type == HIK {
			if v.Model == HIK_DS_TCG225 || v.Model == HIK_DS_2CD9125_KS || v.Model == HIK_DS_TCG2A5_E {
				data := hk_gateway.SetVCLDataReq{}
				vclDataList := hk_gateway.VCLDataList{}
				singlieVCLData := hk_gateway.SingleVCLData{}
				singlieVCLData.RunNum = "0"

				singlieVCLData.ListType = c.ListType
				singlieVCLData.PlateNum = c.CarNum
				singlieVCLData.PlateColor = c.Color
				singlieVCLData.PlateType = c.CarType
				singlieVCLData.CardNo = c.CardNo

				// startTime
				singlieVCLData.StartTime = time.UnixMilli(c.StartTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
				singlieVCLData.EndTime = time.UnixMilli(c.EndTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
				vclDataList.SingleVCLData = append(vclDataList.SingleVCLData, singlieVCLData)
				data.VCLDataList = vclDataList
				err = global.HikGateway.SetVCLData(v.Host, data)
				if err != nil {
					return err
				}
			}

			if v.Model == HIK_DS_TCG2A5_B {
				data := hk_gateway.SetVCLDataReq{}
				vclDataList := hk_gateway.VCLDataList{}
				singlieVCLData := hk_gateway.SingleVCLData{}
				singlieVCLData.RunNum = "0"
				singlieVCLData.ListType = c.ListType
				singlieVCLData.PlateNum = c.CarNum
				singlieVCLData.PlateColor = c.Color
				singlieVCLData.PlateType = c.CarType
				singlieVCLData.CardNo = c.CardNo
				// startTime
				t := time.UnixMilli(c.StartTime).UTC()
				startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
				singlieVCLData.StartTime = startOfDay.Format(time.RFC3339)
				// singlieVCLData.StartTime = time.UnixMilli(c.StartTime).UTC().Format(time.RFC3339)
				t = time.UnixMilli(c.EndTime).UTC()
				endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999000000, time.UTC)
				singlieVCLData.EndTime = endOfDay.Format(time.RFC3339)

				// 将singlieVCLData.EndTime的时间转为最后的时间23:59:59
				vclDataList.SingleVCLData = append(vclDataList.SingleVCLData, singlieVCLData)
				data.VCLDataList = vclDataList
				err = global.HikGateway.TCG2A5EVCLGetCond(v.Host, data)
				if err != nil {
					return err
				}
			}

			if v.Model == HIK_DS_TCG205_E {
				data := hk_gateway.TCG225EVCLGetCondReq{}
				lic := hk_gateway.LicensePlateInfo{}
				lic.PlateColor = TCG205EPlateColor(c.Color)
				lic.PlateType = TCG205EPlateType(c.CarType)
				lic.ListType = TCG205EListType(c.ListType)
				lic.LicensePlate = c.CarNum
				lic.CardNo = c.CardNo
				lic.ID = "1"
				if c.ID == 0 {
					lic.Operation = "new"
				} else {
					lic.Operation = "change"
				}
				lic.CreateTime = time.UnixMilli(c.StartTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
				lic.EffectiveStartDate = time.UnixMilli(c.EndTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
				lic.EffectiveTime = time.UnixMilli(c.EndTime).UTC().Add(8 * time.Hour).Format(time.RFC3339)
				data.LicensePlateInfo = append(data.LicensePlateInfo, lic)
				err = global.HikGateway.TCG225EVCLGetCond(v.Host, data)
				if err != nil {
					return err
				}
			}

		}
	}
	return
}

func (c *CarModel) DeleteByDevice(deviceIds []uint) (err error) {
	var deviceModel []*DeviceModel
	global.TD27_DB.Where("id in (?)", deviceIds).Find(&deviceModel)
	for _, v := range deviceModel {

		if v.Model == HIK_DS_TCG225 ||
			v.Model == HIK_DS_TCG2A5_B ||
			v.Model == HIK_DS_2CD9125_KS ||
			v.Model == HIK_DS_TCG2A5_E {
			data := hk_gateway.VCLDelCondReq{}
			data.CardNo = c.CardNo
			data.PlateColor = c.Color
			data.PlateNum = c.CarNum
			data.PlateType = c.CarType
			data.DelVCLCond = 1
			err = global.HikGateway.VCLDelCond(v.Host, data)
			continue
		}

		if v.Model == HIK_DS_TCG205_E {
			data := hk_gateway.TCG2A5EVCLDelCondReq{}
			data.DeleteAllEnabled = false
			data.LicensePlate = []string{c.CarNum}
			err = global.HikGateway.TCG205EVCLDelCond(v.Host, data)
		}

	}
	return
}
