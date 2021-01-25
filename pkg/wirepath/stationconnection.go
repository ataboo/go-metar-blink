package wirepath

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type StationConnection struct {
	ID     string
	EndID  string
	Length float64
}

func (c *StationConnection) CalculateLength(positions map[string]*sdl.Point) error {
	startPos, ok := positions[c.ID]
	if !ok {
		return fmt.Errorf("failed to get position of station '%s'", c.ID)
	}
	endPos, ok := positions[c.EndID]
	if !ok {
		return fmt.Errorf("failed to get position of station '%s'", c.EndID)
	}

	deltaX := float64(startPos.X - endPos.X)
	deltaY := float64(startPos.Y - endPos.Y)

	c.Length = math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2))

	return nil
}
