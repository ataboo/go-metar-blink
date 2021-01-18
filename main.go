package main

import (
	"os"
	"time"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/metarclient"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
	"github.com/ataboo/go-metar-blink/pkg/virtualmap"
)

func main() {
	appSettings := common.GetAppSettings()

	stationRepo := initStationRepo(appSettings)

	stations, err := stationRepo.GetStations()

	if err != nil {
		common.LogError("failed to get stations: %s", err)
		panic("aborting")
	}

	common.LogInfo("Got stations: %d", len(stations))

	err = virtualmap.ShowMap(stations)
	if err != nil {
		common.LogError("failed to show map: %s", err)
	}

	os.Exit(0)
}

func initStationRepo(settings *common.AppSettings) *stationrepo.StationRepo {
	common.LogInfo("[1.1] Initializing client")
	common.DumpSettingsInfo()
	client, err := metarclient.CreateMetarClient(&metarclient.Settings{
		StationIDs: settings.StationIDs,
		Strategy:   metarclient.MetarStrategy(settings.ClientStrategy),
	})
	if err != nil {
		common.LogError("Failed to start client: %s", err.Error())
	}

	repo := stationrepo.CreateStationRepo(client, &stationrepo.Config{UpdatePeriod: time.Minute * time.Duration(settings.UpdatePeriodMins)})

	return repo
}
