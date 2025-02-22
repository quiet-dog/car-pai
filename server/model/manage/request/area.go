package request

import (
	ComReq "server/model/common/request"
)

type AddArea struct {
	Name string `json:"name" binding:"required"` // 地区名称
	// 备注
	Remark string `json:"remark"`
}

type EditArea struct {
	ComReq.CId
	AddArea
}

type SearchArea struct {
	ComReq.PageInfo
	Name   string `json:"name"`   // 地区名称
	Remark string `json:"remark"` // 备注
}
