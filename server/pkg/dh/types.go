package dh

import "sync"

type Gateway struct {
	Device sync.Map
}

type Client struct {
	Config *Config
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
}

const (
	Post   = "POST"
	Get    = "GET"
	Put    = "PUT"
	Delete = "DELETE"
)

type ReqInitParam struct {
	Url     string
	Query   map[string]string
	Body    interface{}
	Result  interface{}
	Method  string
	Headers map[string]string
}

type ImportCar struct {
	Type        int    `json:"type"`
	IsOverWrite bool   `json:"isOverWrite"`
	CarNo       string `json:"carNo"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Username    string `json:"username"`
}

type Car struct {
	Recno        int    `json:"RecNo"`        // 记录号
	Name         string `json:"Name"`         // TrafficBlackList 黑  TrafficRedList 红
	PlateNumber  string `json:"PlateNumber"`  // 车牌号
	MasterOfCar  string `json:"MasterOfCar"`  // 车主姓名
	PlateColor   string `json:"PlateColor"`   // 车牌颜色
	PlateType    string `json:"PlateType"`    // 车牌类型
	VehicleType  string `json:"VehicleType"`  // 车辆类型
	VehicleColor string `json:"VehicleColor"` // 车辆颜色
	BeginTime    string `json:"BeginTime"`    // 开始时间 2010-05-25 00:00:00
	EndTime      string `json:"EndTime"`      // 结束时间 2010-05-25 00:00:00
}

type DeleteCar struct {
	Name  string `json:"name"`  // TrafficBlackList 黑  TrafficRedList 红
	Recno int    `json:"recno"` // 记录号
}

type GetCarReq struct {
	Name string `json:"name"` // TrafficBlackList 黑  TrafficRedList 红
	// Count       int    `json:"count"`
	StartTime string `json:"startTime"` // 开始时间 2010-05-25 00:00:00
	// EndTime     string `json:"endTime"`   // 结束时间 2010-05-25 00:00:00
	PlateNumber string `json:"plateNumber"`
}

type GetCarRes struct {
	TotalCount int   `json:"totalCount"` // 总记录数
	Found      int   `json:"found"`      // 查找到的记录数
	Records    []Car `json:"records"`    // 记录
}
