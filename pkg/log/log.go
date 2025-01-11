package log

import (
	"fmt"
	"hotwire/pkg/vars"
	"time"
)

func Init() {
	LogTrayChan = make(chan string)
}

func Debug(a ...any) {
	if vars.DebugLogging {
		Normal(a...)
	}
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
