package common

import (
	"io"
	"log"
	"os"
	"path"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func LogInfo(format string, v ...interface{}) {
	if infoLogger == nil {
		initLoggers()
	}

	infoLogger.Printf(format, v)
}

func LogWarn(format string, v ...interface{}) {
	if warningLogger == nil {
		initLoggers()
	}

	warningLogger.Printf(format, v)
}

func LogError(format string, v ...interface{}) {
	if errorLogger == nil {
		initLoggers()
	}

	errorLogger.Printf(format, v)
}

func initLoggers() {
	if os.Getenv("LOG_METHOD") == "file" {
		file, err := os.OpenFile(path.Join(GetProjectRoot(), "logs.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		initLoggersToWriter(file, file)
	} else {
		initLoggersToWriter(os.Stdout, os.Stderr)
	}
}

func initLoggersToWriter(stdOut io.Writer, errOut io.Writer) {
	infoLogger = log.New(stdOut, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(stdOut, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errOut, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
