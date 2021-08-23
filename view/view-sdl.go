package view

import (
	//"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SdlDisplay struct {
	fontPath string
	window   *sdl.Window
	renderer *sdl.Renderer
	surface  *sdl.Surface
	font     *ttf.Font
}

func NewSdlDisplay(fontPath string) *SdlDisplay {
	display := new(SdlDisplay)
	display.fontPath = fontPath
	return display
}

func (display *SdlDisplay) Init() error {
	var err error

	sdl.Init(sdl.INIT_AUDIO | sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)

	display.window, err = sdl.CreateWindow(WINDOW_TITLE, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WINDOW_WIDTH, WINDOW_HEIGHT,
		sdl.WINDOW_SHOWN)

	if err != nil {
		log.Fatal(err)
		return err
	}

	display.renderer, err = sdl.CreateRenderer(display.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
		return err
	}

	display.surface, err = display.window.GetSurface()
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err = ttf.Init(); err != nil {
		log.Fatal(err)
		return err
	}

	display.font, err = ttf.OpenFont(display.fontPath, 24)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (display *SdlDisplay) Clear() error {
	return display.renderer.Clear()
}

func (display *SdlDisplay) Refresh() {
	display.renderer.Present()
}

func (display *SdlDisplay) Terminate() error {
	// display.spritesheet.Destroy()
	if err := display.renderer.Destroy(); err != nil {
		return err
	}
	if err := display.window.Destroy(); err != nil {
		return err
	}
	// display.surface.Destroy()
	display.font.Close()
	ttf.Quit()
	return nil
}

func (display *SdlDisplay) DrawText(text string, pos Coord, color RGBA, hJust Just, vJust Just) error {
	var sdlText *sdl.Surface
	var textTexture *sdl.Texture
	//fmt.Println("Drawing text '" + text + "'")
	var err error
	if sdlText, err = display.font.RenderUTF8Blended(text, sdl.Color{R: color.R, G: color.G, B: color.B, A: color.A}); err != nil {
		return err
	}
	defer sdlText.Free()

	x := getXPos(pos, hJust, sdlText)
	y := getYPos(pos, vJust, sdlText)

	if textTexture, err = display.renderer.CreateTextureFromSurface(sdlText); err != nil {
		return err
	}

	display.renderer.Copy(
		textTexture,
		&sdl.Rect{0, 0, sdlText.W, sdlText.H},
		&sdl.Rect{x, y, sdlText.W, sdlText.H})

	return nil
}

func (display *SdlDisplay) DrawRect(pos Coord, size Rect, bgColor RGBA, fgColor RGBA) error {
	var err error
	var color *RGBA
	var rect sdl.Rect
	color, err = createColor(display.renderer.GetDrawColor())
	if err != nil {
		return err
	}
	if display.renderer.SetDrawColor(fgColor.R, fgColor.G, fgColor.B, fgColor.A); err != nil {
		return err
	}
	rect = sdl.Rect{pos.X - 5, pos.Y - 5, size.W + 10, 5}
	if display.renderer.FillRect(&rect); err != nil {
		return err
	}
	rect = sdl.Rect{pos.X - 5, pos.Y - 5, 5, size.H + 10}
	if display.renderer.FillRect(&rect); err != nil {
		return err
	}
	rect = sdl.Rect{pos.X + size.W, pos.Y - 5, 5, size.H + 10}
	if display.renderer.FillRect(&rect); err != nil {
		return err
	}
	rect = sdl.Rect{pos.X - 5, pos.Y + size.H, size.W + 10, 5}
	if display.renderer.FillRect(&rect); err != nil {
		return err
	}

	if display.renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A); err != nil {
		return err
	}
	rect = sdl.Rect{pos.X, pos.Y, size.W, size.H}
	if display.renderer.FillRect(&rect); err != nil {
		return err
	}
	if display.renderer.SetDrawColor(color.R, color.G, color.B, color.A); err != nil {
		return err
	}

	return nil
}

func (display *SdlDisplay) DrawButton(button Button) error {
	var err error

	if err = display.DrawRect(Coord{button.X, button.Y},
		Rect{button.W, button.H},
		button.BgColor,
		button.BorderColor); err != nil {
		return err
	}
	if display.DrawText(button.Text,
		Coord{button.X + button.W/2, button.Y + button.H/2},
		button.TextColor,
		CENTER, MIDDLE); err != nil {
		return err
	}
	return nil
}

func (display *SdlDisplay) DisplayImage(image Image, pos Coord) error {
	size := image.GetSize()
	src := sdl.Rect{0, 0, size.W, size.H}
	dst := sdl.Rect{pos.X, pos.Y, size.W, size.H}

	sdlImage := image.(*SdlImage)
	if sdlImage.Texture == nil {
		texture, err := display.renderer.CreateTextureFromSurface(sdlImage.Image)
		if err != nil {
			log.Fatal(err)
			return err
		}
		sdlImage.Texture = texture
	}

	return display.renderer.Copy(sdlImage.Texture, &src, &dst)
}

func createColor(r, g, b, a uint8, err error) (*RGBA, error) {
	if err != nil {
		return nil, err
	}
	return &RGBA{r, g, b, a}, nil
}

func getXPos(pos Coord, hJust Just, sdlText *sdl.Surface) int32 {
	var x int32 = pos.X // default: LEFT
	switch hJust {
	case CENTER:
		x = pos.X - (sdlText.W / 2)
	case RIGHT:
		x = pos.X - sdlText.W
	}
	return x
}

func getYPos(pos Coord, vJust Just, sdlText *sdl.Surface) int32 {
	var y int32 = pos.Y - sdlText.H // default: TOP
	switch vJust {
	case BOTTOM:
		y = pos.Y
	case MIDDLE:
		y = pos.Y - (sdlText.H / 2)
	}
	return y
}

type SdlImageLoader struct {
	resourcesBasePath string
}

type SdlImage struct {
	Image   *sdl.Surface
	Texture *sdl.Texture
	size    Rect
}

func CreateSdlImageLoader(resourcesBasePath string) (SdlImageLoader, error) {
	return SdlImageLoader{resourcesBasePath}, nil
}

func (imageLoader SdlImageLoader) Load(imagePath string, size Rect) (Image, error) {
	var err error
	var image *sdl.Surface

	image, err = img.Load(imageLoader.resourcesBasePath + "/" + imagePath)
	if err != nil {
		log.Fatal(err)
	}

	return &SdlImage{image, nil, size}, err
}

func (image *SdlImage) Close() error {
	image.Texture.Destroy()
	image.Image.Free()
	return nil
}

func (image *SdlImage) GetSize() Rect {
	return image.size
}
