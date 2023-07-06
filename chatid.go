package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func chatID(update tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	if config.Debug == -1 || update.Message.Text != "/chatid" {
		return false
	}
	var text string = fmt.Sprintf("会话: %s%s (%d)\n用户: %s (%d)\n", update.Message.Chat.Title, update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.From.String(), update.Message.From.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Printf("发送消息失败: %s", err)
	}
	println(text)
	return true
}
