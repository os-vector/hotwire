package ttr

type Intent struct {
	Intent                 string   `json:"intent"`
	PartialMatchUterrances []string `json:"utterances"`
	ExactMatchUtterances   []string `json:"exactutterances"`
	ExactMatchPreferred    bool     `json:"exactmatchpreferred"`
}

type 

func TextToResponse() {

}
