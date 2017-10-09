package mongo

import (
	"gopkg.in/mgo.v2"
	"log"
)

var (
	mongoSession  *mgo.Session
	MongoDatabase *mgo.Database
)

func init() {
	var err error
	mongoSession, err = mgo.Dial("127.0.0.1:27017/search_engine")
	if err != nil {
		log.Fatal("MongoDB连接失败:", err)
		return
	}
	MongoDatabase = mongoSession.DB("search_engine")
}
