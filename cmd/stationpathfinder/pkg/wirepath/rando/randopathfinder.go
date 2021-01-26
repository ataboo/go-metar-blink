package rando

import (
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/ataboo/go-metar-blink/cmd/stationpathfinder/pkg/wirepath"
	"github.com/ataboo/go-metar-blink/pkg/logger"
)

const (
	PathsPerGeneration = 1000
	KeeperCount        = 1
)

// Uses a randomized approach to select the shortest path generated.  Not very efficient.
type RandoPathfinder struct {
	posMap         map[string]*wirepath.Position
	positions      []*wirepath.Position
	paths          []*StationPath
	semaphore      sync.WaitGroup
	startTime      time.Time
	pathsGenerated int
}

func (f *RandoPathfinder) GetBestPath() []string {
	path := f.paths[0]
	connection := path.Connections[path.StartID]
	ids := make([]string, len(f.positions))
	for i := 0; i < len(f.positions); i++ {
		ids[i] = connection.ID
		connection = path.Connections[connection.EndID]
	}

	return ids
}

func (f *RandoPathfinder) Stats() *wirepath.PathfinderStats {
	return &wirepath.PathfinderStats{
		PathsGenerated: f.pathsGenerated,
		RunTime:        time.Now().Sub(f.startTime),
		ShortestPath:   f.paths[0].Length,
	}
}

func CreateStationPathFinder(positions []*wirepath.Position) (wirepath.PathFinder, error) {
	rand.Seed(time.Now().UnixNano())

	posMap := make(map[string]*wirepath.Position)
	for _, p := range positions {
		posMap[p.Name] = p
	}

	pathFinder := &RandoPathfinder{
		startTime:      time.Now(),
		pathsGenerated: 0,
		semaphore:      sync.WaitGroup{},
		paths:          make([]*StationPath, PathsPerGeneration),
		posMap:         posMap,
		positions:      positions,
	}

	stationIDs := make([]string, len(positions))
	for i, p := range positions {
		stationIDs[i] = p.Name
	}

	for i := 0; i < PathsPerGeneration; i++ {
		pathFinder.paths[i] = CreateStationPath(stationIDs, "CYYJ")
	}

	pathFinder.createRandomPaths(0, PathsPerGeneration)

	return pathFinder, nil
}

func (s *RandoPathfinder) GetPositions() []*wirepath.Position {
	return s.positions
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
			if err := s.paths[idx].CalculateLength(s.posMap); err != nil {
				logger.LogError("pathfinder calc error: %s", err)
			}
			s.semaphore.Done()
		}()
	}
}

func (s *RandoPathfinder) Reset() {
	s.createRandomPaths(0, PathsPerGeneration)
	s.pathsGenerated = 0
	s.startTime = time.Now()
}
