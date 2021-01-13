package animation

import (
	"testing"
	"time"
)

func TestPulseAnimationCosine(t *testing.T) {
	values := make(map[int]Color, 0)
	pulse := CreatePulseAnimation(time.Second*10, 0, ColorWhite, []int{0, 2, 4}).(*PulseAnimation)

	pulse.GetValues(values)
	if len(values) != 3 {
		t.Error("unexpected value count")
	}

	assertMapValuesMatch(values, 0, t)

	pulse.Start()
	pulse.Update(time.Second*1, values)

	if pulse.position != time.Second {
		t.Error("unexpected position", pulse.position)
	}

	assertMapValuesMatch(values, 0x181818, t)

	pulse.Update(time.Second*2, values)
	assertMapValuesMatch(values, 0xa7a7a7, t)
	if pulse.position != time.Second*3 {
		t.Error("unexpected position")
	}

	pulse.Update(time.Second*2, values)
	assertMapValuesMatch(values, 0xffffff, t)

	pulse.Update(time.Second*2, values)
	assertMapValuesMatch(values, 0xa7a7a7, t)

	pulse.Update(time.Second*2, values)
	assertMapValuesMatch(values, 0x181818, t)

	pulse.Update(time.Second, values)
	assertMapValuesMatch(values, 0, t)

	if pulse.position != 0 {
		t.Error("unexpected position")
	}

	pulse.Update(time.Second, values)
	if pulse.position != time.Second {
		t.Error("unexpected position")
	}

	pulse.Stop()

	pulse.Update(time.Second, values)
	assertMapValuesMatch(values, 0x181818, t)

	if pulse.position != time.Second {
		t.Error("unexpected position")
	}
}

func TestPulseAnimationStartStop(t *testing.T) {
	values := make(map[int]Color, 0)
	pulse := CreatePulseAnimation(time.Second*10, 0, 100, []int{0, 2, 4}).(*PulseAnimation)

	pulse.Update(time.Second*1, values)
	if pulse.position != 0 {
		t.Error("unexpected position", pulse.position)
	}

	pulse.Start()
	pulse.Update(time.Second*1, values)
	if pulse.position != time.Second {
		t.Error("unexpected position", pulse.position)
	}

	pulse.Stop()
	pulse.Update(time.Second*1, values)
	if pulse.position != time.Second {
		t.Error("unexpected position", pulse.position)
	}

	pulse.Reset()
	if pulse.position != 0 {
		t.Error("unexpected position", pulse.position)
	}
}

func assertMapValuesMatch(values map[int]Color, expected Color, t *testing.T) {
	if len(values) != 3 {
		t.Error("unexpected value count", len(values))
	}

	if values[0] != expected || values[2] != expected || values[4] != expected {
		t.Errorf("unexpected values R:%s, G:%s, B:%s, expected: %s", values[0].String(), values[2].String(), values[4].String(), expected.String())
	}
}
