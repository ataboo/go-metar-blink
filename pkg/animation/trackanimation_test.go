package animation

import (
	"testing"
	"time"
)

func TestTrackAnimation(t *testing.T) {
	values := make(map[int]Color)
	animation := CreateTrackAnimation(createTestTracks(), 10).(*TrackAnimation)

	animation.GetValues(values)

	assertTrackValues(values, t, "starting value", 0x00FF00, 0)
	animation.Update(time.Millisecond*400, values)
	assertTrackValues(values, t, "stopped 1", 0x00FF00, 0)
	assertTrackPositions(animation.tracks, t, "stopped 1", 0, 0)

	animation.Start()

	animation.Update(time.Millisecond*400, values)
	assertTrackValues(values, t, "4th frame", 0xFF00FF, 0)
	assertTrackPositions(animation.tracks, t, "4th frame", 4, 4)

	animation.Stop()
	animation.Update(time.Millisecond*400, values)
	assertTrackValues(values, t, "stopped 2", 0xFF00FF, 0)
	assertTrackPositions(animation.tracks, t, "stopped 2", 4, 4)

	animation.Start()

	animation.Update(time.Millisecond*100, values)
	assertTrackValues(values, t, "5th frame", 0xCC33CC, 0x400040)
	assertTrackPositions(animation.tracks, t, "5th frame", 5, 5)

	animation.Update(time.Millisecond*50, values)
	assertTrackValues(values, t, "5th frame 1/2 step", 0xCC33CC, 0x400040)
	assertTrackPositions(animation.tracks, t, "5th frame", 5, 5)

	animation.Update(time.Millisecond*50, values)
	assertTrackValues(values, t, "6th frame", 0x996699, 0x800080)
	assertTrackPositions(animation.tracks, t, "6th frame", 6, 6)

	animation.Update(time.Millisecond*400, values)
	assertTrackValues(values, t, "10th frame", 0x00FF00, 0x800080)
	assertTrackPositions(animation.tracks, t, "10th frame", 0, 10)

	animation.Update(time.Millisecond*400, values)
	assertTrackValues(values, t, "14th frame", 0xFF00FF, 0x00FF00)
	assertTrackPositions(animation.tracks, t, "14th frame", 4, 14)
	if animation.frame != 14 {
		t.Error("unnexpected frame", animation.frame)
	}

	animation.Update(time.Millisecond*500, values)
	assertTrackValues(values, t, "19th frame", 0x00FF00, 0x00FF00)
	assertTrackPositions(animation.tracks, t, "19th frame", 9, 14)
	if animation.frame != 19 {
		t.Error("unnexpected frame", animation.frame)
	}

	animation.Reset()
	animation.GetValues(values)
	assertTrackValues(values, t, "reset", 0x00FF00, 0)
	assertTrackPositions(animation.tracks, t, "reset", 0, 0)
}

func assertTrackValues(values map[int]Color, t *testing.T, message string, expectedValues ...Color) {
	for i, val := range values {
		if expectedValues[i] != val {
			t.Errorf("%s | Unnexpected color: %s, expected: %s", message, val, expectedValues[i])
		}
	}
}

func assertTrackPositions(tracks []*Track, t *testing.T, message string, expectedValues ...int) {
	for i, track := range tracks {
		if expectedValues[i] != track.position {
			t.Errorf("%s | Unnexpected track position: %d, expected: %d", message, track.position, expectedValues[i])
		}
	}
}

func createTestTracks() []*Track {
	looping, _ := CreateTrack(10, true, []KeyFrame{
		{0, CreateColor(0, 0xFF, 0)},
		{4, CreateColor(0xFF, 0, 0xFF)},
		{9, CreateColor(0, 0xFF, 0)},
	})
	looping.ChannelIDs = []int{0}
	nonLooping, _ := CreateTrack(15, false, []KeyFrame{
		{0, 0},
		{4, 0},
		{8, CreateColor(0xFF, 0, 0xFF)},
		{12, 0},
		{14, CreateColor(0, 0xFF, 0)},
	})
	nonLooping.ChannelIDs = []int{1}

	return []*Track{looping, nonLooping}
}
