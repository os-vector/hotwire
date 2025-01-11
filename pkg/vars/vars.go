package vars

import (
	"fmt"
	"os"
)

var (
	JdocsFilePath         = "jdocs.json"
	BotGUIDsFilePath      = "botGUIDs.json"
	HotwireConfigFilePath = "config.json"
	PerBotFilePath        = "perBot.json"
)

var DebugLogging bool

func Init() {
	if os.Getenv(DebugLoggingEnv) == "true" {
		fmt.Println("Debug logging is enabled")
		DebugLogging = true
	}
}
