package animation

import "testing"

func TestColorLerping(t *testing.T) {
	start := ColorRed
	end := ColorBlue

	val := lerpColor(start, end, 0, lerpByte)
	assertColorsMatchExpected((0xFF0000), val, t)

	val = lerpColor(start, end, 0.5, lerpByte)
	assertColorsMatchExpected(0x800080, val, t)

	val = lerpColor(start, end, 1, lerpByte)
	assertColorsMatchExpected(0x0000FF, val, t)
}

func TestColorToString(t *testing.T) {
	red := Color(0x123456)
	if red.String() != "0x123456 [R12, G34, B56]" {
		t.Error("unnexpected color string value", red.String())
	}
}

func TestCosineBytePeak(t *testing.T) {
	table := []struct {
		start    byte
		end      byte
		mu       float64
		expected byte
	}{
		{0x00, 0xFF, 0, 0},
		{0x00, 0xFF, .1, 0x18},
		{0x00, 0xFF, .25, 0x7F},
		{0x00, 0xFF, 0.4, 0xE7},
		{0x00, 0xFF, 0.5, 0xFF},
		{0x00, 0xFF, 0.6, 0xE7},
		{0x00, 0xFF, .75, 0x80},
		{0x00, 0xFF, .9, 0x18},
		{0x00, 0xFF, 1, 0},
	}

	for _, row := range table {
		result := cosinePeakByte(row.start, row.end, row.mu)
		if result != row.expected {
			t.Errorf("[0x%x, 0x%x, %f, 0x%x] => 0x%x", row.start, row.end, row.mu, row.expected, result)
		}
	}
}

func TestSharpBytePeak(t *testing.T) {
	table := []struct {
		start    byte
		end      byte
		mu       float64
		expected byte
	}{
		{0x00, 0xFF, 0, 0},
		{0x00, 0xFF, .1, 0x33},
		{0x00, 0xFF, .25, 0x80},
		{0x00, 0xFF, .4, 0xCC},
		{0x00, 0xFF, 0.5, 0xFF},
		{0x00, 0xFF, .6, 0xCC},
		{0x00, 0xFF, .75, 0x80},
		{0x00, 0xFF, .9, 0x33},
		{0x00, 0xFF, 1, 0},
	}

	for _, row := range table {
		result := sharpPeakByte(row.start, row.end, row.mu)
		if result != row.expected {
			t.Errorf("[0x%x, 0x%x, %f, 0x%x] => 0x%x", row.start, row.end, row.mu, row.expected, result)
		}
	}
}

func TestARGB(t *testing.T) {
	table := []struct {
		color    Color
		expected uint32
	}{
		{0xFFFFFF, 0xFFFFFFFF},
		{0x000000, 0xFF000000},
		{0x112233, 0xFF112233},
	}

	for _, row := range table {
		if row.color.ARGB() != row.expected {
			t.Errorf("unnexpected argb: 0x%x => 0x%x, 0x%x", row.color, row.color.ARGB(), row.expected)
		}
	}
}

func TestRGBA(t *testing.T) {
	table := []struct {
		color    Color
		expected uint32
	}{
		{0xFFFFFF, 0xFFFFFFFF},
		{0x000000, 0x000000FF},
		{0x112233, 0x112233FF},
	}

	for _, row := range table {
		if row.color.RGBA() != row.expected {
			t.Errorf("unnexpected argb: 0x%x => 0x%x, 0x%x", row.color, row.color.RGBA(), row.expected)
		}
	}
}

func TestRGB(t *testing.T) {
	table := []struct {
		color    Color
		expected uint32
	}{
		{0xFFFFFF, 0xFFFFFF},
		{0x000000, 0x000000},
		{0x112233, 0x112233},
	}

	for _, row := range table {
		if row.color.RGB() != row.expected {
			t.Errorf("unnexpected argb: 0x%x => 0x%x, 0x%x", row.color, row.color.RGBA(), row.expected)
		}
	}
}

func assertColorsMatchExpected(expected Color, actual Color, t *testing.T) {
	if expected != actual {
		t.Errorf("expected %s instead of %s", expected.String(), actual.String())
	}
}
