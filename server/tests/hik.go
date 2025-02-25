package main

import (
	"fmt"
	"server/pkg/hk_gateway"
	"time"
)

func main() {

	c, err := hk_gateway.NewOnlyHikClinet(hk_gateway.HikConfig{
		Ip:       "http://10.9.11.112",
		Port:     80,
		Username: "admin",
		Password: "admin12345",
	})
	if err != nil {
		fmt.Println("初始化失败")
		panic(err)
	}

	c.UserCheck()
	// for {

	// }
	data := hk_gateway.SetVCLDataReq{
		VCLDataList: hk_gateway.VCLDataList{
			SingleVCLData: []hk_gateway.SingleVCLData{
				{
					ID:         0,
					RunNum:     0,
					ListType:   1,
					PlateNum:   "苏EEEEEEE",
					PlateColor: 0,
					PlateType:  2,
					CardNo:     "123211",
					// 使用CustomTime，强制输出"0000-00-00T00:00:00Z"
					StartTime: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
					EndTime:   time.Now().UTC().Add(time.Hour).Format("2006-01-02T15:04:05Z"),
				},
			},
		},
	}

	err = c.VCLGetCond(data)
	if err != nil {
		fmt.Println("VCLGetCond失败")
		panic(err)
	}

}
