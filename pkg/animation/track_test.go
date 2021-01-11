package animation

import "testing"

func TestLoopingTrackStepping(t *testing.T) {
	track := createExampleTrack()

	if track.Value() != 50 {
		t.Error("unnexpected value", track.Value())
	}

	ok := track.Step()
	if !ok {
		t.Error("expected true response")
	}

	if track.GetPosition() != 1 {
		t.Error("unnexpected position", track.position)
	}

	if track.Value() != 47 {
		t.Error("unnexpected value", track.Value())
	}

	track.Seek(49)
	if track.Value() != 53 {
		t.Error("unnexpected value", track.Value())
	}

	ok = track.Step()
	if !ok {
		t.Error("expected ok")
	}

	if track.Value() != 50 {
		t.Error("unnexpected value", track.Value())
	}
}

// func TestNonLoopingTrackValues() {
// 	// track := createExampleTrack()

// }

func TestTrackSeek(t *testing.T) {
	track := createExampleTrack()

	err := track.Seek(-1)
	if err == nil {
		t.Error("expected error")
	}

	err = track.Seek(50)
	if err == nil {
		t.Error("expected error")
	}

	err = track.Seek(0)
	if err != nil {
		t.Error(err)
	}

	if track.GetPosition() != 0 {
		t.Error("unnexpected position", track.position)
	}

	err = track.Seek(49)
	if err != nil {
		t.Error(err)
	}

	if track.GetPosition() != 49 {
		t.Error("unnexpected position", track.position)
	}
}

func createExampleTrack() Track {
	keyFrames := []KeyFrame{
		{
			Position: 10,
			Value:    20,
		},
		{
			Position: 20,
			Value:    60,
		},
		{
			Position: 30,
			Value:    40,
		},
		{
			Position: 40,
			Value:    80,
		},
	}

	track := Track{
		Looping:    true,
		ChannelIDs: []int{},
	}

	track.SetKeyFrames(keyFrames, 50)

	return track
}
