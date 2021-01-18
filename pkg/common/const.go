package common

import (
	"os"
	"path"
)

const (
	FlightRuleVFR   = "VFR"
	FlightRuleSVFR  = "SVFR"
	FlightRuleIFR   = "IFR"
	FlightRuleLIFR  = "LIFR"
	FlightRuleError = "Error"
)

func GetProjectRoot() string {
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
