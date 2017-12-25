package models

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"uu/HtmlWorker"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ChapterDetail struct {
	Content string
	Author  string
	Title   string
}

func GetChaptersForNovelWithId(id string, page int) (*[]Chapter, error) {
	novel, err := getNovelWithId(id)
	if nil != err {
		return nil, err
	} else if nil != novel {
		return novel.getChapters(page)
	} else {
		return nil, NewModelError(-1, "没有对应的小说")
	}
}

func GetChapterDetailWithURL(url string) (*ChapterDetail, error) {
	chapter, _ := getChapterContentWithURLFromRedis(url)
	if chapter != nil {
		return chapter, nil
	}
	return getChapterContentWithURL(url)
}

func getChapterContentWithURL(url string) (*ChapterDetail, error) {
	chapter := ChapterDetail{}
	action1 := HtmlWorker.NewAction("div.h1title", func(s *goquery.Selection) {
		chapter.Title = s.Find("#timu").Text()
		chapter.Author = s.Find("span.author").Text()
	})
	action2 := HtmlWorker.NewAction("div.contentbox", func(s *goquery.Selection) {
		text := s.Text()
		removeList := []string{"UU看书",
			"www.uukanshu.cｏm",
			"www.uuｋａnsｈu.ｃom",
			"添加更新提醒，有最新章节时，将会发送邮件到您的邮箱。",
			"(adsbygoogle = window.adsbygoogle || []).push({});",
			"如果喜欢《修真聊天群》，请把网址发给您的朋友。收藏本页请按  Ctrl + D，为方便下次阅读也可把本书添加到桌面，添加桌面请猛击这里。"}
		i := 0
		for i < len(removeList) {
			text = strings.Replace(text, removeList[i], "", -1)
			i++
		}
		text = strings.TrimSpace(text)
		chapter.Content = text
	})
	worker := HtmlWorker.New(url, []HtmlWorker.WorkerAction{action1, action2})

	worker.CookieStrig = "lastread=11356%3D0%3D%7C%7C17203%3D0%3D%7C%7C17151%3D0%3D%7C%7C482%3D0%3D%7C%7C55516%3D10981%3D%u7B2C8%u7AE0%20%u5C38%u53D8; ASP.NET_SessionId=fm1nai0bstdsevx2zoxva3vh; _ga=GA1.2.1243761825.1494000552; _gid=GA1.2.779825662.1496043539; fcip=111"
	worker.Encoder = func(s []byte) ([]byte, error) {
		reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
		d, e := ioutil.ReadAll(reader)
		if e != nil {
			return nil, e
		}
		return d, nil
	}
	var e error
	worker.OnFail = func(err error) {
		e = err
	}
	worker.Run()
	if e == nil {
		saveChapterToRedis(url, chapter)
	}
	return &chapter, e
}

func saveChapterToRedis(url string, chapter ChapterDetail) error {
	serialized, err := json.Marshal(chapter)
	if nil != err {
		return err
	}
	content := string(serialized)
	err = redisClient.Set(url, content, 0).Err()
	return err
}

func getChapterContentWithURLFromRedis(url string) (*ChapterDetail, error) {
	val, err := redisClient.Get(url).Result()
	if err != nil {
		return nil, err
	}
	chaper := ChapterDetail{}
	err = json.Unmarshal([]byte(val), &chaper)
	if nil != err {
		return nil, err
	} else {
		return &chaper, nil
	}
}

func (n *Novel) getChapters(page int) (*[]Chapter, error) {
	c := n.getChapterCollection()
	if c == nil {
		return nil, NewModelError(-1, "没有找到对应的章节")
	} else {
		query := c.Find(bson.M{"cateurl": n.URL}).Sort("index").Skip(page * 20).Limit(20)
		list := make([]Chapter, 10)
		err := query.All(&list)
		if nil != err {
			return nil, err
		} else {
			return &list, nil
		}
	}
}

func (n *Novel) getChapterCollection() *mgo.Collection {
	i := 0
	for i < len(chapterCollections) {
		j := 0
		for j < len(chapterCollections[i].URLs) {
			if 0 == strings.Compare(chapterCollections[i].URLs[j], n.URL) {
				return db.C(chapterCollections[i].CollectionName)
			}
			j++
		}
		i++
	}
	return nil
}
