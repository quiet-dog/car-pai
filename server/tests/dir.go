package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func gbkToUtf8(gbkData []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(gbkData), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Data), nil
}

func main() {
	watchDir := "/home/zks/Projects/private/car-pai/server/tests"
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
					return
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

				} else if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("文件修改:", event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("文件/目录删除:", event.Name)
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("文件/目录重命名:", event.Name)
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

	fmt.Println("监听目录:", watchDir)

	// 保持主 goroutine 运行
	select {}
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
