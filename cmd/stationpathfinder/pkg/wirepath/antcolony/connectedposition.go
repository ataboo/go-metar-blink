package antcolony

import (
	"math"

	"github.com/ataboo/go-metar-blink/cmd/stationpathfinder/pkg/wirepath"
)

// connectedPositions represents a position for travel.
type connectedPosition struct {
	position   *wirepath.Position
	neighbours []*neighbour
}

// neighbour represents a path between a connectedPosition and another one.
type neighbour struct {
	distance       float64
	posIdx         int
	pheromoneLevel float64
}

// weightedDistance takes the pheromone level of a neighbour to skew the distance.
func (n *neighbour) weightedDistance() float64 {
	return (1 - n.pheromoneLevel) * n.distance
}

// findNeighbourWithPosIdx finds the neighbour with the matching position index.
func (c *connectedPosition) findNeighbourWithPosIdx(posIdx int) *neighbour {
	for _, n := range c.neighbours {
		if n.posIdx == posIdx {
			return n
		}
	}

	panic("failed to find neighbour matching idx")
}

// decayPheromones reduces the pheromone level for all neighbours in the connectedPosition.
func (c *connectedPosition) decayPheromones(pheromoneDecay float64) {
	for _, n := range c.neighbours {
		n.pheromoneLevel = math.Max(0, n.pheromoneLevel-pheromoneDecay)
	}
}

// createConnectedPosition creates a new connected position centered on the positions specified by thisIdx, with neighbours to all the other positions.
func createConnectedPosition(positions []*wirepath.Position, thisIdx int) *connectedPosition {
	thisPos := positions[thisIdx]
	connected := &connectedPosition{
		position:   thisPos,
		neighbours: make([]*neighbour, len(positions)-1),
	}

	neighbourIdx := 0
	for i, p := range positions {
		if i == thisIdx {
			continue
		}
		connected.neighbours[neighbourIdx] = &neighbour{
			distance: math.Sqrt(math.Pow(float64(thisPos.X)-float64(p.X), 2) + math.Pow(float64(thisPos.Y)-float64(p.Y), 2)),
			posIdx:   i,
		}
		neighbourIdx++
	}

	return connected
}
