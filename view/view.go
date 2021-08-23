package view

import (
	"github.com/mcustiel/read-game/events"
)

const WINDOW_TITLE string = "Learning Go/SDL"
const WINDOW_WIDTH, WINDOW_HEIGHT int32 = 800, 600

type Coord struct {
	X int32
	Y int32
}

type Just int

const (
	LEFT Just = iota
	CENTER
	RIGHT
	TOP
	MIDDLE
	BOTTOM
)

type Rect struct {
	W int32
	H int32
}

type Button struct {
	Coord
	Rect
	Text        string
	BgColor     RGBA
	BorderColor RGBA
	TextColor   RGBA
	OnClick     func(event events.Event, args ...interface{}) error
}

type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Display interface {
	Init() error
	DrawText(text string, pos Coord, color RGBA, hJust Just, vJust Just) error
	DrawRect(pos Coord, size Rect, bgColor RGBA, fgColor RGBA) error
	DrawButton(button Button) error
	DisplayImage(image Image, pos Coord) error
	Refresh()
	Clear() error
	Terminate() error
}

type ImageLoader interface {
	Load(imagePath string, size Rect) (Image, error)
}

type Image interface {
	GetSize() Rect
	Close() error
}
