package main

import (
	"os"
	"path/filepath"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

func main() {
	goMetarRoot, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	os.Setenv("GO_METAR_BLINK_ROOT", goMetarRoot)

	pathfinder, err := CreateStationPathFinder()
	if err != nil {
		common.LogError("error creating pathfinder: %s", err)
		panic("aborting")
	}

	if err := pathfinder.RunRound(); err != nil {
		common.LogError("error running round: %s", err)
		panic("aborting")
	}

	visualizer, err := CreatePathFindingVisualization(pathfinder)
	if err != nil {
		common.LogError("failed to create pathfinding visualization: %s", err)
		panic("aborting")
	}

	count := 0
	for {
		if err := pathfinder.RunRound(); err != nil {
			common.LogError("error running round: %s", err)
			panic("aborting")
		}
		if count > 100 {
			count = 0

			if err := visualizer.Update(pathfinder); err != nil {
				if _, ok := err.(*common.MapQuitError); ok {
					visualizer.Dispose()
					common.LogInfo("map quit")
					os.Exit(0)
				}
				common.LogError("error updating visualization: %s", err)
				panic("aborting")
			}
		}
		count++
	}
}
