package antcolony

import (
	"github.com/ataboo/go-metar-blink/pkg/wirepath"
)

type Config struct {
	AntCount        int
	PheromoneWeight float64
	Positions       []*Position
}

type Position struct {
	Name string
	X    int
	Y    int
}

type AntColony struct {
	ants       []*Ant
	config     *Config
	pheromones [][]float64
}

type tourQueue []*Position

func (t tourQueue) RemoveAt(idx int) {
	t = append(t[:idx], t[idx+1:]...)
}

type Ant struct {
	unvisited tourQueue
	tour      []*Position
	position  int
}

func CreateAntColonyPathFinder(config *Config) (wirepath.PathFinder, error) {
	colony := &AntColony{
		config: config,
	}

	if err := colony.init(); err != nil {
		return nil, err
	}

	return *&colony, nil
}

func (a *AntColony) init() error {
	posCount := len(a.config.Positions)
	a.pheromones = make([][]float64, posCount)
	for i := range a.pheromones {
		a.pheromones[i] = make([]float64, posCount)
	}

	a.ants = make([]*Ant, a.config.AntCount)
}

func (a *AntColony) GetBestPath() *wirepath.StationPath {

}

func (a *AntColony) RunRound() error {

}

func (a *AntColony) Stats() *wirepath.PathfinderStats {

}
