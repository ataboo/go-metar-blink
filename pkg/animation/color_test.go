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

func assertColorsMatchExpected(expected Color, actual Color, t *testing.T) {
	if expected != actual {
		t.Errorf("expected %s instead of %s", expected.String(), actual.String())
	}
}
