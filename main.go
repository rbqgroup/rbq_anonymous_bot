//go:generate goversioninfo -icon=icon.ico -manifest=main.exe.manifest
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

type ChatObj struct {
	ID    int64
	Title string
}

func main() {
	fmt.Println("rbq_anonymous_bot v1.1.0")
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
	bot, err := tgbotapi.NewBotAPIWithClient(config.Apikey, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		log.Printf("連線伺服器出現問題: %s\n", err)
		return
	}

	bot.Debug = config.Debug != -1
	if config.Debug != -1 {
		log.Printf("已開啟除錯模式")
	}
	log.Printf("已登入 %s", bot.Self.UserName)

	go getUpdates(bot)

	signalch := make(chan os.Signal, 1)
	signal.Notify(signalch, os.Interrupt, os.Kill)
	signal := <-signalch
	fmt.Println("收到系統訊號: ", signal)
	if signal == os.Interrupt || signal == os.Kill {
		bot.StopReceivingUpdates()
		client.CloseIdleConnections()
		fmt.Println("終止 BOT")
		os.Exit(0)
	}
}
