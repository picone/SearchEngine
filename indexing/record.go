package indexing

import (
	//"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

type indexRecord struct {
	lock sync.RWMutex
	records []bson.ObjectId
}

func (record *indexRecord) ExistRecord(value bson.ObjectId) (exist bool) {
	exist = false
	for _, mValue := range record.records {
		if mValue == value {
			exist = true
			return
		}
	}
	return
}

type IndexRecords = sync.Map

/*func (records *IndexRecords) Marshal() ([]byte, error) {
	result := IndexStorage{}
	records.Range(func(key, value interface{}) bool {
		record := &IndexStorageRecord{
			Key: key.(string),
		}
		r := value.(indexRecord).records
		record.Value = []string(r)
		result.Records = append(result.Records, record)
		return true
	})
	return proto.Marshal(&result)
}

func (records *IndexRecords) UnMarshal(data []byte) (*IndexStorage, error) {
	storage := IndexStorage{}
	if err := proto.Unmarshal(data, &storage); err != nil {
		return nil, err
	}
	return &storage, nil
}*/
