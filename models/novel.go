package models

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Novel struct {
	URL      string
	Title    string
	Author   string
	Summary  string
	Coverimg string
	Tags     []string
}

type Chapter struct {
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

func init() {
	dbSession, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	db = dbSession.DB(dbName)
	novelCollection = db.C(novelCollectionName)
	chapterCollectionCollection = db.C(chapterCCName)

	updateChapterCC()
}

/*
	model action
*/
func (n *Novel) GetChapters(page int) []Novel {
	return []Novel{}
}

func GetNovelsOnPage(page int) ([]Novel, error) {
	query := db.C(novelCollectionName).Find(bson.M{}).Skip(page * 10).Limit(10)
	_, err := query.Count()
	if err == nil {
		list := make([]Novel, 10)
		err = query.All(&list)
		return list, err
	} else {
		return nil, err
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
