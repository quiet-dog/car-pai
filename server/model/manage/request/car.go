package request

import (
	ComReq "server/model/common/request"
)

type AddCar struct {
	CarNum    string `json:"carNum" binding:"required"` // 车牌号
	Name      string `json:"name"`                      // 车主姓名
	Phone     string `json:"phone"`                     // 车主电话
	StartTime int64  `json:"startTime"`                 // 开始时间
	EndTime   int64  `json:"endTime"`                   // 结束时间
	AreaIDs   []uint `json:"areaIds"`                   // 地区ID
	Remark    string `json:"remark"`                    // 备注
	Color     string `json:"color"`                     // 颜色
	CarType   string `json:"carType"`                   // 车辆类型
	ListType  string `json:"listType"`                  // 名单类型 0白 1 黑
	CardNo    string `json:"cardNo"`                    // 卡号
}

type EditCar struct {
	ComReq.CId
	AddCar
}

type SearchCar struct {
	ComReq.PageInfo
	CarNum    string `json:"carNum"`    // 车牌号
	Name      string `json:"name"`      // 车主姓名
	Phone     string `json:"phone"`     // 车主电话
	StartTime int64  `json:"startTime"` // 开始时间
	EndTime   int64  `json:"endTime"`   // 结束时间
	AreaID    uint   `json:"areaId"`    // 地区ID
}
