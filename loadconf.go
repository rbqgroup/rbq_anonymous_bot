package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var config ConfigFile

type ConfigFile struct {
	Ver       int8              `json:"ver"`
	Debug     bool              `json:"debug"`
	Proxy     string            `json:"proxy"`
	Apikey    string            `json:"apikey"`
	Timeout   int               `json:"timeout"`
	Whitelist []int64           `json:"whitelist"`
	To        map[string]string `json:"to"`
	Nitter    string            `json:"nitterHost"`
}

func cmdTChat(cmd string) (bool, string) {
	if len(cmd) == 0 {
		return false, ""
	}
	for k, v := range config.To {
		if cmd == k {
			return strings.HasPrefix(v, "C"), v[1:]
		}
	}
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
	if config.Ver != 1 {
		log.Println("配置檔案版本不符: ", config.Ver)
		return false
	}
	return true
}
