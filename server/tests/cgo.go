package main

import (
	"fmt"
	"regexp"
	"time"
)

func fixInvalidTime(timeStr string) (string, error) {
	// 匹配 "T" 之后的时间部分 (hh:mm:ss)，捕获 ss 部分
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}):(\d{2,})Z`)

	// 使用正则替换不合规的秒数为 00
	fixedTimeStr := re.ReplaceAllString(timeStr, "${1}:00Z")

	// 再次尝试解析
	_, err := time.Parse(time.RFC3339, fixedTimeStr)
	if err != nil {
		return "", fmt.Errorf("修正失败: %v", err)
	}
	return fixedTimeStr, nil
}

func main() {
	invalidTime := "2016-01-01T01:01:122Z"
	fixedTime, err := fixInvalidTime(invalidTime)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fixedTime)
}
