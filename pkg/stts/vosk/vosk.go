package vosk

import (
	"hotwire/pkg/log"
	"hotwire/pkg/stt"
	"path/filepath"

	vosk "github.com/kercre123/vosk-api/go"
)

var model *vosk.VoskModel
var rec *vosk.VoskRecognizer

type VoskSTT struct {
	stt.STTProcessor
}

func NewVoskSTT() VoskSTT {
	return VoskSTT{}
}

func (v VoskSTT) Name() string {
	return "Vosk"
}

func (v VoskSTT) Load(path string) error {
	var err error
	model, err = vosk.NewModel(filepath.Join(path, "en-US"))
	if err != nil {
		return err
	}
	rec, err = vosk.NewRecognizer(model, 16000)
	if err != nil {
		return err
	}
	return nil
}

func (v VoskSTT) Unload() error {
	return nil
}

func (v VoskSTT) Process(stream stt.SpeechStream) (string, error) {
	for {
		chunk, err := stream.Read()
		if err != nil {
			return "", err
		}
		rec.AcceptWaveform(chunk)
		if stream.DetectEndOfSpeech(chunk) {
			break
		}
	}
	log.Debug(rec.FinalResult())
	rec.Reset()
	return "", nil
}
