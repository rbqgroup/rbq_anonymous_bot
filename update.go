package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var modeString []string = []string{"0空白", "1文字", "2圖片", "3影片", "4動畫", "5多圖組"}

func getUpdates(bot *tgbotapi.BotAPI) {
	var u tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	u.Timeout = config.Timeout
	var updates tgbotapi.UpdatesChannel = bot.GetUpdatesChan(u)

	var timers map[string]*time.Ticker = make(map[string]*time.Ticker)
	var medias map[string][]interface{} = make(map[string][]interface{})
	var tousrs map[string]string = make(map[string]string)

	for update := range updates {
		if update.Message == nil || chatID(update, bot) {
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
		var mode int8 = 0

		var msg tgbotapi.Chattable
		// var fromUser ChatObj
		// var toUser ChatObj
		// var fromChat ChatObj
		// var toChat ChatObj
		var toChat string = ""
		var toChatID int64 = -1
		var toChannel bool = false
		var text string = ""

		// if update.Message.IsCommand() {
		// 	var userCommand string = update.Message.Command()
		// 	var userCommandArg string = update.Message.CommandArguments()
		// 	println("收到指令: ", userCommand, userCommandArg)
		// }

		// fromChat = ChatObj{ID: update.Message.Chat.ID, Title: update.Message.Chat.UserName}
		// fromUser = ChatObj{ID: update.Message.From.ID, Title: update.Message.From.UserName}
		text = update.Message.Text
		if len(update.Message.Caption) > 0 {
			text = update.Message.Caption
		}

		log.Printf("收到來自會話 %s(%d) 裡 %s(%d) 的訊息：%s", update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.ID, text)
		if len(update.Message.MediaGroupID) > 0 {
			log.Println("多圖組: ", update.Message.MediaGroupID)
		}

		if (len(text) > 0 && text[0] == '/') || update.Message.IsCommand() {
			var textUnit []string = strings.Split(text, " ")
			var cmd string = textUnit[0]
			textUnit = textUnit[1:]
			text = strings.Join(textUnit, " ")
			toChannel, toChat = cmdTChat(cmd)
			toChatID, _ = strconv.ParseInt(toChat, 10, 64)
			var log = fmt.Sprintf("已指定收件人: %d", toChatID)
			if toChannel {
				log += " (頻道)"
			}
			println(log)
		}
		text = filterTwitterURL(text)
		var isMediaGroup = len(update.Message.MediaGroupID) > 0
		var fileID tgbotapi.FileID
		if update.Message.Photo != nil || update.Message.Video != nil || update.Message.Animation != nil {
			var photo tgbotapi.InputMediaPhoto
			var video tgbotapi.InputMediaVideo
			var animation tgbotapi.InputMediaAnimation
			if update.Message.Video != nil {
				mode = 3
				fileID = tgbotapi.FileID(update.Message.Video.FileID)
				video = tgbotapi.NewInputMediaVideo(fileID)
				if medias[update.Message.MediaGroupID] == nil {
					text = config.HeadVideo + text
					video.Caption = text
				}
			} else if update.Message.Photo != nil {
				mode = 2
				fileID = tgbotapi.FileID(update.Message.Photo[0].FileID)
				photo = tgbotapi.NewInputMediaPhoto(fileID)
				if medias[update.Message.MediaGroupID] == nil {
					text = config.HeadPhoto + text
					photo.Caption = text
				}
			} else if update.Message.Animation != nil {
				mode = 4
				fileID = tgbotapi.FileID(update.Message.Animation.FileID)
				animation = tgbotapi.NewInputMediaAnimation(fileID)
				if medias[update.Message.MediaGroupID] == nil {
					text = config.HeadAnimation + text
					animation.Caption = text
				}
			}
			println("收到的資訊型別: ", modeString[mode])
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
			text = config.HeadText + text
			println("收到的資訊型別: ", modeString[mode])
			if len(config.Nitter) > 0 && tweetGETchk(text) {
				println("有且僅有一個推特連結，開始解析。")
				go tweetPush(update, bot, text, toChannel, toChat)
				continue
			}
		}
		if toChatID == -1 {
			if config.Debug {
				toChatID = update.Message.Chat.ID
			} else {
				continue
			}
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
				log.Printf("向 %d 傳送 %s类型 訊息失敗: %s\n", toChatID, modeString[mode], err)
				health(false)
			} else {
				log.Printf("已向 %d 傳送 %s类型 訊息: %s\n", toChatID, modeString[mode], text)
				health(true)
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
					var mediaGroup []interface{} = medias[MediaGroupID]
					if len(mediaGroup) > 0 {
						to, _ := strconv.ParseInt(tousrs[MediaGroupID], 10, 64)
						var mediaGroupMsg tgbotapi.MediaGroupConfig = tgbotapi.NewMediaGroup(to, mediaGroup)
						msg = mediaGroupMsg
						if _, err := bot.Send(msg); err != nil {
							log.Printf("向 %d 傳送多圖訊息失敗: %s", to, err)
							health(false)
						} else {
							log.Printf("已向 %d 傳送多圖訊息: %s", to, text)
							health(true)
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
