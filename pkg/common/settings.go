package common

import (
	"flag"
	"io/ioutil"
	"path"
	"strings"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type AppSettings struct {
	StationIDs        []string `json:"station_ids"`
	ClientStrategy    string   `json:"client_strategy"`
	WindyThresholdKts float32  `json:"windy_threshold_kts"`
	UpdatePeriodMins  int      `json:"update_period_mins"`
	LoggingDir        string   `json:"logging_dir"`
	LoggingMethod     string   `json:"logging_method"`
	LoggingLevel      string   `json:"logging_level"`
	CacheDir          string   `json:"cache_dir"`
}

var _appSettings *AppSettings

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

	LogDebug("\tActive Station IDs: %s", strings.Join(settings.StationIDs, ", "))
	LogDebug("\tClient Strategy: %s", settings.ClientStrategy)
	LogDebug("\tWindyThresholdKts: %f", settings.WindyThresholdKts)
	LogDebug("\tUpdatePeriodMins: %d", settings.UpdatePeriodMins)
}

func inTestEnvironment() bool {
	return flag.Lookup("test.v") != nil
}

func mustLoadAppSettings() {
	settingsRaw, err := ioutil.ReadFile(path.Join(GetProjectRoot(), "settings.json"))
	if err != nil {
		panic("failed to read 'settings.json': " + err.Error())
	}

	_appSettings = &AppSettings{}
	err = json5.Unmarshal(settingsRaw, &_appSettings)
	if err != nil {
		panic("failed to parse 'settings.json': " + err.Error())
	}

	CurrentLogLevel = MustParseLogLevel(_appSettings.LoggingLevel)

	validateStationIds()
}

func validateStationIds() {
	keyMap := make(map[string]bool)
	for _, id := range _appSettings.StationIDs {
		if _, ok := keyMap[id]; ok {
			LogWarn("station '%s' is found more than once in settings", id)
		}

		keyMap[id] = true
	}

	LogInfo("loaded %d stations", len(_appSettings.StationIDs))
}
