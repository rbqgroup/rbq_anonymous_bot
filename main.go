//go:generate goversioninfo -icon=icon.ico -manifest=main.exe.manifest
package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

var startTime time.Time = time.Now()
var dataCounts []int64 = []int64{0, 0, 0} //收發錯

type ChatObj struct {
	ID    int64
	Title string
}

const botvar string = "v1.2.2"

var bot *tgbotapi.BotAPI

func main() {
	log.Println("rbq_anonymous_bot " + botvar)
	if !loadConfig() {
		return
	}
	client := &http.Client{}
	if len(config.Proxy) > 0 {
		log.Printf("代理伺服器: %s\n", config.Proxy)
		tgProxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			log.Printf("代理伺服器地址配置錯誤: %s\n", err)
		}
		tgDialer, err := proxy.FromURL(tgProxyURL, proxy.Direct)
		if err != nil {
			log.Printf("代理伺服器錯誤: %s\n", err)
		}
		tgTransport := &http.Transport{
			Dial: tgDialer.Dial,
		}
		client.Transport = tgTransport
	}
	var err error
	bot, err = tgbotapi.NewBotAPIWithClient(config.Apikey, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		log.Printf("連線伺服器出現問題: %s\n", err)
		return
	}

	bot.Debug = config.Debug != -1
	if bot.Debug {
		log.Printf("已開啟除錯模式\n")
	}
	log.Printf("已登入 %s %s (%s)\n", bot.Self.FirstName, bot.Self.LastName, bot.Self.UserName)
	initNitter()

	go getUpdates(bot)

	signalch := make(chan os.Signal, 1)
	signal.Notify(signalch, os.Interrupt, os.Kill)
	signal := <-signalch
	log.Println("收到系統訊號: ", signal)
	if signal == os.Interrupt || signal == os.Kill {
		bot.StopReceivingUpdates()
		client.CloseIdleConnections()
		log.Println("終止 BOT")
		os.Exit(0)
	}
}

func logCaches(lines []string) {
	if len(lines) == 0 {
		return
	}
	for _, line := range lines {
		log.Println(line)
	}
}

func in(s string, arr []string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
