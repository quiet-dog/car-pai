package manage

import (
	"server/global"
	"server/model/authority"

	"gorm.io/gorm"
)

type AreaModel struct {
	global.TD27_MODEL
	Name string `json:"name" gorm:"not null;comment:地区名称" binding:"required"` // 地区名称
	// 备注
	Remark  string                 `json:"remark" gorm:"comment:备注"`
	Users   []*authority.UserModel `json:"users" gorm:"many2many:user_area;"` // 用户
	UserIds []uint                 `json:"userIds" gorm:"-"`
}

func (a *AreaModel) AfterFind(tx *gorm.DB) (err error) {
	tx.Table("user_area").Select("user_model_id").Find(&a.UserIds)
	return
}
