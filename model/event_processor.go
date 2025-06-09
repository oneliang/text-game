package model

type EventProcessor interface {
	Process(event Event) error
}
