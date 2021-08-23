package audio

type Audio struct {
	resourcesBasePath string
}

type AudioLoader interface {
	LoadAudioFile(audioFile string) (AudioFile, error)
}

type AudioFile interface {
	Play(volume int, channel int) error
	Close() error
}
