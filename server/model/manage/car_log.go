package manage

import "server/global"

type CarLogModel struct {
	global.TD27_MODEL
	CarNum    string       `json:"carNum" gorm:"comment:车牌号"`     // 车牌号
	DeviceId  uint         `json:"deviceId" gorm:"comment:设备ID"`  // 设备ID
	Device    *DeviceModel `json:"device"`                        // 设备
	Uri       string       `json:"uri" gorm:"comment:图片地址"`       // 图片地址
	PlateType string       `json:"plateType" gorm:"comment:车牌类型"` // 车牌类型
	SubTime   int64        `json:"subTime" gorm:"comment:提交时间"`   // 提交时间
}
