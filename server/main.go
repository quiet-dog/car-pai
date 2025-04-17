package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"server/core"
	"server/global"
	"server/initialize"
)

func main() {
	global.TD27_VP = core.Viper() // 初始化viper
	global.TD27_LOG = core.Zap()  // 初始化zap日志
	zap.ReplaceGlobals(global.TD27_LOG)
	global.TD27_DB = initialize.Gorm() // gorm连接数据库
	initialize.Redis()                 // 初始化redis
	fmt.Println("初始化redis完成")
	global.TD27_CRON = initialize.InitCron() // 初始化cron
	fmt.Println("初始化cron完成")
	initialize.CheckCron() // start cron entry, if exists
	fmt.Println("初始化CheckCron完成")
	if global.TD27_DB == nil {
		global.TD27_LOG.Error("mysql连接失败，退出程序")
		fmt.Println("连接数据启动服务222")
		os.Exit(127)
	} else {
		fmt.Println("连接数据启动服务44")
		initialize.RegisterTables(global.TD27_DB) // 初始化表
		// 程序结束前关闭数据库链接
		fmt.Println("连接数据启动服务333")
		db, _ := global.TD27_DB.DB()
		fmt.Println("连接数据启动服务")
		defer db.Close()
	}
	go initialize.WatchFile()
	initialize.InitHikGateway() // 初始化海康威视
	initialize.InitDh()

	core.RunServer()
}
