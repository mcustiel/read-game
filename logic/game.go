package logic

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/mcustiel/read-game/audio"
	"github.com/mcustiel/read-game/events"
	"github.com/mcustiel/read-game/timing"
	"github.com/mcustiel/read-game/view"
	"github.com/mcustiel/read-game/words"

	"errors"
	"fmt"
	"log"

	"math/rand"
)

const BUTTON_HEIGHT = 40
const BUTTON_WIDTH = 80

type HiddenWord struct {
	button *view.Button
	word   *words.Word
}

type SylabeOption struct {
	button *view.Button
	sylabe *words.Sylabe
}

type BlankPlace struct {
	button *view.Button
	sylabe *words.Sylabe
	filled bool
}

type Game struct {
	data          words.GameDataAccessor
	timer         timing.FrameRateController
	eventScanner  events.EventScanner
	display       view.Display
	audioLoader   audio.AudioLoader
	imageLoader   view.ImageLoader
	hiddenWord    HiddenWord
	options       []SylabeOption
	targets       []BlankPlace
	currentSylabe int
	playing       bool
	previousWord  *words.Word
}

var audioCache map[string]audio.AudioFile
var bgImage view.Image

func NewGame(data words.GameDataAccessor,
	framerateController timing.FrameRateController,
	eventScanner events.EventScanner,
	display view.Display,
	audioLoader audio.AudioLoader,
	imageLoader view.ImageLoader) *Game {

	game := new(Game)
	game.data = data
	game.timer = framerateController
	game.eventScanner = eventScanner
	game.display = display
	game.audioLoader = audioLoader
	game.imageLoader = imageLoader
	game.playing = false
	game.previousWord = nil

	return game
}

func getOptionEventHandler() func(events.Event, ...interface{}) error {
	return func(event events.Event, args ...interface{}) error {
		data := event.GetEventData()
		index := args[1].(int)
		game := args[0].(*Game)
		option := game.options[index]
		button := option.button
		x := data["x"].(int32)
		y := data["y"].(int32)
		if !(x > button.X && x < button.X+button.W && y > button.Y && y < button.Y+button.H) {
			return nil
		}
		hiddenWord := game.hiddenWord
		fmt.Println("Clicked option " + option.sylabe.Sylabe)
		fmt.Println(option.sylabe, hiddenWord.word.Sylabes[game.currentSylabe])
		if option.sylabe == hiddenWord.word.Sylabes[game.currentSylabe] {
			fmt.Println("They are equal")
			target := game.targets[game.currentSylabe]
			fmt.Printf("Target sylabe: %s\n", target.sylabe.Sylabe)
			target.button.Text = target.sylabe.Sylabe
			target.button.BgColor = view.RGBA{255, 255, 255, 255}
			option.button.BgColor = view.RGBA{0, 255, 0, 255}
			option.button.TextColor = view.RGBA{128, 128, 0, 255}
			target.filled = true
			game.currentSylabe = game.currentSylabe + 1
			audioCache["magic"].Play(25, 1)
			audioCache[target.sylabe.AudioFile].Play(100, 2)
		}
		if game.currentSylabe == len(hiddenWord.word.Sylabes) {
			game.playing = false
			game.hiddenWord.button.BgColor = view.RGBA{255, 255, 255, 255}
			game.hiddenWord.button.Text = game.hiddenWord.word.Word
		}
		return nil
	}
}

func getTargetEventHandler() func(events.Event, ...interface{}) error {
	return func(event events.Event, args ...interface{}) error {
		index := args[1].(int)
		game := args[0].(*Game)
		target := game.targets[index]
		button := target.button
		data := event.GetEventData()
		x := data["x"].(int32)
		y := data["y"].(int32)
		if !(x > button.X && x < button.X+button.W && y > button.Y && y < button.Y+button.H) {
			return nil
		}
		fmt.Println("Clicked target rectange: " + target.sylabe.Sylabe)
		if button.Text != "" {
			audioCache[target.sylabe.AudioFile].Play(100, 2)
		}
		return nil
	}
}

func getHiddenWordEventHandler() func(events.Event, ...interface{}) error {
	return func(event events.Event, args ...interface{}) error {
		game := args[0].(*Game)
		hiddenWord := game.hiddenWord
		button := hiddenWord.button
		data := event.GetEventData()
		x := data["x"].(int32)
		y := data["y"].(int32)
		if x > button.X && x < button.X+button.W && y > button.Y && y < button.Y+button.H {
			fmt.Println("Clicked main button: " + button.Text)
			audioCache[hiddenWord.word.AudioFile].Play(100, 2)
		}
		return nil
	}
}

func (game *Game) init() error {
	var err error
	game.currentSylabe = 0
	windowFourth := view.WINDOW_HEIGHT / 4
	bgImage, err = game.imageLoader.Load("images/backgrounds/child-drawing.jpg", view.Rect{800, 600})
	word, err := game.data.GetRandomWord([]*words.Word{game.previousWord})
	if err != nil {
		log.Fatal(err)
		return err
	}
	audioCache = make(map[string]audio.AudioFile)

	sylabes, err := game.data.GetUnusedSylabes(5-len(word.Sylabes), word)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for _, s := range word.Sylabes {
		sylabes = append(sylabes, s)
	}

	rand.Shuffle(len(sylabes), func(i, j int) {
		sylabes[i], sylabes[j] = sylabes[j], sylabes[i]
	})

	audioFile, err := game.audioLoader.LoadAudioFile("audio/fx/sfx-magic2.ogg")
	if err != nil {
		log.Fatal(err)
		return err
	}
	audioCache["magic"] = audioFile
	game.options = make([]SylabeOption, 0)
	game.targets = make([]BlankPlace, 0)
	optionEventHandler := getOptionEventHandler()

	var button *view.Button
	targetEventHandler := getTargetEventHandler()
	separationPx := int32(50)
	buttonTotalWidth := int32(BUTTON_WIDTH + separationPx)
	offset := (view.WINDOW_WIDTH - int32(len(sylabes))*BUTTON_WIDTH - (int32(len(sylabes))-1)*separationPx) / 2
	for i := 0; i < len(sylabes); i++ {
		button = &view.Button{
			Text: sylabes[i].Sylabe,
			Coord: view.Coord{
				X: int32(offset + int32(i)*buttonTotalWidth),
				Y: view.WINDOW_HEIGHT - windowFourth},
			Rect: view.Rect{
				W: BUTTON_WIDTH,
				H: BUTTON_HEIGHT},
			BgColor:     view.RGBA{255, 255, 255, 255},
			BorderColor: view.RGBA{0, 64, 255, 255},
			TextColor:   view.RGBA{0, 64, 255, 255},
			OnClick:     optionEventHandler}
		game.options = append(game.options, SylabeOption{button, sylabes[i]})

		if sylabes[i].AudioFile != "" {
			audioFile, err := game.audioLoader.LoadAudioFile(sylabes[i].AudioFile)
			if err != nil {
				log.Fatal(err)
				return err
			}
			audioCache[sylabes[i].AudioFile] = audioFile
		}
	}

	targetEventHandler = getTargetEventHandler()
	separationPx = int32(25)
	buttonTotalWidth = int32(BUTTON_WIDTH + separationPx)
	offset = (view.WINDOW_WIDTH - int32(len(word.Sylabes))*BUTTON_WIDTH - (int32(len(word.Sylabes))-1)*separationPx) / 2
	for i := 0; i < len(word.Sylabes); i++ {
		button = &view.Button{
			Text: "",
			Coord: view.Coord{
				X: int32(offset + int32(i)*buttonTotalWidth),
				Y: 2 * windowFourth},
			Rect: view.Rect{
				W: BUTTON_WIDTH,
				H: BUTTON_HEIGHT},
			BgColor:     view.RGBA{0, 0, 200, 200},
			BorderColor: view.RGBA{0, 255, 0, 255},
			TextColor:   view.RGBA{64, 64, 0, 0},
			OnClick:     targetEventHandler}
		game.targets = append(game.targets, BlankPlace{button, word.Sylabes[i], false})
	}

	button = &view.Button{
		// Text: word.Word,
		Text: "???????",
		Coord: view.Coord{
			X: (view.WINDOW_WIDTH - 200) / 2,
			Y: windowFourth},
		Rect: view.Rect{
			W: 200,
			H: BUTTON_HEIGHT},
		BgColor:     view.RGBA{0, 255, 0, 255},
		BorderColor: view.RGBA{0, 64, 255, 255},
		TextColor:   view.RGBA{0, 64, 255, 255},
		OnClick:     getHiddenWordEventHandler()}
	game.hiddenWord = HiddenWord{button, word}
	if game.hiddenWord.word.AudioFile == "" {
		return errors.New("Audio file is empty for word" + game.hiddenWord.word.Word)
	}
	audioFile, err = game.audioLoader.LoadAudioFile(game.hiddenWord.word.AudioFile)
	audioCache[game.hiddenWord.word.AudioFile] = audioFile
	return nil
}

func (game *Game) end() error {
	for _, data := range audioCache {
		data.Close()
	}
	audioCache = nil
	bgImage.Close()
	return nil
}

func (game *Game) Play() {
	var eventList []events.Event

	game.init()
	defer game.end()

	exit := false

	againButton := view.Button{
		// Text: word.Word,
		Text: "->",
		Coord: view.Coord{
			X: (view.WINDOW_WIDTH - 200) / 2,
			Y: view.WINDOW_HEIGHT - BUTTON_HEIGHT*1.5},
		Rect: view.Rect{
			W: 200,
			H: BUTTON_HEIGHT},
		BgColor:     view.RGBA{0, 255, 0, 255},
		BorderColor: view.RGBA{0, 64, 255, 255},
		TextColor:   view.RGBA{0, 64, 255, 255},
		OnClick: func(event events.Event, args ...interface{}) error {
			game := args[0].(*Game)
			button := args[1].(view.Button)
			data := event.GetEventData()
			x := data["x"].(int32)
			y := data["y"].(int32)
			if !(x > button.X && x < button.X+button.W && y > button.Y && y < button.Y+button.H) {
				return nil
			}
			fmt.Println("Clicked again button")
			game.init()
			game.playing = true
			return nil
		}}

	for !exit {

		game.playing = true
		for game.playing {
			game.display.Clear()
			game.display.DisplayImage(bgImage, view.Coord{0, 0})

			eventList = game.eventScanner.GetEvents()
			for i := 0; i < len(eventList); i++ {
				event := eventList[i]
				if event.IsQuit() {
					game.playing = false
					exit = true
				} else if event.IsMouseUp() {
					game.hiddenWord.button.OnClick(event, game, game.hiddenWord)
					for i := 0; i < len(game.options); i++ {
						game.options[i].button.OnClick(event, game, i)
					}
					for i := 0; i < len(game.targets); i++ {
						game.targets[i].button.OnClick(event, game, i)
					}
				}
			}
			for i := 0; i < len(game.options); i++ {
				game.display.DrawButton(*game.options[i].button)
			}
			for i := 0; i < len(game.targets); i++ {
				game.display.DrawButton(*game.targets[i].button)
			}
			game.display.DrawButton(*game.hiddenWord.button)
			if !game.playing {
				game.display.DrawButton(againButton)
			}
			game.display.Refresh()

			game.timer.WaitFrameRate()
		}

		for !game.playing && !exit {
			eventList = game.eventScanner.GetEvents()
			for i := 0; i < len(eventList); i++ {
				event := eventList[i]
				if event.IsQuit() {
					game.playing = false
					exit = true
				} else if event.IsMouseUp() {
					againButton.OnClick(event, game, againButton)
				}
			}
		}
	}
}
