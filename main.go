package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	var u tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	u.Timeout = config.Timeout
	var updates tgbotapi.UpdatesChannel = bot.GetUpdatesChan(u)

	var timers map[string]*time.Ticker = make(map[string]*time.Ticker)
	var medias map[string][]interface{} = make(map[string][]interface{})
	var tousrs map[string]string = make(map[string]string)

	for update := range updates {
		// for whitelist
		if update.Message == nil {
			continue
		}
		var isOK bool = false
		for _, id := range config.Whitelist {
			println(id, update.Message.Chat.ID)
			if update.Message.Chat.ID == id {
				isOK = true
			}
		}
		if !isOK {
			continue
		}
		// 0: 無訊息 1: 文字訊息 2: 圖片訊息 3: 影片訊息 4: 音訊訊息 5: 檔案訊息
		var mode = 0

		var msg tgbotapi.Chattable
		// var fromUser ChatObj
		// var toUser ChatObj
		// var fromChat ChatObj
		// var toChat ChatObj
		var toChat string = ""
		var toChatID int64 = 0
		var toChannel bool = false
		var text string = ""
		// for i := 0; i < updatesLen; i++ {
		// 	update := updates
		log.Printf("收到來自會話 %s(%d) 裡 %s(%d) 的訊息：%s | %s", update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.ID, update.Message.Text, update.Message.CommandArguments())
		if len(update.Message.MediaGroupID) > 0 {
			log.Println("多圖組: ", update.Message.MediaGroupID)
		}

		// if update.Message.IsCommand() {
		// 	switch update.Message.Command() {
		// 	case "c2":
		// 		toChannel = config.C2
		// 		msg.Text = "c2" + update.Message.CommandArguments()
		// }

		// fromChat = ChatObj{ID: update.Message.Chat.ID, Title: update.Message.Chat.UserName}
		// fromUser = ChatObj{ID: update.Message.From.ID, Title: update.Message.From.UserName}
		text = update.Message.Text
		if len(update.Message.Caption) > 0 {
			text = update.Message.Caption
		}
		if (len(text) > 0 && text[0] == '/') || update.Message.IsCommand() {
			var textUnit []string = strings.Split(text, " ")
			var cmd string = textUnit[0]
			textUnit = textUnit[1:]
			text = strings.Join(textUnit, " ")
			toChannel, toChat = cmdTChat(cmd)
			toChatID, _ = strconv.ParseInt(toChat, 10, 64)
		}
		var isMediaGroup = len(update.Message.MediaGroupID) > 0
		var fileID tgbotapi.FileID
		if update.Message.Photo != nil || update.Message.Video != nil || update.Message.Animation != nil {
			var photo tgbotapi.InputMediaPhoto
			var video tgbotapi.InputMediaVideo
			var animation tgbotapi.InputMediaAnimation
			if update.Message.Photo != nil {
				mode = 2
				fileID = tgbotapi.FileID(update.Message.Photo[0].FileID)
				photo = tgbotapi.NewInputMediaPhoto(fileID)
				if medias[update.Message.MediaGroupID] == nil {
					photo.Caption = text
				}
			} else if update.Message.Video != nil {
				mode = 3
				fileID = tgbotapi.FileID(update.Message.Video.FileID)
				video = tgbotapi.NewInputMediaVideo(fileID)
				if medias[update.Message.MediaGroupID] == nil {
					video.Caption = text
				}
			} else if update.Message.Animation != nil {
				mode = 4
				fileID = tgbotapi.FileID(update.Message.Animation.FileID)
				animation = tgbotapi.NewInputMediaAnimation(fileID)
				if medias[update.Message.MediaGroupID] == nil {
					animation.Caption = text
				}
			}
			println("mode", mode)
			if isMediaGroup {
				var nMedia []interface{} = make([]interface{}, 0)
				if toChannel || len(tousrs[update.Message.MediaGroupID]) == 0 {
					tousrs[update.Message.MediaGroupID] = toChat
				}
				if mode == 2 {
					if medias[update.Message.MediaGroupID] != nil {
						nMedia = append(medias[update.Message.MediaGroupID], photo)
					} else {
						nMedia = append(nMedia, photo)
					}
				} else if mode == 3 {
					if medias[update.Message.MediaGroupID] != nil {
						nMedia = append(medias[update.Message.MediaGroupID], video)
					} else {
						nMedia = append(nMedia, video)
					}
				} else if mode == 4 {
					if medias[update.Message.MediaGroupID] != nil {
						nMedia = append(medias[update.Message.MediaGroupID], animation)
					} else {
						nMedia = append(nMedia, animation)
					}
				}
				medias[update.Message.MediaGroupID] = nMedia
				// println("新增媒體", update.Message.MediaGroupID, len(medias[update.Message.MediaGroupID]))
			}
		} else if len(text) > 0 {
			mode = 1
		}
		if !isMediaGroup {
			switch mode {
			case 1:
				if toChannel {
					msg = tgbotapi.NewMessageToChannel(toChat, text)
				} else {
					msg = tgbotapi.NewMessage(toChatID, text)
				}
			case 2:
				if toChannel {
					var photoMsg tgbotapi.PhotoConfig = tgbotapi.NewPhotoToChannel(toChat, fileID)
					photoMsg.Caption = text
					msg = photoMsg
				} else {
					var photoMsg tgbotapi.PhotoConfig = tgbotapi.NewPhoto(toChatID, fileID)
					photoMsg.Caption = text
					msg = photoMsg
				}
			case 3:
				var videoMsg tgbotapi.VideoConfig = tgbotapi.NewVideo(toChatID, fileID)
				videoMsg.Caption = text
				msg = videoMsg
			case 4:
				var animationMsg tgbotapi.AnimationConfig = tgbotapi.NewAnimation(toChatID, fileID)
				animationMsg.Caption = text
				msg = animationMsg
			default:
				return
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("向 %d 傳送 类型%d 訊息失敗: %s\n", toChatID, mode, err)
			} else {
				log.Printf("已向 %d 傳送 类型%d 訊息: %s\n", toChatID, mode, text)
			}
		} else {
			if timers[update.Message.MediaGroupID] == nil {
				newTicker := time.NewTicker(3 * time.Second)
				timers[update.Message.MediaGroupID] = newTicker
				var MediaGroupID = update.Message.MediaGroupID
				go func() {
					<-newTicker.C
					// println("提交媒體", update.Message.MediaGroupID, len(medias[update.Message.MediaGroupID]))
					timers[MediaGroupID].Stop()
					if len(medias[MediaGroupID]) > 0 {
						to, _ := strconv.ParseInt(tousrs[MediaGroupID], 10, 64)
						var photoMsg tgbotapi.MediaGroupConfig = tgbotapi.NewMediaGroup(to, medias[MediaGroupID])
						msg = photoMsg
						if _, err := bot.Send(msg); err != nil {
							log.Printf("向 %d 傳送多圖訊息失敗: %s", to, err)
						} else {
							log.Printf("已向 %d 傳送多圖訊息: %s", to, text)
						}
					}
					delete(timers, MediaGroupID)
					delete(medias, MediaGroupID)
					delete(tousrs, MediaGroupID)
				}()
			}
		}
	}
}
