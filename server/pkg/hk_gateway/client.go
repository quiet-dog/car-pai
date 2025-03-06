package hk_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type HikConfig struct {
	Ip       string
	Port     int
	Username string
	Password string
}

type HikClient struct {
	hikConfig HikConfig
	client    *resty.Client
	isConnect bool
	// 上下文
	ctx context.Context
	// 取消长连接的上下文
	cancel context.CancelFunc
	// 长连接状态
	longConnect bool
	// 长连接通道
	longConnectCtx       context.Context
	longCancel           context.CancelFunc
	doorStatus           interface{}
	doorConnectCtx       context.Context
	doorCancel           context.CancelFunc
	doorEventLogs        []EventInfo
	doorEventLogsCtx     context.Context
	doorEventLogsCancel  context.CancelFunc
	doorPersonInfoCount  *GetPersonInfoCountRes
	doorAcsEventTotalNum *GetAcsEventTotalNumRes
	doorAcsEventLogs     *GetAcsEventRes
	doorThreeCtx         context.Context
	doorThreeCancel      context.CancelFunc
}

func NewHikClinet(conf HikConfig) (hikClient *HikClient, err error) {
	hikClient = &HikClient{
		hikConfig: conf,
		client:    resty.New().SetBaseURL(fmt.Sprintf("http://%s:%d", conf.Ip, conf.Port)).SetTimeout(3 * time.Second),
		isConnect: false,
	}

	// 开启长连接
	hikClient.isConnect = true
	if _, err = hikClient.UserCheck(); err != nil {
	}

	if _, err = hikClient.GetDeviceInfo(); err != nil {
	}

	// 主线程结束时关闭长连接
	hikClient.longConnectCtx, hikClient.longCancel = context.WithCancel(context.Background())
	hikClient.doorEventLogsCtx, hikClient.doorEventLogsCancel = context.WithCancel(context.Background())
	hikClient.doorThreeCtx, hikClient.doorThreeCancel = context.WithCancel(context.Background())
	time.Sleep(100 * time.Millisecond)
	go hikClient.WatchLongConnect()
	go hikClient.WatchDoorStatus()
	go hikClient.WatchDoorEventLogs()
	go hikClient.WatchThreeDoorInfo()

	hikClient.doorConnectCtx, hikClient.doorCancel = context.WithCancel(context.Background())

	return hikClient, err
}

func NewOnlyHikClinet(conf HikConfig) (hikClient *HikClient, err error) {
	hikClient = &HikClient{
		hikConfig: conf,
		client:    resty.New().SetBaseURL(fmt.Sprintf("%s:%d", conf.Ip, conf.Port)).SetTimeout(3 * time.Second),
		isConnect: false,
	}
	// hikClient.isConnect = true

	// hikClient.doorConnectCtx, hikClient.doorCancel = context.WithCancel(context.Background())

	return hikClient, err
}

func newHikClinet(conf HikConfig) (hikClient *HikClient, err error) {
	hikClient = &HikClient{
		hikConfig: conf,
		client:    resty.New().SetBaseURL(fmt.Sprintf("%s:%d", conf.Ip, conf.Port)).SetTimeout(3 * time.Second),
		isConnect: false,
	}

	// 开启长连接
	// if _, err = hikClient.UserCheck(); err != nil {
	// }
	hikClient.isConnect = true

	return hikClient, err
}

func (c *HikClient) Close() {
	// c.longCancel()
}

func getDoorEventLogs(c *HikClient, index int, startTime string, endTime string, list *[]EventInfo) {

	params := GetAcsEventReq{}
	params.AcsEventCond = AcsEventCond{
		SearchID:             uuid.New().String(),
		SearchResultPosition: index * 10,
		MaxResults:           30,
		StartTime:            startTime,
		EndTime:              endTime,
		EmployeeNoString:     "",
	}

	result, err := c.GetAcsEvent(params)
	if err != nil || result == nil {
		return
	}
	for _, v := range result.AcsEvent.InfoList {
		if v.PictureURL != "" && len(*list) < 10 {
			*list = append(*list, v)
		}
		if len(*list) == 10 {
			return
		}
	}
	if len(*list) < 10 {
		index++
		getDoorEventLogs(c, index, startTime, endTime, list)
	}

}

func (c *HikClient) WatchThreeDoorInfo() {
	var getDoorInfo = func() {
		currentTime := time.Now()
		// 开始时间 格式 2024-06-05T10:57:02+08:00
		startTime := currentTime.Add(-time.Hour * 24 * 30).Format("2006-01-02T15:04:05+08:00")
		// 结束时间
		endTime := currentTime.Format("2006-01-02T15:04:05+08:00")
		personInfoCount, _ := c.GetPersonInfoCount()
		acsEventNum, _ := c.GetAcsEventTotalNum(GetAcsEventTotalNumReq{
			AcsEventTotalNumCond: AcsEventTotalNumCond{},
		})
		acsEventLog, _ := c.GetAcsEvent(GetAcsEventReq{
			AcsEventCond: AcsEventCond{
				SearchID:             uuid.NewString(),
				SearchResultPosition: 0,
				MaxResults:           5,
				StartTime:            startTime,
				EndTime:              endTime,
			},
		})
		if acsEventLog != nil {
			c.doorAcsEventLogs = acsEventLog
		}
		if personInfoCount != nil {
			c.doorPersonInfoCount = personInfoCount
		}
		if acsEventNum != nil {
			c.doorAcsEventTotalNum = acsEventNum
		}
	}
	getDoorInfo()
	for {
		select {
		case <-c.doorThreeCtx.Done():
			{
				return
			}
		case <-time.After(10 * time.Second):
			{
				getDoorInfo()
			}
		}
	}
}

func (c *HikClient) WatchDoorEventLogs() {
	currentTime := time.Now()
	// 开始时间 格式 2024-06-05T10:57:02+08:00
	startTime := currentTime.Add(-time.Hour * 24).Format("2006-01-02T15:04:05+08:00")
	// 结束时间
	endTime := currentTime.Format("2006-01-02T15:04:05+08:00")
	list := []EventInfo{}
	getDoorEventLogs(c, 1, startTime, endTime, &list)
	c.doorEventLogs = list
	for {
		select {
		case <-c.doorEventLogsCtx.Done():
			{
				return
			}
		case <-time.After(10 * time.Second):
			{
				currentTime := time.Now()
				// 开始时间 格式 2024-06-05T10:57:02+08:00
				startTime := currentTime.Add(-time.Hour * 24).Format("2006-01-02T15:04:05+08:00")
				// 结束时间
				endTime := currentTime.Format("2006-01-02T15:04:05+08:00")
				list := []EventInfo{}
				getDoorEventLogs(c, 1, startTime, endTime, &list)
				c.doorEventLogs = list
			}
		}
	}

}

func (c *HikClient) WatchLongConnect() {
	param := ReqInitParam{
		Url:    "/ISAPI/Security/userCheck",
		Query:  nil,
		Body:   nil,
		Result: nil,
		Method: Get,
	}
	userCheck := UserCheckRes{}
	if 0 < userCheck.RetryLoginTime && userCheck.RetryLoginTime < 3 {
		c.isConnect = false
	}

	// 1分钟计时器
	timer := time.NewTicker(1 * time.Minute)
	defer timer.Stop()
	timerT := time.NewTicker(3 * time.Second)
	defer timerT.Stop()
	var reload = func() {
		resClient := resty.New()
		// param.Headers["Authorization"] = auth

		req := resClient.SetBaseURL("http://"+c.hikConfig.Ip).SetTimeout(3*time.Second).R().SetDigestAuth(c.hikConfig.Username, c.hikConfig.Password)
		req = req.SetQueryParams(param.Query).SetBody(param.Body).SetHeaders(param.Headers)

		var resp *resty.Response
		var err error
		switch param.Method {
		case Post:
			{
				resp, err = req.Post(param.Url)
			}
		case Get:
			{
				resp, err = req.Get(param.Url)
			}
		case Put:
			{
				resp, err = req.Put(param.Url)
			}
		case Delete:
			{
				resp, err = req.Delete(param.Url)
			}
		}
		if err != nil {
			c.isConnect = false
			timer.Reset(1 * time.Minute)
			timerT.Stop()
			return
		}

		switch resp.Header().Get("Content-Type") {
		case "application/xml":
			{
				if resp.StatusCode() != 200 {
					// 先验证是否经过摘要认证
					c.isConnect = false
					if err = xml.Unmarshal(resp.Body(), &userCheck); err != nil {
						fmt.Printf("xml 解析失败 %s", err.Error())
					}
					timer.Reset(1 * time.Minute)
					timerT.Stop()
					return
				} else {
					c.isConnect = true
				}

			}
		default:
			{
				c.isConnect = true
			}
		}
	}
	// 开始计时
	for {
		select {
		case <-timerT.C:
			{
				reload()
			}
		case <-timer.C:
			{
				reload()
				if c.isConnect {
					timerT.Reset(3 * time.Second)
				}
			}

		case <-c.longConnectCtx.Done():
			{
				c.isConnect = false
				return
			}
		}
	}
}

func (c *HikClient) WatchDoorStatus() {
	result, err := c.GetAcsWorkStatus()
	if err != nil {
		c.doorStatus = nil
	}
	c.doorStatus = result
	for {
		select {
		case <-c.doorConnectCtx.Done():
			{
				return
			}
		case <-time.After(1 * time.Second):
			{
				result, err := c.GetAcsWorkStatus()
				if err != nil {
					continue
				}
				c.doorStatus = result
			}
		}
	}
}

func (c *HikClient) UserCheck() (result *UserCheckRes, err error) {
	req := ReqInitParam{
		Url:    "/ISAPI/Security/userCheck",
		Query:  nil,
		Body:   nil,
		Result: &result,
		Method: Get,
	}

	if err = c.Do(req); err != nil {
		return
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.isConnect = true
	return
}

func (c *HikClient) AddPersonInfo(person AddPersonInfoReq) (result *ErrorMsg, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/UserInfo/Record",
		Query: map[string]string{
			"format": "json",
		},
		Body:   person,
		Result: &result,
		Method: Post,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 设置人脸信息
// /ISAPI/Intelligent/FDLib/FDSetUp
func (c *HikClient) SetFaceInfo(reqBody SetFaceInfo) (result *ErrorMsg, err error) {

	req := ReqInitParam{
		Url: "/ISAPI/Intelligent/FDLib/FDSetUp",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Put,
	}

	if err = c.DoByte(req); err != nil {
		return
	}
	return
}

// 修改人员信息
// ISAPI/AccessControl/UserInfo/Modify?format=json
func (c *HikClient) ModifyPersonInfo(person AddPersonInfoReq) (result *ErrorMsg, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/UserInfo/Modify",
		Query: map[string]string{
			"format": "json",
		},
		Body:   person,
		Result: &result,
		Method: Put,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 删除人员信息
// ISAPI/AccessControl/UserInfo/Delete?format=json
func (c *HikClient) DeletePersonInfo(reqBody DelPersonInfo) (result *ErrorMsg, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/UserInfo/Delete",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Put,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 获取人员列表
// /ISAPI/AccessControl/UserInfo/Search
func (c *HikClient) GetPersonInfoList(reqBody GetPersonInfoList) (result *GetPersonInfoListRes, err error) {

	body := map[string]interface{}{
		"UserInfoSearchCond": map[string]interface{}{},
	}
	reqData := body["UserInfoSearchCond"].(map[string]interface{})
	reqData["searchID"] = reqBody.UserInfoSearchCond.SearchID
	reqData["searchResultPosition"] = reqBody.UserInfoSearchCond.SearchResultPosition
	reqData["maxResults"] = reqBody.UserInfoSearchCond.MaxResults
	if reqBody.UserInfoSearchCond.FuzzySearch != "" {
		reqData["fuzzySearch"] = reqBody.UserInfoSearchCond.FuzzySearch
	}

	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/UserInfo/Search",
		Query: map[string]string{
			"format": "json",
		},
		Body:   body,
		Result: &result,
		Method: Post,
	}

	if err := c.Do(req); err != nil {
		return nil, err
	}
	return
}

// 日志查询
// /ISAPI/ContentMgmt/logSearch
func (c *HikClient) LogSearch(reqBody LogSearch) (result *LogSearchRes, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/ContentMgmt/logSearch",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Post,
	}

	if err := c.Do(req); err != nil {
		return nil, err
	}
	return result, nil
}

// 获取人员数量信息
// /ISAPI/AccessControl/UserInfo/Count
func (c *HikClient) GetPersonInfoCount() (result *GetPersonInfoCountRes, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/UserInfo/Count",
		Query: map[string]string{
			"format": "json",
		},
		Body:   nil,
		Result: &result,
		Method: Get,
	}

	if err := c.Do(req); err != nil {
		return nil, err
	}

	return
}

// 查询设备中已有的人脸数量及人脸信息
// /ISAPI/Intelligent/FDLib/Count?format=json
func (c *HikClient) GetFaceInfoCount() (result *GetFaceInfoCountRes, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/Intelligent/FDLib/Count",
		Query: map[string]string{
			"format": "json",
		},
		Body:   nil,
		Result: &result,
		Method: Get,
	}

	if err := c.Do(req); err != nil {
		return nil, err
	}
	return result, nil
}

// 查询指定或全部人员的卡数量
// /ISAPI/AccessControl/CardInfo/Count?format=json
func (c *HikClient) GetCardInfoCount() (result *GetCardInfoCountRes, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/CardInfo/Count",
		Query: map[string]string{
			"format": "json",
		},
		Body:   nil,
		Result: &result,
		Method: Get,
	}

	if err := c.Do(req); err != nil {
		return nil, err
	}
	return
}

// 查询门禁事件条数。
// /ISAPI/AccessControl/AcsEventTotalNum?format=json
func (c *HikClient) GetAcsEventTotalNum(reqBody GetAcsEventTotalNumReq) (result *GetAcsEventTotalNumRes, err error) {
	body := make(map[string]interface{})
	body["AcsEventTotalNumCond"] = make(map[string]interface{})
	subBody := body["AcsEventTotalNumCond"].(map[string]interface{})
	subBody["major"] = reqBody.AcsEventTotalNumCond.Major
	subBody["minor"] = reqBody.AcsEventTotalNumCond.Minor

	if reqBody.AcsEventTotalNumCond.CardNo != "" {
		subBody["cardNo"] = reqBody.AcsEventTotalNumCond.CardNo
	}
	if reqBody.AcsEventTotalNumCond.Name != "" {
		subBody["name"] = reqBody.AcsEventTotalNumCond.Name
	}
	if reqBody.AcsEventTotalNumCond.EmployeeNoString != "" {
		subBody["employeeNoString"] = reqBody.AcsEventTotalNumCond.EmployeeNoString
	}
	if reqBody.AcsEventTotalNumCond.BeginSerialNo != 0 {
		subBody["beginSerialNo"] = reqBody.AcsEventTotalNumCond.BeginSerialNo
	}
	if reqBody.AcsEventTotalNumCond.EndSerialNo != 0 {
		subBody["endSerialNo"] = reqBody.AcsEventTotalNumCond.EndSerialNo
	}

	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/AcsEventTotalNum",
		Query: map[string]string{
			"format": "json",
		},
		Body:   body,
		Result: &result,
		Method: Post,
	}

	if err := c.Do(req); err != nil {
		return nil, err
	}
	return
}

// 获取门禁主机工作状态
// /ISAPI/AccessControl/AcsWorkStatus?format=json
func (c *HikClient) GetAcsWorkStatus() (result *GetAcsWorkStatusRes, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/AcsWorkStatus",
		Query: map[string]string{
			"format": "json",
		},
		Body:   nil,
		Result: &result,
		Method: Get,
	}

	if err = c.Do(req); err != nil {
		return nil, err
	}
	return result, nil
}

// 添加卡信息
// /ISAPI/AccessControl/CardInfo/Record
func (c *HikClient) AddCardInfo(reqBody AddCardInfoReq) (result *ErrorMsg, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/CardInfo/Record",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Post,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 获取用户的图片
func (c *HikClient) GetPicture(reqBody GetPictureReq) (result []byte, err error) {
	// 创建新的缓冲区
	resultBuf := bytes.NewBuffer([]byte{})
	req := ReqInitParam{
		Url:    reqBody.URL,
		Query:  nil,
		Body:   nil,
		Result: resultBuf,
		Method: Post,
	}

	if err = c.Do(req); err != nil {
		return
	}

	result = make([]byte, resultBuf.Len())
	if _, err = resultBuf.Read(result); err != nil {
		return
	}

	return result, nil
}

// 删除卡信息
// /ISAPI/AccessControl/CardInfo/Delete?format=json
func (c *HikClient) DeleteCardInfo(reqBody DelCardInfoReq) (result *ErrorMsg, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/CardInfo/Delete",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Put,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 获取卡信息
// /ISAPI/AccessControl/CardInfo/Search?format=json
func (c *HikClient) GetCardInfo(reqBody GetCardInfoReq) (result *GetCardInfoRes, err error) {

	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/CardInfo/Search",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Post,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 查询门禁事件
// /ISAPI/AccessControl/AcsEvent?format=json
func (c *HikClient) GetAcsEvent(reqBody GetAcsEventReq) (result *GetAcsEventRes, err error) {

	body := make(map[string]interface{})
	body["AcsEventCond"] = map[string]interface{}{
		"searchID":             reqBody.AcsEventCond.SearchID,
		"searchResultPosition": reqBody.AcsEventCond.SearchResultPosition,
		"maxResults":           reqBody.AcsEventCond.MaxResults,
		"major":                reqBody.AcsEventCond.Major,
		"minor":                reqBody.AcsEventCond.Minor,
		"timeReverseOrder":     true,
	}

	if reqBody.AcsEventCond.StartTime != "" {
		body["AcsEventCond"].(map[string]interface{})["startTime"] = reqBody.AcsEventCond.StartTime
	}

	if reqBody.AcsEventCond.EndTime != "" {
		body["AcsEventCond"].(map[string]interface{})["endTime"] = reqBody.AcsEventCond.EndTime
	}
	if reqBody.AcsEventCond.CardNo != "" {
		// 删除空字段
		body["AcsEventCond"].(map[string]interface{})["cardNo"] = reqBody.AcsEventCond.CardNo
	}
	if reqBody.AcsEventCond.Name != "" {
		// 删除空字段
		body["AcsEventCond"].(map[string]interface{})["name"] = reqBody.AcsEventCond.Name
	}
	if reqBody.AcsEventCond.EmployeeNoString != "" {
		// 删除空字段
		body["AcsEventCond"].(map[string]interface{})["employeeNoString"] = reqBody.AcsEventCond.EmployeeNoString
	}

	if reqBody.AcsEventCond.BeginSerialNo != 0 {
		// 删除空字段
		body["AcsEventCond"].(map[string]interface{})["beginSerialNo"] = reqBody.AcsEventCond.BeginSerialNo
	}

	if reqBody.AcsEventCond.EndSerialNo != 0 {
		// 删除空字段
		body["AcsEventCond"].(map[string]interface{})["endSerialNo"] = reqBody.AcsEventCond.EndSerialNo
	}

	if reqBody.AcsEventCond.Major != 0 {
		// 删除空字段
		body["AcsEventCond"].(map[string]interface{})["major"] = reqBody.AcsEventCond.Major
	}

	if reqBody.AcsEventCond.Minor != 0 {
		// 删除空字段
		body["AcsEventCond"].(map[string]interface{})["minor"] = reqBody.AcsEventCond.Minor
	}

	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/AcsEvent",
		Query: map[string]string{
			"format": "json",
		},
		Body:   body,
		Result: &result,
		Method: Post,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 配置门禁主机参数
// /ISAPI/AccessControl/AcsCfg?format=json
func (c *HikClient) SetAcsCfg(reqBody SetAcsCfgReq) (result *ErrorMsg, err error) {
	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/AcsCfg",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Put,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// /ISAPI/AccessControl/AcsEvent/StorageCfg?format=json
func (c *HikClient) SetStorageCfg(reqBody SetStorageCfgReq) (result *ErrorMsg, err error) {
	body := make(map[string]interface{})
	body["EventStorageCfg"] = make(map[string]interface{})
	subBody := body["EventStorageCfg"].(map[string]interface{})
	if reqBody.EventStorageCfg.CheckTime != "" {
		subBody["checkTime"] = reqBody.EventStorageCfg.CheckTime
	}
	if reqBody.EventStorageCfg.Mode != "" {
		subBody["mode"] = reqBody.EventStorageCfg.Mode
	}
	if reqBody.EventStorageCfg.Period != 0 {
		subBody["period"] = reqBody.EventStorageCfg.Period
	}

	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/AcsEvent/StorageCfg",
		Query: map[string]string{
			"format": "json",
		},
		Body:   body,
		Result: &result,
		Method: Put,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 配置人员及凭证显示参数
// /ISAPI/AccessControl/userAndRightShow?format=json
func (c *HikClient) GetUserAndRightShow(reqBody GetUserAndRightShowReq) (result *ErrorMsg, err error) {
	body := make(map[string]interface{})
	body["showAuthenticationList"] = reqBody.ShowAuthenticationList
	if reqBody.ShowCardNo != "" {
		body["showCardNo"] = reqBody.ShowCardNo
	}
	if reqBody.ShowDuration != 0 {
		body["showDuration"] = reqBody.ShowDuration
	}

	req := ReqInitParam{
		Url: "/ISAPI/AccessControl/userAndRightShow",
		Query: map[string]string{
			"format": "json",
		},
		Body:   reqBody,
		Result: &result,
		Method: Put,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 远程控门
// PUT/ISAPI/AccessControl/RemoteControl/door/<doorID>
func (c *HikClient) RemoteControlDoor(reqBody RemoteControlDoorReq) (result *ErrorMsg, err error) {

	type RemoteControlDoorInt struct {
		XMLName xml.Name `xml:"RemoteControlDoor"`
		Cmd     string   `xml:"cmd"`
	}
	body := RemoteControlDoorInt{
		Cmd: reqBody.RemoteControlDoor.Cmd,
	}

	xmlData, err := xml.Marshal(body)
	if err != nil {
		return
	}

	fmt.Println(reqBody.DoorID)
	req := ReqInitParam{
		Url:    fmt.Sprintf("/ISAPI/AccessControl/RemoteControl/door/%s", reqBody.DoorID.DoorID),
		Query:  map[string]string{},
		Body:   xmlData,
		Result: &result,
		Method: Put,
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 获取设备信息参数
// GET/ISAPI/System/deviceInfo
func (c *HikClient) GetDeviceInfo() (result *DeviceInfo, err error) {
	req := ReqInitParam{
		Url:    "/ISAPI/System/deviceInfo",
		Query:  nil,
		Body:   nil,
		Result: &result,
		Method: Get,
	}

	if err = c.Do(req); err != nil {
		return
	}
	return result, nil
}

// 设置车拍
func (c *HikClient) SetVCLData(reqBody SetVCLDataReq) (err error) {

	xmlData, err := xml.Marshal(reqBody)
	if err != nil {
		return
	}
	xmlWithHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlData))
	fmt.Println("xmlWithHer", string(xmlWithHeader))
	req := ReqInitParam{
		Url: "/ISAPI/ITC/Entrance/VCL",
		// Query: map[string]string{
		// 	"format": "json",
		// },
		Query:  nil,
		Body:   xmlWithHeader,
		Result: nil,
		Method: Put,
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}

	if reqBody.VCLDataList.SingleVCLData[0].CreateTime != "" {
		req.Url = "ISAPI/Traffic/channels/1licensePlateAuditData/record"
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 修改
func (c *HikClient) VCLGetCond(reqBody SetVCLDataReq) (err error) {
	xmlData, err := xml.Marshal(reqBody)
	if err != nil {
		return
	}

	xmlWithHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlData))
	fmt.Println("请求信息", string(xmlWithHeader))
	req := ReqInitParam{
		Url: "/ISAPI/ITC/Entrance/VCL",
		// Query: map[string]string{
		// 	"format": "json",
		// },
		Body:   xmlWithHeader,
		Result: nil,
		Method: Put,
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}

	if err = c.Do(req); err != nil {
		return
	}
	return
}

func (c *HikClient) TCG225EVCLGetCond(reqBody TCG225EVCLGetCondReq) (err error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return
	}

	fmt.Println("请求信息", string(jsonData))

	req := ReqInitParam{
		Url: "/ISAPI/Traffic/channels/1/licensePlateAuditData/record",
		// Query: map[string]string{
		// 	"format": "json",
		// },
		Body:   jsonData,
		Result: nil,
		Method: Put,
		Headers: map[string]string{
			"Content-Type": "application/json,charset=utf-8",
		},
	}

	// if reqBody.VCLDataList.SingleVCLData[0].CreateTime != "" {
	// 	req.Url = "ISAPI/Traffic/channels/1licensePlateAuditData/record"
	// }
	err = c.Do(req)
	return
}

func (c *HikClient) VCLDelCond(reqBody VCLDelCondReq) (err error) {
	xmlData, err := xml.Marshal(reqBody)
	if err != nil {
		return
	}
	xmlWithHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlData))
	fmt.Println("请求信息", string(xmlWithHeader))
	req := ReqInitParam{
		Url: "/ISAPI/ITC/Entrance/VCL",
		// Query: map[string]string{
		// 	"format": "json",
		// },
		Body:   xmlWithHeader,
		Result: nil,
		Method: Delete,
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}
	if err = c.Do(req); err != nil {
		return
	}
	return
}

// 获取车牌列表
func (c *HikClient) VCLGetList(reqBody VCLGetListReq) (res *VCLGetListRes, err error) {
	xmlData, err := xml.Marshal(reqBody)
	if err != nil {
		return
	}

	xmlWithHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlData))
	req := ReqInitParam{
		Url:    "/ISAPI/ITC/Entrance/VCL",
		Method: Post,
		Body:   xmlWithHeader,
		Result: &res,
	}
	err = c.Do(req)
	return
}

func (c *HikClient) TCG225EVCLDelCond(reqBody TCG225EVCLDelCondReq) (err error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return
	}

	req := ReqInitParam{
		Url:    "/ISAPI/Traffic/channels/1/DelLicensePlateAuditData",
		Body:   jsonData,
		Result: nil,
		Method: Put,
		Headers: map[string]string{
			"Content-Type": "application/json,charset=utf-8",
		},
	}

	err = c.Do(req)
	return
}

func (c *HikClient) VCLSearch(reqBody CMSearchDescriptionReq) (result *CMSearchResult, err error) {
	xmlData, err := xml.Marshal(reqBody)
	if err != nil {
		return
	}
	xmlWithHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlData))
	fmt.Println("请求信息", string(xmlWithHeader))
	req := ReqInitParam{
		Url:    "/ISAPI/ITC/ContentMgmt/logSearch",
		Body:   xmlWithHeader,
		Result: &result,
		Method: Post,
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}
	if err = c.Do(req); err != nil {
		return
	}
	return
}

func (c *HikClient) ZhuaPaiSeach(reqBody ZhuaPaiCMSearchDescription) (result *CMSearchResult, err error) {
	xmlData, err := xml.Marshal(reqBody)
	if err != nil {
		return
	}
	xmlWithHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlData))
	// xmlWithHeader = []byte("<?xml version='1.0' encoding='utf-8'?><CMSearchDescription><searchID>CB1A4C1D-AC20-0001-81AD-153015A011BD</searchID><trackIDList><trackID>120</trackID></trackIDList><timeSpanList><timeSpan><startTime>2025-01-25T00:00:00Z</startTime><endTime>2025-02-25T23:59:59Z</endTime><laneNumber></laneNumber><carType>all</carType><illegalType>all</illegalType></timeSpan></timeSpanList><contentTypeList><contentType>metadata</contentType></contentTypeList><maxResults>40</maxResults><searchResultPostion>280</searchResultPostion><metadataList><metadataDescriptor>//recordType.meta.hikvision.com/timing</metadataDescriptor></metadataList></CMSearchDescription>")
	// fmt.Println("请求信息", string(xmlWithHeader))
	req := ReqInitParam{
		Url:    "/ISAPI/ITC/ContentMgmt/search",
		Body:   xmlWithHeader,
		Result: &result,
		Method: Post,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}
	if err = c.Do(req); err != nil {
		return
	}
	return
}
