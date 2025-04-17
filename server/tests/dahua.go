package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"server/pkg/dh"
	"strings"
)

func main() {

	// ITC436-PW9H-Z
	c := dh.New(dh.Config{
		Username: "admin",
		Password: "Kygs12345",
		Host:     "192.168.1.108",
		Port:     "80",
	})

	l, err := c.GetCar(dh.GetCarReq{
		PlateNumber: "粤A12345",
		Name:        "TrafficBlackList",
	})
	if err != nil {
		panic(err)
	}
	for _, v := range l.Records {
		err = c.Delete(dh.DeleteCar{
			Name:  "TrafficBlackList",
			Recno: v.Recno,
		})
		if err != nil {
			panic(err)
		}
	}

	// url := "http://192.168.1.108/cgi-bin/magicBox.cgi?action=getLanguageCaps"
	// username := "admin"
	// password := "Kygs12345"

	// c := resty.New()
	// var authParams map[string]string

	// // 中间件：捕获 401 响应并解析 WWW-Authenticate
	// c.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
	// 	if resp.StatusCode() == 401 {
	// 		authHeader := resp.Header().Get("WWW-Authenticate")
	// 		if authHeader != "" {
	// 			authParams = parseDigestHeader(authHeader)
	// 		}
	// 	}
	// 	return nil
	// })

	// // 中间件：在请求前添加自定义 Authorization
	// c.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
	// 	if authParams != nil {
	// 		// 自定义参数
	// 		nc := "00000001"           // 可递增
	// 		cnonce := generateCnonce() // 自定义 cnonce
	// 		uri := "/cgi-bin/magicBox.cgi?action=getLanguageCaps"
	// 		method := req.Method

	// 		// 计算 response
	// 		response := calculateResponse(username, password, authParams["realm"], authParams["nonce"], nc, cnonce, authParams["qop"], method, uri)

	// 		// 构造 Authorization header
	// 		auth := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s", opaque="%s"`,
	// 			username, authParams["realm"], authParams["nonce"], uri, authParams["qop"], nc, cnonce, response, authParams["opaque"])
	// 		req.SetHeader("Authorization", auth)
	// 	}
	// 	return nil
	// })

	// resp, err := c.R().
	// 	SetHeader("User-Agent", "client/1.0").
	// 	Get(url)

	// if err != nil {
	// 	fmt.Println("请求失败:", err)
	// 	return
	// }

	// // 如果需要手动重试（视情况而定）
	// if resp.StatusCode() == 401 && authParams != nil {
	// 	resp, err = c.R().
	// 		SetHeader("User-Agent", "client/1.0").
	// 		Get(url)
	// 	if err != nil {
	// 		fmt.Println("重试请求失败:", err)
	// 		return
	// 	}
	// }

	// // 输出结果
	// fmt.Println("状态码:", resp.StatusCode())
	// fmt.Println("响应内容:", resp.String())

	// 创建 HTTP 客户端
	// client := &http.Client{}

	// // 1. 发送初始请求
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	fmt.Println("创建请求失败:", err)
	// 	return
	// }
	// req.Header.Set("User-Agent", "client/1.0")

	// resp, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println("初始请求失败:", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// // 检查是否收到 401
	// if resp.StatusCode != http.StatusUnauthorized {
	// 	fmt.Println("意外的状态码:", resp.StatusCode)
	// 	return
	// }

	// // 2. 解析 WWW-Authenticate 头
	// authHeader := resp.Header.Get("WWW-Authenticate")
	// if authHeader == "" {
	// 	fmt.Println("未找到 WWW-Authenticate 头")
	// 	return
	// }

	// // 提取参数
	// params := parseDigestHeader(authHeader)
	// realm := params["realm"]
	// nonce := params["nonce"]
	// qop := params["qop"]
	// opaque := params["opaque"]

	// // 3. 生成客户端参数
	// nc := "00000001"           // 请求计数
	// cnonce := generateCnonce() // 随机客户端 nonce
	// uri := "/cgi-bin/magicBox.cgi?action=getLanguageCaps"
	// method := "GET"

	// // 计算 response
	// response := calculateResponse(username, password, realm, nonce, nc, cnonce, qop, method, uri)

	// // 4. 构造 Authorization header
	// auth := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s", opaque="%s"`,
	// 	username, realm, nonce, uri, qop, nc, cnonce, response, opaque)

	// // 5. 发送带认证的请求
	// req, err = http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	fmt.Println("创建认证请求失败:", err)
	// 	return
	// }
	// req.Header.Set("User-Agent", "client/1.0")
	// req.Header.Set("Authorization", auth)

	// resp, err = client.Do(req)
	// if err != nil {
	// 	fmt.Println("认证请求失败:", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// // 输出结果
	// fmt.Println("状态码:", resp.StatusCode)
	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("响应内容:", string(body))
}

// 解析 WWW-Authenticate 头
func parseDigestHeader(header string) map[string]string {
	params := make(map[string]string)
	header = strings.TrimPrefix(header, "Digest ")
	parts := strings.Split(header, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			key := kv[0]
			value := strings.Trim(kv[1], `"`)
			params[key] = value
		}
	}
	return params
}

// 生成随机 cnonce
func generateCnonce() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// 计算 Digest response
func calculateResponse(username, password, realm, nonce, nc, cnonce, qop, method, uri string) string {
	// HA1 = MD5(username:realm:password)
	ha1 := md5.Sum([]byte(username + ":" + realm + ":" + password))
	ha1Str := hex.EncodeToString(ha1[:])

	// HA2 = MD5(method:uri)
	ha2 := md5.Sum([]byte(method + ":" + uri))
	ha2Str := hex.EncodeToString(ha2[:])

	// response = MD5(HA1:nonce:nc:cnonce:qop:HA2)
	resp := md5.Sum([]byte(ha1Str + ":" + nonce + ":" + nc + ":" + cnonce + ":" + qop + ":" + ha2Str))
	return hex.EncodeToString(resp[:])
}
