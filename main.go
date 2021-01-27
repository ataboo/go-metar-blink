package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/engine"
	"github.com/ataboo/go-metar-blink/pkg/logger"
	"github.com/ataboo/go-metar-blink/pkg/metarclient"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
)

func main() {
	common.GetAppSettings()
	runMainApp()

	os.Exit(0)
}

func runMainApp() {
	stationRepo := initStationRepo(common.GetAppSettings())

	engine, err := engine.CreateEngine(stationRepo, common.GetAppSettings())
	if err != nil {
		logger.LogError("failed to create engine, %s", err)
		panic("aborting")
	}

	err = engine.Start()
	if err != nil {
		logger.LogError("failed to start engine: %s", err)
		panic("aborting")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-engine.DoneSubscribe():
		logger.LogDebug("engine done")
		break
	case sigVal := <-sigs:
		logger.LogDebug("received signal: %d", sigVal)
	}

	engine.Dispose()
}

func initStationRepo(settings *common.AppSettings) *stationrepo.StationRepo {
	logger.LogInfo("[1.1] Initializing client")
	common.DumpSettingsInfo()
	client, err := metarclient.CreateMetarClient(&metarclient.Settings{
		StationIDs: settings.StationIDs,
		Strategy:   metarclient.MetarStrategy(settings.ClientStrategy),
	})
	if err != nil {
		logger.LogError("Failed to start client: %s", err.Error())
	}

	repo := stationrepo.CreateStationRepo(client, &stationrepo.Config{StationIDs: settings.StationIDs})

	return repo
}
