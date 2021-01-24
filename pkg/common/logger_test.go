package common

import (
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := InitLoggersToTestWriter()

	_, err := ParseLogLevel("extreme")
	if err == nil {
		t.Error("expected error in parsing log level")
	}

	CurrentLogLevel, err = ParseLogLevel("DeBuG")
	if err != nil {
		t.Error(err)
	}
	LogDebug("%s", "debug")
	LogInfo("%s", "info")
	LogWarn("%s", "warn")
	LogError("%s", "error")

	assertLogLines(t, logger, "debug", "info", "warn", "error")

	CurrentLogLevel, err = ParseLogLevel("InFo")
	if err != nil {
		t.Error(err)
	}
	logger.Lines = make([]string, 0)

	LogDebug("%s", "debug")
	LogInfo("%s", "info")
	LogWarn("%s", "warn")
	LogError("%s", "error")

	assertLogLines(t, logger, "info", "warn", "error")

	CurrentLogLevel, err = ParseLogLevel("WaRnInG")
	if err != nil {
		t.Error(err)
	}
	logger.Lines = make([]string, 0)

	LogDebug("%s", "debug")
	LogInfo("%s", "info")
	LogWarn("%s", "warn")
	LogError("%s", "error")

	assertLogLines(t, logger, "warn", "error")

	CurrentLogLevel, err = ParseLogLevel("ErRoR")
	if err != nil {
		t.Error(err)
	}
	logger.Lines = make([]string, 0)

	LogDebug("%s", "debug")
	LogInfo("%s", "info")
	LogWarn("%s", "warn")
	LogError("%s", "error")

	assertLogLines(t, logger, "error")
}

func assertLogLines(t *testing.T, writer *TestLogWriter, lines ...string) {
	if len(writer.Lines) != len(lines) {
		t.Errorf("unnexpected log line count %d, %d", len(writer.Lines), len(lines))
		return
	}

	for i, line := range lines {
		if strings.HasSuffix(writer.Lines[i], line) {
			t.Errorf("expected line %d to be %s, found %s", i, line, writer.Lines[i])
		}
	}
}
