package wirepath

import (
	"time"
)

type Position struct {
	Name string
	X    int
	Y    int
}

type PathFinder interface {
	// GetBestPath gets the best path as station names ordered.
	GetBestPath() []string
	// Stats gets statistics.
	Stats() *PathfinderStats
	// RunRound simulates a number of ants and increments the pheromones.
	RunRound() error
	// GetPositions gets all the positions.
	GetPositions() []*Position
	// Reset the pathfinder.
	Reset()
}

type PathfinderStats struct {
	PathsGenerated int
	RunTime        time.Duration
	ShortestPath   float64
}
