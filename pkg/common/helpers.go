package common

import (
	"fmt"
	"math"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/ataboo/go-metar-blink/pkg/animation"
	"github.com/ataboo/go-metar-blink/pkg/logger"
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

func initLoggingFromSettings(appSettings *AppSettings) error {
	level, err := logger.ParseLogLevel(_appSettings.LoggingLevel)
	if err != nil {
		return err
	}
	logger.CurrentLogLevel = level

	switch strings.ToLower(appSettings.LoggingMethod) {
	case logger.LoggingMethodSingleFile:
		logger.InitFileLogging(appSettings.LoggingDir, "go-metar-blink", false)
	case logger.LoggingMethodMultiFile:
		logger.InitFileLogging(appSettings.LoggingDir, "go-metar-blink", true)
	case logger.LoggingMethodStdio:
		// logger defaults to stdio.
	default:
		return fmt.Errorf("unsupported logging method: " + appSettings.LoggingMethod)
	}

	return nil
}

func GetLocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iFace := range ifaces {
		if strings.HasPrefix(iFace.Name, "wl") {
			addresses, err := iFace.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addresses {
				ipNet, ok := addr.(*net.IPNet)

				if !ok || ipNet.IP.To4() == nil || ipNet.IP.IsLoopback() {
					continue
				}

				return ipNet.IP, nil
			}
		}
	}

	return nil, nil
}
