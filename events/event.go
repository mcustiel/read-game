package events

type Event interface {
	IsQuit() bool
	IsMouseMotion() bool
	IsMouseDown() bool
	IsMouseUp() bool
	IsMouseWheel() bool
	GetEventData() EventData
}

type EventScanner interface {
	GetEvents() []Event
}

type EventData map[string]interface{}
