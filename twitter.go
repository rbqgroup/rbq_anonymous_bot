package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Tweet struct {
	Success   bool
	URL       string
	NitterURL string
	Fullname  string
	Username  string
	Time      string
	Comments  int
	Retweets  int
	Quotes    int
	Likes     int
	Content   string
	Images    []string
	Videos    []string
	MediaNum  int8
}

func tweetPush(update tgbotapi.Update, bot *tgbotapi.BotAPI, text string, toChannel bool, toChat string) {
	var tweet Tweet = tweetGET(text)
	if !tweet.Success {
		println("推特解析失敗。")
		return
	}
	toChatID, _ := strconv.ParseInt(toChat, 10, 64)
	var timeStr string = timeFormat(tweet.Time)
	var statText string = fmt.Sprintf("💬%d  🔁%d  🔀%d  💖%d  %s", tweet.Comments, tweet.Retweets, tweet.Quotes, tweet.Likes, timeStr)
	var enter = "\n"
	if strings.Contains(tweet.Content, "\n") {
		enter += "\n"
	}
	text = text + enter + tweet.Content + enter + statText
	var msg tgbotapi.Chattable
	var mode = 0
	if tweet.MediaNum == 0 {
		// 文字
		text = config.HeadText + text
		if toChannel {
			msg = tgbotapi.NewMessageToChannel(toChat, text)
		} else {
			msg = tgbotapi.NewMessage(toChatID, text)
		}
		mode = 1
	} else if tweet.MediaNum == 1 {
		// 單一附件
		if len(tweet.Images) > 0 {
			var file tgbotapi.FileURL = tgbotapi.FileURL(tweet.Images[0])
			var photoMsg tgbotapi.PhotoConfig
			if toChannel {
				photoMsg = tgbotapi.NewPhotoToChannel(toChat, file)
			} else {
				photoMsg = tgbotapi.NewPhoto(toChatID, file)
			}
			text = config.HeadPhoto + text
			photoMsg.Caption = text
			msg = photoMsg
			mode = 2
		} else if len(tweet.Videos) > 0 {
			var file tgbotapi.FileURL = tgbotapi.FileURL(tweet.Videos[0])
			var videoMsg tgbotapi.VideoConfig = tgbotapi.NewVideo(toChatID, file)
			text = config.HeadVideo + text
			videoMsg.Caption = text
			msg = videoMsg
			mode = 3
		}
	} else if tweet.MediaNum > 1 {
		// 多附件
		var files []interface{} = []interface{}{}
		var isCaption = false
		if len(tweet.Videos) > 0 {
			if len(config.HeadVideo) > 0 {
				text = config.HeadVideo + text
			} else if len(config.HeadPhoto) > 0 {
				text = config.HeadPhoto + text
			}
		}
		for _, v := range tweet.Images {
			var file tgbotapi.FileURL = tgbotapi.FileURL(v)
			var photo tgbotapi.InputMediaPhoto = tgbotapi.NewInputMediaPhoto(file)
			if !isCaption {
				photo.Caption = text
				isCaption = true
			}
			files = append(files, photo)
		}
		for _, v := range tweet.Videos {
			var file tgbotapi.FileURL = tgbotapi.FileURL(v)
			var video tgbotapi.InputMediaVideo = tgbotapi.NewInputMediaVideo(file)
			if !isCaption {
				video.Caption = text
				isCaption = true
			}
			files = append(files, video)
		}
		var mediaGroupMsg tgbotapi.MediaGroupConfig = tgbotapi.NewMediaGroup(toChatID, files)
		msg = mediaGroupMsg
		mode = 5
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("向 %d 傳送 %s类型 訊息失敗: %s\n", toChatID, modeString[mode], err)
		health(false)
	} else {
		log.Printf("已向 %d 傳送 %s类型 訊息: %s\n", toChatID, modeString[mode], text)
		health(true)
	}
}

func timeFormat(timeStr string) string {
	var layout string = "Jan 02, 2006 · 3:04 PM UTC" // UTC = GMT + 0
	nTime, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("時間格式化失敗", timeStr, err)
		return timeStr
	}
	nTime = nTime.Add(time.Hour * time.Duration(config.TimeZone))
	layout = "2006年1月2日 PM3:04"
	var newStr string = nTime.Format(layout)
	var th int = nTime.Hour()
	var thC string = "AM"
	if th >= 12 {
		thC = "PM"
	}
	if th >= 0 && th < 6 {
		newStr = strings.Replace(newStr, thC, " 夜间", 1)
	} else if th >= 6 && th < 8 {
		newStr = strings.Replace(newStr, thC, " 早晨", 1)
	} else if th >= 8 && th < 11 {
		newStr = strings.Replace(newStr, thC, " 上午", 1)
	} else if th >= 11 && th < 13 {
		newStr = strings.Replace(newStr, thC, " 中午", 1)
	} else if th >= 13 && th < 18 {
		newStr = strings.Replace(newStr, thC, " 下午", 1)
	} else if th >= 18 && th < 24 {
		newStr = strings.Replace(newStr, thC, " 晚上", 1)
	}
	if len(newStr) == 0 {
		return timeStr
	}
	var timeZoneStr string = ""
	if config.TimeZone < 0 {
		timeZoneStr = fmt.Sprintf("%d", config.TimeZone)
	} else if config.TimeZone > 0 {
		timeZoneStr = fmt.Sprintf("+%d", config.TimeZone)
	}
	newStr += " (GMT" + timeZoneStr + ")"
	var timeEmoji []string = strings.Split("🕛🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚", "")
	var timeEmojiN string = timeEmoji[3]
	if th >= 12 {
		timeEmojiN = timeEmoji[th-12]
	} else {
		timeEmojiN = timeEmoji[th]
	}
	return timeEmojiN + newStr
}

func tweetGETchk(url string) bool {
	if !strings.Contains(url, "twitter.com/") {
		return false
	}
	if len(strings.Split(url, " ")) != 1 {
		return false
	}
	return true
}

func tweetGET(url string) Tweet {
	var tweet Tweet = Tweet{
		Success:   false,
		URL:       url,
		NitterURL: "",
		Fullname:  "",
		Username:  "",
		Time:      "",
		Comments:  0,
		Retweets:  0,
		Quotes:    0,
		Likes:     0,
		Content:   "",
		Images:    []string{},
		Videos:    []string{},
	}
	if url == "" {
		flag.Usage()
		return tweet
	}
	println("目標連結:", tweet.URL)
	if strings.Contains(tweet.URL, "twitter.com/") {
		println("這是一個推特連結，開始清理額外引數。")
		tweet.URL = strings.Split(tweet.URL, "?")[0]
		println("目標連結:", tweet.URL)
		println("選擇 Nitter Node:", config.Nitter)
		tweet.NitterURL = strings.Replace(tweet.URL, "twitter.com/", config.Nitter+"/", 1)
		println("正在載入:", tweet.NitterURL)
		res, err := http.Get(tweet.NitterURL)
		if err != nil {
			println("載入失敗:", tweet.NitterURL)
			return tweet
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			println("載入失敗:", res.StatusCode, res.Status)
			return tweet
		}
		println("載入成功:", res.StatusCode, res.Status)
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			println("解析資料失敗:", err)
		}
		println("解析資料成功，解析推文...")
		var tweetUsernames []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .tweet-header .username").Nodes
		if len(tweetUsernames) == 0 {
			println("解析推文作者帳號失敗")
			return tweet
		} else {
			tweet.Username = tweetUsernames[0].FirstChild.Data
			println("推文作者帳號:", tweet.Username)
		}
		var tweetFullnames []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .tweet-header .fullname").Nodes
		if len(tweetFullnames) == 0 {
			println("解析推文作者暱稱失敗")
			tweet.Fullname = tweetUsernames[0].FirstChild.Data
			println("使用帳號名:", tweet.Username)
		} else {
			tweet.Fullname = tweetFullnames[0].FirstChild.Data
			println("推文作者暱稱:", tweet.Fullname)
		}
		var tweetTimes []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .tweet-published").Nodes
		if len(tweetTimes) == 0 {
			println("解析推文時間失敗")
			return tweet
		} else {
			tweet.Time = tweetTimes[0].FirstChild.Data
			println("推文時間:", tweet.Time)
		}
		var tweetStats []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .tweet-stats .icon-container").Nodes
		if len(tweetStats) == 0 {
			println("解析推文統計失敗")
			return tweet
		} else {
			var tweetStatTitle []string = []string{"回覆", "轉推", "引用", "喜歡"}
			var tweetStatNum []int = []int{0, 0, 0, 0}
			var i int = 0
			for _, tweetStat := range tweetStats {
				if i >= 4 {
					break
				}
				num, err := strconv.Atoi(strings.Replace(strings.Replace(tweetStat.LastChild.Data, " ", "", -1), ",", "", -1))
				if err != nil {
					num = 0
				}
				tweetStatNum[i] = num
				println(tweetStatTitle[i], ":", tweetStatNum[i])
				i++
			}
			tweet.Comments = tweetStatNum[0]
			tweet.Retweets = tweetStatNum[1]
			tweet.Quotes = tweetStatNum[2]
			tweet.Likes = tweetStatNum[3]
		}
		var tweetContents []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .tweet-content").Nodes
		if len(tweetContents) == 0 {
			println("解析推文內容失敗")
			return tweet
		} else {
			var tweetContent *html.Node = tweetContents[0]
			for c := tweetContent.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					tweet.Content += c.Data
				} else if c.Type == html.ElementNode {
					if c.Data == "a" {
						tweet.Content += c.FirstChild.Data
					}
				}
			}
			println("推文內容:", tweet.Content)
		}
		var images []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .attachments .attachment .still-image").Nodes
		if len(images) > 0 {
			println("推文附圖:")
			for _, image := range images {
				for _, attr := range image.Attr {
					if attr.Key == "href" {
						var imageURL string = fmt.Sprintf("https://%s%s\n", config.Nitter, attr.Val)
						tweet.MediaNum++
						tweet.Images = append(tweet.Images, imageURL)
						log.Println(imageURL)
					}
				}
			}
		}
		var videos []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .video-container video").Nodes
		if len(videos) > 0 {
			println("推文附影片:")
			for _, video := range videos {
				for _, attr := range video.Attr {
					if attr.Key == "data-url" {
						var videoURL string = fmt.Sprintf("https://%s%s\n", config.Nitter, attr.Val)
						tweet.MediaNum++
						tweet.Videos = append(tweet.Videos, videoURL)
						log.Println(videoURL)
					}
				}
			}
		}
	}
	tweet.Success = true
	return tweet
}
