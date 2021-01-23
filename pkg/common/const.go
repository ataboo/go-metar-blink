package common

import (
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	FlightRuleVFR     = "VFR"
	FlightRuleSVFR    = "SVFR"
	FlightRuleMVFR    = "MVFR"
	FlightRuleIFR     = "IFR"
	FlightRuleLIFR    = "LIFR"
	FlightRuleUnknown = "Unknown"
	FlightRuleError   = "Error"
)

type MapQuitError struct{}

func (e *MapQuitError) Error() string { return "this map is no longer running" }

func GetProjectRoot() string {
	exec, _ := os.Executable()
	fmt.Printf("Exec: %s\n", exec)

	if !strings.HasSuffix(exec, ".test") {
		return path.Dir(exec)
	}

	return stepUpFromTestToProjectRoot()
}

func stepUpFromTestToProjectRoot() string {
	tryCount := 0
	dir, _ := os.Getwd()

	for tryCount < 10 {
		tryCount++
		if path.Base(dir) == "go-metar-blink" {
			return dir
		}

		dir = path.Join(dir, "../")
	}

	panic("failed to find project root")
}

func GetResourcesRoot() string {
	return path.Join(GetProjectRoot(), "resources")
}
