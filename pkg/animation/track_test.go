package animation

import (
	"errors"
	"testing"
)

func TestLoopingTrackStepping(t *testing.T) {
	track := createExampleTrack(true)

	if !track.IsLooping() {
		t.Error("unnexpected value")
	}

	if track.Value() != 50 {
		t.Error("unnexpected value", track.Value())
	}

	ok := track.Step(1)
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

	ok = track.Step(1)
	if !ok || track.position != 0 {
		t.Error("failed to step")
	}

	if track.Value() != 50 {
		t.Error("unnexpected value", track.Value())
	}
}

func TestNonLoopingTrackValues(t *testing.T) {
	track := createExampleTrack(false)

	if track.IsLooping() {
		t.Error("unnexpected value")
	}

	if track.Value() != 20 {
		t.Error("unnexpected track value", track.Value())
	}

	track.Step(1)
	if track.Value() != 20 {
		t.Error("unnexpected value", track.Value())
	}

	track.Seek(15)
	if track.Value() != 40 {
		t.Error("unnexpected value", track.Value())
	}

	err := track.Seek(40)
	if err != nil || track.GetPosition() != 40 {
		t.Error("failed to set position")
	}

	if track.Value() != 80 {
		t.Error("unnexpected value", track.Value())
	}

	track.Seek(49)
	if track.Value() != 80 {
		t.Error("unnexpected error", track.Value())
	}

	ok := track.Step(1)
	if ok || track.position != 49 {
		t.Error("failed to not step")
	}
}

func TestTrackSeek(t *testing.T) {
	track := createExampleTrack(true)

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

func TestNormalizeShortKeyFrames(t *testing.T) {
	track, err := CreateTrack(16, true, []KeyFrame{
		{0, 5},
	})
	if err != nil {
		t.Error(err)
	}

	if len(track.keyFrames) != 2 {
		t.Error("failed to make end key frame")
	}

	if track.keyFrames[0].Position != 0 || track.keyFrames[0].Value != 5 {
		t.Error("unnexpected initial keyframe")
	}

	if track.keyFrames[1].Position != 15 || track.keyFrames[1].Value != 5 {
		t.Error("unnexpected final keyframe")
	}

	err = track.SetKeyFrames([]KeyFrame{}, 16, true)
	if err != nil {
		t.Error(err)
	}

	if len(track.keyFrames) != 2 {
		t.Error("failed to make keyframes")
	}

	if track.keyFrames[0].Position != 0 || track.keyFrames[0].Value != 0 {
		t.Error("unnexpected initial keyframe")
	}

	if track.keyFrames[1].Position != 15 || track.keyFrames[1].Value != 0 {
		t.Error("unnexpected final keyframe")
	}
}

func TestKeyFramesWithSamePositionRejected(t *testing.T) {
	_, err := CreateTrack(10, true, []KeyFrame{
		{5, 5},
		{5, 10},
	})

	if err == nil {
		t.Error("expected error")
	}
}

func TestSingleFrameAnimationValid(t *testing.T) {
	track, err := CreateTrack(1, false, []KeyFrame{
		{
			Position: 0,
			Value:    42,
		},
	})

	if err != nil {
		t.Error(err)
	}

	track.Step(1)

	if track.GetPosition() != 0 {
		t.Error("unnexpected position", track.GetPosition())
	}

	if track.GetLength() != 1 {
		t.Error("unnexpected length", track.GetLength())
	}
}

func TestPanicIfLerpSameFrame(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	track := Track{
		looping:  true,
		position: 0,
		length:   10,
	}

	track.lerpFrameValues(KeyFrame{5, 0}, KeyFrame{5, 5}, 0)
}

func TestTrackPanicsWithNonNormalizedKeyframes(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	track := Track{
		keyFrames: []KeyFrame{{1, 2}},
	}

	track.Value()
}

func TestSetLength(t *testing.T) {
	_, err := CreateTrack(-1, true, []KeyFrame{})
	if err == nil {
		t.Error(errors.New("expected error"))
	}

	track, err := CreateTrack(10, true, []KeyFrame{})
	if err != nil {
		t.Error(err)
	}

	if track.GetLength() != 10 {
		t.Error("unnexpected length")
	}

	err = track.SetLength(-1)
	if err == nil {
		t.Error(errors.New("expected error"))
	}

	track.Seek(9)

	err = track.SetLength(5)
	if err != nil {
		t.Error(err)
	}

	if track.position != 4 {
		t.Error("unnexpected position")
	}
}

func createExampleTrack(looping bool) *Track {
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

	track, _ := CreateTrack(50, looping, keyFrames)

	return track
}
