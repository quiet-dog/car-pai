package request

import ComReq "server/model/common/request"

type AddDevice struct {
	Host        string `json:"host" binding:"required"` // 地区名称
	Port        string `json:"port" binding:"required"` // 端口号
	HikUsername string `json:"hikUsername"`             // 用户名
	HikPassword string `json:"hikPassword"`             // 密码
	DhUsername  string `json:"dhUsername"`              // 用户名
	DhPassword  string `json:"dhPassword"`              // 密码
	Remark      string `json:"remark"`                  // 备注
	AreaId      uint   `json:"areaId"`                  // 地区ID
	Rtsp        string `json:"rtsp"`                    // RTSP地址
	Type        string `json:"type"`                    // 设备类型
	Model       string `json:"model"`
}

type EditDevice struct {
	ComReq.CId
	AddDevice
}

type SearchDevice struct {
	ComReq.PageInfo
	AreaId uint `json:"areaId"` // 地区ID
}

type RemoteControl struct {
	ComReq.CIds
	LockStatus string `json:"lockStatus"` // 锁定状态
}
