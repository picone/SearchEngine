package documents

import (
	"sync"
	"time"
	"github.com/huichen/sego"
	"github.com/picone/SearchEngine/utils/mongo"
	mSegment "github.com/picone/SearchEngine/utils/segment"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	IndexingCollection *mgo.Collection
	indexingLock       = sync.Map{}
)

type Indexing struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Keyword   string        `bson:"keyword" json:"keyword,omitempty"`
	Pages     []mgo.DBRef   `bson:"pages" json:"pages,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func init() {
	IndexingCollection = mongo.MongoDatabase.C("indexing")
	IndexingCollection.EnsureIndex(mgo.Index{
		Key:        []string{"keyword"},
		Sparse:     false,
		Unique:     true,
		Background: true,
	})
}

func IndexingAdd(keyword string, id bson.ObjectId) {
	var l *sync.Mutex
	i := Indexing{}
	if lock, ok := indexingLock.Load(keyword); ok {
		l = lock.(*sync.Mutex)
	} else {
		l = &sync.Mutex{}
		indexingLock.Store(keyword, l)
	}
	l.Lock()
	err := IndexingCollection.Find(bson.M{"keyword": keyword}).One(&i)
	if err == nil {
		l.Unlock()
		IndexingCollection.UpdateId(i.Id, bson.M{
			"$addToSet": bson.M{
				"pages": mgo.DBRef{Collection: PageCollection.Name, Id: id},
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		})
	} else {
		i.Id = bson.NewObjectId()
		i.Keyword = keyword
		i.Pages = []mgo.DBRef{
			{Collection: PageCollection.Name, Id: id},
		}
		i.CreatedAt = time.Now()
		i.UpdatedAt = i.CreatedAt
		IndexingCollection.Insert(i)
		l.Unlock()
	}
}

func IndexingFind(keyword string) ([]bson.ObjectId, bool) {
	segments := sego.SegmentsToSlice(mSegment.GetSegmenter().Segment([]byte(keyword)), true)
	indexes := []Indexing{}
	idsMap := make(map[bson.ObjectId]bool)
	if err := IndexingCollection.Find(bson.M{"keyword": bson.M{"$in": segments}}).All(&indexes); err == nil {
		for _, i := range indexes {
			for _, p := range i.Pages {
				idsMap[p.Id.(bson.ObjectId)] = true
			}
		}
	}
	if l := len(idsMap); l == 0 {
		return nil, false
	} else {
		idsSlice := make([]bson.ObjectId, l)
		i := 0
		for id := range idsMap {
			idsSlice[i] = id
			i++
		}
		return idsSlice, true
	}
}
