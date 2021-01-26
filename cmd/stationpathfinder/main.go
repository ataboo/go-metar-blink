package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ataboo/go-metar-blink/cmd/stationpathfinder/pkg/visualization"
	"github.com/ataboo/go-metar-blink/cmd/stationpathfinder/pkg/wirepath"
	"github.com/ataboo/go-metar-blink/cmd/stationpathfinder/pkg/wirepath/antcolony"
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/logger"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

func loadPositions() ([]*wirepath.Position, error) {
	bytes, err := common.LoadCachedFile("station_screen_pos.json")
	if err != nil {
		logger.LogError("failed to load cached station screen positions")
		return nil, err
	}

	screenPositions := make(map[string]*sdl.Point)
	err = json5.Unmarshal(bytes, &screenPositions)

	positions := make([]*wirepath.Position, len(screenPositions))
	idx := 0
	for id, p := range screenPositions {
		positions[idx] = &wirepath.Position{
			Name: id,
			X:    int(p.X),
			Y:    int(p.Y),
		}
		idx++
	}

	return positions, err
}

type LoggedScore struct {
	Distance float64  `json:"distance"`
	Stations []string `json:"stations"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	startTime := time.Now()

	os.Mkdir("paths", logger.LoggingDirPermission)

	goMetarRoot, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	os.Setenv("GO_METAR_BLINK_ROOT", goMetarRoot)

	positions, err := loadPositions()
	if err != nil {
		logger.LogError("failed to load positions: %s", err)
		panic("aborting")
	}

	pathfinder, err := antcolony.CreateAntColonyPathFinder(&antcolony.Config{
		AntCount:                4,
		MaxPheromoneFactor:      0.6,
		PheromoneSpreadForward:  0.01,
		PheromoneSpreadBackward: 0.01,
		PheromoneDecay:          0.005,
		Positions:               positions,
	})
	if err != nil {
		logger.LogError("error creating pathfinder: %s", err)
		panic("aborting")
	}

	if err := pathfinder.RunRound(); err != nil {
		logger.LogError("error running round: %s", err)
		panic("aborting")
	}

	visualizer, err := visualization.CreatePathFindingVisualization(pathfinder)
	if err != nil {
		logger.LogError("failed to create pathfinding visualization: %s", err)
		panic("aborting")
	}

	bestCount := 0
	bestScore := math.MaxFloat64
	var bestPath []string = nil

	for {

		for i := 0; i < 4000; i++ {
			if err := pathfinder.RunRound(); err != nil {
				logger.LogError("error running round: %s", err)
				panic("aborting")
			}
			if i%100 == 0 {
				if err := visualizer.Update(pathfinder, bestPath, bestScore); err != nil {
					if _, ok := err.(*visualization.MapQuitError); ok {
						visualizer.Dispose()
						logger.LogInfo("map quit")
						os.Exit(0)
					}
					logger.LogError("error updating visualization: %s", err)
					panic("aborting")
				}
			}
		}

		stats := pathfinder.Stats()
		if stats.ShortestPath < bestScore {
			logger.LogInfo("Found new Best!: %f", stats.ShortestPath)
			bestPath = pathfinder.GetBestPath()
			bestScore = stats.ShortestPath

			data, err := json5.MarshalIndent(LoggedScore{Distance: bestScore, Stations: bestPath}, "", "  ")
			if err != nil {
				panic(err)
			}
			fileName := fmt.Sprintf("paths/best_path.%s.%d.json", startTime.Format("2006-01-15_15:04"), bestCount)
			ioutil.WriteFile(fileName, data, logger.LoggingFilePermission)
			bestCount++
		}

		pathfinder.Reset()
	}
}
