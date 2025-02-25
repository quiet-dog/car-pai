package hk_gateway

import "encoding/xml"

type GateWayConnect struct {
	HikConfig HikConfig
	IsConnect bool
}

type UserCheckRes struct {
	XMLNAME           xml.Name `xml:"userCheck"`
	StatusValue       int      `xml:"statusValue"`
	StatusString      string   `xml:"statusString"`
	IsDefaultPassword bool     `xml:"isDefaultPassword"`
	IsRiskPassword    bool     `xml:"isRiskPassword"`
	IsActivated       bool     `xml:"isActivated"`
	ResidualValidity  int      `xml:"residualValidity"`
	LockStatus        string   `xml:"lockStatus"`
	UnlockTime        int      `xml:"unlockTime"`
	RetryLoginTime    int      `xml:"retryLoginTime"`
}

type ResponseStatusXML struct {
	XMLNAME       xml.Name `xml:"ResponseStatus"`
	RequestURL    string   `xml:"requestURL" json:"requestURL"`
	StatusCode    int      `xml:"statusCode" json:"statusCode"`
	StatusString  string   `xml:"statusString" json:"statusString"`
	SubStatusCode string   `xml:"subStatusCode" json:"subStatusCode"`
	ErrorCode     int      `xml:"errorCode" json:"errorCode"`
	ErrorMsg      string   `xml:"errorMsg" json:"errorMsg"`
}

type ErrorMsg struct {
	StatusCode    int    `json:"statusCode"`
	StatusString  string `json:"statusString"`
	SubStatusCode string `json:"subStatusCode"`
	ErrorCode     int    `json:"errorCode"`
	ErrorMsg      string `json:"errorMsg"`
}

type ResponseStatus struct {
	XMLName       xml.Name `xml:"ResponseStatus"` // 根元素名称
	Version       string   `xml:"version,attr"`   // version作为属性
	RequestURL    string   `xml:"requestURL"`     // 请求URL
	StatusCode    int      `xml:"statusCode"`     // 状态码
	StatusString  string   `xml:"statusString"`   // 状态描述
	SubStatusCode string   `xml:"subStatusCode"`  // 子状态码
}

type Msg struct {
	Ip   string
	Data Result
	Type string
	Msg  string
}

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type VCLData struct {
	XMLName      xml.Name     `xml:"VCLData"`
	Version      string       `xml:"version,attr"`
	CurrentUpNum int          `xml:"currentUpNum"`
	TotalNum     int          `xml:"totalNum"`
	VCLDataList  VCLDataListO `xml:"VCLDataList"`
}

// VCLDataList 包含多个singleVCLData的列表
type VCLDataListO struct {
	SingleVCLData []SingleVCLDataO `xml:"singleVCLData"`
}

// SingleVCLData 单个车辆数据
type SingleVCLDataO struct {
	ID         int    `xml:"id"`
	RunNum     int    `xml:"runNum"`
	ListType   int    `xml:"listType"`
	PlateType  int    `xml:"plateType"`
	PlateColor int    `xml:"plateColor"`
	PlateNum   string `xml:"plateNum"`
	CardNo     string `xml:"cardNo"`
	StartTime  string `xml:"startTime"`
	EndTime    string `xml:"endTime"`
}
