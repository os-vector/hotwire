package vars

import (
	"os"
	"path/filepath"
)

var (
	RootPath = "./storage"

	JdocsFilePath         = "jdocs.json"
	BotGUIDsFilePath      = "botGUIDs.json"
	HotwireConfigFilePath = "hotwireConfig.json"
	PerBotFilePath        = "perBotConfig.json"
	SessionCertsStorage   = "session-certs"

	CertPath = "cert.pem"
	KeyPath  = "cert.key"
)

func init() {
	if _, err := os.Stat(RootPath); err != nil {
		os.Mkdir(RootPath, 0777)
	}
	JdocsFilePath = filepath.Join(RootPath, JdocsFilePath)
	BotGUIDsFilePath = filepath.Join(RootPath, BotGUIDsFilePath)
	HotwireConfigFilePath = filepath.Join(RootPath, HotwireConfigFilePath)
	PerBotFilePath = filepath.Join(RootPath, PerBotFilePath)
	SessionCertsStorage = filepath.Join(RootPath, SessionCertsStorage)

	CertPath = filepath.Join(RootPath, CertPath)
	KeyPath = filepath.Join(RootPath, KeyPath)
}
