package main

import (
	"fmt"
	"time"
)

func main() {
	t, err := time.Parse(time.RFC3339, "2034-02-01T01:01:02Z")
	if err != nil {
		fmt.Println("初始化失败")
		panic(err)
	}
	// 转为时间戳
	fmt.Println(t.UnixMilli())
	fmt.Println(time.UnixMilli(t.UnixMilli()).UTC().Format(time.RFC3339))
}
