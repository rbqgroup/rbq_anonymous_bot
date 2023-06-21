package main

import (
	"fmt"
	"os"
	"time"
)

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func health(isOK bool) {
	var path = config.HealthCheck
	var healthFile *os.File = nil
	var err error = nil
	if isOK {
		if checkFileIsExist(path) {
			healthFile, err = os.OpenFile(path, os.O_WRONLY, 0666)
			if err != nil {
				println("健康檢查記錄檔案開啟失敗:", err)
			}
		} else {
			healthFile, err = os.Create(path)
			if err != nil {
				println("健康檢查記錄檔案建立失敗:", err)
			}
		}
		var timeUnix int64 = time.Now().Unix()
		_, err = healthFile.Write([]byte(fmt.Sprintf("%d", timeUnix)))
		if err != nil {
			println("健康檢查記錄檔案寫入失敗:", err)
		}
		err = healthFile.Close()
		if err != nil {
			println("健康檢查記錄檔案關閉失敗:", err)
		}
	} else {
		if checkFileIsExist(path) {
			err = os.Remove(path)
			if err != nil {
				println("健康檢查記錄檔案刪除失敗:", err)
			}
		}
	}
}
