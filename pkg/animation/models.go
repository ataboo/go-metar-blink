package animation

import (
	"errors"
	"math"
	"sort"
)

type KeyFrame struct {
	Position int
	Value    uint64
}

type Track struct {
	Looping    bool
	ChannelIDs []int
	position   int
	length     int
	keyFrames  []KeyFrame
}

type Animation struct {
	Tracks []Track
}

func (t *Track) Step() bool {
	if t.position >= t.length-1 {
		if !t.Looping {
			return false
		}

		t.position = 0
	} else {
		t.position++
	}

	return true
}

func (t *Track) GetPosition() int {
	return t.position
}

func (t *Track) Seek(position int) error {
	if position < 0 || position >= t.length {
		return errors.New("position out of range")
	}
	t.position = position

	return nil
}

func (t *Track) SetLength(length int) error {
	if length < 0 {
		return errors.New("length out of range")
	}

	t.length = length
	if t.position > length-1 {
		t.position = length - 1
	}

	return nil
}

func (t *Track) GetLength() int {
	return t.length
}

func (t *Track) SetKeyFrames(keyFrames []KeyFrame, length int) error {
	err := t.SetLength(length)
	if err != nil {
		return err
	}

	sort.Slice(keyFrames, func(i int, j int) bool {
		return keyFrames[i].Position < keyFrames[j].Position
	})

	t.keyFrames = t.normalizeKeyFrames(keyFrames, length)

	return nil
}

func (t *Track) normalizeKeyFrames(keyFrames []KeyFrame, length int) []KeyFrame {
	if len(keyFrames) < 2 {
		fixedValue := uint64(0)
		if len(keyFrames) == 1 {
			fixedValue = keyFrames[0].Value
		}

		return []KeyFrame{
			{Position: 0, Value: fixedValue},
			{Position: length - 1, Value: fixedValue},
		}
	}

	lastFrame := keyFrames[len(keyFrames)-1]
	firstFrame := keyFrames[0]

	if firstFrame.Position != 0 {
		startKeyFrame := KeyFrame{
			Position: 0,
			Value:    t.lerpFrameValues(lastFrame, firstFrame, 0),
		}

		keyFrames = append([]KeyFrame{startKeyFrame}, keyFrames...)
	}

	if lastFrame.Position != length-1 {
		endKeyFrame := KeyFrame{
			Position: length - 1,
			Value:    t.lerpFrameValues(lastFrame, keyFrames[0], length-1),
		}

		keyFrames = append(keyFrames, endKeyFrame)
	}

	return keyFrames
}

func (t *Track) Value() uint64 {
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

func (t *Track) lerpFrameValues(startFrame KeyFrame, endFrame KeyFrame, position int) uint64 {
	duration := (endFrame.Position - startFrame.Position + t.length) % t.length
	normalizedPosition := (position - startFrame.Position + t.length) % t.length
	progress := float64(normalizedPosition) / float64(duration)

	value := math.Round(float64(startFrame.Value)*(1-progress) + float64(endFrame.Value)*progress)

	return uint64(value)
}
