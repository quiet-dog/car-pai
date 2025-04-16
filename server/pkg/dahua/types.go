package dahua

type TollgateInfoReq struct {
	Picture Picture `json:"Picture"`
}

type LoginReq struct {
	Method  string `json:"method"`
	Session string `json:"session"`
	Info    Info   `json:"info"`
}

type Info struct {
	UserName string `json:"userName"`
	Password string `json:"passWord"`
}

type Picture struct {
	NormalPic  NormalPic   `json:"NormalPic"`
	CombinPic  []NormalPic `json:"CombinPic"`
	CutouPic   NormalPic   `json:"CutouPic"`
	VehiclePic NormalPic   `json:"VehiclePic"`
	FacePic    []Face      `json:"FacePic"`
	Plate      Plate       `json:"Plate"`
	Vehicle    Vehicle     `json:"Vehicle"`
	SnapInfo   SnapInfo    `json:"SnapInfo"`
}

type NormalPic struct {
	PicName string `json:"PicName"`
	Content string `json:"Content"`
}

type Face struct {
	PicType int32 `json:"PicType"`
	NormalPic
}

type Plate struct {
	IsExit      bool    `json:"IsExit"`
	PlateNumber string  `json:"PlateNumber"`
	ADR         string  `json:"ADR"`
	PlateColor  string  `json:"PlateColor"`
	PlateType   string  `json:"PlateType"`
	Confidence  int32   `json:"Confidence"`
	BoundingBox []int32 `json:"BoundingBox"`
	UploadNum   int32   `json:"UploadNum"`
	Channel     int32   `json:"Channel"`
	Region      string  `json:"Region"`
}

type Vehicle struct {
	VehicleColor       string  `json:"VehicleColor"`
	VehicleSign        string  `json:"VehicleSign"`
	VehicleType        string  `json:"VehicleType"`
	VehocleSeries      string  `json:"VehocleSeries"`
	Speed              int32   `json:"Speed"`
	VehicleBoundingBox []int32 `json:"VehicleBoundingBox"`
}

type SnapInfo struct {
	TriggerSource    string `json:"TriggerSource"`
	SnapTime         string `json:"SnapTime"`
	AccurateTime     string `json:"AccurateTime"`
	TimeZone         int32  `json:"TimeZone"`
	DSTTune          int32  `json:"DSTTune"`
	SnapAddress      string `json:"SnapAddress"`
	LanNo            int32  `json:"LanNo"`
	Direction        string `json:"Direction"`
	OpenStrobe       bool   `json:"OpenStrobe"`
	AllowUser        bool   `json:"AllowUser"`
	AllowUserEndTime string `json:"AllowUserEndTime"`
	BlockUser        bool   `json:"BlockUser"`
	BlockUserEndTime string `json:"BlockUserEndTime"`
	DefenceCode      string `json:"DefenceCode"`
	DeviceID         string `json:"DeviceID"`
	InCarPeopleNum   int32  `json:"InCarPeopleNum"`
}

type DeviceInfo struct {
	DeviceName   string `json:"DeviceName"`
	DeviceModel  string `json:"DeviceModel"`
	DeviceType   string `json:"DeviceType"`
	Manufacturer string `json:"Manufacturer"`
	IPAddress    string `json:"IPAddress"`
	IPv6Address  string `json:"IPv6Address"`
	MACAddress   string `json:"MACAddress"`
	DeviceID     string `json:"DeviceID"`
}
