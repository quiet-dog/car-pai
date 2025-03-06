package manage

import "server/global"

type CarLogModel struct {
	global.TD27_MODEL
	CarNum    string `json:"carNum" gorm:"not null;comment:车牌号"`     // 车牌号
	DeviceId  uint   `json:"deviceId" gorm:"not null;comment:设备ID"`  // 设备ID
	Uri       string `json:"uri" gorm:"not null;comment:图片地址"`       // 图片地址
	PlateType string `json:"plateType" gorm:"not null;comment:车牌类型"` // 车牌类型
}
