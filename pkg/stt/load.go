package stt

import (
	"hotwire/pkg/vars"
	"path/filepath"
)

var SelectedSTTProcessor STTProcessor

func SetSTT(s STTProcessor) {
	SelectedSTTProcessor = s
	SelectedSTTProcessor.Load(filepath.Join(vars.STTStorage, SelectedSTTProcessor.Name()))
}
