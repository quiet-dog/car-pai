package dh

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

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
