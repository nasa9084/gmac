package log

import (
	"log"
	"os"
	"sync"
)

var (
	mu            sync.Mutex
	isVerbose     bool
	verboseLogger *log.Logger
)

func init() {
	verboseLogger = log.New(os.Stderr, "", log.LstdFlags)
}

func SetVerbose(v bool) {
	mu.Lock()
	defer mu.Unlock()

	isVerbose = v
}

var (
	Fatal = log.Fatal

	Print   = log.Print
	Printf  = log.Printf
	Println = log.Println
)

func Vprint(v ...interface{}) {
	if isVerbose {
		verboseLogger.Print(v...)
	}
}

func Vprintf(format string, v ...interface{}) {
	if isVerbose {
		verboseLogger.Printf(format, v...)
	}
}
