package wirepath

import (
	"time"
)

type PathFinder interface {
	GetBestPath() *StationPath
	Stats() *PathfinderStats
	RunRound() error
}

type PathfinderStats struct {
	PathsGenerated int
	RunTime        time.Duration
	ShortestPath   float64
}
