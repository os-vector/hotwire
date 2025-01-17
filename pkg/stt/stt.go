package stt

import (
	"errors"
	"hotwire/pkg/audioproc"
	"hotwire/pkg/log"
	"hotwire/pkg/vtt"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/maxhawkins/go-webrtcvad"
)

type STTProcessor interface {
	Name() string

	Load(path string) error
	Unload() error

	Process(stream SpeechStream) (string, error)

	MultiLanguage() bool
	SupportedLanguages() []string

	DownloadableModels() bool
	ModelURLs() map[string]string
}

type SpeechStream struct {
	Device string

	stream interface{}

	inactiveFrames int
	activeFrames   int
	vadInst        *webrtcvad.VAD

	firstChunk []byte

	isPastFirstChunk bool

	audioProc *audioproc.AudioProcessor
}

func NewSpeechStream(req interface{}) SpeechStream {
	var s SpeechStream
	if str, ok := req.(*vtt.IntentRequest); ok {
		s = SpeechStream{
			Device:     str.Device,
			stream:     str.Stream,
			firstChunk: str.FirstReq.InputAudio,
		}
	} else if str, ok := req.(*vtt.KnowledgeGraphRequest); ok {
		s = SpeechStream{
			Device:     str.Device,
			stream:     str.Stream,
			firstChunk: str.FirstReq.InputAudio,
		}
	} else if str, ok := req.(*vtt.IntentGraphRequest); ok {
		s = SpeechStream{
			Device:     str.Device,
			stream:     str.Stream,
			firstChunk: str.FirstReq.InputAudio,
		}
	} else if _, ok := req.(string); ok {
		s = SpeechStream{
			Device: "test",
		}
	} else {
		log.Error("req type is unknown (Init)")
	}
	var err error
	s.vadInst, err = webrtcvad.New()
	if err != nil {
		log.Error("failed to create new VAD instance (Init): ", err)
	}
	s.vadInst.SetMode(3)
	s.audioProc, err = audioproc.NewAudioProcessor(16000, 200, 1)
	if err != nil {
		log.Error("failed to create new audio processor (Init): ", err)
	}
	return s
}

func (s *SpeechStream) Read() ([]byte, error) {
	// returns next chunk in voice stream as pcm
	var chunkData []byte
	if str, ok := s.stream.(pb.ChipperGrpc_StreamingIntentServer); ok {
		if !s.isPastFirstChunk {
			chunkData = s.firstChunk
			s.isPastFirstChunk = true
		} else {
			chunk, err := str.Recv()
			if err != nil {
				log.Error(err)
				return nil, err
			}
			chunkData = chunk.InputAudio
		}
	} else if str, ok := s.stream.(pb.ChipperGrpc_StreamingIntentGraphServer); ok {
		if !s.isPastFirstChunk {
			chunkData = s.firstChunk
			s.isPastFirstChunk = true
		} else {
			chunk, err := str.Recv()
			if err != nil {
				log.Error(err)
				return nil, err
			}
			chunkData = chunk.InputAudio
		}
	} else if str, ok := s.stream.(pb.ChipperGrpc_StreamingKnowledgeGraphServer); ok {
		if !s.isPastFirstChunk {
			chunkData = s.firstChunk
			s.isPastFirstChunk = true
		} else {
			chunk, err := str.Recv()
			if err != nil {
				log.Error(err)
				return nil, err
			}
			chunkData = chunk.InputAudio
		}
	} else {
		log.Error("invalid type")
		return nil, errors.New("invalid type")
	}
	decoded := s.audioProc.ProcessAudio(chunkData)
	return decoded, nil
}

func (s *SpeechStream) DetectEndOfSpeech(chunk []byte) bool {
	inactiveNumMax := 23
	for _, chunk := range audioproc.SplitIntoFrames(chunk, 320) {
		active, err := s.vadInst.Process(16000, chunk)
		if err != nil {
			log.Error("VAD err:")
			log.Error(err)
			return true
		}
		if active {
			s.activeFrames = s.activeFrames + 1
			s.inactiveFrames = 0
		} else {
			s.inactiveFrames = s.inactiveFrames + 1
		}
		if s.inactiveFrames >= inactiveNumMax && s.activeFrames > 18 {
			log.Important("(Bot " + s.Device + ") End of speech detected.")
			return true
		}
	}
	if s.activeFrames < 5 {
		return false
	}
	return false
}
