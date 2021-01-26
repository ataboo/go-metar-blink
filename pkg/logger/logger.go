package logger

import (
	"errors"
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

func (l LogLevel) String() string {
	switch l {
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

func ParseLogLevel(levelStr string) (LogLevel, error) {
	switch strings.ToLower(levelStr) {
	case "error":
		return LogLevelError, nil
	case "warning":
		return LogLevelWarn, nil
	case "info":
		return LogLevelInfo, nil
	case "debug":
		return LogLevelDebug, nil
	default:
		fmt.Println("unsupported log level: " + levelStr)
		return LogLevelError, errors.New("unsupported log level")
	}
}

func LogDebug(format string, v ...interface{}) {
	if CurrentLogLevel > LogLevelDebug {
		return
	}

	assertLoggersInitialized()
	debugLogger.Output(2, fmt.Sprintf(format, v...))
}

func LogInfo(format string, v ...interface{}) {
	if CurrentLogLevel > LogLevelInfo {
		return
	}

	assertLoggersInitialized()
	infoLogger.Output(2, fmt.Sprintf(format, v...))
}

func LogWarn(format string, v ...interface{}) {
	if CurrentLogLevel > LogLevelWarn {
		return
	}

	assertLoggersInitialized()
	warningLogger.Output(2, fmt.Sprintf(format, v...))
}

func LogError(format string, v ...interface{}) {
	assertLoggersInitialized()
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}

func assertLoggersInitialized() {
	if debugLogger != nil && infoLogger != nil && warningLogger != nil && errorLogger != nil {
		return
	}

	InitLoggersToWriter(os.Stdout, os.Stderr)
}

func InitFileLogging(loggingDir string, baseName string, multiFile bool) {
	err := os.MkdirAll(loggingDir, LoggingDirPermission)
	if err != nil {
		panic(err)
	}

	var fileName string
	if multiFile {
		fileName = fmt.Sprintf("%s.%s.log", baseName, time.Now().Format("2006-01-02_150405"))
	} else {
		fileName = fmt.Sprintf("%s.log", baseName)
	}
	file, err := os.OpenFile(path.Join(loggingDir, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, LoggingFilePermission)
	if err != nil {
		panic(err)
	}

	InitLoggersToWriter(file, file)
}

func InitLoggersToWriter(stdOut io.Writer, errOut io.Writer) {
	debugLogger = log.New(stdOut, LogLevel(LogLevelDebug).String()+": ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(stdOut, LogLevel(LogLevelInfo).String()+":  ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(stdOut, LogLevel(LogLevelWarn).String()+":  ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errOut, LogLevel(LogLevelError).String()+": ", log.Ldate|log.Ltime|log.Lshortfile)

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
