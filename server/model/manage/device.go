package manage

import (
	"fmt"
	"server/global"
	"server/pkg/hk_gateway"
	"strconv"

	"gorm.io/gorm"
)

type DeviceModel struct {
	global.TD27_MODEL
	Host        string     `json:"host" gorm:"not null;comment:主机地址" binding:"required"`         // 主机地址
	Port        string     `json:"port" gorm:"not null;comment:端口号" binding:"required"`          // 端口号
	HikUsername string     `json:"hikUsername" gorm:"not null;comment:海康用户名" binding:"required"` // 用户名
	HikPassword string     `json:"hikPassword" gorm:"not null;comment:海康密码" binding:"required"`  // 密码
	DhUsername  string     `json:"dhUsername" gorm:"not null;comment:大华用户名" binding:"required"`  // 用户名
	DhPassword  string     `json:"dhPassword" gorm:"not null;comment:大华密码" binding:"required"`   // 密码
	Remark      string     `json:"remark" gorm:"comment:备注"`                                     // 备注
	AreaId      uint       `json:"areaId" gorm:"not null;comment:地区ID" binding:"required"`       // 地区ID
	Type        string     `json:"type" gorm:"not null;comment:设备类型" binding:"required"`         // 设备类型
	Rtsp        string     `json:"rtsp" gorm:"not null;comment:RTSP地址" binding:"required"`       // RTSP地址
	Area        *AreaModel `json:"area"`                                                         // 地区
	// 型号
	Model string `json:"model" gorm:"not null;comment:型号" binding:"required"`
}

func (d *DeviceModel) AfterCreate(tx *gorm.DB) (err error) {
	if d.Type == "海康" {
		port, err := strconv.Atoi(d.Port)
		if err != nil {
			return fmt.Errorf("端口号错误")
		}
		if err = global.HikGateway.RegisterHikGateway(hk_gateway.HikConfig{
			Ip:       d.Host,
			Port:     port,
			Username: d.HikUsername,
			Password: d.HikPassword,
		}); err != nil {
			// 删除设备
			tx.Delete(d)
			return err
		}
	}

	return
}
