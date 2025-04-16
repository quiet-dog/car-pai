package initialize

import (
	"fmt"
	"os"
	"regexp"
	"server/global"
	"server/model/manage/request"
	"server/service"
	"time"

	"github.com/xuri/excelize/v2"
)

type Car struct {
	CarNum    string
	Color     string
	CarType   string
	Type      string
	CardNum   string
	StartTime string
	EndTime   string
}

type Source struct {
	Excel  string
	Layout string
	Sheet  string
	AreaId uint
}

var columnMap = map[string]string{
	"车牌号码":   "CarNum",
	"车牌颜色":   "Color",
	"车牌类型":   "CarType",
	"名单类型":   "Type",
	"卡号":     "CardNum",
	"有效开始时间": "StartTime",
	"有效结束时间": "EndTime",
}

var columnMap1 = map[string]string{
	"车牌号码": "CarNum",
	"车牌颜色": "Color",
	"车牌类型": "CarType",
	"停车类型": "Type",
	"卡号":   "CardNum",
	"开始时间": "StartTime",
	"截止时间": "EndTime",
}

// 解析 Excel 并映射到结构体
func parseExcelToStruct(filePath, sheet string) ([]Car, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel 内容不足")
	}

	// 获取表头并建立映射
	header := rows[0]
	colIndex := make(map[int]string)
	for i, colName := range header {
		if sheet == "SurveilSheet" || sheet == "sheet1" {
			if field, ok := columnMap[colName]; ok {
				colIndex[i] = field
			}
		} else {
			if field, ok := columnMap1[colName]; ok {
				colIndex[i] = field
			}
		}

	}

	var users []Car
	for _, row := range rows[1:] {
		user := Car{}
		for i, value := range row {
			if field, ok := colIndex[i]; ok {
				switch field {
				case "CarNum":
					user.CarNum = value
				case "Color":
					user.Color = value
				case "CarType":
					user.CarType = value
				case "Type":
					user.Type = value
				case "CardNum":
					user.CardNum = value
				case "StartTime":
					user.StartTime = value
				case "EndTime":
					user.EndTime = value

				}
			}
		}
		users = append(users, user)
	}
	return users, nil
}

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

func Dao() {

	global.TD27_DB.Exec("delete from car_models")
	global.TD27_DB.Exec("delete from car_area")
	global.TD27_DB.Exec("delete from car_device")
	sou := []Source{
		{Excel: "gaotie.xlsx", Layout: time.RFC3339, Sheet: "SurveilSheet", AreaId: 1},
		{Excel: "aoti.xlsx", Layout: time.RFC3339, Sheet: "SurveilSheet", AreaId: 6},
		{Excel: "zhong.xlsx", Layout: time.RFC3339, Sheet: "SurveilSheet", AreaId: 4},
		{Excel: "zongzhanqian.xlsx", Layout: time.RFC3339, Sheet: "sheet1", AreaId: 3},
		{Excel: "zongzhanrukou.xlsx", Layout: time.RFC3339, Sheet: "SurveilSheet", AreaId: 16},
		{Excel: "chengdong.xlsx", Layout: "2006-01-02", Sheet: "车卡资料", AreaId: 5},
		{Excel: "zhencun.xlsx", Layout: "2006-01-02", Sheet: "车卡资料", AreaId: 15},
	}

	carList := []*request.AddCar{}
	ser := service.ServiceGroupApp.Manage.CarService
	// gaotie := 0
	// aoti := 0
	// zhong := 0
	// zongzhanqian := 0
	// zongzhanrukou := 0
	// chengdong := 0
	// zhencun := 0
	for _, o := range sou {
		f, err := parseExcelToStruct(o.Excel, o.Sheet)
		// if o.AreaId == 1 {
		// 	gaotie = len(f)
		// } else if o.AreaId == 6 {
		// 	aoti = len(f)
		// } else if o.AreaId == 4 {
		// 	zhong = len(f)
		// } else if o.AreaId == 3 {
		// 	zongzhanqian = len(f)
		// } else if o.AreaId == 16 {
		// 	zongzhanrukou = len(f)
		// } else if o.AreaId == 5 {
		// 	chengdong = len(f)
		// } else if o.AreaId == 15 {
		// 	zhencun = len(f)
		// }
		// areaIds := []uint{o.AreaId}
		if err != nil {
			fmt.Println(err, o.Sheet, o.Excel)
			os.Exit(0)
			return
		}

		for i, v := range f {
			// layout := "2006-01-02" time.RFC3339
			var t, s time.Time
			fmt.Println("颜色", v.Color)
			if v.StartTime != "" {
				// 去前面的2006-01-02
				v.StartTime = v.StartTime[:10]
			} else {
				v.StartTime = "2006-01-02"
			}

			if v.EndTime != "" {
				// 去前面的2006-01-02
				v.EndTime = v.EndTime[:10]
				fmt.Println("去结束前面的2006-01-02", v.EndTime)
			} else {
				v.EndTime = "2034-02-01"
			}

			layout := "2006-01-02"
			t, err = time.Parse(layout, v.StartTime)
			if err != nil {
				fmt.Println("解析时间错误:", err, i, v.CarNum, v.StartTime)
				os.Exit(0)
				return
			}
			s, err = time.Parse(layout, v.EndTime)
			if err != nil {
				fmt.Println("解析时间错误:", err, i, v.CarNum)
				os.Exit(0)
				return
			}

			color := "0"
			switch v.Color {
			case "蓝色":
				color = "0"
			case "黄色":
				color = "1"
			case "白色":
				color = "2"
			case "黑色":
				color = "3"
			case "绿色":
				color = "4"
			case "其它", "其他":
				color = "255"
			default:
				color = "255"
			}

			catType := "0"
			switch v.CarType {
			case "标准民用用车与军车":
				catType = "0"
			case "02式民用车牌", "92式民用车牌类型":
				catType = "1"
			case "武警车":
				catType = "2"
			case "警车":
				catType = "3"
			case "民用车双行尾牌":
				catType = "4"
			case "使馆车牌":
				catType = "5"
			case "农用车牌":
				fmt.Println("====")
				catType = "6"
			case "摩托车牌":
				catType = "7"
			case "新能源车车牌":
				catType = "8"
			default:
				catType = "1"
			}

			churu := "0"

			switch v.Type {
			case "白名单":
				churu = "0"
			case "黑名单":
				churu = "1"
			}
			v.CarType = catType
			v.Color = color
			v.Type = churu

			// if err = global.TD27_DB.Where("car_num = ? and car_type = ? and color = ?", v.CarNum, v.CarType, v.Color).First(&manage.CarModel{}).Error; err == nil {
			// 	continue
			// }

			isExit := false
			targetCar := &request.AddCar{}
			for _, car := range carList {
				if car.CarNum == v.CarNum && v.CarType == car.CarType && v.Color == car.Color {
					isExit = true
					targetCar = car
					break
				}
			}
			if !isExit {
				req := request.AddCar{
					CarNum:    v.CarNum,
					Color:     v.Color,
					CarType:   v.CarType,
					ListType:  v.Type,
					CardNo:    v.CardNum,
					StartTime: t.UnixMilli(),
					EndTime:   s.UnixMilli(),
					AreaIDs:   []uint{o.AreaId},
				}
				carList = append(carList, &req)
			} else {
				carPaiIsExit := false
				for _, num := range targetCar.AreaIDs {
					if num == o.AreaId {
						carPaiIsExit = true
						break
					}
				}
				if !carPaiIsExit {
					targetCar.AreaIDs = append(targetCar.AreaIDs, o.AreaId)
				}
			}
		}
	}

	// fmt.Println("高铁:", gaotie, "奥体:", aoti, "中:", zhong, "总站前:", zongzhanqian, "总站入口:", zongzhanrukou, "城东:", chengdong, "镇村:", zhencun)
	// os.Exit(0)

	// targetGaotie := 0
	// targetAoti := 0
	// targetZhong := 0
	// targetZongzhanqian := 0
	// targetZongzhanrukou := 0
	// targetChengdong := 0
	// targetZhencun := 0

	// l, _ := os.OpenFile("car7.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// defer l.Close()
	// // 创建一个新的 Excel 文件
	// f := excelize.NewFile()

	// // 需要创建的 sheet 名称
	// sheets := []string{"奥体", "高铁", "中关村", "总站前", "总站入口", "城东", "镇村"}
	// headers := []string{"车牌号码", "车牌颜色", "车辆类型", "名单类型", "开始时间", "结束时间"}

	// for i, sheet := range sheets {
	// 	sheetName := sheet
	// 	if i == 0 {
	// 		f.SetSheetName("Sheet1", sheetName)
	// 	} else {
	// 		f.NewSheet(sheetName)
	// 	}

	// 	// 设置第一行的标题
	// 	for colIdx, header := range headers {
	// 		cell := fmt.Sprintf("%s1", string(rune('A'+colIdx)))
	// 		f.SetCellValue(sheetName, cell, header)
	// 	}

	// }
	// 高铁行 := 0
	// 奥体行 := 0
	// 中关村行 := 0
	// 总站前行 := 0
	// 总站入口行 := 0
	// 城东行 := 0
	// 镇村行 := 0
	// for _, car := range carList {
	// 	switch car.CarType {
	// 	case "0":
	// 		car.CarType = "标准民用用车与军车"
	// 	case "1":
	// 		car.CarType = "02式民用车牌"
	// 	case "2":
	// 		car.CarType = "武警车"
	// 	case "3":
	// 		car.CarType = "警车"
	// 	case "4":
	// 		car.CarType = "民用车双行尾牌"
	// 	case "5":
	// 		car.CarType = "使馆车牌"
	// 	case "6":
	// 		car.CarType = "农用车牌"
	// 	case "7":
	// 		car.CarType = "摩托车牌"
	// 	case "8":
	// 		car.CarType = "新能源车车牌"
	// 	default:
	// 		car.CarType = "02式民用车牌"
	// 	}

	// 	switch car.Color {
	// 	case "0":
	// 		car.Color = "蓝色"
	// 	case "1":
	// 		car.Color = "黄色"
	// 	case "2":
	// 		car.Color = "白色"
	// 	case "3":
	// 		car.Color = "黑色"
	// 	case "4":
	// 		car.Color = "绿色"
	// 	case "255":
	// 		car.Color = "其它"
	// 	default:
	// 		car.Color = "其它"
	// 	}
	// 	switch car.ListType {
	// 	case "0":
	// 		car.ListType = "白名单"
	// 	case "1":
	// 		car.ListType = "黑名单"
	// 	}

	// 	// fmt.Println("车牌号:", car.CarNum, "颜色:", car.Color, "类型:", car.CarType, "名单类型:", car.ListType, "开始时间:", car.StartTime, "结束时间:", car.EndTime)
	// 	for _, id := range car.AreaIDs {
	// 		switch id {
	// 		case 1:
	// 			f.SetCellValue("高铁", fmt.Sprintf("A%d", 高铁行+2), car.CarNum)
	// 			f.SetCellValue("高铁", fmt.Sprintf("B%d", 高铁行+2), car.Color)
	// 			f.SetCellValue("高铁", fmt.Sprintf("C%d", 高铁行+2), car.CarType)
	// 			f.SetCellValue("高铁", fmt.Sprintf("D%d", 高铁行+2), car.ListType)
	// 			f.SetCellValue("高铁", fmt.Sprintf("E%d", 高铁行+2), time.UnixMilli(car.StartTime).UTC().Format(time.RFC3339))
	// 			f.SetCellValue("高铁", fmt.Sprintf("F%d", 高铁行+2), time.UnixMilli(car.EndTime).UTC().Format(time.RFC3339))
	// 			高铁行++
	// 		case 6:
	// 			f.SetCellValue("奥体", fmt.Sprintf("A%d", 奥体行+2), car.CarNum)
	// 			f.SetCellValue("奥体", fmt.Sprintf("B%d", 奥体行+2), car.Color)
	// 			f.SetCellValue("奥体", fmt.Sprintf("C%d", 奥体行+2), car.CarType)
	// 			f.SetCellValue("奥体", fmt.Sprintf("D%d", 奥体行+2), car.ListType)
	// 			f.SetCellValue("奥体", fmt.Sprintf("E%d", 奥体行+2), time.UnixMilli(car.StartTime).UTC().Format(time.RFC3339))
	// 			f.SetCellValue("奥体", fmt.Sprintf("F%d", 奥体行+2), time.UnixMilli(car.EndTime).UTC().Format(time.RFC3339))
	// 			奥体行++
	// 		case 4:
	// 			f.SetCellValue("中关村", fmt.Sprintf("A%d", 中关村行+2), car.CarNum)
	// 			f.SetCellValue("中关村", fmt.Sprintf("B%d", 中关村行+2), car.Color)
	// 			f.SetCellValue("中关村", fmt.Sprintf("C%d", 中关村行+2), car.CarType)
	// 			f.SetCellValue("中关村", fmt.Sprintf("D%d", 中关村行+2), car.ListType)
	// 			f.SetCellValue("中关村", fmt.Sprintf("E%d", 中关村行+2), time.UnixMilli(car.StartTime).UTC().Format(time.RFC3339))
	// 			f.SetCellValue("中关村", fmt.Sprintf("F%d", 中关村行+2), time.UnixMilli(car.EndTime).UTC().Format(time.RFC3339))
	// 			中关村行++
	// 		case 3:
	// 			f.SetCellValue("总站前", fmt.Sprintf("A%d", 总站前行+2), car.CarNum)
	// 			f.SetCellValue("总站前", fmt.Sprintf("B%d", 总站前行+2), car.Color)
	// 			f.SetCellValue("总站前", fmt.Sprintf("C%d", 总站前行+2), car.CarType)
	// 			f.SetCellValue("总站前", fmt.Sprintf("D%d", 总站前行+2), car.ListType)
	// 			f.SetCellValue("总站前", fmt.Sprintf("E%d", 总站前行+2), time.UnixMilli(car.StartTime).UTC().Format(time.RFC3339))
	// 			f.SetCellValue("总站前", fmt.Sprintf("F%d", 总站前行+2), time.UnixMilli(car.EndTime).UTC().Format(time.RFC3339))
	// 			总站前行++
	// 		case 16:
	// 			f.SetCellValue("总站入口", fmt.Sprintf("A%d", 总站入口行+2), car.CarNum)
	// 			f.SetCellValue("总站入口", fmt.Sprintf("B%d", 总站入口行+2), car.Color)
	// 			f.SetCellValue("总站入口", fmt.Sprintf("C%d", 总站入口行+2), car.CarType)
	// 			f.SetCellValue("总站入口", fmt.Sprintf("D%d", 总站入口行+2), car.ListType)
	// 			f.SetCellValue("总站入口", fmt.Sprintf("E%d", 总站入口行+2), time.UnixMilli(car.StartTime).UTC().Format(time.RFC3339))
	// 			f.SetCellValue("总站入口", fmt.Sprintf("F%d", 总站入口行+2), time.UnixMilli(car.EndTime).UTC().Format(time.RFC3339))
	// 			总站入口行++
	// 		case 5:
	// 			f.SetCellValue("城东", fmt.Sprintf("A%d", 城东行+2), car.CarNum)
	// 			f.SetCellValue("城东", fmt.Sprintf("B%d", 城东行+2), car.Color)
	// 			f.SetCellValue("城东", fmt.Sprintf("C%d", 城东行+2), car.CarType)
	// 			f.SetCellValue("城东", fmt.Sprintf("D%d", 城东行+2), car.ListType)
	// 			f.SetCellValue("城东", fmt.Sprintf("E%d", 城东行+2), time.UnixMilli(car.StartTime).UTC().Format(time.RFC3339))
	// 			f.SetCellValue("城东", fmt.Sprintf("F%d", 城东行+2), time.UnixMilli(car.EndTime).UTC().Format(time.RFC3339))
	// 			城东行++
	// 		case 15:
	// 			f.SetCellValue("镇村", fmt.Sprintf("A%d", 镇村行+2), car.CarNum)
	// 			f.SetCellValue("镇村", fmt.Sprintf("B%d", 镇村行+2), car.Color)
	// 			f.SetCellValue("镇村", fmt.Sprintf("C%d", 镇村行+2), car.CarType)
	// 			f.SetCellValue("镇村", fmt.Sprintf("D%d", 镇村行+2), car.ListType)
	// 			f.SetCellValue("镇村", fmt.Sprintf("E%d", 镇村行+2), time.UnixMilli(car.StartTime).UTC().Format(time.RFC3339))
	// 			f.SetCellValue("镇村", fmt.Sprintf("F%d", 镇村行+2), time.UnixMilli(car.EndTime).UTC().Format(time.RFC3339))
	// 			镇村行++

	// 		}
	// 	}
	// }

	// f.SetActiveSheet(0)
	// if err := f.SaveAs("car.xlsx"); err != nil {
	// 	fmt.Println(err)
	// }
	// os.Exit(0)

	for i, car := range carList {
		deviceIds := []uint{}
		global.TD27_DB.Table("device_models").Where("area_id in (?)", car.AreaIDs).Pluck("id", &deviceIds)
		car.DeviceIDs = deviceIds
		ser.CreateCar(*car)
		// if err != nil {
		// 	if strings.Contains(err.Error(), "车牌号已存在") {
		// 		l.Write([]byte(fmt.Sprintf("车牌号已存在: %d %s %v\n", i, car.CarNum, car.AreaIDs)))
		// 		continue
		// 	}

		// 	l.Write([]byte(fmt.Sprintf("创建车牌号失败: %d %s %v\n", i, car.CarNum, car.AreaIDs)))
		// 	continue
		// 	// fmt.Println("创建车牌号失败", err, car.CarNum)
		// 	// os.Exit(0)
		// 	// return
		// }
		// for _, num := range car.AreaIDs {
		// if num == 1 {
		// 	targetGaotie++
		// } else if num == 6 {
		// 	targetAoti++
		// } else if num == 4 {
		// 	targetZhong++
		// } else if num == 3 {
		// 	targetZongzhanqian++
		// } else if num == 16 {
		// 	targetZongzhanrukou++
		// } else if num == 5 {
		// 	targetChengdong++
		// } else if num == 15 {
		// 	targetZhencun++
		// }
		// }

		// fmt.Println("高铁:", gaotie, "奥体:", aoti, "中:", zhong, "总站前:", zongzhanqian, "总站入口:", zongzhanrukou, "城东:", chengdong, "镇村:", zhencun)
		// fmt.Println("高铁:", targetGaotie, "奥体:", targetAoti, "中:", targetZhong, "总站前:", targetZongzhanqian, "总站入口:", targetZongzhanrukou, "城东:", targetChengdong, "镇村:", targetZhencun)
		fmt.Println("完成:", car.CarNum, i, len(carList), i+1/len(carList))
		// time.Sleep(500 * time.Millisecond)
	}

}
