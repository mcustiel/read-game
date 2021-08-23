package timing

import "github.com/veandco/go-sdl2/sdl"

type SdlFramerateController int

func NewSdlFrameRateController(framesPerSecond int) SdlFramerateController {
	var controller SdlFramerateController = SdlFramerateController(framesPerSecond)
	return controller
}

func (frameRate SdlFramerateController) WaitFrameRate() {
	sdl.Delay(uint32(1000 / frameRate))
}

func (frameRate SdlFramerateController) WaitMillis(millis uint32) {
	sdl.Delay(millis)
}
