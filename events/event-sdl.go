package events

import (
	// "fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type SdlEvent struct {
	isQuit        bool
	isMouseMotion bool
	isMouseDown   bool
	isMouseUp     bool
	isMouseWheel  bool

	eventData map[string]interface{}
}

type SdlEventScanner struct{}

func (evscanner *SdlEventScanner) GetEvents() []Event {
	var events []Event = make([]Event, 0)
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		events = append(events, newSdlEvent(event))
	}
	return events
}

func NewEventScanner() *SdlEventScanner {
	s := new(SdlEventScanner)
	return s
}

func newSdlEvent(sdlEvent sdl.Event) *SdlEvent {
	ev := new(SdlEvent)

	ev.isQuit = false
	ev.isMouseMotion = false
	ev.isMouseDown = false
	ev.isMouseUp = false
	ev.isMouseWheel = false

	switch t := sdlEvent.(type) {
	case *sdl.QuitEvent:
		ev.isQuit = true
		ev.eventData = nil
	case *sdl.MouseMotionEvent:
		//fmt.Printf("MOUSE MOTION: %+v\n", t)
		ev.isMouseMotion = true
		ev.eventData = EventData{
			"mouse":  t.Which,
			"xMov":   t.XRel,
			"yMov":   t.YRel,
			"xStart": t.X,
			"yStart": t.Y,
			"button": t.State,
		}
	case *sdl.MouseButtonEvent:
		//fmt.Printf("MOUSE BUTTON: %+v\n", t)
		ev.isMouseDown = t.State == sdl.PRESSED
		ev.isMouseUp = t.State == sdl.RELEASED
		ev.eventData = EventData{
			"mouse":  t.Which,
			"button": t.Button,
			"x":      t.X,
			"y":      t.Y,
			"clicks": t.Clicks,
		}
	case *sdl.MouseWheelEvent:
		//fmt.Printf("MOUSE WHEEL: %+v\n", t)
		ev.isMouseWheel = true
		var multiplier int32
		if t.Direction == sdl.MOUSEWHEEL_FLIPPED {
			multiplier = -1
		} else {
			multiplier = 1
		}
		ev.eventData = EventData{
			"mouse":     t.Which,
			"xMov":      t.X * multiplier,
			"yMov":      t.Y * multiplier,
			"direction": t.Direction,
		}

	}
	return ev
}

func (e *SdlEvent) IsQuit() bool {
	return e.isQuit
}

func (e *SdlEvent) IsMouseMotion() bool {
	return e.isMouseMotion
}

func (e *SdlEvent) IsMouseDown() bool {
	return e.isMouseDown
}

func (e *SdlEvent) IsMouseUp() bool {
	return e.isMouseUp
}

func (e *SdlEvent) IsMouseWheel() bool {
	return e.isMouseWheel
}

func (e *SdlEvent) GetEventData() EventData {
	return e.eventData
}
