package dahua

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	Post   = "POST"
	Get    = "GET"
	Put    = "PUT"
	Delete = "DELETE"
)

func (c *Client) ToolgateInfo(q TollgateInfoReq) (err error) {
	req := ReqInitParam{
		Url:    "/NotificationInfo/TollgateInfo",
		Method: Post,
		Body:   q,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	err = c.do(req)
	return
}

func (c *Client) DeviceInfo() (err error) {
	req := ReqInitParam{
		Url:    "/NotificationInfo/DeviceInfo",
		Method: Post,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	err = c.do(req)
	return
}

// func (c *Client) Create() {
// 	req := ReqInitParam{
// 		Url:    "/RPC2",
// 		Method: Post,
// 		Headers: map[string]string{
// 			"Content-Type": "application/json",
// 		},
// 		Body: map[string]interface{}{
// 			"id":      1,
// 			"menthod": "RecordUpdater.insert",
// 			"params": map[string]interface{}{
// 				"record": map[string]interface{}{
// 					"BeginTime":       "2025-03-26 00:00:00",
// 					"CancelTime":      "2025-03-26 00:00:00",
// 					"DepartMent":      " ",
// 					"MasterMent":      "123",
// 					"PlateNumber":     "苏E12345",
// 					"TelephoneNumber": "",
// 				},
// 			},
// 			"session": c.,
// 		},
// 	}
// 	err = c.do(req)
// }

func (c *Client) Login(q LoginReq) (err error) {
	q.Method = "/sp/login"
	q.Info.UserName = c.Username
	q.Info.Password = c.Password
	q.Session = "1ccc"
	req := ReqInitParam{
		Url:    "/sp/login",
		Body:   q,
		Method: Post,
	}

	err = c.do(req)

	return
}

func (c *Client) GetPublicKey() (err error) {
	req := ReqInitParam{
		Url:    "/evo-apigw/evo-oauth/1.0.0/oauth/public-key",
		Method: Get,
	}
	err = c.do(req)
	return
}

func (c *Client) Keeplive() (err error) {
	req := ReqInitParam{
		Url:    "/NotificationInfo/KeepAlive",
		Method: Post,
	}
	err = c.do(req)
	return
}

func RSAEncrypt(origData string, publicKey string) (string, error) {
	publicKeyRAS := ""
	if !strings.Contains(publicKey, "----") {
		publicKeyRAS = "-----BEGIN RSA PUBLIC KEY-----\n" + publicKey + "\n-----END RSA PUBLIC KEY-----"
	}
	block, _ := pem.Decode([]byte(publicKeyRAS))
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", errors.New("public key error")
	}
	pub := pubInterface.(*rsa.PublicKey)

	aimByte, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
	str := base64.StdEncoding.EncodeToString(aimByte)
	return str, nil
}

// func (c *Client) TestLogin() (err error) {
// 	md5Pass := md5.Sum([]byte(c.Password))
// 	data := map[string]interface{}{
// 		"grant_type": "password",
// 		"username":   c.Username,
// 		"password":   string(md5Pass[:]),
// 	}
// 	req := ReqInitParam{
// 		Url:    "/evo-apigw/evo-oauth/1.0.0/oauth/extend/token",
// 		Body:   nil,
// 		Method: Post,
// 	}

// 	err = c.do(req)
// 	return
// }
