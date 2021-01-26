package antcolony

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/ataboo/go-metar-blink/cmd/stationpathfinder/pkg/wirepath"
)

// Config sets the parameters for the pathfinder.
type Config struct {
	AntCount int
	// Maximum weight in decision making that a pheromone may hold vs shortest distance.
	MaxPheromoneFactor float64
	// Pheromone left in the direction of travel at all the paths the best ant took.
	PheromoneSpreadForward float64
	// Pheromone left in the reverse direction at all the paths the best ant took.
	PheromoneSpreadBackward float64
	// Pheromone amount lost when the fastest and did not take a path.
	PheromoneDecay float64
	Positions      []*wirepath.Position
}

// AntColony models a number of ants leaving pheromone trails to find the shortest path to all stations.
// https://en.wikipedia.org/wiki/Ant_colony_optimization_algorithms
type AntColony struct {
	bestAnt        *ant
	ants           []*ant
	config         *Config
	positions      []*connectedPosition
	toursGenerated int
	startTime      time.Time
}

// CreateAntColonyPathFinder constructs a new AntColony.
func CreateAntColonyPathFinder(config *Config) (wirepath.PathFinder, error) {
	colony := &AntColony{
		config:         config,
		toursGenerated: 0,
		startTime:      time.Now(),
	}

	if err := colony.init(); err != nil {
		return nil, err
	}

	return colony, nil
}

func (a *AntColony) init() error {
	posCount := len(a.config.Positions)

	a.positions = make([]*connectedPosition, posCount)
	for i := 0; i < posCount; i++ {
		a.positions[i] = createConnectedPosition(a.config.Positions, i)
	}

	a.ants = make([]*ant, a.config.AntCount)
	for i := 0; i < a.config.AntCount; i++ {
		a.ants[i] = &ant{}
	}

	return nil
}

func (a *AntColony) Reset() {
	for _, p := range a.positions {
		for _, n := range p.neighbours {
			n.pheromoneLevel = 0
		}
	}

	a.startTime = time.Now()
	a.toursGenerated = 0
	a.bestAnt = nil
}

func (a *AntColony) GetBestPath() []string {
	bestAnt := a.bestAnt
	ids := make([]string, len(a.positions))
	for i := 0; i < len(a.positions); i++ {
		ids[i] = a.positions[bestAnt.tour[i]].position.Name
	}

	return ids
}

func (a *AntColony) GetPositions() []*wirepath.Position {
	return a.config.Positions
}

func (a *AntColony) RunRound() error {
	posCount := len(a.positions)
	semaphore := sync.WaitGroup{}
	semaphore.Add(len(a.ants))
	for _, thisAnt := range a.ants {
		go func(thisAnt *ant) {
			thisAnt.reset(len(a.config.Positions), rand.Intn(len(a.config.Positions)-1))
			for i := 0; i < posCount-1; i++ {
				thisAnt.moveToNextPosition(a.positions)
			}

			startNeighbour := a.positions[thisAnt.tour[len(thisAnt.tour)-1]].findNeighbourWithPosIdx(thisAnt.tour[0])
			thisAnt.travelled += startNeighbour.distance
			thisAnt.weightedTravelled += startNeighbour.weightedDistance()

			semaphore.Done()
		}(thisAnt)
	}

	semaphore.Wait()

	sort.Slice(a.ants, func(i, j int) bool {
		return a.ants[i].travelled < a.ants[j].travelled
	})

	if a.bestAnt == nil || a.bestAnt.travelled > a.ants[0].travelled {
		a.bestAnt = a.ants[0]
		a.ants[0] = &ant{}
	}

	a.incrementPheromones(a.bestAnt)

	a.toursGenerated += a.config.AntCount

	return nil
}

func (a *AntColony) Stats() *wirepath.PathfinderStats {
	return &wirepath.PathfinderStats{
		PathsGenerated: a.toursGenerated,
		RunTime:        time.Now().Sub(a.startTime),
		ShortestPath:   a.bestAnt.travelled,
	}
}

func (a *AntColony) incrementPheromones(thisAnt *ant) {
	for i := 0; i < len(a.positions); i++ {
		nextIdx := i + 1
		thisPos := a.positions[thisAnt.tour[i]]
		thisPos.decayPheromones(a.config.PheromoneDecay)
		if i == len(a.positions)-1 {
			nextIdx = 0
		}

		nextPos := a.positions[thisAnt.tour[nextIdx]]

		forwardNeighbour := thisPos.findNeighbourWithPosIdx(thisAnt.tour[nextIdx])
		backwardNeighbour := nextPos.findNeighbourWithPosIdx(thisAnt.tour[i])

		forwardNeighbour.pheromoneLevel = math.Min(forwardNeighbour.pheromoneLevel+a.config.PheromoneSpreadForward, a.config.MaxPheromoneFactor)
		backwardNeighbour.pheromoneLevel = math.Min(backwardNeighbour.pheromoneLevel+a.config.PheromoneSpreadBackward, a.config.MaxPheromoneFactor)
	}
}
