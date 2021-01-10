package common

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type AppSettings struct {
	StationIDs        []string `json:"station_ids"`
	ClientStrategy    string   `json:"client_strategy"`
	WindyThresholdKts float32  `json:"windy_threshold_kts"`
	UpdatePeriodMins  int      `json:"update_period_mins"`
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

func mustLoadAppSettings() {
	settingsRaw, err := ioutil.ReadFile(path.Join(GetProjectRoot(), "settings.json"))
	if err != nil {
		LogError("failed to read 'settings.json': %s", err.Error())
		panic("aborting")
	}

	_appSettings = &AppSettings{}
	err = json.Unmarshal(settingsRaw, &_appSettings)
	if err != nil {
		LogError("failed to parse 'settings.json': %s", err.Error())
		panic("aborting")
	}
}
