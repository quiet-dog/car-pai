package manage

import (
	"fmt"
	"server/global"
	"server/pkg/dh"
	"server/pkg/hk_gateway"
	"strconv"

	"gorm.io/gorm"
)

const (
	HIK_DS_TCG225     = "DS-TCG225"
	HIK_DS_TCG2A5_E   = "DS-TCG2A5-E"
	HIK_DS_TCG2A5_B   = "DS-TCG2A5-B"
	HIK_DS_TCG205_E   = "DS-TCG205-E"
	HIK_DS_2CD9125_KS = "DS-2CD9125-KS"
	DH_ITC436_PW9H_Z  = "ITC436-PW9H-Z"
	HIK               = "海康"
	DH                = "大华"
)

type DeviceModel struct {
	global.TD27_MODEL
	Host        string     `json:"host" gorm:"comment:主机地址"`         // 主机地址
	Port        string     `json:"port" gorm:"comment:端口号"`          // 端口号
	HikUsername string     `json:"hikUsername" gorm:"comment:海康用户名"` // 用户名
	HikPassword string     `json:"hikPassword" gorm:"comment:海康密码"`  // 密码
	DhUsername  string     `json:"dhUsername" gorm:"comment:大华用户名"`  // 用户名
	DhPassword  string     `json:"dhPassword" gorm:"comment:大华密码"`   // 密码
	Remark      string     `json:"remark" gorm:"comment:备注"`         // 备注
	AreaId      uint       `json:"areaId" gorm:"comment:地区ID"`       // 地区ID
	Type        string     `json:"type" gorm:"comment:设备类型"`         // 设备类型
	Rtsp        string     `json:"rtsp" gorm:"comment:RTSP地址"`       // RTSP地址
	Area        *AreaModel `json:"area"`                             // 地区
	// 型号
	Model string `json:"model" gorm:"comment:型号"`
}

func (d *DeviceModel) AfterCreate(tx *gorm.DB) (err error) {
	if d.Type == HIK {
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

		carList := []*CarModel{}
		tx.Table("car_area").Where("area_model_id = ?", d.AreaId).Find(&carList)
		for _, v := range carList {
			if err = tx.Save(&v).Error; err != nil {
				return err
			}
		}
	}

	if d.Type == DH {
		global.DhGateway.AddDevice(d.ID, dh.Config{
			Host:     d.Host,
			Port:     d.Port,
			Username: d.DhUsername,
			Password: d.DhPassword,
		})
	}

	return
}
