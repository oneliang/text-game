package model

const (
	EVENT_NONE  Event = 0
	EVENT_UP    Event = 1
	EVENT_DOWN  Event = 2
	EVENT_LEFT  Event = 3
	EVENT_RIGHT Event = 4

	EVENT_CONFIRM Event = 5
	EVENT_CANCEL  Event = 6
	EVENT_MENU    Event = 7
)

type Event uint32
