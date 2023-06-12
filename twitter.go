package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
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
	flag.StringVar(&url, "d", "", "Twitter URL")
	flag.Parse()
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
		var tweetFullnames []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .tweet-header .fullname").Nodes
		if len(tweetFullnames) == 0 {
			println("解析推文作者暱稱失敗")
			return tweet
		} else {
			tweet.Fullname = tweetFullnames[0].FirstChild.Data
			println("推文作者暱稱:", tweet.Fullname)
		}
		var tweetUsernames []*html.Node = doc.Find(".main-tweet .timeline-item .tweet-body .tweet-header .username").Nodes
		if len(tweetUsernames) == 0 {
			println("解析推文作者帳號失敗")
			return tweet
		} else {
			tweet.Username = tweetUsernames[0].FirstChild.Data
			println("推文作者帳號:", tweet.Username)
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
					println("解析推文統計失敗", tweetStat.LastChild.Data, err)
					return tweet
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
