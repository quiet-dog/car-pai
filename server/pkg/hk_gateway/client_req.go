package hk_gateway

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

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

func (c *HikClient) Do(param ReqInitParam) (err error) {

	// if param.Url != "/ISAPI/Security/userCheck" && !c.isConnect {
	// 	return errors.New("未连接")
	// }

	resClient := resty.New()
	if param.Headers == nil {
		param.Headers = make(map[string]string)
	}

	req := resClient.SetBaseURL(c.hikConfig.Ip).SetDebug(true).SetTimeout(3*time.Second).R().SetDigestAuth(c.hikConfig.Username, c.hikConfig.Password)
	req = req.SetQueryParams(param.Query).SetBody(param.Body).SetHeaders(param.Headers)

	var resp *resty.Response
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
		return
	}

	// 打印响应的内容
	// 判断响应的内容格式
	switch resp.Header().Get("Content-Type") {
	case "application/xml":
		{
			if resp.StatusCode() != 200 {
				userCheck := UserCheckRes{}
				// 先验证是否经过摘要认证
				if err = xml.Unmarshal(resp.Body(), &userCheck); err != nil {
					fmt.Printf("xml 解析失败 %s", err.Error())
					return fmt.Errorf("xml 解析失败 %s", err.Error())
				}

				if userCheck.StatusValue != 0 && userCheck.StatusValue != 200 {
					err = fmt.Errorf("StatusString = %s, StatusValue = %d, RetryLoginTime = %d", userCheck.StatusString, userCheck.StatusValue, userCheck.RetryLoginTime)
					return
				}
				if userCheck.StatusString == "Unauthorized" {
					c.isConnect = false
				}

				return fmt.Errorf("StatusString = %s, StatusValue = %d, RetryLoginTime = %d", userCheck.StatusString, userCheck.StatusValue, userCheck.RetryLoginTime)
			}
			if param.Result == nil {
				res := ResponseStatus{}
				if res.StatusCode == 200 {
					return
				}
				if err := xml.Unmarshal(resp.Body(), &res); err != nil {
					return errors.New("xml json 解析失败12" + err.Error())
				}
				if res.StatusCode != 0 && res.StatusCode != 1 && res.SubStatusCode != "OK" {
					return fmt.Errorf("错误1")
				}
				return
			}
			fmt.Println("打印响应的内容", string(resp.Body()))
			if err = xml.Unmarshal(resp.Body(), param.Result); err != nil {
				fmt.Printf("xml 解析失败 验证后  %s", err.Error())
				return fmt.Errorf("xml 解析失败 验证后1  %s", err.Error())
			}

		}
	case "application/json":
		{
			fmt.Println("打印响应的内容xxxxx", string(resp.Body()))
			res := ResponseStatus{}
			if err := json.Unmarshal(resp.Body(), &res); err != nil {
				return errors.New("xml json 解析失败" + err.Error())
			}

			if res.StatusCode != 0 && res.StatusCode != 1 && res.SubStatusCode != "OK" {
				return fmt.Errorf("xml json 解析失败 验证后  %s %d %s", res.StatusString, res.StatusCode, res.SubStatusCode)
			}

			if err = xml.Unmarshal(resp.Body(), param.Result); err != nil {
				return errors.New("xml json 解析失败 验证后2" + err.Error() + resp.Status())
			}

			// if resp.StatusCode() != 200 {
			// errMsg := ErrorMsg{}
			// if err = json.Unmarshal(resp.Body(), &errMsg); err != nil {
			// 	return errors.New("json 解析失败" + err.Error())
			// }

			// 判断返回的json格式是不是这个格式
			// if errMsg.StatusCode != 0 && errMsg.StatusCode != 1 && errMsg.SubStatusCode != "OK" {
			// 	err = fmt.Errorf("ErrorCode = %d, ErrorMsg = %s, StatusString = %s, StatusCode = %d, SubStatusCode = %s", errMsg.ErrorCode, errMsg.ErrorMsg, errMsg.StatusString, errMsg.StatusCode, errMsg.SubStatusCode)
			// 	return
			// }

			// if err = json.Unmarshal(resp.Body(), param.Result); err != nil {
			// 	return errors.New("json 解析失败 验证后2" + err.Error() + resp.Status())
			// }

		}
	case "image/jpeg":
		{
			param.Result.(*bytes.Buffer).Write(resp.Body())
		}
	default:
		{
			if resp.StatusCode() != 200 {
				return errors.New("请求失败")
			}
		}
	}
	// c.isConnect = true

	return
}

func (c *HikClient) DoByte(param ReqInitParam) (err error) {
	if !c.isConnect && param.Url != "/ISAPI/Security/userCheck" {
		return
	}

	req := c.client.R().SetQueryParams(param.Query).SetDigestAuth(c.hikConfig.Username, c.hikConfig.Password)
	img := param.Body.(SetFaceInfo).Img

	var imgCont multipart.File

	if img != nil {
		if imgCont, err = img.Open(); err != nil {
			return
		}
		defer imgCont.Close()
	}

	FaceDataRecordStr, err := json.Marshal(param.Body.(SetFaceInfo).FaceDataRecord)
	if err != nil {
		return
	}

	req = req.SetHeader("Content-Type", "multipart/form-data").SetQueryParams(param.Query).SetMultipartFields(&resty.MultipartField{
		Param:    "FaceDataRecord",
		Reader:   strings.NewReader(string(FaceDataRecordStr)),
		FileName: "",
	}, &resty.MultipartField{
		Param:       "img",
		Reader:      imgCont,
		FileName:    img.Filename,
		ContentType: "image/png",
	})

	var resp *resty.Response
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
		return
	}

	// 判断响应的内容格式
	switch resp.Header().Get("Content-Type") {
	case "application/xml":
		{
			if resp.StatusCode() != 200 {
				userCheck := UserCheckRes{}
				// 先验证是否经过摘要认证
				if err = xml.Unmarshal(resp.Body(), &userCheck); err != nil {
					fmt.Printf("xml 解析失败 %s", err.Error())
					return fmt.Errorf("xml 解析失败 %s", err.Error())
				}

				if userCheck.StatusValue != 0 && userCheck.StatusValue != 200 {
					err = fmt.Errorf("StatusString = %s, StatusValue = %d, RetryLoginTime = %d", userCheck.StatusString, userCheck.StatusValue, userCheck.RetryLoginTime)
					return
				}

			}

			if err = xml.Unmarshal(resp.Body(), param.Result); err != nil {
				fmt.Printf("xml 解析失败 验证后  %s", err.Error())
				return fmt.Errorf("xml 解析失败 验证后1  %s", err.Error())
			}
		}
	case "application/json":
		{
			// if resp.StatusCode() != 200 {
			errMsg := ErrorMsg{}
			if err = json.Unmarshal(resp.Body(), &errMsg); err != nil {
				return errors.New("json 解析失败" + err.Error())
			}

			// 判断返回的json格式是不是这个格式
			if errMsg.StatusCode != 0 && errMsg.StatusCode != 1 && errMsg.SubStatusCode != "OK" {
				err = fmt.Errorf("ErrorCode = %d, ErrorMsg = %s, StatusString = %s, StatusCode = %d, SubStatusCode = %s", errMsg.ErrorCode, errMsg.ErrorMsg, errMsg.StatusString, errMsg.StatusCode, errMsg.SubStatusCode)
				return
			}

			if err = json.Unmarshal(resp.Body(), param.Result); err != nil {
				return errors.New("json 解析失败 验证后2" + err.Error() + resp.Status())
			}
		}
	case "image/jpeg":
		{
			param.Result.(*bytes.Buffer).Write(resp.Body())
		}
	}
	return
}
