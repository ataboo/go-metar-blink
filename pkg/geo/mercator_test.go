package geo

import (
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

func TestCoordinateEquality(t *testing.T) {
	table := []struct {
		coordA   Coordinate
		coordB   Coordinate
		expected bool
	}{
		{Coordinate{0, 0, 0}, Coordinate{0, 0, 0}, true},
		{Coordinate{1, 2, 3}, Coordinate{1, 2, 3}, true},
		{Coordinate{1, 2, 3}, Coordinate{0, 2, 3}, false},
		{Coordinate{1, 2, 3}, Coordinate{1, 0, 3}, false},
		{Coordinate{1, 2, 3}, Coordinate{1, 2, 0}, false},
	}

	for _, row := range table {
		if row.coordA.Equal(&row.coordB) != row.expected {
			t.Errorf("%+v, %+v => %t, %t", row.coordA, row.coordB, row.coordA.Equal(&row.coordB), row.expected)
		}
	}
}

func TestParseDecimalCoord(t *testing.T) {
	coord, err := ParseDecimalCoordinate("1.234", "5.678")
	if err != nil {
		t.Error(err)
	}

	if !coord.Equal(&Coordinate{Latitude: 1.234, Longitude: 5.678, Altitude: 0}) {
		t.Error("unnexpected value")
	}

	_, err = ParseDecimalCoordinate("invalid", "5.678")
	if err == nil {
		t.Error("expected error")
	}

	_, err = ParseDecimalCoordinate("1.234", "invalid")
	if err == nil {
		t.Error("expected error")
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
