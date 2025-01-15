package vars

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

var (
	BotJdocs   []botJdoc
	jdocsMutex sync.Mutex
)

type AJdoc struct {
	DocVersion     uint64 `json:"doc_version,omitempty"`
	FmtVersion     uint64 `json:"fmt_version,omitempty"`
	ClientMetadata string `json:"client_metadata,omitempty"`
	JsonDoc        string `json:"json_doc,omitempty"`
}

type botJdoc struct {
	Thing string `json:"thing"`
	Name  string `json:"name"`
	Jdoc  AJdoc  `json:"jdoc"`
}

func loadJdocsFile() error {
	if _, err := os.Stat(JdocsFilePath); os.IsNotExist(err) {
		return nil
	}
	f, err := os.Open(JdocsFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var data []botJdoc
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}
	BotJdocs = data
	return nil
}

func saveJdocsFile() error {
	f, err := os.Create(JdocsFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(BotJdocs)
}

func Thingifier(esn string) string {
	esn = strings.ToLower(strings.TrimSpace(esn))
	if !strings.HasPrefix(esn, "vic:") {
		esn = "vic:" + esn
	}
	return esn
}

func WriteJdoc(thing, name string, j AJdoc) error {
	jdocsMutex.Lock()
	defer jdocsMutex.Unlock()
	found := false
	for i := range BotJdocs {
		if BotJdocs[i].Thing == thing && BotJdocs[i].Name == name {
			BotJdocs[i].Jdoc = j
			found = true
			break
		}
	}
	if !found {
		BotJdocs = append(BotJdocs, botJdoc{Thing: thing, Name: name, Jdoc: j})
	}
	return saveJdocsFile()
}

func ReadJdoc(thing, name string) (AJdoc, error) {
	jdocsMutex.Lock()
	defer jdocsMutex.Unlock()
	for _, v := range BotJdocs {
		if v.Thing == thing && v.Name == name {
			return v.Jdoc, nil
		}
	}
	return AJdoc{}, errors.New("not found")
}
