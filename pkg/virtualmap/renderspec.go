package virtualmap

import (
	"errors"
	"math"

	"github.com/ataboo/go-metar-blink/pkg/geo"
	"github.com/ataboo/go-metar-blink/pkg/logger"
)

type StationRenderSpec struct {
	center       *geo.Coordinate
	mercatorSpec *geo.MercatorSpec
	widthDeg     float64
	heightDeg    float64
	paddingPx    int
	scaleFactor  float64
	imgWidthPx   int
	imgHeightPx  int
}

func (s *StationRenderSpec) mercatorWidthRad() float64 {
	return s.widthDeg * geo.DegToRad
}

func (s *StationRenderSpec) mercatorHeightRad() float64 {
	northCoord := geo.Coordinate{
		Latitude:  s.center.Latitude + s.heightDeg/2,
		Longitude: 0,
	}
	latRads, _ := northCoord.MercatorPosition(s.mercatorSpec)

	return latRads * 2
}

func CreateRenderSpec(coordinates []*geo.Coordinate, imgWidthPx int, imgHeightPx int, paddingPx int) *StationRenderSpec {
	spec := &StationRenderSpec{
		imgWidthPx:  imgWidthPx,
		imgHeightPx: imgHeightPx,
		paddingPx:   paddingPx,
	}

	spec.computeCenterAndDimensions(coordinates)
	scale, err := spec.computeScale()
	if err != nil {
		logger.LogError("failed to compute scale: %s", err)
	}

	spec.scaleFactor = scale

	return spec
}

func (s *StationRenderSpec) ProjectCoordinate(c *geo.Coordinate) (x float64, y float64) {
	latRads, longRads := c.MercatorPosition(s.mercatorSpec)

	y = math.Round(float64(s.imgHeightPx/2) - latRads*s.scaleFactor)
	x = math.Round(float64(s.imgWidthPx/2) + longRads*s.scaleFactor)

	return x, y
}

func (s *StationRenderSpec) computeScale() (scale float64, err error) {
	if s.heightDeg == 0 && s.widthDeg == 0 {
		return 0, errors.New("stations must have some positional range")
	}
	heightScale := float64(s.imgHeightPx-2*s.paddingPx) / s.mercatorHeightRad()
	widthScale := float64(s.imgWidthPx-2*s.paddingPx) / s.mercatorWidthRad()

	if s.heightDeg == 0 || widthScale < heightScale {
		return widthScale, nil
	}

	return heightScale, nil
}

func (s *StationRenderSpec) computeCenterAndDimensions(coordinates []*geo.Coordinate) {
	latMin := 90.0
	latMax := -90.0
	longMin := 360.0
	longMax := 0.0

	// Get the min/max latitudes and min/max longitude normalized as a positive angle
	for _, coordinate := range coordinates {
		latMin = math.Min(latMin, coordinate.Latitude)
		latMax = math.Max(latMax, coordinate.Latitude)

		posLong := math.Mod(coordinate.Longitude+360.0, 360.0)
		longMin = math.Min(posLong, longMin)
		longMax = math.Max(posLong, longMax)
	}

	// Get center lat and long as simple average.
	centerLat := (latMin + latMax) / 2
	centerLong := (longMax + longMin) / 2
	height := latMax - latMin
	width := longMax - longMin

	// If we're going the long way around, rotate center by 180.
	if longMax-longMin > 180 {
		centerLong = math.Mod(centerLong+180, 360)
		width = 360 - width
	}

	// Normalize longitude back to  +/- 180.
	if centerLong > 180 {
		centerLong -= 360
	}

	s.center = &geo.Coordinate{
		Latitude:  centerLat,
		Longitude: centerLong,
	}
	s.widthDeg = width
	s.heightDeg = height
	s.mercatorSpec = &geo.MercatorSpec{
		LatCenter:     centerLat,
		LongCenter:    centerLong,
		LatitudeScale: 1.4,
	}
}
