package manage

import "server/global"

type AreaModel struct {
	global.TD27_MODEL
	Name string `json:"name" gorm:"not null;comment:地区名称" binding:"required"` // 地区名称
	// 备注
	Remark string `json:"remark" gorm:"comment:备注"`
}
