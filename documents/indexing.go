package documents

import (
	"ChienHo/SearchEngine/utils/mongo"
	mSegment "ChienHo/SearchEngine/utils/segment"
	"github.com/huichen/sego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
	"log"
)

var (
	IndexingCollection *mgo.Collection
	addCollectionLock  sync.Mutex
)

type Indexing struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Keyword     string        `bson:"keyword" json:"keyword,omitempty"`
	Pages       []mgo.DBRef   `bson:"pages" json:"pages,omitempty"`
	SearchTimes uint64        `bson:"search_times" json:"search_times,omitempty"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
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

func (indexing *Indexing) Add(id bson.ObjectId) {
	indexing.UpdatedAt = time.Now()
	addCollectionLock.Lock()
	n, err := IndexingCollection.Find(bson.M{"keyword": indexing.Keyword}).Count()
	if err == nil && n > 0 {
		addCollectionLock.Unlock()
		IndexingCollection.UpdateId(indexing.Id, bson.M{
			"pages": bson.M{"$addToSet": mgo.DBRef{
				Collection: PageCollection.Name,
				Id:         id,
			}},
		})
	} else {
		indexing.Id = bson.NewObjectId()
		indexing.Pages = []mgo.DBRef{
			{Collection: PageCollection.Name, Id: id},
		}
		indexing.CreatedAt = indexing.UpdatedAt
		IndexingCollection.Insert(indexing)
		addCollectionLock.Unlock()
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
