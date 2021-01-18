package geo

import (
	"math"
	"testing"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

func TestZeroedMercatorProjection(t *testing.T) {
	spec := &MercatorSpec{
		LatCenter:  0,
		LongCenter: 0,
	}

	table := []struct {
		Lat  float64
		Long float64
		xExp float64
		yExp float64
	}{
		{0, 0, 0, 0},
		{30, 180, 0.549306, 3.141593},
		{-30, -180, -0.549306, 3.141593},
		{60, 90, 1.316958, 1.570796},
		{-60, -90, -1.316958, -1.570796},
	}

	for _, row := range table {
		assertMercatorPos(Coordinate{row.Lat, row.Long, 0}, spec, t, row.xExp, row.yExp)
	}
}

func TestNormalizeAngle(t *testing.T) {
	table := []struct {
		angle    float64
		expected float64
	}{
		{math.Pi + 0.5, 0.5 - math.Pi},
		{-math.Pi - 0.5, math.Pi - 0.5},
	}

	for _, row := range table {
		normalized := common.NormalizePlusMinusPi(row.angle)
		if !common.Similar(normalized, row.expected) {
			t.Errorf("Got %f instead of %f", normalized, row.expected)
		}
	}
}

func TestOffsetMercatorProjection(t *testing.T) {
	spec := &MercatorSpec{
		LatCenter:  20,
		LongCenter: -140,
	}

	table := []struct {
		Lat  float64
		Long float64
		xExp float64
		yExp float64
	}{
		{20, -140, 0, 0},
		{-10, -170, -.549306, -0.523599},
		{50, -110, .549306, 0.523599},
		{-5, -90, -0.450875, 0.872665},
		{20, 170, 0, -0.872665},
	}

	for _, row := range table {
		assertMercatorPos(Coordinate{row.Lat, row.Long, 0}, spec, t, row.xExp, row.yExp)
	}
}

func assertMercatorPos(coord Coordinate, spec *MercatorSpec, t *testing.T, latExpected, longExpected float64) {
	latRads, longRads := coord.MercatorPosition(spec)

	if !common.Similar(latRads, latExpected) || !common.Similar(longRads, longExpected) {
		t.Errorf(
			"Coordinate: (%f, %f), Expected: (%f, %f), got: (%f, %f)",
			coord.Latitude,
			coord.Longitude,
			latExpected,
			longExpected,
			latRads,
			longRads,
		)
	}
}
