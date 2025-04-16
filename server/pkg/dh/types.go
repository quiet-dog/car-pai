package dh

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
	Name         string `json:"name"`         // TrafficBlackList 黑  TrafficRedList 红
	PlateNumber  string `json:"plateNumber"`  // 车牌号
	MasterOfCar  string `json:"masterOfCar"`  // 车主姓名
	PlateColor   string `json:"plateColor"`   // 车牌颜色
	PlateType    string `json:"plateType"`    // 车牌类型
	VehicleType  string `json:"vehicleType"`  // 车辆类型
	VehicleColor string `json:"vehicleColor"` // 车辆颜色
	BeginTime    string `json:"beginTime"`    // 开始时间 2010-05-25 00:00:00
	EndTime      string `json:"endTime"`      // 结束时间 2010-05-25 00:00:00
}
