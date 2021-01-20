package animation

import "time"

// TrackAnimation is an animation containing a number of track sequences.
type TrackAnimation struct {
	running bool
	tracks  []*Track
	fps     int
	runTime time.Duration
	frame   uint64
}

// CreateTrackAnimation creates a new TrackAnimation
func CreateTrackAnimation(tracks []*Track, fps int) Animation {
	return &TrackAnimation{
		running: false,
		tracks:  tracks,
		fps:     fps,
	}
}

// Reset sets all the tracks to their first frame.
func (a *TrackAnimation) Reset() {
	a.forEachTrack(func(track *Track) {
		track.Seek(0)
	})

	a.runTime = 0
	a.frame = 0
}

// Update advances the tracks by a single frame and sets the new channel value in the `values` map.
func (a *TrackAnimation) Update(delta time.Duration, values map[int]Color) {
	stepCount := 0
	if a.running {
		a.runTime += delta
		targetFrame := uint64(a.runTime.Seconds() * float64(a.fps))
		stepCount = int(targetFrame - a.frame)
		a.frame = targetFrame
	}

	a.getValuesFromAllTracks(values, stepCount)
}

func (a *TrackAnimation) Step(values map[int]Color) {
	if a.running {
		a.getValuesFromAllTracks(values, 1)
	}
}

// Start allows advancing of the tracks on Update.
func (a *TrackAnimation) Start() {
	a.running = true
}

// Stop stops the advancing of the tracks on Update.
func (a *TrackAnimation) Stop() {
	a.running = false
}

// GetValues reads the values for all tracks mapped to the appropriate channel.
func (a *TrackAnimation) GetValues(values map[int]Color) {
	a.getValuesFromAllTracks(values, 0)
}

func (a *TrackAnimation) getValuesFromAllTracks(values map[int]Color, steps int) {
	a.forEachTrack(func(track *Track) {
		if steps > 0 {
			track.Step(steps)
		}
		for _, trackChan := range track.ChannelIDs {
			values[trackChan] = track.Value()
		}
	})
}

func (a *TrackAnimation) forEachTrack(action func(track *Track)) {
	for _, track := range a.tracks {
		action(track)
	}
}
