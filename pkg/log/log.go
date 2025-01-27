package log

import (
	"fmt"
	"hotwire/pkg/vars"
	"os"
	"time"

	"github.com/fatih/color"
)

var DebugLogging bool
var SuperDebugLogging bool

func init() {
	LogTrayChan = make(chan string)
	if os.Getenv(vars.DebugLoggingEnv) == "true" {
		fmt.Println("Debug logging is enabled")
		DebugLogging = true
	}
	if os.Getenv(vars.SuperDebugLoggingEnv) == "true" {
		fmt.Println("Super debug logging is enabled")
		SuperDebugLogging = true
	}
}

func SuperDebug(a ...any) {
	if SuperDebugLogging {
		color.Cyan(fmt.Sprint(a...))
	}
}

func Debug(a ...any) {
	if DebugLogging {
		Normal(a...)
	}
}

func Important(a ...any) {
	color.Green(fmt.Sprint(a...))
}

func Error(a ...any) {
	color.Red(fmt.Sprint(a...))
}

var debugLogging bool = true
var LogList string
var LogArray []string

var LogTrayList string
var LogTrayArray []string
var LogTrayChan chan string

func GetLogTrayChan() chan string {
	return LogTrayChan
}

// things which should be seen on the command line
func Normal(a ...any) {
	LogUIFull(a...)
	fmt.Println(a...)
}

func UI(a ...any) {
	LogArray = append(LogArray, time.Now().Format("2006.01.02 15:04:05")+": "+fmt.Sprint(a...)+"\n")
	if len(LogArray) >= 50 {
		LogArray = LogArray[1:]
	}
	LogList = ""
	for _, b := range LogArray {
		LogList = LogList + b
	}
}

func LogUIFull(a ...any) {
	LogTrayArray = append(LogTrayArray, time.Now().Format("2006.01.02 15:04:05")+": "+fmt.Sprint(a...)+"\n")
	if len(LogTrayArray) >= 200 {
		LogTrayArray = LogTrayArray[1:]
	}
	LogTrayList = ""
	for _, b := range LogTrayArray {
		LogTrayList = LogTrayList + b
	}
	select {
	case LogTrayChan <- time.Now().Format("2006.01.02 15:04:05") + ": " + fmt.Sprint(a...) + "\n":
	default:
	}
}
