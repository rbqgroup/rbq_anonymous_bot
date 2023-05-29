package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ConfigFile struct {
	Ver       int8    `json:"ver"`
	Proxy     string  `json:"proxy"`
	Apikey    string  `json:"apikey"`
	Whitelist []int64 `json:"whitelist"`
	G1        string  `json:"g1"`
	C2        string  `json:"c2"`
	C25       string  `json:"c25"`
	C3        string  `json:"c3"`
}

func cmdTChat(cmd string) string {
	switch cmd {
	case "/c2":
		return config.C2
	case "/c25":
		return config.C25
	case "/c3":
		return config.C3
	case "/g1":
		return config.G1
	default:
		return ""
	}
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
