package main

import (
	"fmt"
	"server/pkg/hk_gateway"
	"time"
)

func main() {

	c, err := hk_gateway.NewOnlyHikClinet(hk_gateway.HikConfig{
		Ip:       "http://10.9.0.243",
		Port:     80,
		Username: "admin",
		Password: "admin12345",
		// Ip:       "http://192.168.5.9",
		// Port:     80,
		// Username: "admin",
		// Password: "147258369q",
	})
	if err != nil {
		fmt.Println("初始化失败")
		panic(err)
	}

	// r, err := c.GetDeviceInfo()
	// if err != nil {
	// 	fmt.Println("GetDeviceInfo失败")
	// 	panic(err)
	// }
	// fmt.Println(r)
	// for {

	// }
	data := hk_gateway.SetVCLDataReq{
		VCLDataList: hk_gateway.VCLDataList{
			SingleVCLData: []hk_gateway.SingleVCLData{
				{
					ID:         "",
					RunNum:     "0",
					ListType:   "0",
					PlateNum:   "苏EEEEEEW",
					PlateColor: "0",
					PlateType:  "2",
					CardNo:     "",
					Operation:  "new",
					// 使用CustomTime，强制输出"0000-00-00T00:00:00Z"
					StartTime: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
					EndTime:   time.Now().UTC().Add(time.Hour).Format("2006-01-02T15:04:05Z"),
					// CreateTime:         time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02T15:04:05Z"),
					// EffectiveTime:      time.Now().UTC().Add(time.Hour * 9).Format("2006-01-02T15:04:05Z"),
					// EffectiveStartDate: time.Now().UTC().Format("2006-01-02"),
				},
			},
		},
	}

	// data := hk_gateway.TCG225EVCLGetCondReq{}
	// data.LicensePlateInfo = append(data.LicensePlateInfo, hk_gateway.LicensePlateInfo{
	// 	ID:                 "1",
	// 	PlateColor:         "blue",
	// 	PlateType:          "92TypeCivil",
	// 	ListType:           "allowList",
	// 	LicensePlate:       "苏EEEEEEB",
	// 	CardNo:             "",
	// 	CardID:             "",
	// 	Operation:          "new",
	// 	CreateTime:         time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02T15:04:05Z"),
	// 	EffectiveTime:      time.Now().UTC().Add(time.Hour * 9).Format("2006-01-02T15:04:05Z"),
	// 	EffectiveStartDate: time.Now().UTC().Format("2006-01-02"),
	// })

	// data := hk_gateway.VCLDelCondReq{
	// 	DelVCLCond: 1,
	// 	PlateNum:   "苏EEEEEEE",
	// 	PlateColor: 0,
	// 	PlateType:  2,
	// 	CardNo:     "123211",
	// }

	// data := hk_gateway.TCG225EVCLDelCondReq{
	// 	DeleteAllEnabled: false,
	// 	CompoundCond: hk_gateway.CompoundCond{
	// 		PlateColor:   "blue",
	// 		LicensePlate: "苏EEEEEEB",
	// 	},
	// }
	err = c.TCG2A5EVCLGetCond(data)
	if err != nil {
		fmt.Println("VCLGetCond失败")
		panic(err)
	}

	// data := hk_gateway.ZhuaPaiCMSearchDescription{
	// 	SearchID: uuid.New().String(),
	// 	TrackIDList: hk_gateway.ZhuaPaiTrackIDList{
	// 		TrackIDs: []int{120},
	// 	},
	// 	ContentTypeList: hk_gateway.ZhuaPaiContentTypeList{
	// 		ContentType: []string{"video"},
	// 	},
	// 	MaxResults:           40,
	// 	SearchResultPosition: 280,
	// 	MetadataList: hk_gateway.MetadataList{
	// 		MetadataDescriptor: []string{
	// 			"//recordType.meta.hikvision.com/timing",
	// 		},
	// 	},
	// }
	// r, err := c.ZhuaPaiSeach(data)
	// if err != nil {
	// 	fmt.Println("VCLSearch失败")
	// 	panic(err)
	// }
	// fmt.Println(r)

}
