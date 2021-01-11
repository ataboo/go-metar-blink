package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var (
	debugLogger     *log.Logger
	infoLogger      *log.Logger
	warningLogger   *log.Logger
	errorLogger     *log.Logger
	CurrentLogLevel LogLevel = LogLevelInfo
)

const (
	LogLevelDebug           = LogLevel(iota)
	LogLevelInfo            = LogLevel(iota)
	LogLevelWarn            = LogLevel(iota)
	LogLevelError           = LogLevel(iota)
	LoggingDirPermission    = 0777
	LoggingFilePermission   = 0666
	LoggingMethodSingleFile = "single-file"
	LoggingMethodMultiFile  = "multi-file"
	LoggingMethodStdio      = "console"
)

var _ io.Writer = (*TestLogWriter)(nil)

type LogLevel int

func (l *LogLevel) String() string {
	switch *l {
	case LogLevelDebug:
		return "Debug"
	case LogLevelInfo:
		return "Info"
	case LogLevelWarn:
		return "Warn"
	case LogLevelError:
		return "Error"
	default:
		panic("log level not supported")
	}
}

func MustParseLogLevel(levelStr string) LogLevel {
	switch strings.ToLower(_appSettings.LoggingLevel) {
	case "error":
		return LogLevelError
	case "warning":
		return LogLevelWarn
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	default:
		panic("unsupported log level: " + _appSettings.LoggingLevel)
	}
}

func LogDebug(format string, v ...interface{}) {
	if CurrentLogLevel > LogLevelDebug {
		return
	}

	assertLoggersInitialized()
	debugLogger.Printf(format, v...)
}

func LogInfo(format string, v ...interface{}) {
	if CurrentLogLevel > LogLevelInfo {
		return
	}

	assertLoggersInitialized()
	infoLogger.Printf(format, v...)
}

func LogWarn(format string, v ...interface{}) {
	if CurrentLogLevel > LogLevelWarn {
		return
	}

	assertLoggersInitialized()
	warningLogger.Printf(format, v...)
}

func LogError(format string, v ...interface{}) {
	assertLoggersInitialized()
	errorLogger.Printf(format, v...)
}

func assertLoggersInitialized() {
	if infoLogger != nil && warningLogger != nil && errorLogger != nil {
		return
	}

	if inTestEnvironment() {
		panic("loggers should be initialized when testing")
	}

	appSettings := GetAppSettings()

	switch strings.ToLower(appSettings.LoggingMethod) {
	case LoggingMethodSingleFile:
		initFileLogging(appSettings, false)
	case LoggingMethodMultiFile:
		initFileLogging(appSettings, true)
	case LoggingMethodStdio:
		InitLoggersToWriter(os.Stdout, os.Stderr)
	default:
		panic("unsupported logging method: " + appSettings.LoggingMethod)
	}
}

func initFileLogging(appSettings *AppSettings, multiFile bool) {
	err := os.MkdirAll(appSettings.LoggingDir, LoggingDirPermission)
	if err != nil {
		panic(err)
	}

	var fileName string
	if multiFile {
		fileName = fmt.Sprintf("go-metar-blink.%s.log", time.Now().Format("2006-01-02_150405"))
	} else {
		fileName = fmt.Sprintf("go-metar-blink.log")
	}
	file, err := os.OpenFile(path.Join(appSettings.LoggingDir, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, LoggingFilePermission)
	if err != nil {
		panic(err)
	}

	InitLoggersToWriter(file, file)
}

func InitLoggersToWriter(stdOut io.Writer, errOut io.Writer) {
	debugLogger = log.New(stdOut, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(stdOut, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(stdOut, "WARN:  ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errOut, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	LogDebug("initialized loggers to level '%s'", CurrentLogLevel.String())
}

func InitLoggersToTestWriter() *TestLogWriter {
	writer := TestLogWriter{
		Lines: make([]string, 0),
	}

	InitLoggersToWriter(&writer, &writer)

	return &writer
}

type TestLogWriter struct {
	Lines []string
}

func (w *TestLogWriter) Write(p []byte) (n int, err error) {
	if w.Lines == nil {
		w.Lines = make([]string, 0)
	}

	w.Lines = append(w.Lines, string(p))

	return len(p), nil
}
