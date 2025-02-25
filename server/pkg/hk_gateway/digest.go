package hk_gateway

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type handleBody = func(data Result) bool

type Result struct {
	Data []byte
	Type string
}

// func (h *HikClient) StartLongConnect(handle handleBody) (err error) {
// 	if !h.isConnect {
// 		return
// 	}

// 	// h.ctx, h.cancel = context.WithCancel(context.Background())
// 	err = h.newHikClientLongConnect(handle)
// 	return
// 	//  结束newHikClientLongConnect的线程
// }

func (h *HikClient) newHikDisage(param ReqInitParam) (authValue string, err error) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	url := fmt.Sprintf("http://%s%s", h.hikConfig.Ip, param.Url)
	req, err := http.NewRequest(param.Method, url, nil)
	// 添加query
	query := req.URL.Query()
	for k, v := range param.Query {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Connection", "Keep-Alive ")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		result := UserCheckRes{}
		if err = xml.Unmarshal(body, &result); err != nil {
			return
		}

		if result.StatusString == "Unauthorized" {
			authDig := resp.Header.Get("Www-Authenticate")
			if strings.HasPrefix(authDig, "Digest") {
				authInfo := parseHeader(authDig)
				cnonce, _ := generateCnonce()
				response := calculateResponse(h.hikConfig.Username, h.hikConfig.Password, authInfo["realm"], param.Method, req.URL.RequestURI(), authInfo["qop"], authInfo["nonce"], "00000001", cnonce)
				authValue = fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s", opaque="%s"`, "admin", authInfo["realm"], authInfo["nonce"], req.URL.RequestURI(), authInfo["qop"], "00000001", cnonce, response, authInfo["opaque"])
				return
			}
		}
	}
	return
}

func (h *HikClient) newHikClientLongConnect(handle handleBody) (err error) {

	uri := "/ISAPI/Event/notification/alertStream"
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s%s", h.hikConfig.Ip, uri), nil)
	req.Header.Add("Connection", "Keep-Alive ")

	resp, err := client.Do(req)
	if err != nil {
		h.longConnect = false
		return err
	}
	defer resp.Body.Close()

	// 不确定body的长度,读取body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.longConnect = false
		return
	}

	result := UserCheckRes{}
	if err = xml.Unmarshal([]byte(body), &result); err != nil {
		h.longConnect = false
		return
	}

	// file, _ := os.OpenFile(fmt.Sprintf("./%ssteam.txt", h.hikConfig.Ip), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if result.StatusString == "Unauthorized" {
		authDig := resp.Header.Get("Www-Authenticate")
		// 获取realm的值
		if strings.HasPrefix(authDig, "Digest") {
			authInfo := parseHeader(authDig)
			cnonce, _ := generateCnonce()
			response := calculateResponse(h.hikConfig.Username, h.hikConfig.Password, authInfo["realm"], "GET", uri, authInfo["qop"], authInfo["nonce"], "00000001", cnonce)
			authValue := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s", opaque="%s"`, "admin", authInfo["realm"], authInfo["nonce"], uri, authInfo["qop"], "00000001", cnonce, response, authInfo["opaque"])

			req.Header.Add("Authorization", authValue)
			client.Timeout = 0
			if resp, err = client.Do(req); err != nil {
				h.longConnect = false
				return
			}

			switch resp.Header.Get("Content-Type") {
			case "application/xml":
				{

					defer resp.Body.Close()
					h.longConnect = false
					body, err = io.ReadAll(resp.Body)
					if err != nil {
						h.longConnect = false
						return err
					}
					result := ResponseStatusXML{}
					if err = xml.Unmarshal([]byte(body), &result); err != nil {
						h.longConnect = false
						return err
					}

					err = fmt.Errorf("statusCode %d, statusString %s, subStatusCode %s, errorCode %d, errorMsg %s", result.StatusCode, result.StatusString, result.SubStatusCode, result.ErrorCode, result.ErrorMsg)
					h.longConnect = false
					return err
				}
			default:
				{
					defer resp.Body.Close()
					for {
						result := Result{}
						streamBuf := make([]byte, 1024*2)
						_, err := resp.Body.Read(streamBuf)
						// 获取本次的流所有数据

						if err != nil {
							continue
						}
						if strings.Contains(string(streamBuf), "--MIME_boundary") && strings.Contains(string(streamBuf), `Content-Type: application/json; charset="UTF-8"`) {
							startIndex := strings.Index(string(streamBuf), "{")
							endIndex := strings.LastIndex(string(streamBuf), "}")
							if startIndex == -1 || endIndex == -1 {
								continue
							}
							if startIndex >= endIndex {
								continue
							}
							streamBuf = streamBuf[startIndex : endIndex+1]
							result.Data = streamBuf
							result.Type = "application/json"
							h.longConnect = true
							// file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", streamBuf))
							if !handle(result) {
								return err
							}
						}

						// file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", streamBuf))
						// result := Result{}
						// result.Data = streamBuf
						// result.Type = resp.Header.Get("Content-Type")
						// // 获取Content-Type：后面的值
						// startIndex := strings.Index(string(streamBuf), "{")
						// endIndex := strings.LastIndex(string(streamBuf), "}")
						// if startIndex == -1 || endIndex == -1 {
						// 	continue
						// }
						// streamBuf = streamBuf[startIndex : endIndex+1]
						// h.longConnect = true
						// file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", streamBuf))
						// if !handle(result) {
						// 	return err
						// }
					}

					// mediatype, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))

					// if err != nil {
					// 	h.longConnect = false
					// 	return err
					// }

					// if !strings.HasPrefix(mediatype, "multipart/mixed") {
					// 	h.longConnect = false
					// 	err = errors.New("No boundary found in Content-Type ssss")
					// 	return err
					// }

					// boundary, ok := params["boundary"]
					// if !ok {
					// 	h.longConnect = false
					// 	err = errors.New("No boundary found in Content-Type")
					// 	return err
					// }

					// reader := multipart.NewReader(resp.Body, boundary)

					// for {
					// 	result := Result{}

					// 	part, err := reader.NextPart()
					// 	if err == io.EOF {
					// 		file.WriteString(fmt.Sprintf("data is EOF\n"))
					// 		continue
					// 	}
					// 	if err != nil {
					// 		result.Type = ""
					// 		result.Data = []byte(err.Error())
					// 		file.WriteString(fmt.Sprintf("data is EOF error\n"))
					// 		continue
					// 	}

					// 	header := part.Header.Get("Content-Type")
					// 	partData, _ := io.ReadAll(part)

					// 	if len(partData) == 0 {
					// 		continue
					// 	}
					// 	file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", partData))

					// 	result.Data = partData
					// 	result.Type = header

					// 	h.longConnect = true
					// 	if !handle(result) {
					// 		return err
					// 	}
					// }

				}
			}
		}
	} else {
		h.longConnect = false
	}
	return
}

func (h *HikClient) NewHikClientLongConnect(handle handleBody) (err error) {

	uri := "/ISAPI/Event/notification/alertStream"
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", h.hikConfig.Ip, uri), nil)
	req.Header.Add("Connection", "Keep-Alive ")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("错误1")
		h.longConnect = false
		return err
	}
	defer resp.Body.Close()

	// 不确定body的长度,读取body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("错误2")
		h.longConnect = false
		return
	}

	result := UserCheckRes{}
	if err = xml.Unmarshal([]byte(body), &result); err != nil {
		h.longConnect = false
		fmt.Println("错误3")
		return
	}

	// file, _ := os.OpenFile(fmt.Sprintf("./%ssteam.txt", h.hikConfig.Ip), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if result.StatusString == "Unauthorized" {
		authDig := resp.Header.Get("Www-Authenticate")
		fmt.Println("=========1111==============")
		// 获取realm的值
		if strings.HasPrefix(authDig, "Digest") {
			authInfo := parseHeader(authDig)
			cnonce, _ := generateCnonce()
			response := calculateResponse(h.hikConfig.Username, h.hikConfig.Password, authInfo["realm"], "GET", uri, authInfo["qop"], authInfo["nonce"], "00000001", cnonce)
			authValue := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s", opaque="%s"`, "admin", authInfo["realm"], authInfo["nonce"], uri, authInfo["qop"], "00000001", cnonce, response, authInfo["opaque"])
			fmt.Println("=========33333==============")
			req.Header.Add("Authorization", authValue)
			client.Timeout = 0
			if resp, err = client.Do(req); err != nil {
				fmt.Println("=========44444==============")
				h.longConnect = false
				return
			}

			switch resp.Header.Get("Content-Type") {
			case "application/xml":
				{

					fmt.Println("=========5555==============")
					defer resp.Body.Close()
					h.longConnect = false
					body, err = io.ReadAll(resp.Body)
					fmt.Println("======body", string(body))
					if err != nil {
						fmt.Println("=========6666==============")
						h.longConnect = false
						return err
					}
					result := ResponseStatusXML{}
					if err = xml.Unmarshal([]byte(body), &result); err != nil {
						fmt.Println("=========7777==============")
						h.longConnect = false
						return err
					}

					err = fmt.Errorf("statusCode %d, statusString %s, subStatusCode %s, errorCode %d, errorMsg %s", result.StatusCode, result.StatusString, result.SubStatusCode, result.ErrorCode, result.ErrorMsg)
					h.longConnect = false
					return err
				}
			default:
				{
					fmt.Println("=========88888==============")
					defer resp.Body.Close()
					for {
						result := Result{}
						streamBuf := make([]byte, 1024*2)
						_, err := resp.Body.Read(streamBuf)
						// 获取本次的流所有数据
						if err != nil {
							continue
						}
						if strings.Contains(string(streamBuf), "--MIME_boundary") && strings.Contains(string(streamBuf), `Content-Type: application/json; charset="UTF-8"`) {
							startIndex := strings.Index(string(streamBuf), "{")
							endIndex := strings.LastIndex(string(streamBuf), "}")
							if startIndex == -1 || endIndex == -1 {
								continue
							}
							if startIndex >= endIndex {
								continue
							}
							streamBuf = streamBuf[startIndex : endIndex+1]
							result.Data = streamBuf
							result.Type = "application/json"
							h.longConnect = true
							// file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", streamBuf))
							if !handle(result) {
								return err
							}
						}
					}

				}
			}
		}
	} else {
		fmt.Println("=========2222==============")
		h.longConnect = false
	}
	return
}

func (h *HikClient) newBuFangHikClientLongConnect(handle handleBody) (err error) {

	uri := "/ISAPI/Event/notification/alertStream"
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s%s", h.hikConfig.Ip, uri), nil)
	req.Header.Add("Connection", "Keep-Alive ")

	resp, err := client.Do(req)
	if err != nil {
		h.longConnect = false
		return err
	}
	defer resp.Body.Close()

	// 不确定body的长度,读取body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.longConnect = false
		return
	}

	result := UserCheckRes{}
	if err = xml.Unmarshal([]byte(body), &result); err != nil {
		h.longConnect = false
		return
	}

	// file, _ := os.OpenFile(fmt.Sprintf("./%ssteam.txt", h.hikConfig.Ip), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if result.StatusString == "Unauthorized" {
		authDig := resp.Header.Get("Www-Authenticate")
		// 获取realm的值
		if strings.HasPrefix(authDig, "Digest") {
			authInfo := parseHeader(authDig)
			cnonce, _ := generateCnonce()
			response := calculateResponse(h.hikConfig.Username, h.hikConfig.Password, authInfo["realm"], "GET", uri, authInfo["qop"], authInfo["nonce"], "00000001", cnonce)
			authValue := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s", opaque="%s"`, "admin", authInfo["realm"], authInfo["nonce"], uri, authInfo["qop"], "00000001", cnonce, response, authInfo["opaque"])

			req.Header.Add("Authorization", authValue)
			client.Timeout = 0
			if resp, err = client.Do(req); err != nil {
				h.longConnect = false
				return
			}

			switch resp.Header.Get("Content-Type") {
			case "application/xml":
				{

					defer resp.Body.Close()
					h.longConnect = false
					body, err = io.ReadAll(resp.Body)
					if err != nil {
						h.longConnect = false
						return err
					}
					result := ResponseStatusXML{}
					if err = xml.Unmarshal([]byte(body), &result); err != nil {
						h.longConnect = false
						return err
					}

					err = fmt.Errorf("statusCode %d, statusString %s, subStatusCode %s, errorCode %d, errorMsg %s", result.StatusCode, result.StatusString, result.SubStatusCode, result.ErrorCode, result.ErrorMsg)
					h.longConnect = false
					return err
				}
			default:
				{
					defer resp.Body.Close()
					for {
						result := Result{}
						streamBuf := make([]byte, 1024*2)
						_, err := resp.Body.Read(streamBuf)
						// 获取本次的流所有数据

						if err != nil {
							continue
						}
						if strings.Contains(string(streamBuf), "--MIME_boundary") && strings.Contains(string(streamBuf), `Content-Type: application/json; charset="UTF-8"`) {
							startIndex := strings.Index(string(streamBuf), "{")
							endIndex := strings.LastIndex(string(streamBuf), "}")
							if startIndex == -1 || endIndex == -1 {
								continue
							}
							if startIndex >= endIndex {
								continue
							}
							streamBuf = streamBuf[startIndex : endIndex+1]
							result.Data = streamBuf
							result.Type = "application/json"
							h.longConnect = true
							// file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", streamBuf))
							if !handle(result) {
								return err
							}
						}

						// file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", streamBuf))
						// result := Result{}
						// result.Data = streamBuf
						// result.Type = resp.Header.Get("Content-Type")
						// // 获取Content-Type：后面的值
						// startIndex := strings.Index(string(streamBuf), "{")
						// endIndex := strings.LastIndex(string(streamBuf), "}")
						// if startIndex == -1 || endIndex == -1 {
						// 	continue
						// }
						// streamBuf = streamBuf[startIndex : endIndex+1]
						// h.longConnect = true
						// file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", streamBuf))
						// if !handle(result) {
						// 	return err
						// }
					}

					// mediatype, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))

					// if err != nil {
					// 	h.longConnect = false
					// 	return err
					// }

					// if !strings.HasPrefix(mediatype, "multipart/mixed") {
					// 	h.longConnect = false
					// 	err = errors.New("No boundary found in Content-Type ssss")
					// 	return err
					// }

					// boundary, ok := params["boundary"]
					// if !ok {
					// 	h.longConnect = false
					// 	err = errors.New("No boundary found in Content-Type")
					// 	return err
					// }

					// reader := multipart.NewReader(resp.Body, boundary)

					// for {
					// 	result := Result{}

					// 	part, err := reader.NextPart()
					// 	if err == io.EOF {
					// 		file.WriteString(fmt.Sprintf("data is EOF\n"))
					// 		continue
					// 	}
					// 	if err != nil {
					// 		result.Type = ""
					// 		result.Data = []byte(err.Error())
					// 		file.WriteString(fmt.Sprintf("data is EOF error\n"))
					// 		continue
					// 	}

					// 	header := part.Header.Get("Content-Type")
					// 	partData, _ := io.ReadAll(part)

					// 	if len(partData) == 0 {
					// 		continue
					// 	}
					// 	file.WriteString(fmt.Sprintf("接收到的数据\n%s\n", partData))

					// 	result.Data = partData
					// 	result.Type = header

					// 	h.longConnect = true
					// 	if !handle(result) {
					// 		return err
					// 	}
					// }

				}
			}
		}
	} else {
		h.longConnect = false
	}
	return
}

func parseHeader(header string) map[string]string {
	authInfo := make(map[string]string)
	parts := strings.Split(header, " ")
	for _, part := range parts[1:] {
		parts := strings.SplitN(part, "=", 2)
		key := strings.Trim(parts[0], " ,")
		value := strings.Trim(parts[1], `"`)
		// 去除value的,号
		value = strings.ReplaceAll(value, ",", "")
		value = strings.ReplaceAll(value, "\"", "")
		authInfo[key] = value
	}
	return authInfo
}

func calculateResponse(username, password, realm, method, uri, qop, nonce, nc, cnonce string) string {
	h1 := md5.Sum([]byte(username + ":" + realm + ":" + password))
	h2 := md5.Sum([]byte(method + ":" + uri))
	response := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%x:%s:%s:%s:%s:%x", h1, nonce, nc, cnonce, qop, h2))))
	return response
}

func generateCnonce() (string, error) {
	// 生成随机字节序列
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// 将随机字节序列编码为十六进制字符串
	cnonce := hex.EncodeToString(randomBytes)
	return cnonce, nil
}
