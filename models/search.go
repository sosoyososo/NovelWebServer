package models

import (
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type State int

const (
	searTypeTitle   State = iota // value --> 0
	searTypeAuthor               // value --> 1
	searTypeSummary              // value --> 2
	searTypeTag                  // value --> 3
)

func SearchNovelsByAuthor(key string, page int) (*[]Novel, error) {
	return searchNovels(key, searTypeAuthor, page)
}
func SearchNovelsByTitle(key string, page int) (*[]Novel, error) {
	return searchNovels(key, searTypeTitle, page)
}
func SearchNovelsBySummary(key string, page int) (*[]Novel, error) {
	return searchNovels(key, searTypeSummary, page)
}
func SearchNovelsByTags(key string, page int) (*[]Novel, error) {
	return searchNovels(key, searTypeTag, page)
}

func searchNovels(key string, s State, page int) (*[]Novel, error) {
	searchKey := ".*" + key + ".*"
	searchField := "title"
	switch s {
	case searTypeAuthor:
		searchField = "author"
	case searTypeSummary:
		searchField = "summary"
	case searTypeTag:
		searchField = "tags"
	}

	findSel := bson.M{searchField: bson.M{"$regex": bson.RegEx{searchKey, ""}}}
	query := db.C(novelCollectionName).Find(findSel).Skip(page * 20).Limit(20)
	_, err := query.Count()
	if err == nil {
		list := make([]Novel, 10)
		err = query.All(&list)
		removeList := []string{"UU看书",
			"https://www.uukanshu.cｏm",
			"https://www.uukanshu.com",
			"www.uukanshu.cｏm",
			"www.uukanshu.com",
			"www.uuｋａnsｈu.ｃom",
			"添加更新提醒，有最新章节时，将会发送邮件到您的邮箱。",
			"(adsbygoogle = window.adsbygoogle || []).push({});",
			"如果喜欢《修真聊天群》，请把网址发给您的朋友。收藏本页请按  Ctrl + D，为方便下次阅读也可把本书添加到桌面，添加桌面请猛击这里。"}
		j := 0
		for j < len(list) {
			i := 0
			for i < len(removeList) {
				list[j].Summary = strings.Replace(list[j].Summary, removeList[i], "", -1)
				i++
			}

			j++
		}
		return &list, err
	} else {
		return nil, err
	}
}
