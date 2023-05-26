package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ConfigFile struct {
	Ver         int8   `json:"ver"`
	Proxy       string `json:"proxy"`
	Apikey      string `json:"apikey"`
	TestChat    int64  `json:"testchat"`
	TestChannel string `json:"testchannel"`
	Ch2         string `json:"ch2"`
	Ch25        string `json:"ch25"`
	Ch3         string `json:"ch3"`
}

var config ConfigFile

func loadConfig() bool {
	f, err := os.OpenFile("config.json", os.O_RDONLY, 0600)
	if err != nil {
		log.Println("開啟配置檔案失敗: ", err)
		return false
	} else {
		contentByte, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println("讀取配置檔案失敗: ", err)
			return false
		} else {
			err = json.Unmarshal(contentByte, &config)
			if err != nil {
				log.Println("解析配置檔案失敗: ", err)
				return false
			}
		}
	}
	if config.Ver != 1 {
		log.Println("配置檔案版本不符: ", config.Ver)
		return false
	}
	return true
}
