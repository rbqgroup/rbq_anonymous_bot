package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/mem"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func chatID(update tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	if update.Message.Text != "/"+bot.Self.UserName && update.Message.Text != "/"+bot.Self.UserName+"@"+bot.Self.UserName {
		return false
	}
	var runTime time.Duration = time.Since(startTime)
	var permissions string = "标准"
	for _, id := range config.Whitelist {
		println(id, update.Message.From.ID)
		if update.Message.From.ID == id {
			permissions = "管理员"
		}
	}
	v, _ := mem.VirtualMemory()
	var failPercent int64 = 0
	var total = dataCounts[1] + dataCounts[2]
	if total > 0 {
		failPercent = dataCounts[1] * 100 / total
	}
	var currentTime time.Time = time.Now()
	var lines []string = []string{
		fmt.Sprintf("%s %s %s", bot.Self.FirstName, bot.Self.LastName, botvar),
		"[系统状态]",
		fmt.Sprintf("服务器时间: %d-%02d-%02d %02d:%02d:%02d", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second()),
		fmt.Sprintf("运行时间: %d天 %d小时 %d分钟 %d秒", int(runTime.Hours()/24), int(runTime.Hours())%24, int(runTime.Minutes())%60, int(runTime.Seconds())%60),
		fmt.Sprintf("内存: %d MB / %d MB (%d %%)", v.Used/1024/1024, v.Total/1024/1024, int(v.UsedPercent)),
		"[BOT 状态]",
		fmt.Sprintf("收取: %d  发送: %d  失败: %d (%d %%)", dataCounts[0], dataCounts[1], dataCounts[2], failPercent),
		fmt.Sprintf("会话: %s%s (%d)", update.Message.Chat.Title, update.Message.Chat.UserName, update.Message.Chat.ID),
		fmt.Sprintf("用户: %s (%d)", update.Message.From.String(), update.Message.From.ID),
		fmt.Sprintf("权限: %s", permissions),
		nitterInfo(),
		"        本 BOT 具有超级绒力。",
	}
	var text string = strings.Join(lines, "\n")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		dataCounts[2]++
		log.Printf("发送消息失败: %s", err)
	} else {
		dataCounts[1]++
	}
	println(text)
	return true
}
