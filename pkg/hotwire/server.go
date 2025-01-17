package hotwire

import (
	"hotwire/pkg/stt"
)

type Server struct{}

// New returns a new server
func New(sttProc stt.STTProcessor) (*Server, error) {

	stt.SetSTT(sttProc)

	return &Server{}, nil
}
