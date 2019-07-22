package event

type EventStore interface {
	Save(event Event) error
}
