package main

import (
	"os"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/metarclient"
)

func main() {
	appSettings := common.GetAppSettings()

	common.LogInfo("[1.1] Initializing client")
	common.DumpSettingsInfo()
	_, err := metarclient.CreateMetarClient(&metarclient.Settings{
		StationIDs: appSettings.StationIDs,
		Strategy:   metarclient.MetarStrategy(appSettings.ClientStrategy),
	})
	if err != nil {
		common.LogError("Failed to start client: %s", err.Error())
	}

	common.LogInfo("Hello!")

	os.Exit(0)
}
