package input

import "github.com/veandco/go-sdl2/sdl"

func KeyboardState() *Input {
	input := NewInput()
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		setInputDataFromSdlEvent(event, input)
	}
	//state := sdl.GetKeyboardState()

	return input
}

func setInputDataFromSdlEvent(event sdl.Event, input *Input) {
	switch event.(type) {
	case *sdl.QuitEvent:
		input.exit = true
	}
}
