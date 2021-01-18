package geo

import (
	"math"
	"strconv"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

const (
	EarthRadius        = float64(6371)
	EarthCircumference = float64(40030)
	DegToRad           = float64(math.Pi / 180.0)
	RadToDeg           = float64(180.0 / math.Pi)
)

type MercatorSpec struct {
	LatCenter  float64
	LongCenter float64
}

func (m *MercatorSpec) Center() *Coordinate {
	return &Coordinate{
		Latitude:  m.LatCenter,
		Longitude: m.LongCenter,
	}
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

func (c *Coordinate) Equal(other *Coordinate) bool {
	return c.Latitude == other.Latitude && c.Longitude == other.Longitude
}

func ParseDecimalCoordinate(latStr, longStr string) (coord *Coordinate, err error) {
	coord = &Coordinate{}

	coord.Latitude, err = strconv.ParseFloat(latStr, 64)
	if err != nil {
		return nil, err
	}

	coord.Longitude, err = strconv.ParseFloat(longStr, 64)
	if err != nil {
		return nil, err
	}

	return coord, err
}

func (c *Coordinate) MercatorPosition(spec *MercatorSpec) (latRads float64, longRads float64) {
	center := spec.Center()

	longRads = common.NormalizePlusMinusPi((c.Longitude - center.Longitude) * DegToRad)

	latRads = math.Log(math.Tan(math.Pi/4 + DegToRad*(c.Latitude-center.Latitude)/2))
	return latRads, longRads
}
