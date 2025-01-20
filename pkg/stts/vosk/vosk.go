package vosk

import (
	"encoding/json"
	"hotwire/pkg/log"
	"hotwire/pkg/stt"
	"path/filepath"

	vosk "github.com/kercre123/vosk-api/go"
)

var model *vosk.VoskModel

type VoskSTT struct {
	stt.STTProcessor
}

func NewVoskSTT() VoskSTT {
	return VoskSTT{}
}

func (v VoskSTT) Name() string {
	return "Vosk"
}

// path is ./storage/stt/Vosk
func (v VoskSTT) Load(path string) error {
	var err error
	model, err = vosk.NewModel(filepath.Join(path, "en-US"))
	if err != nil {
		return err
	}
	log.Normal("Vosk en-US model loaded")
	return nil
}

func (v VoskSTT) Unload() error {
	return nil
}

func (v VoskSTT) Process(stream stt.SpeechStream) (string, error) {
	rec, err := vosk.NewRecognizer(model, 16000)
	if err != nil {
		return "", err
	}
	log.Debug("recognizer created")
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
	var resp map[string]string
	json.Unmarshal([]byte(rec.FinalResult()), &resp)
	log.Important("VOSK result for device " + stream.Device + ": " + resp["text"])
	rec.Free()
	return resp["text"], nil
}
