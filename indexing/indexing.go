package indexing

import (
	"gopkg.in/mgo.v2/bson"
)

var indexRecords IndexRecords

func init() {

}

//通过关键词增加索引
func Add(key string, value bson.ObjectId) {
	if record, ok := indexRecords.Load(key); ok {
		record.(*indexRecord).lock.Lock()
		defer record.(*indexRecord).lock.Unlock()
		if !record.(*indexRecord).ExistRecord(value) {
			record.(*indexRecord).records = append(record.(*indexRecord).records, value)
		}
	} else {
		indexRecords.Store(key, &indexRecord{
			records: []bson.ObjectId{value},
		})
	}
}

//查找对应关键词包含的文档ID
func Find(key string) ([]bson.ObjectId, bool) {
	records, ok := indexRecords.Load(key)
	if ok {
		return records.(*indexRecord).records, true
	} else {
		return nil, false
	}
}
