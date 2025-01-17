package hotwire

import (
	"hotwire/pkg/log"
	"hotwire/pkg/stt"
	"hotwire/pkg/vtt"

	chippergrpc2 "github.com/digital-dream-labs/api/go/chipperpb"
)

func (s *Server) ProcessIntentGraph(req *vtt.IntentGraphRequest) (*vtt.IntentGraphResponse, error) {
	log.Debug("a bot!")
	strm := stt.NewSpeechStream(req)
	_, err := stt.SelectedSTTProcessor.Process(strm)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Stream.Send(
		&chippergrpc2.IntentGraphResponse{
			IsFinal:      true,
			ResponseType: chippergrpc2.IntentGraphMode_INTENT,
			IntentResult: &chippergrpc2.IntentResult{
				Action:    "intent_greeting_hello",
				QueryText: "bruh",
			},
		},
	)
	return &vtt.IntentGraphResponse{
		Intent: &chippergrpc2.IntentGraphResponse{
			IntentResult: &chippergrpc2.IntentResult{
				Action:    "intent_greeting_hello",
				QueryText: "bruh",
			},
		},
	}, nil
}

func (s *Server) ProcessIntent(req *vtt.IntentRequest) (*vtt.IntentResponse, error) {
	log.Debug("a bot!")
	strm := stt.NewSpeechStream(req)
	_, err := stt.SelectedSTTProcessor.Process(strm)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Stream.Send(
		&chippergrpc2.IntentResponse{
			IsFinal: true,
			IntentResult: &chippergrpc2.IntentResult{
				Action:    "intent_greeting_hello",
				QueryText: "bruh",
			},
		},
	)
	return &vtt.IntentResponse{
		Intent: &chippergrpc2.IntentResponse{
			IsFinal: true,
			IntentResult: &chippergrpc2.IntentResult{
				Action:    "intent_greeting_hello",
				QueryText: "bruh",
			},
		},
	}, nil
}

func (s *Server) ProcessKnowledgeGraph(req *vtt.KnowledgeGraphRequest) (*vtt.KnowledgeGraphResponse, error) {
	return nil, nil
}
