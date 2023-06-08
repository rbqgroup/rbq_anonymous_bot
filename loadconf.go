package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ConfigFile struct {
	Ver       int8    `json:"ver"`
	Debug     bool    `json:"debug"`
	Proxy     string  `json:"proxy"`
	Apikey    string  `json:"apikey"`
	Timeout   int     `json:"timeout"`
	Whitelist []int64 `json:"whitelist"`
	G         string  `json:"g"`
	G18       string  `json:"g18"`
	C2        string  `json:"c2"`
	C25       string  `json:"c25"`
	C3        string  `json:"c3"`
	GY        string  `json:"gy"`
}

func cmdTChat(cmd string) (bool, string) {
	switch cmd {
	case "/c2":
		return chatType(config.C2)
	case "/c25":
		return chatType(config.C25)
	case "/c3":
		return chatType(config.C3)
	case "/g":
		return chatType(config.G)
	case "/g18":
		return chatType(config.G18)
	case "/gy":
		return chatType(config.GY)
	default:
		return false, ""
	}
}

func chatType(chatConfItem string) (bool, string) {
	return strings.HasPrefix(chatConfItem, "C"), chatConfItem[1:]
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
