package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var modeString []string = []string{"0空白", "1文字", "2圖片", "3影片", "4動畫", "5多圖組", "6文件"}

func getUpdates(bot *tgbotapi.BotAPI) {
	var u tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	u.Timeout = config.Timeout
	var updates tgbotapi.UpdatesChannel = bot.GetUpdatesChan(u)

	var timers map[string]*time.Ticker = make(map[string]*time.Ticker)
	var medias map[string][]interface{} = make(map[string][]interface{})
	var tousrs map[string]int64 = make(map[string]int64)

	for update := range updates {
		dataCounts[0]++
		if update.Message == nil || chatID(update, bot) {
			continue
		}
		var text string = update.Message.Text
		if len(update.Message.Caption) > 0 {
			text = update.Message.Caption
		}
		var isCommand bool = (len(text) > 0 && text[0] == '/')
		if isCommand {
			var isOK bool = false
			for _, id := range config.Whitelist {
				if update.Message.From.ID == id {
					isOK = true
				}
			}
			if !isOK {
				log.Println("未授權的使用者使用命令: ", update.Message.From.ID)
				continue
			}
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
		var defaultTo bool = false

		// if update.Message.IsCommand() {
		// 	var userCommand string = update.Message.Command()
		// 	var userCommandArg string = update.Message.CommandArguments()
		// 	log.Println("收到指令: ", userCommand, userCommandArg)
		// }

		// fromChat = ChatObj{ID: update.Message.Chat.ID, Title: update.Message.Chat.UserName}
		// fromUser = ChatObj{ID: update.Message.From.ID, Title: update.Message.From.UserName}
		var logCache []string = []string{}
		var message *tgbotapi.Message = update.Message
		logCache = append(logCache, fmt.Sprintf("收到來自會話 %s(%d) 裡 %s(%d) 的訊息：%s", update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.ID, text))
		_, inTousrs := tousrs[message.MediaGroupID]
		if update.Message.ReplyToMessage != nil { //  && update.Message.ReplyToMessage.From.ID == bot.Self.ID
			message = update.Message.ReplyToMessage
			if len(message.Caption) > 0 {
				text = text + " " + message.Caption
			}
			logCache = append(logCache, fmt.Sprintf("需要轉發的訊息來自 %s (%d): %s\n", message.From.UserName, message.From.ID, message.Text))
		}
		var isMediaGroup bool = len(message.MediaGroupID) > 0
		if isMediaGroup {
			log.Println("多圖組: ", message.MediaGroupID)
		}
		if isCommand { //  || update.Message.IsCommand()
			var textUnit []string = strings.Split(text, " ")
			var cmd string = textUnit[0]
			textUnit = textUnit[1:]
			text = strings.Join(textUnit, " ")
			toChannel, toChat = cmdTChat(cmd)
			if len(toChat) == 0 {
				continue
			}
			toChatID, _ = strconv.ParseInt(toChat, 10, 64)
			var logt = fmt.Sprintf("已指定收件人: %d", toChatID)
			if toChannel {
				logt += " (頻道)"
			}
			logCache = append(logCache, logt)
		} else if !inTousrs && update.Message.Chat.ID == update.Message.From.ID && config.DefTo != -1 {
			toChatID = config.DefTo
			logCache = append(logCache, fmt.Sprintf("已指定為預設收件人: %d", toChatID))
			defaultTo = true
		}
		text = filterTwitterURL(text)
		var fileID tgbotapi.FileID
		if message.Photo != nil || message.Video != nil || message.Animation != nil || message.Document != nil {
			var photo tgbotapi.InputMediaPhoto
			var video tgbotapi.InputMediaVideo
			var animation tgbotapi.InputMediaAnimation
			var file tgbotapi.InputMediaDocument
			if message.Video != nil {
				mode = 3
				fileID = tgbotapi.FileID(message.Video.FileID)
				video = tgbotapi.NewInputMediaVideo(fileID)
				if medias[message.MediaGroupID] == nil {
					text = config.Head.Video + text
					video.Caption = text
				}
			} else if message.Photo != nil {
				mode = 2
				fileID = tgbotapi.FileID(message.Photo[0].FileID)
				photo = tgbotapi.NewInputMediaPhoto(fileID)
				if medias[message.MediaGroupID] == nil {
					text = config.Head.Photo + text
					photo.Caption = text
				}
			} else if message.Animation != nil {
				mode = 4
				fileID = tgbotapi.FileID(message.Animation.FileID)
				animation = tgbotapi.NewInputMediaAnimation(fileID)
				if medias[message.MediaGroupID] == nil {
					text = config.Head.Animation + text
					animation.Caption = text
				}
			} else if message.Document != nil {
				mode = 6
				fileID = tgbotapi.FileID(message.Document.FileID)
				file = tgbotapi.NewInputMediaDocument(fileID)
				if medias[message.MediaGroupID] == nil {
					text = config.Head.Document + text
					file.Caption = text
				}
			}
			logCache = append(logCache, fmt.Sprintf("收到的資訊型別: %s", modeString[mode]))
			if isMediaGroup {
				var nMedia []interface{} = make([]interface{}, 0)
				if toChannel || !inTousrs {
					tousrs[message.MediaGroupID] = toChatID
				}
				if mode == 2 {
					if medias[message.MediaGroupID] != nil {
						nMedia = append(medias[message.MediaGroupID], photo)
					} else {
						nMedia = append(nMedia, photo)
					}
				} else if mode == 3 {
					if medias[message.MediaGroupID] != nil {
						nMedia = append(medias[message.MediaGroupID], video)
					} else {
						nMedia = append(nMedia, video)
					}
				} else if mode == 4 {
					if medias[message.MediaGroupID] != nil {
						nMedia = append(medias[message.MediaGroupID], animation)
					} else {
						nMedia = append(nMedia, animation)
					}
				} else if mode == 6 {
					if medias[message.MediaGroupID] != nil {
						nMedia = append(medias[message.MediaGroupID], file)
					} else {
						nMedia = append(nMedia, file)
					}
				}
				medias[update.Message.MediaGroupID] = nMedia
				// log.Println("新增媒體", update.Message.MediaGroupID, len(medias[update.Message.MediaGroupID]))
			}
		} else if len(text) > 0 {
			mode = 1
			text = config.Head.Text + text
			logCache = append(logCache, fmt.Sprintf("收到的資訊型別: %s", modeString[mode]))
			if update.Message.Chat.ID == update.Message.From.ID && len(config.Nitter) > 0 && tweetGETchk(text) {
				logCache = append(logCache, "有且僅有一個推特連結，開始解析。")
				logCaches(logCache)
				go tweetPush(update, bot, text, toChannel, toChat)
				continue
			}
		}
		if toChatID == -1 {
			if update.Message.From.ID == config.Debug {
				toChatID = config.Debug
			} else {
				continue
			}
		}
		if defaultTo && (mode != 2 && mode != 3 && !inTousrs || (len(text) == 0 || !strings.Contains(text, "http") || !strings.Contains(text, "://") || !strings.Contains(text, ".") || !strings.Contains(text, "/"))) && inTousrs {
			logCache = append(logCache, fmt.Sprintf("無效投稿: mode%d: %s\n", mode, text))
			toChatID = update.Message.Chat.ID
			medias = make(map[string][]interface{})
			isMediaGroup = false
			mode = 1
			toChannel = false
			text = "直接发送给我图片或者视频，同时标注来源链接，可以进行投稿。\n投稿会经过编辑审查后，才会发布到频道。\n当前发送的内容无效喵。"
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
			case 6:
				var fileMsg tgbotapi.DocumentConfig = tgbotapi.NewDocument(toChatID, fileID)
				fileMsg.Caption = text
				msg = fileMsg
			default:
				return
			}
			logCaches(logCache)
			if _, err := bot.Send(msg); err != nil {
				dataCounts[2]++
				log.Printf("向 %d 傳送 %s类型 訊息失敗[N]: %s\n", toChatID, modeString[mode], err)
				health(false)
			} else {
				dataCounts[1]++
				log.Printf("已向 %d 傳送 %s类型 訊息: %s\n", toChatID, modeString[mode], text)
				health(true)
			}
		} else {
			var MediaGroupID string = update.Message.MediaGroupID
			if !inTousrs {
				tousrs[MediaGroupID] = toChatID
			}
			if timers[MediaGroupID] == nil {
				newTicker := time.NewTicker(5 * time.Second)
				timers[MediaGroupID] = newTicker
				go func() {
					<-newTicker.C
					// logCache = append(logCache, fmt.Sprintf("提交媒體 %s %d", MediaGroupID, len(medias[MediaGroupID])))
					timers[MediaGroupID].Stop()
					var mediaGroup []interface{} = medias[MediaGroupID]
					if len(mediaGroup) > 0 {
						var mediaGroupMsg tgbotapi.MediaGroupConfig = tgbotapi.NewMediaGroup(tousrs[MediaGroupID], mediaGroup)
						msg = mediaGroupMsg
						logCaches(logCache)
						if _, err := bot.Send(msg); err != nil {
							dataCounts[2]++
							log.Printf("向 %d 傳送多圖訊息失敗[G]: %s", tousrs[MediaGroupID], err)
							health(false)
						} else {
							dataCounts[1]++
							log.Printf("已向 %d 傳送多圖訊息: %s", tousrs[MediaGroupID], text)
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
