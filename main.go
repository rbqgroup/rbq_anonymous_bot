package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

type ChatObj struct {
	ID    int64
	Title string
}

func main() {

	fmt.Println("絨！")
	socks5 := "socks5://localhost:23332"
	client := &http.Client{}
	tgProxyURL, err := url.Parse(socks5)
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

	bot, err := tgbotapi.NewBotAPIWithClient(apikey, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		log.Printf("連線伺服器出現問題: %s\n", err)
		return
	}

	bot.Debug = true

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

func getUpdates(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	timers := make(map[string]*time.Ticker)
	medias := make(map[string][]interface{})
	for update := range updates {
		var mode = 0 // 0: 無訊息 1: 文字訊息 2: 圖片訊息 3: 影片訊息 4: 音訊訊息 5: 檔案訊息
		var msg tgbotapi.Chattable
		var fromUser ChatObj
		// var toUser ChatObj
		// var fromChat ChatObj
		// var toChat ChatObj
		var text string = ""
		// for i := 0; i < updatesLen; i++ {
		// 	update := updates
		log.Print(">>>>>>>>>>>>")
		if update.Message == nil { // 過濾非訊息型別
			continue
		}
		log.Printf("收到來自會話 %s(%d) 裡 %s(%d) 的訊息：%s", update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.ID, update.Message.Text)
		log.Println("組: ", update.Message.MediaGroupID)
		// fromChat = ChatObj{ID: update.Message.Chat.ID, Title: update.Message.Chat.UserName}
		fromUser = ChatObj{ID: update.Message.From.ID, Title: update.Message.From.UserName}
		text = update.Message.Text
		mode = 0
		var isMediaGroup = len(update.Message.MediaGroupID) > 0
		var fileID tgbotapi.FileID
		if update.Message.Photo != nil {
			fileID = tgbotapi.FileID(update.Message.Photo[0].FileID)
			println(fileID, update.Message.Caption)
			var photo tgbotapi.InputMediaPhoto = tgbotapi.NewInputMediaPhoto(fileID)
			if isMediaGroup {
				var nMedia []interface{} = make([]interface{}, 0)
				if medias[update.Message.MediaGroupID] != nil {
					nMedia = append(medias[update.Message.MediaGroupID], photo)
				} else {
					nMedia = append(nMedia, photo)
				}
				medias[update.Message.MediaGroupID] = nMedia
				// println("新增媒體", update.Message.MediaGroupID, len(medias[update.Message.MediaGroupID]))
			}
			mode = 2
			text = update.Message.Caption
		} else {
			mode = 1
		}
		if !isMediaGroup {
			switch mode {
			case 1:
				msg = tgbotapi.NewMessage(testchat, text)
			case 2:
				var photoMsg tgbotapi.PhotoConfig = tgbotapi.NewPhoto(testchat, fileID)
				photoMsg.Caption = text
				msg = photoMsg
			default:
				return
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("傳送訊息失敗: %s", err)
			} else {
				log.Printf("已向 %s(%d) 傳送訊息: %s", fromUser.Title, fromUser.ID, text)
			}
		} else {
			if timers[update.Message.MediaGroupID] == nil {
				newTicker := time.NewTicker(5 * time.Second)
				timers[update.Message.MediaGroupID] = newTicker
				go func() {
					<-newTicker.C
					// println("提交媒體", update.Message.MediaGroupID, len(medias[update.Message.MediaGroupID]))
					timers[update.Message.MediaGroupID].Stop()
					if len(medias[update.Message.MediaGroupID]) > 0 {
						var photoMsg tgbotapi.MediaGroupConfig = tgbotapi.NewMediaGroup(testchat, medias[update.Message.MediaGroupID])
						msg = photoMsg
						if _, err := bot.Send(msg); err != nil {
							log.Printf("傳送訊息失敗: %s", err)
						} else {
							log.Printf("已向 %s(%d) 傳送訊息: %s", fromUser.Title, fromUser.ID, text)
						}
					}
					delete(timers, update.Message.MediaGroupID)
					delete(medias, update.Message.MediaGroupID)
				}()
			}
		}
	}
}
