package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var config ConfigFile

type ConfigHead struct {
	Video     string `json:"video"`
	Animation string `json:"animation"`
	Photo     string `json:"photo"`
	Document  string `json:"document"`
	Text      string `json:"text"`
}

type ConfigFile struct {
	Ver         int8              `json:"ver"`
	Debug       int64             `json:"debug"`
	HealthCheck string            `json:"healthcheck"`
	TimeZone    int8              `json:"timezone"`
	Proxy       string            `json:"proxy"`
	Apikey      string            `json:"apikey"`
	Timeout     int               `json:"timeout"`
	Whitelist   []int64           `json:"whitelist"`
	To          map[string]string `json:"to"`
	DefTo       int64             `json:"defto"`
	Nitter      []string          `json:"nitterHost"`
	Head        ConfigHead        `json:"head"`
}

func cmdTChat(cmd string) (bool, string) {
	if len(cmd) == 0 {
		return false, ""
	}
	var cmdU []string = strings.Split(cmd[1:], "@")
	if len(cmdU) > 1 && cmdU[1] != bot.Self.UserName {
		return false, ""
	}
	for k, v := range config.To {
		if cmdU[0] == k {
			return strings.HasPrefix(v, "C"), v[1:]
		}
	}
	log.Println("找不到预设发送目标: ", cmd)
	return false, ""
}

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
	if config.Ver != 2 {
		log.Println("配置檔案版本不符: ", config.Ver)
		return false
	}
	return true
}
