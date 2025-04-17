package dh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

func New(conf Config) *Client {
	c := &Client{
		Config: &conf,
	}

	return c
}

func (c *Client) MagicBox() {
	err := c.do(ReqInitParam{
		Method: Get,
		Url:    "/cgi-bin/magicBox.cgi?action=getLanguageCaps",
		Query:  nil,
		Body:   nil,
	})
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
}

func (c *Client) Import(body ImportCar) {
	body.IsOverWrite = false
	body.StartTime = "2017/1/1 0:00"
	body.EndTime = "2037/12/31 0:00"
	body.Type = 0
	body.Username = "admin"
	body.CarNo = "粤A12345"

	d := &bytes.Buffer{}
	writer := multipart.NewWriter(d)
	boundary := "my-custom-boundary"
	writer.SetBoundary(boundary)

	jsond := map[string]interface{}{
		"type":        body.Type,
		"isOverWrite": body.IsOverWrite,
	}
	jsbytes, _ := json.Marshal(jsond)

	jsonPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":   {"application/json"},
		"Content-Length": {fmt.Sprintf("%d", len(jsbytes))}, // {"type": 0, "isOverWrite": true} 的长度
	})
	if err != nil {
		fmt.Println("创建 JSON 部分失败:", err)
		return
	}

	_, err = jsonPart.Write([]byte(jsbytes))
	if err != nil {
		fmt.Println("写入 JSON 数据失败:", err)
		return
	}

	csvData := fmt.Sprintf("开始时间,结束时间,车主姓名,车牌号\n%s,%s,%s,%s\n", body.StartTime, body.EndTime, body.Username, body.CarNo)

	csvBytes := []byte(csvData)
	csvPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":   {"application/octet-stream"},
		"Content-Length": {fmt.Sprintf("%d", len(csvBytes))}, // 动态计算长度
	})
	if err != nil {
		fmt.Println("创建 CSV 部分失败:", err)
		return
	}
	_, err = csvPart.Write(csvBytes)
	if err != nil {
		fmt.Println("写入 CSV 数据失败:", err)
		return
	}

	// 关闭 writer
	writer.Close()

	// fmt.Println("请求的内容:", d.String())
	// 将buff转为buty
	err = c.do(ReqInitParam{
		Method: Post,
		Url:    "/cgi-bin/api/ImExport/importData",
		Query:  nil,
		Body:   d.String(),
		Headers: map[string]string{
			"Content-Type": "multipart/x-mixed-replace; boundary=" + boundary,
		},
	})
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
}

func (c *Client) Insert(car Car) (err error) {
	query := map[string]string{}
	if car.Name != "" {
		query["name"] = car.Name
	}
	if car.PlateNumber != "" {
		query["PlateNumber"] = car.PlateNumber
	}
	if car.MasterOfCar != "" {
		query["MasterOfCar"] = car.MasterOfCar
	}
	if car.PlateColor != "" {
		query["PlateColor"] = car.PlateColor
	}
	if car.PlateType != "" {
		query["PlateType"] = car.PlateType
	}
	if car.VehicleType != "" {
		query["VehicleType"] = car.VehicleType
	}
	if car.VehicleColor != "" {
		query["VehicleColor"] = car.VehicleColor
	}
	if car.BeginTime != "" {
		query["BeginTime"] = car.BeginTime
	}
	if car.EndTime != "" {
		query["CancelTime"] = car.EndTime
	}
	query["action"] = "insert"
	uri := fmt.Sprintf("action=insert&name=%s&PlateNumber=%s&MasterOfCar=%s&PlateColor=%s&PlateType=%s&VehicleType=%s&VehicleColor=%s&BeginTime=%s&EndTime=%s",
		car.Name, car.PlateNumber, car.MasterOfCar, car.PlateColor, car.PlateType, car.VehicleType, car.VehicleColor, car.BeginTime, car.EndTime)

	err = c.do(ReqInitParam{
		Method: Get,
		Url:    "/cgi-bin/recordUpdater.cgi?" + uri,
	})
	return
}

func (c *Client) Update(car Car) (err error) {
	if car.Recno == 0 {
		return fmt.Errorf("recno 不能为空")
	}
	query := map[string]string{}
	query["action"] = "update"
	if car.Name != "" {
		query["name"] = car.Name
	}
	if car.PlateNumber != "" {
		query["PlateNumber"] = car.PlateNumber
	}
	if car.MasterOfCar != "" {
		query["MasterOfCar"] = car.MasterOfCar
	}
	if car.PlateColor != "" {
		query["PlateColor"] = car.PlateColor
	}
	if car.PlateType != "" {
		query["PlateType"] = car.PlateType
	}
	if car.VehicleType != "" {
		query["VehicleType"] = car.VehicleType
	}
	if car.VehicleColor != "" {
		query["VehicleColor"] = car.VehicleColor
	}
	if car.BeginTime != "" {
		query["BeginTime"] = car.BeginTime
	}
	if car.EndTime != "" {
		query["CancelTime"] = car.EndTime
	}
	query["recno"] = fmt.Sprintf("%d", car.Recno)
	err = c.do(ReqInitParam{
		Method: Get,
		Url:    "/cgi-bin/recordUpdater.cgi",
		Query:  query,
	})
	return
}

func (c *Client) Delete(car DeleteCar) (err error) {
	if car.Recno == 0 {
		return fmt.Errorf("recno 不能为空")
	}
	query := map[string]string{}
	query["action"] = "remove"
	query["recno"] = fmt.Sprintf("%d", car.Recno)
	query["name"] = car.Name
	err = c.do(ReqInitParam{
		Method: Get,
		Url:    "/cgi-bin/recordUpdater.cgi",
		Query:  query,
	})
	return
}

func (c *Client) GetCar(req GetCarReq) (result *GetCarRes, err error) {
	// query := map[string]string{}
	// query["action"] = "find"
	// if req.Name != "" {
	// 	query["name"] = req.Name
	// }
	// if req.Count != 0 {
	// 	query["count"] = fmt.Sprintf("%d", req.Count)
	// }
	// if req.StartTime != "" {
	// 	query["StartTime"] = req.StartTime
	// }
	// if req.EndTime != "" {
	// 	query["EndTime"] = req.EndTime
	// }
	uri := fmt.Sprintf("action=find&name=%s&condition.PlateNumber=%s",
		req.Name, req.PlateNumber)
	bResult := ""
	err = c.do(ReqInitParam{
		Method: Get,
		Url:    "/cgi-bin/recordFinder.cgi?" + uri,
		// Query:  query,
		Result: &bResult,
	})
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	result, err = parseResponse(bResult)
	return
}

func parseResponse(body string) (*GetCarRes, error) {
	lines := strings.Split(body, "\n")
	resp := &GetCarRes{}
	recordMap := map[int]*Car{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		// 处理 found
		if key == "found" {
			found, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			resp.Found = found
			continue
		}

		// 处理 records
		if strings.HasPrefix(key, "records[") {
			idxStart := strings.Index(key, "[") + 1
			idxEnd := strings.Index(key, "]")
			if idxStart < 0 || idxEnd < 0 || idxEnd <= idxStart {
				continue
			}
			indexStr := key[idxStart:idxEnd]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				continue
			}
			field := key[idxEnd+2:] // skip ].  (eg: ].CreateTime)

			// 创建或获取记录对象
			rec, ok := recordMap[index]
			if !ok {
				rec = &Car{}
				recordMap[index] = rec
			}

			// 填充字段
			switch field {
			case "BeginTime":
				rec.BeginTime = value
			case "CancelTime":
				rec.EndTime = value
			// case "ControlledType":
			// 	rec.ControlledType = value
			// case "CreateTime":
			// 	rec.BeginTime, _ = strconv.ParseInt(value, 10, 64)
			case "MasterOfCar":
				rec.MasterOfCar = value
			case "PlateColor":
				rec.PlateColor = value
			case "PlateNumber":
				rec.PlateNumber = value
			case "PlateType":
				rec.PlateType = value
			case "RecNo":
				rec.Recno, _ = strconv.Atoi(value)
			case "VehicleColor":
				rec.VehicleColor = value
			case "VehicleType":
				rec.VehicleType = value
			}
		}
	}

	// 按序组装记录
	for i := 0; i < len(recordMap); i++ {
		if rec, ok := recordMap[i]; ok {
			resp.Records = append(resp.Records, *rec)
		}
	}

	return resp, nil
}

func (c *Client) do(req ReqInitParam) (err error) {
	client := resty.New().
		// SetHeader("User-Agent", "client/1.0").
		SetBaseURL(fmt.Sprintf("http://%s:%s", c.Config.Host, c.Config.Port)).
		SetDebug(true)
	var authParams map[string]string

	// 中间件：捕获 401 响应并解析 WWW-Authenticate
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		if resp.StatusCode() == 401 {
			authHeader := resp.Header().Get("WWW-Authenticate")
			if authHeader != "" {
				authParams = parseDigestHeader(authHeader)
			}
		}
		return nil
	})

	// 中间件：在请求前添加自定义 Authorization
	client.OnBeforeRequest(func(rc *resty.Client, req *resty.Request) error {
		if authParams != nil {
			// 自定义参数
			nc := "00000001"           // 可递增
			cnonce := generateCnonce() // 自定义 cnonce
			uri := req.URL
			method := req.Method

			// 计算 response
			response := calculateResponse(c.Config.Username, c.Config.Password, authParams["realm"], authParams["nonce"], nc, cnonce, authParams["qop"], method, uri)

			// 构造 Authorization header
			auth := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s", opaque="%s"`,
				c.Config.Username, authParams["realm"], authParams["nonce"], uri, authParams["qop"], nc, cnonce, response, authParams["opaque"])
			req.SetHeader("Authorization", auth)
		}
		return nil
	})

	d := client.R().SetBody(req.Body).SetQueryParams(req.Query).SetHeaders(req.Headers).SetResult(req.Result)

	var resp *resty.Response

	switch req.Method {
	case Post:
		resp, err = d.Post(req.Url)
	case Get:
		resp, err = d.Get(req.Url)
	case Put:
		resp, err = d.Put(req.Url)
	case Delete:
		resp, err = d.Delete(req.Url)
	}

	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}

	// 如果需要手动重试（视情况而定）
	if resp.StatusCode() == 401 && authParams != nil {

		switch req.Method {
		case Post:
			resp, err = d.Post(req.Url)
		case Get:
			resp, err = d.Get(req.Url)
		case Put:
			resp, err = d.Put(req.Url)
		case Delete:
			resp, err = d.Delete(req.Url)
		}
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode())
	}

	fmt.Println("状态码:", resp.String())
	// req.Result是 &string ,将相应的结果赋值给result
	if req.Result != nil {
		result, ok := req.Result.(*string)
		if ok {
			*result = resp.String()
		}
	}
	return
}
