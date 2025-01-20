package ttr

import (
	"encoding/json"
	"hotwire/pkg/log"
	"hotwire/pkg/vars"
	"os"
	"path/filepath"
)

var idata IntentData

type IntentData struct {
	Intents []Intent
	Lclztn  LocalizationDetails
}

type Intent struct {
	Intent                 string   `json:"intent"`
	PartialMatchUterrances []string `json:"utterances"`
	ExactMatchUtterances   []string `json:"exactutterances"`
	ExactMatchPreferred    bool     `json:"exactmatchpreferred"`
}

type LocalizationDetails struct {
	Weather struct {
		In               []string `json:"in"`
		Forecast         []string `json:"forecast"`
		Tomorrow         []string `json:"tomorrow"`
		DayAfterTomorrow []string `json:"dayaftertomorrow"`
		Tonight          []string `json:"tonight"`
		ThisAfternoon    []string `json:"thisafternoon"`
	} `json:"weather"`
	EyeColor struct {
		Purple []string `json:"purple"`
		Blue   []string `json:"blue"`
		Yellow []string `json:"yellow"`
		Teal   []string `json:"teal"`
		Green  []string `json:"green"`
		Orange []string `json:"orange"`
	} `json:"eyecolor"`
	Volume struct {
		Low        []string `json:"low"`
		MediumLow  []string `json:"mediumlow"`
		Medium     []string `json:"medium"`
		MediumHigh []string `json:"mediumhigh"`
		High       []string `json:"high"`
	} `json:"volume"`
	NameIs []string `json:"nameis"`
}

func Load(language string) error {
	iData, err := os.Open(filepath.Join(vars.IntentDataPath, language))
	if err != nil {
		log.Error("failed to open intent data file:", err)
		return err
	}
	if json.NewDecoder(iData).Decode(&idata) != nil {
		return err
	}
	return nil
}

func TextToResponse() {

}
