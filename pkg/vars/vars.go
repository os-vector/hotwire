package vars

import (
	"os"
	"path/filepath"
)

var (
	UserData = "./storage"

	JdocsFilePath         = "jdocs.json"
	SavedRobotsFilePath   = "savedRobots.json"
	HotwireConfigFilePath = "hotwireConfig.json"
	PerBotFilePath        = "perBotConfig.json"
	SessionCertsStorage   = "session-certs"
	STTStorage            = "stt"

	IntentDataPath = "intent-data"

	CertPath = "cert.pem"
	KeyPath  = "cert.key"
)

func init() {
	if _, err := os.Stat(UserData); err != nil {
		os.Mkdir(UserData, 0777)
	}
	JdocsFilePath = filepath.Join(UserData, JdocsFilePath)
	SavedRobotsFilePath = filepath.Join(UserData, SavedRobotsFilePath)
	HotwireConfigFilePath = filepath.Join(UserData, HotwireConfigFilePath)
	PerBotFilePath = filepath.Join(UserData, PerBotFilePath)
	SessionCertsStorage = filepath.Join(UserData, SessionCertsStorage)
	STTStorage = filepath.Join(UserData, STTStorage)

	CertPath = filepath.Join(UserData, CertPath)
	KeyPath = filepath.Join(UserData, KeyPath)
}
