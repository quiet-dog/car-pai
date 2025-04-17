package main

import (
	"fmt"
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

var columnMap = map[string]string{
	"车牌号码": "CarNum",
	"车牌颜色": "Color",
	"车辆类型": "CarType",
	"名单类型": "Type",
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
		if field, ok := columnMap[colName]; ok {
			colIndex[i] = field
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

func main() {
	c, err := parseExcelToStruct("chengdong.xlsx", "车卡资料")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range c {
		layout := "2006-01-02"
		t, err := time.Parse(layout, v.StartTime)
		if err != nil {
			fmt.Println("解析时间错误:", err, v.CarNum)
			return
		}
		s, err := time.Parse(layout, v.EndTime)
		if err != nil {
			// fmt.Println("解析时间错误:", err, i, v.CarNum)
			return
		}
		fmt.Println(t.UnixMilli(), s.UnixMilli())
	}
}
