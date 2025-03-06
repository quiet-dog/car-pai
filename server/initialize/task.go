package initialize

import (
	"server/global"
	"server/model/manage"
	"server/pkg/hk_gateway"
	"time"
)

func Hik() {
	hikModels := []*manage.DeviceModel{}
	global.TD27_DB.Where("type = ?", "海康").Find(&hikModels)
	for _, v := range hikModels {
		areaModel := manage.AreaModel{}
		if err := global.TD27_DB.Where("id = ?", v.AreaId).First(&areaModel).Error; err != nil {
			continue
		}

		client, err := hk_gateway.NewOnlyHikClinet(hk_gateway.HikConfig{
			Ip:   v.Host,
			Port: 80,
		})
		if err != nil {
			continue
		}
		size := 10
		total := 0
		res, err := client.VCLGetList(hk_gateway.VCLGetListReq{
			StartOffSet: 0,
			GetVCLNum:   size,
		})
		if err != nil {
			continue
		}
		if res != nil {
			total = res.TotalNum
			for i := 0; i < total; i += size {
				res, err := client.VCLGetList(hk_gateway.VCLGetListReq{
					StartOffSet: i,
					GetVCLNum:   size,
				})
				if err != nil {
					break
				}
				if res != nil {
					for _, v := range res.VCLDataList.SingleVCLData {
						var carModel manage.CarModel
						global.TD27_DB.Where("car_num = ?", v.PlateNum).First(&carModel)
						if carModel.ID == 0 {
							t, err := time.Parse(time.RFC3339, v.StartTime)
							if err != nil {
								continue
							}

							e, err := time.Parse(time.RFC3339, v.EndTime)
							if err != nil {
								continue
							}

							global.TD27_DB.Create(&manage.CarModel{
								Name:      "未知",
								CarNum:    v.PlateNum,
								Color:     v.PlateColor,
								CarType:   v.PlateType,
								CardNo:    v.CardNo,
								StartTime: t.UnixMilli(),
								EndTime:   e.UnixMilli(),
								Areas:     []*manage.AreaModel{&areaModel},
							})
						}
					}
				}
			}
		}
	}
}

// func InitTask() {
// 	gocron.Every(1).Second().Do()
// }
