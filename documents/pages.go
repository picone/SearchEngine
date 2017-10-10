package documents

import (
	"ChienHo/SearchEngine/utils/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var PageCollection *mgo.Collection

type Page struct {
	Id          bson.ObjectId `bson:"_id"`
	Url         string        `bson:"url"`
	Domain      string        `bson:"domain"`
	Title       string        `bson:"title"`
	Keyword     string        `bson:"keyword,omitempty" json:"keyword,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Content     string        `bson:"content" json:"content,omitempty"`
	Rank        uint64        `bson:"rank" json:"rank,omitempty"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}

func init() {
	PageCollection = mongo.MongoDatabase.C("pages")
	PageCollection.EnsureIndex(mgo.Index{
		Key:        []string{"url"},
		Sparse:     false,
		Unique:     true,
		Background: true,
	})
}
