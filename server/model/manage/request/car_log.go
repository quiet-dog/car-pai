package request

import ComReq "server/model/common/request"

type SearchCarLog struct {
	ComReq.PageInfo
	CarNum string `json:"carNum"` // 车牌号
	Color  string `json:"color"`  // 颜色
	AreaId uint   `json:"areaId"` // 区域ID
}
