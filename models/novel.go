package models

import (
	"strings"

	"github.com/go-redis/redis"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Novel struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	URL      string
	Title    string
	Author   string
	Summary  string
	Coverimg string
	Tags     []string
}

type Chapter struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Cateurl string
	Title   string
	URL     string
	Index   int
}

type ChapterCollection struct {
	URLs           []string //collectionName 对应的collection存储内容的小说列表
	CollectionName string
}

var (
	db                          *mgo.Database
	novelCollection             *mgo.Collection
	chapterCollectionCollection *mgo.Collection //保存章节表信息
	chapterCollections          []ChapterCollection
)

const (
	dbAddress           = "127.0.0.1:27017"
	dbName              = "novel"
	novelCollectionName = "novels"
	chapterCCName       = "chaptercollections"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)

func init() {
	dbSetupOnInit()
	err := redisClient.Ping().Err()
	if nil != err {
		panic(err)
	}
}

/*
	model action
*/
func GetNovelsOnPage(page int) ([]Novel, error) {
	query := db.C(novelCollectionName).Find(bson.M{}).Skip(page * 20).Limit(20)
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
		return list, err
	} else {
		return nil, err
	}
}

func getNovelWithURL(url string) (*Novel, error) {
	query := novelCollection.Find(bson.M{"url": url})
	c, err := query.Count()
	if nil != err {
		return nil, err
	} else if c <= 0 {
		return nil, NewModelError(-1, "没有找到小说")
	} else {
		novel := Novel{}
		err = query.One(&novel)
		return &novel, err
	}
}

func getNovelWithId(id string) (*Novel, error) {
	query := novelCollection.FindId(bson.ObjectIdHex(id))
	c, err := query.Count()
	if nil != err {
		return nil, err
	} else if c <= 0 {
		return nil, NewModelError(-1, "没有找到小说")
	} else {
		novel := Novel{}
		err = query.One(&novel)
		return &novel, err
	}
}

func SearchAuthorNovel(author string) []Novel {
	return []Novel{}
}

func SearchTagNovel(author string) []Novel {
	return []Novel{}
}

/*
db action
*/
func dbSetupOnInit() {
	dbSession, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	db = dbSession.DB(dbName)
	novelCollection = db.C(novelCollectionName)
	chapterCollectionCollection = db.C(chapterCCName)

	updateChapterCC()
}
func updateChapterCC() {
	chapterCollections = []ChapterCollection{}
	iter := chapterCollectionCollection.Find(bson.M{}).Iter()
	defer iter.Close()
	chapterCollection := ChapterCollection{}
	for iter.Next(&chapterCollection) {
		chapterCollections = append(chapterCollections, chapterCollection)
	}
}

func (n *Novel) GetChaptersCollection(page int) *mgo.Collection {
	i := 0
	for i < len(chapterCollections) {
		j := 0
		for j < len(chapterCollections[i].URLs) {
			if chapterCollections[i].URLs[j] == n.URL {
				return db.C(chapterCollections[i].CollectionName)
			}
			j++
		}
		i++
	}
	return nil
}
