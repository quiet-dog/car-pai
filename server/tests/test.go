package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	d := "20250312133715169"
	fmt.Println(hikTimeFormatMilli(d))
	fmt.Println(time.UnixMilli(hikTimeFormatMilli(d)).UTC().Format("2006-01-02 15:04:05.000"))
}
func hikTimeFormatMilli(timeSr string) int64 {
	year, _ := strconv.Atoi(timeSr[:4])
	month, _ := strconv.Atoi(timeSr[4:6])
	day, _ := strconv.Atoi(timeSr[6:8])
	hour, _ := strconv.Atoi(timeSr[8:10])
	minute, _ := strconv.Atoi(timeSr[10:12])
	second, _ := strconv.Atoi(timeSr[12:14])
	millisecond, _ := strconv.Atoi(timeSr[14:])
	t := time.Date(year, time.Month(month), day, hour, minute, second, millisecond*1_000_000, time.UTC)
	return t.UTC().UnixMilli()
}
