package animation

import (
	"math"
	"time"
)

// Animation represents playback of a sequence of values.
type Animation interface {
	Reset()
	Start()
	Stop()
	Update(delta time.Duration, values map[int]Color)
	GetValues(values map[int]Color)
}

type byteInterpolation func(start, end byte, mu float64) byte

func lerpColor(start Color, end Color, mu float64, interpFunc byteInterpolation) Color {
	return CreateColor(
		interpFunc(start.R(), end.R(), mu),
		interpFunc(start.G(), end.G(), mu),
		interpFunc(start.B(), end.B(), mu),
	)
}

func lerpByte(start byte, end byte, mu float64) byte {
	return byte(math.Round(float64(start)*(1-mu) + float64(end)*mu))
}

func sharpPeakByte(start byte, end byte, mu float64) byte {
	factor := 1 - 2*math.Abs(mu-0.5)

	return byte(math.Round(float64(start)*(1-factor) + float64(end)*factor))
}

func cosinePeakByte(startVal byte, endVal byte, mu float64) byte {
	factor := (1 - math.Cos(mu*math.Pi*2)) / 2

	return byte(math.Round(float64(startVal)*(1-factor) + float64(endVal)*factor))
}
