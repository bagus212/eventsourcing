package event

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"temp/booking/event"
	"time"
)

type Event struct {
	EventID   string `json:"event_id,omitempty"   bson:"_id,omitempty"`
	EventName string `json:"event_name,omitempty" bson:"event_name,omitempty"`
	EventType string `json:"event_type,omitempty" bson:"event_type,omitempty"`
	DataID    string `json:"data_id,omitempty"    bson:"data_id,omitempty"`
	DataName  string `json:"data_name,omitempty"  bson:"data_name,omitempty"`
	Data      string `json:"data,omitempty"       bson:"data,omitempty"`
	CreateBy  string `json:"create_by,omitempty"  bson:"create_by,omitempty"`
	TimeStamp string `json:"timestamp,omitempty"  bson:"timestamp,omitempty"`
}

type EventCreator struct {
	DataID   string
	DataName string
	CreateBy string
}

func (ec EventCreator) Create(eventName, eventType string, data interface{}) (event.Event, error) {
	eventData, err := json.Marshal(data)
	if err != nil {
		return event.Event{}, err
	}
	return event.Event{
		EventID:   ec.generateEventID(),
		EventName: eventName,
		EventType: eventType,
		DataID:    ec.DataID,
		DataName:  ec.DataName,
		CreateBy:  ec.CreateBy,
		Data:      string(eventData),
		TimeStamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (ec EventCreator) generateEventID() string {
	minIDLength := 10000000
	additionalIDLength := 1000000
	rand.Seed(time.Now().UnixNano())
	id := strconv.Itoa(minIDLength + rand.Intn(additionalIDLength))
	return id
}
