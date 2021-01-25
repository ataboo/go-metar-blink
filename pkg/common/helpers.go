package common

import (
	"math"
	"os"
	"path"
	"strconv"

	"github.com/ataboo/go-metar-blink/pkg/animation"
)

func Similar(a float64, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

func NormalizePlusMinusPi(angleRads float64) float64 {
	twoPi := math.Pi * 2
	positiveReduced := math.Mod(math.Mod(angleRads, twoPi)+twoPi, twoPi)

	if positiveReduced > math.Pi {
		positiveReduced -= twoPi
	}

	return positiveReduced
}

func ParseByteHexString(strval string) (byte, error) {
	byteVal, err := strconv.ParseUint(strval, 0, 8)

	return byte(byteVal), err
}

func ParseColorHexString(strVal string) (animation.Color, error) {
	intVal, err := strconv.ParseUint(strVal, 0, 32)

	return animation.Color(intVal), err
}

func GetProjectRoot() string {
	if root, ok := os.LookupEnv("GO_METAR_BLINK_ROOT"); ok {
		return root
	}

	if inTestEnvironment() {
		return stepUpFromTestToProjectRoot()
	}

	exec, _ := os.Executable()
	return path.Dir(exec)
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
