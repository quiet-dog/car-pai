package manage

import "server/global"

type CarModel struct {
	global.TD27_MODEL
	CarNum    string       `json:"carNum" gorm:"not null;comment:车牌号" binding:"required"`    // 车牌号
	Name      string       `json:"name" gorm:"not null;comment:车主姓名"`                        // 车主姓名
	Phone     string       `json:"phone" gorm:"not null;comment:车主电话"`                       // 车主电话
	StartTime int64        `json:"startTime" gorm:"not null;comment:开始时间"`                   // 开始时间
	EndTime   int64        `json:"endTime" gorm:"not null;comment:结束时间"`                     // 结束时间
	DeviceID  uint         `json:"deviceId" gorm:"not null;comment:设备ID" binding:"required"` // 设备ID
	Device    *DeviceModel `json:"device"`                                                   // 设备
	Remark    string       `json:"remark" gorm:"comment:备注"`                                 // 备注
	Areas     []*AreaModel `json:"areas" gorm:"many2many:car_area;"`
}
