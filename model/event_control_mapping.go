package model

var (
	EVENT_KEY_CODE_MAPPING = map[Event]byte{
		EVENT_UP:      'w',
		EVENT_DOWN:    's',
		EVENT_LEFT:    'a',
		EVENT_RIGHT:   'd',
		EVENT_CONFIRM: 'h',
		EVENT_CANCEL:  'j',
		EVENT_MENU:    'm',
	}

	KEY_CODE_EVENT_MAPPING = map[byte]Event{
		'w': EVENT_UP,
		's': EVENT_DOWN,
		'a': EVENT_LEFT,
		'd': EVENT_RIGHT,
		'h': EVENT_CONFIRM,
		'j': EVENT_CANCEL,
		'm': EVENT_MENU,
	}
)
