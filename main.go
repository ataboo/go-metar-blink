package main

import (
	"os"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/engine"
	"github.com/ataboo/go-metar-blink/pkg/metarclient"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
)

func main() {
	appSettings := common.GetAppSettings()
	stationRepo := initStationRepo(appSettings)

	engine, err := engine.CreateEngine(stationRepo, appSettings)
	if err != nil {
		common.LogError("failed to create engine, %s", err)
		panic("aborting")
	}

	err = engine.Start()
	if err != nil {
		common.LogError("failed to start engine: %s", err)
		panic("aborting")
	}

	select {
	case <-engine.DoneSubscribe():
		break
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

	repo := stationrepo.CreateStationRepo(client, &stationrepo.Config{StationIDs: settings.StationIDs})

	return repo
}
