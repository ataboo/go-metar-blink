package animation

import (
	"errors"
	"sort"
)

// Track is a sequence of for playback that will give the value according to keyframes.
type Track struct {
	looping    bool
	ChannelIDs []int
	position   int
	length     int
	keyFrames  []KeyFrame
}

// KeyFrame is a position and value for a track to interpolate between.
type KeyFrame struct {
	Position int
	Value    Color
}

// CreateTrack create a new track.
func CreateTrack(length int, looping bool, keyFrames []KeyFrame) (*Track, error) {
	track := Track{
		position: 0,
	}

	err := track.SetKeyFrames(keyFrames, length, looping)
	if err != nil {
		return nil, err
	}

	return &track, nil
}

// Step steps the track forward to the next frame.
// A looping track will return to 0 after the last frame.
// A non-looping track will stop at the last frame, and return false.
func (t *Track) Step(count int) bool {
	finalPos := (t.position + count)
	stepped := true

	if t.looping {
		finalPos %= t.length
	} else if finalPos > t.length-1 {
		finalPos = t.length - 1
		stepped = t.position != t.length-1
	}

	t.Seek(finalPos)

	return stepped
}

// GetPosition Gets the current position of this track in frames.
func (t *Track) GetPosition() int {
	return t.position
}

// Seek moves the track to a specific frame.
func (t *Track) Seek(position int) error {
	if position < 0 || position >= t.length {
		return errors.New("position out of range")
	}
	t.position = position

	return nil
}

// SetLength sets the length of the track in frames.
func (t *Track) SetLength(length int) error {
	if length < 1 {
		return errors.New("length out of range")
	}

	t.length = length
	if t.position > length-1 {
		t.position = length - 1
	}

	return nil
}

// GetLength gets the length of the track in frames.
func (t *Track) GetLength() int {
	return t.length
}

// IsLooping returns whether the track restarts when it reaches the end.
func (t *Track) IsLooping() bool {
	return t.looping
}

// SetKeyFrames sets the key frames for this track, length, and if it is looping.
// New keyframes will be interpolated at the first and last keyframe depending on the length and looping behaviour.
func (t *Track) SetKeyFrames(keyFrames []KeyFrame, length int, looping bool) error {
	t.looping = looping

	err := t.SetLength(length)
	if err != nil {
		return err
	}

	t.keyFrames, err = t.normalizeKeyFrames(keyFrames, length)
	if err != nil {
		return err
	}

	return nil
}

// Value gets this track's value at the current position.
func (t *Track) Value() Color {
	if len(t.keyFrames) < 2 {
		panic("key frames must be normalized")
	}

	var startKeyIdx = -1
	var endKeyIdx = -1
	for i, key := range t.keyFrames {
		if key.Position == t.position {
			return key.Value
		}

		if key.Position > t.position {
			endKeyIdx = i
			break
		}

		startKeyIdx = i
	}

	return t.lerpFrameValues(t.keyFrames[startKeyIdx], t.keyFrames[endKeyIdx], t.position)
}

func (t *Track) normalizeKeyFrames(keyFrames []KeyFrame, length int) ([]KeyFrame, error) {
	sort.Slice(keyFrames, func(i int, j int) bool {
		return keyFrames[i].Position < keyFrames[j].Position
	})

	if len(keyFrames) < 2 {
		fixedValue := Color(0)
		if len(keyFrames) == 1 {
			fixedValue = keyFrames[0].Value
		}

		return []KeyFrame{
			{Position: 0, Value: fixedValue},
			{Position: length - 1, Value: fixedValue},
		}, nil
	}

	for i := 0; i < len(keyFrames)-1; i++ {
		if keyFrames[i].Position == keyFrames[i+1].Position {
			return nil, errors.New("key frames may not be at the same position")
		}
	}

	lastFrame := keyFrames[len(keyFrames)-1]
	firstFrame := keyFrames[0]

	if firstFrame.Position != 0 {
		frameZeroValue := firstFrame.Value

		if t.looping {
			frameZeroValue = t.lerpFrameValues(lastFrame, firstFrame, 0)
		}

		startKeyFrame := KeyFrame{
			Position: 0,
			Value:    frameZeroValue,
		}

		keyFrames = append([]KeyFrame{startKeyFrame}, keyFrames...)
	}

	if lastFrame.Position != length-1 {
		endFrameValue := lastFrame.Value
		if t.looping {
			endFrameValue = t.lerpFrameValues(lastFrame, keyFrames[0], length-1)
		}

		endKeyFrame := KeyFrame{
			Position: length - 1,
			Value:    endFrameValue,
		}

		keyFrames = append(keyFrames, endKeyFrame)
	}

	return keyFrames, nil
}

func (t *Track) lerpFrameValues(startFrame KeyFrame, endFrame KeyFrame, position int) Color {
	duration := (endFrame.Position - startFrame.Position + t.length) % t.length
	if duration == 0 {
		panic("start and end frame cannot be at the same position.")
	}

	normalizedPosition := (position - startFrame.Position + t.length) % t.length
	progress := float64(normalizedPosition) / float64(duration)

	return lerpColor(startFrame.Value, endFrame.Value, progress, lerpByte)
}
