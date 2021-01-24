package animation

import (
	"time"
)

// PulseAnimation loops between 2 values.
type PulseAnimation struct {
	value     Color
	running   bool
	start     Color
	end       Color
	position  time.Duration
	period    time.Duration
	channels  []int
	interFunc byteInterpolation
	fps       int
}

// CreatePulseAnimation creates a new pulse animation.
func CreatePulseAnimation(period time.Duration, start Color, end Color, channels []int, fps int) Animation {
	return &PulseAnimation{
		value:     start,
		running:   false,
		start:     start,
		end:       end,
		period:    period,
		position:  time.Duration(0),
		channels:  channels,
		interFunc: cosinePeakByte,
		fps:       fps,
	}
}

// Reset starts the animation from the beginning.
func (a *PulseAnimation) Reset() {
	a.position = time.Duration(0)
}

// Start starts the animation.
func (a *PulseAnimation) Start() {
	a.running = true
}

// Stop stops the animation at the current position.
func (a *PulseAnimation) Stop() {
	a.running = false
}

// Update ticks the animation foward and gets the channel values in a map.
func (a *PulseAnimation) Update(delta time.Duration, values map[int]Color) {
	if a.running {
		a.position = (a.position + delta) % a.period
	}

	a.GetValues(values)
}

func (a *PulseAnimation) Step(values map[int]Color) {
	a.Update(time.Second/time.Duration(a.fps), values)
}

// GetValues gets the values for each channel.
func (a *PulseAnimation) GetValues(values map[int]Color) {
	value := lerpColor(a.start, a.end, float64(a.position)/float64(a.period), a.interFunc)
	for _, channel := range a.channels {
		values[channel] = value
	}
}
