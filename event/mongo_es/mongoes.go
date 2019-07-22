package mongo_es

import (
	"eventsourcing/event"

	"gopkg.in/mgo.v2"
)

func NewMongoEventStore(col *mgo.Collection) MongoEventStore {
	return MongoEventStore{
		Col: col,
	}
}

type MongoEventStore struct {
	Col *mgo.Collection
}

func (es MongoEventStore) Save(event event.Event) error {
	return es.Col.Insert(event)
}
