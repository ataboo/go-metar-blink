package rando

import (
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/wirepath"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

const (
	PathsPerGeneration = 1000
	KeeperCount        = 1
)

type RandoPathfinder struct {
	screenPositions map[string]*sdl.Point
	paths           []*wirepath.StationPath
	semaphore       sync.WaitGroup
	startTime       time.Time
	pathsGenerated  int
}

func (f *RandoPathfinder) GetBestPath() *wirepath.StationPath {
	return f.paths[0]
}

func (f *RandoPathfinder) Stats() *wirepath.PathfinderStats {
	return &wirepath.PathfinderStats{
		PathsGenerated: f.pathsGenerated,
		RunTime:        time.Now().Sub(f.startTime),
		ShortestPath:   f.paths[0].Length,
	}
}

func CreateStationPathFinder() (wirepath.PathFinder, error) {
	rand.Seed(time.Now().UnixNano())

	pathFinder := &RandoPathfinder{
		startTime:      time.Now(),
		pathsGenerated: 0,
		semaphore:      sync.WaitGroup{},
		paths:          make([]*wirepath.StationPath, PathsPerGeneration),
	}
	bytes, err := common.LoadCachedFile("station_screen_pos.json")
	if err != nil {
		common.LogError("failed to load cached station screen positions")
		return nil, err
	}

	screenPositions := make(map[string]*sdl.Point)
	err = json5.Unmarshal(bytes, &screenPositions)

	pathFinder.screenPositions = screenPositions

	stationIDs := make([]string, len(screenPositions))
	idx := 0
	for id := range screenPositions {
		stationIDs[idx] = id
		idx++
	}

	//TODO load cached paths

	for i := 0; i < PathsPerGeneration; i++ {
		pathFinder.paths[i] = wirepath.CreateStationPath(stationIDs, "CYYJ")
	}

	pathFinder.createRandomPaths(0, PathsPerGeneration)

	return pathFinder, nil
}

func (s *RandoPathfinder) RunRound() error {
	if err := s.cullAndSeedPaths(); err != nil {
		return err
	}

	s.semaphore.Wait()

	sort.Slice(s.paths, func(i, j int) bool {
		return s.paths[i].Length < s.paths[j].Length
	})

	return nil
}

func (s *RandoPathfinder) cullAndSeedPaths() error {
	keeperCount := KeeperCount
	newSeedCount := PathsPerGeneration - KeeperCount

	s.createRandomPaths(keeperCount, newSeedCount)

	return nil
}

func (s *RandoPathfinder) createRandomPaths(startIdx int, count int) {
	s.semaphore.Add(count)
	s.pathsGenerated += count

	for i := startIdx; i < PathsPerGeneration; i++ {
		idx := i
		go func() {
			s.paths[idx].Randomize()
			if err := s.paths[idx].CalculateLength(s.screenPositions); err != nil {
				common.LogError("pathfinder calc error: %s", err)
			}
			s.semaphore.Done()
		}()
	}
}
