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

/*
func sharpPeakByte(start byte, end byte, mu float64) byte {
	factor := 1 - 2*math.Abs(mu-0.5)

	return byte(math.Round(float64(start)*(1-factor) + float64(end)*factor))
}

func cosinePeakByte(startVal byte, endVal byte, mu float64) byte {
	factor := (1 - math.Cos(mu*math.Pi*2)) / 2

	return byte(math.Round(float64(startVal)*(1-factor) + float64(endVal)*factor))
}
*/

func assertColorsMatchExpected(expected Color, actual Color, t *testing.T) {
	if expected != actual {
		t.Errorf("expected %s instead of %s", expected.String(), actual.String())
	}
}
