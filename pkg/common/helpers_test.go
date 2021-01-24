package common

import (
	"math"
	"testing"
)

func TestNormalizePlusMinusPi(t *testing.T) {
	table := []struct {
		angleRads float64
		expected  float64
	}{
		{0, 0},
		{math.Pi * 2, 0},
		{-math.Pi * 2, 0},
		{math.Pi * 4, 0},
		{-math.Pi * 4, 0},
		{math.Pi / 2, math.Pi / 2},
		{5 * math.Pi / 2, math.Pi / 2},
		{-5 * math.Pi / 2, -math.Pi / 2},
		{3 * math.Pi / 2, -math.Pi / 2},
		{-3 * math.Pi / 2, math.Pi / 2},
		{math.Pi, math.Pi},
		{-math.Pi, math.Pi},
	}

	for _, row := range table {
		normalized := NormalizePlusMinusPi(row.angleRads)
		if normalized != row.expected {
			t.Errorf("%f => %f, %f", row.angleRads, normalized, row.expected)
		}
	}
}

func TestSimilar(t *testing.T) {
	table := []struct {
		first    float64
		second   float64
		expected bool
	}{
		{0, 9e-7, true},
		{0, 1e-6, false},
		{0, -9e-7, true},
		{0, -9e-6, false},
		{10, 10 + 9e-7, true},
		{10, 10 + 1e-6, true},
		{10, 10 - 9e-7, true},
		{10, 10 - 1e-6, true},
	}

	for _, row := range table {
		if row.expected != Similar(row.first, row.second) {
			t.Errorf("%f, %f => %t, %t", row.first, row.second, Similar(row.first, row.second), row.expected)
		}
	}
}
