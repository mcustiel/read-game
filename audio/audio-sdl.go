package audio

import (
	"errors"
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/mix"
)

func Init() error {
	if err := mix.Init(mix.INIT_OGG); err != nil {
		log.Println(err)
		return err
	}

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Quit() {
	mix.CloseAudio()
	mix.Quit()
}

type SdlAudio struct {
	Audio
}

type SdlAudioFile struct {
	audioData *mix.Chunk
}

func NewSdlAudio(resourcesBasePath string) *SdlAudio {
	audio := new(SdlAudio)
	audio.resourcesBasePath = resourcesBasePath
	return audio
}

func CreateAudioFile(resourcePath string) (SdlAudioFile, error) {
	fmt.Printf("Loading file %s\n", resourcePath)
	music, err := mix.LoadWAV(resourcePath)
	if err != nil {
		log.Println(err)
		return SdlAudioFile{nil}, err
	}

	return SdlAudioFile{music}, nil
}

func (audioFile SdlAudioFile) Play(volume int, channel int) error {
	if volume < 0 || volume > 100 {
		return errors.New(fmt.Sprintf("Invalid volume received: %d. Expected a value between 0 and 100", volume))
	}
	if channel < 1 || channel == 0 || channel > 2 {
		return errors.New(fmt.Sprintf("Invalid channel received: %d. Expected a value in the set (-1, 1, 2)", channel))
	}
	mix.Volume(channel, mix.MAX_VOLUME/100*volume)
	if _, err := audioFile.audioData.Play(channel, 0); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (audioFile SdlAudioFile) Close() error {
	audioFile.audioData.Free()
	return nil
}

func (audio *SdlAudio) LoadAudioFile(audioFile string) (AudioFile, error) {
	resourcePath := audio.resourcesBasePath + "/" + audioFile
	fmt.Printf("Loading file %s\n", resourcePath)
	return CreateAudioFile(resourcePath)
}
