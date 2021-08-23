package main

import (
	// "database/sql"
	"fmt"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mcustiel/read-game/audio"
	"github.com/mcustiel/read-game/events"
	"github.com/mcustiel/read-game/logic"
	"github.com/mcustiel/read-game/timing"
	"github.com/mcustiel/read-game/view"
	"github.com/mcustiel/read-game/words"
)

func main() {
	fmt.Println("Starting...")

	var err error
	var display *view.SdlDisplay
	var audioPlayer *audio.SdlAudio

	currentExec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return
	}
	execDir := path.Dir(currentExec)
	display = view.NewSdlDisplay(execDir + "/resources/Roboto-Medium.ttf")
	err = display.Init()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer display.Terminate()

	loader := words.NewSqliteLoader(execDir + "/resources/db.sqlite")

	data, err := loader.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = audio.Init()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer audio.Quit()

	resourcesPath := execDir + "/resources"
	imageLoader, err := view.CreateSdlImageLoader(resourcesPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	audioPlayer = audio.NewSdlAudio(resourcesPath)
	game := logic.NewGame(data,
		timing.NewSdlFrameRateController(30),
		events.NewEventScanner(),
		display,
		audioPlayer,
		imageLoader)

	game.Play()
}
