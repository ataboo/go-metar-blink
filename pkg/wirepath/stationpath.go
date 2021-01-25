package wirepath

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type StationPath struct {
	Connections map[string]*StationConnection
	Length      float64
	StartID     string
}

func CreateStationPath(stationIDs []string, startID string) *StationPath {
	path := &StationPath{
		Connections: make(map[string]*StationConnection, len(stationIDs)),
		Length:      0,
		StartID:     startID,
	}

	for _, id := range stationIDs {
		path.Connections[id] = &StationConnection{
			ID:     id,
			EndID:  "",
			Length: 0,
		}
	}

	return path
}

func (p *StationPath) Randomize() {
	ids := make([]string, len(p.Connections))
	idx := 0
	for k := range p.Connections {
		ids[idx] = k
		idx++
	}

	rand.Shuffle(len(ids), func(i, j int) {
		tmp := ids[j]
		ids[j] = ids[i]
		ids[i] = tmp
	})

	for i, id := range ids {
		if i == len(ids)-1 {
			p.Connections[id].EndID = p.Connections[ids[0]].ID
		} else {
			p.Connections[id].EndID = p.Connections[ids[i+1]].ID
		}
	}
}

func (p StationPath) String() string {
	ids, err := p.GetIDsOrdered()
	if err != nil {
		return fmt.Sprintf("failed to render: %s", err)
	}

	return strings.Join(ids, ", ")
}

func (p *StationPath) GetIDsOrdered() ([]string, error) {
	ids := make([]string, len(p.Connections))

	connection := p.Connections[p.StartID]
	for i := 0; i < len(p.Connections); i++ {
		ids[i] = connection.ID
		if connection.EndID == "" {
			return nil, fmt.Errorf("connection '%s' has no end id", connection.ID)
		}
		connection = p.Connections[connection.EndID]
	}

	return ids, nil
}

func (p *StationPath) CalculateLength(positions map[string]*sdl.Point) error {
	sum := 0.0

	for _, c := range p.Connections {
		if err := c.CalculateLength(positions); err != nil {
			return err
		}
		sum += c.Length
	}

	p.Length = sum

	return nil
}
