package model

import (
	"errors"
	"fmt"
	"github.com/oneliang/util-golang/concurrent"
)

type EventExecutor struct {
	queue *concurrent.ResourceQueueThread[Event]
}

func NewEventExecutor(eventProcessor EventProcessor, realStopCallback func()) (*EventExecutor, error) {
	if eventProcessor == nil {
		return nil, errors.New("eventProcessor can not be nil")
	}
	eventExecutor := &EventExecutor{}
	eventExecutor.queue = concurrent.NewResourceQueueThread[Event](func(event Event) {
		err := eventProcessor.Process(event)
		if err != nil {
			fmt.Println(fmt.Sprintf("%v", err))
		}
	}, realStopCallback)
	eventExecutor.queue.Start()
	return eventExecutor, nil

}

func (this *EventExecutor) PostEvent(event Event) {
	this.queue.AddResource(event)
}

func (this *EventExecutor) Stop() {
	this.queue.Stop()
}
