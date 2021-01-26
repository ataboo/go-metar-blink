package common

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ataboo/go-metar-blink/pkg/logger"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

const (
	PiBootAppSettingsPath = "/boot/go-metar-blink.settings.json"
	PiBootPanicErrorPath  = "/boot/go-metar-blink.panic.log"
)

var _appSettings *AppSettings

type AppSettings struct {
	StationIDs       []string           `json:"station_ids"`
	ClientStrategy   string             `json:"client_strategy"`
	UpdatePeriodMins int                `json:"update_period_mins"`
	LoggingDir       string             `json:"logging_dir"`
	LoggingMethod    string             `json:"logging_method"`
	LoggingLevel     string             `json:"logging_level"`
	CacheDir         string             `json:"cache_dir"`
	Colors           *ColorThemeStrings `json:"colors"`
	colorsParsed     *ColorTheme
}

func (a *AppSettings) GetParsedColors() *ColorTheme {
	if a.colorsParsed == nil {
		errors := make(map[string]string)
		a.colorsParsed = a.Colors.ParseColors(errors)
	}

	return a.colorsParsed
}

func SetAppSettings(settings *AppSettings) {
	_appSettings = settings
}

func GetAppSettings() *AppSettings {
	if _appSettings == nil {
		mustLoadAppSettings()
	}

	return _appSettings
}

func DumpSettingsInfo() {
	settings := GetAppSettings()

	logger.LogDebug("\tActive Station IDs: %s", strings.Join(settings.StationIDs, ", "))
	logger.LogDebug("\tClient Strategy: %s", settings.ClientStrategy)
	logger.LogDebug("\tUpdatePeriodMins: %d", settings.UpdatePeriodMins)
	logger.LogDebug("\tLoggingDir: %s", settings.LoggingDir)
	logger.LogDebug("\tLoggingMethod: %s", settings.LoggingMethod)
	logger.LogDebug("\tLoggingLevel: %s", settings.LoggingLevel)
	logger.LogDebug("\tCacheDir: %s", settings.CacheDir)
	logger.LogDebug("\tColors")
	logger.LogDebug("\t\tVFR: %s", settings.Colors.VFR)
	logger.LogDebug("\t\tSVFR: %s", settings.Colors.SVFR)
	logger.LogDebug("\t\tIFR: %s", settings.Colors.IFR)
	logger.LogDebug("\t\tLIFR: %s", settings.Colors.LIFR)
	logger.LogDebug("\t\tError: %s", settings.Colors.Error)
	logger.LogDebug("\t\tBrightness: %s", settings.Colors.Brightness)
}

func inTestEnvironment() bool {
	return flag.Lookup("test.v") != nil
}

func mustLoadAppSettings() {
	errors := make(map[string]string)

	settingsRaw, err := loadRawSettingsFile()
	if err != nil {
		errors["settings"] = "failed to load settings: " + err.Error()
		logErrorsAndPanic(errors)
	}

	_appSettings = &AppSettings{}
	err = json5.Unmarshal(settingsRaw, &_appSettings)
	if err != nil {
		errors["settings"] = "failed to parse settings.json: " + err.Error()
		logErrorsAndPanic(errors)
	}

	validateSettings(_appSettings, errors)

	if len(errors) > 0 {
		logErrorsAndPanic(errors)
	}

	initLoggingFromSettings(_appSettings)
}

func loadRawSettingsFile() ([]byte, error) {
	if runtime.GOOS == "arm" {
		if _, err := os.Stat(PiBootAppSettingsPath); err != nil {
			fmt.Println("loading appsettings from boot partition")
			return ioutil.ReadFile(PiBootAppSettingsPath)
		}
	}

	fmt.Println("loading appsettings from project root")
	return ioutil.ReadFile(path.Join(GetProjectRoot(), "settings.json"))
}

func logErrorsAndPanic(errors map[string]string) {
	defer panic("aborting")

	if inTestEnvironment() {
		for field, err := range errors {
			fmt.Printf("settings error: %s|%s", field, err)
		}

		return
	}

	var filePath string

	if runtime.GOOS == "arm" {
		filePath = PiBootPanicErrorPath
	} else {
		filePath = path.Join(GetProjectRoot(), "panic.log")
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, logger.LoggingFilePermission)
	if err != nil {
		fmt.Println("failed to write panic log")
	} else {
		for field, err := range errors {
			line := fmt.Sprintf("%s|%s\n", field, err)
			f.WriteString(line)
			fmt.Printf("settings error: %s", line)
		}
		f.Close()
	}
}

func validateSettings(settings *AppSettings, errors map[string]string) {
	switch settings.ClientStrategy {
	case AviationWeatherMetarStrategy:
		break
	default:
		errors["ClientStrategy"] = "invalid client strategy"
	}

	settings.colorsParsed = settings.Colors.ParseColors(errors)

	switch settings.LoggingMethod {
	case logger.LoggingMethodStdio:
	case logger.LoggingMethodMultiFile:
	case logger.LoggingMethodSingleFile:
		break
	default:
		errors["LoggingMethod"] = "invalid logging method"
	}

	if _, err := logger.ParseLogLevel(settings.LoggingLevel); err != nil {
		errors["LoggingLevel"] = "invalid logging level"
	}

	if settings.UpdatePeriodMins < 1 {
		errors["UpdatePeriodMins"] = "update period must be atleast 1 minute"
	}

	validateStationIds(errors)
}

func validateStationIds(errors map[string]string) {
	keyMap := make(map[string]bool)
	for _, id := range _appSettings.StationIDs {
		if _, ok := keyMap[id]; ok {
			logger.LogWarn("station '%s' is found more than once in settings", id)
		}

		keyMap[id] = true
	}
}
