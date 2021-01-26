package antcolony

import "errors"

// ant visits all positions, tracking its path as a tour.
type ant struct {
	visited           map[int]bool
	tour              []int
	step              int
	travelled         float64
	weightedTravelled float64
}

// reset resets the ants values and starts it at the position specified.
func (a *ant) reset(posCount int, startingPosIdx int) {
	a.tour = make([]int, posCount)
	a.visited = make(map[int]bool)
	a.step = 0
	a.travelled = 0
	a.weightedTravelled = 0
	a.visited[startingPosIdx] = true
	a.tour[0] = startingPosIdx
}

// currentPosition gets the index of the ant's current position.
func (a *ant) currentPosition() int {
	return a.tour[a.step]
}

func (a *ant) hasVisitedPos(posIdx int) bool {
	visited, ok := a.visited[posIdx]

	return ok && visited
}

// moveToNextPosition moves the ant to the next unvisited location favouring close neighbours and pheromone levels.
func (a *ant) moveToNextPosition(positions []*connectedPosition) error {
	currentPos := positions[a.currentPosition()]

	var bestNeighbour *neighbour = nil
	bestWeightedDistance := float64(0)
	for _, n := range currentPos.neighbours {
		if a.hasVisitedPos(n.posIdx) {
			continue
		}

		weightedDistance := n.weightedDistance()
		if bestNeighbour == nil || bestWeightedDistance > weightedDistance {
			bestNeighbour = n
			bestWeightedDistance = weightedDistance
		}
	}

	if bestNeighbour == nil {
		return errors.New("fell through move to next without finding a neighbour")
	}

	a.step++
	a.tour[a.step] = bestNeighbour.posIdx
	a.travelled += bestNeighbour.distance
	a.weightedTravelled += bestWeightedDistance
	a.visited[bestNeighbour.posIdx] = true

	return nil
}
