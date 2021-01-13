package animation

import "testing"

func TestTrackAnimation(t *testing.T) {
	tracks := createTestTracks()
	CreateTrackAnimation(tracks, 10)

}

func createTestTracks() []*Track {
	looping, _ := CreateTrack(10, true, []KeyFrame{
		{0, 0},
		{4, 100},
		{9, 0},
	})
	nonLooping, _ := CreateTrack(15, false, []KeyFrame{
		{0, 0},
		{4, 0},
		{8, 100},
		{12, 0},
		{14, 0},
	})

	return []*Track{looping, nonLooping}
}
