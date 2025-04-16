package initialize

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"server/global"
	"server/model/manage"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func WatchFile() {

	watchDir := global.TD27_CONFIG.Ftp.Watch
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 监听事件处理
	go func() {
		for {

			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				// 处理不同类型的文件事件
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("文件/目录创建:", event.Name)

					handleNewFile(event.Name)

					// 如果新创建的是目录，则递归添加到 watcher
					fileInfo, err := os.Stat(event.Name)
					if err == nil && fileInfo.IsDir() {
						fmt.Println("新建目录，添加监听:", event.Name)
						addAllDirs(watcher, event.Name)
					}

					// 判断文件名字是否是gbk编码

					if err == nil && !fileInfo.IsDir() && strings.Contains(event.Name, "jpg") {
						fmt.Println("开始处理文件:", event.Name)
						fileName := filepath.Base(event.Name)
						dirPath := filepath.Dir(event.Name)
						ipName := filepath.Base(dirPath)
						if !utf8.Valid([]byte(fileName)) {
							utf8Name, err := gbkToUtf8([]byte(fileName))
							if err != nil {
								continue
							}
							filePath := fmt.Sprintf("%s/%s", dirPath, utf8Name)
							if err = os.Rename(event.Name, filePath); err != nil {
								fmt.Println("更改文件名失败:", err)
								continue
							}
						}

						var deviceModel manage.DeviceModel
						if err = global.TD27_DB.Where("host = (?)", fmt.Sprintf("http://%s", ipName)).First(&deviceModel).Error; err != nil {
							fmt.Println("设备不存在:", ipName)
							continue
						}

						if deviceModel.Type == "海康" {
							fmt.Println("海康设备:", ipName)
							utfFileName, err := gbkToUtf8([]byte(fileName))
							if err != nil {
								fmt.Println("转码失败:", err)
								continue
							}
							fmt.Println("转码后文件名:", utfFileName)

							info := strings.Split(utfFileName, "_")
							var carLog manage.CarLogModel
							carLog.DeviceId = deviceModel.ID
							carLog.CarNum = strings.Split(info[1], ".")[0]
							carLog.SubTime = hikTimeFormatMilli(info[0])
							carLog.Uri = fmt.Sprintf("%s/%s", ipName, utfFileName)
							carLog.PlateType = "2"
							if err = global.TD27_DB.Create(&carLog).Error; err != nil {
								fmt.Println("插入失败:", err)
							}
						}

					}

				} else if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("文件修改:", event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("文件/目录删除:", event.Name)
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					// fileName := filepath.Base(event.Name)
					// dirPath := filepath.Dir(event.Name)
					// ipName := filepath.Base(dirPath)

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("错误:", err)

			}
		}
	}()

	// 添加要监听的目录及其所有子目录
	addAllDirs(watcher, watchDir)

	// 监听上下文,如果程序退出的话,结束监听
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
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

func gbkToUtf8(gbkData []byte) (string, error) {

	if utf8.Valid(gbkData) {
		return string(gbkData), fmt.Errorf("不是gbk编码")
	}

	reader := transform.NewReader(bytes.NewReader(gbkData), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Data), nil
}

// 递归添加目录到 watcher
func addAllDirs(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err := watcher.Add(path)
			if err != nil {
				log.Println("无法监听目录:", path, err)
			} else {
				fmt.Println("已添加监听:", path)
			}
		}
		return nil
	})
}

// 处理新文件的逻辑
func handleNewFile(filename string) {
	// 等待文件稳定
	time.Sleep(500 * time.Millisecond)

	info, err := os.Stat(filename)
	if err != nil {
		fmt.Println("无法获取文件信息:", err)
		return
	}

	// 确保是文件而不是目录
	if !info.IsDir() {
		fmt.Println("新文件已创建:", filename)
	}
}
