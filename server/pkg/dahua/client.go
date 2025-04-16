package dahua

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Client struct {
	Config
	clinet *resty.Client
}

type ReqInitParam struct {
	Url     string
	Query   map[string]string
	Body    interface{}
	Result  interface{}
	Method  string
	Headers map[string]string
}

func (c *Client) do(req ReqInitParam) (err error) {

	client := c.clinet.R()
	client.SetBody(req.Body).
		SetQueryParams(req.Query).
		SetHeaders(req.Headers)
	// SetResult(req.Result)
	// SetDigestAuth(c.Username, c.Password)
	r := &resty.Response{}
	switch req.Method {
	case "POST":
		r, err = client.Post(req.Url)
	case "GET":
		r, err = client.Get(req.Url)
	case "PUT":
		r, err = client.Put(req.Url)
	case "DELETE":
		r, err = client.Delete(req.Url)
	}
	// 打印响应的内容
	fmt.Println("结果", r.String(), err)
	return
}
