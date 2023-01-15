package types

import (
	"fmt"
	"testing"
)

func expectPulse(r *Rhythm, t *testing.T, step uint64, ppq, expectedBeat, expectedStrength int) {
	beat, strength := r.Pulse(step, ppq)
	if beat != expectedBeat || strength != expectedStrength {
		t.Errorf("got (%v, %v) for step %v, expected (%v, %v)", beat, strength, step, expectedBeat, expectedStrength)
	}
}

// Print out the pulse so we can see the actual pattern.
func dumpPulse(r *Rhythm, ppq int) {
	fmt.Println("actual pulse: ----------------------------------")
	length := r.Meter.Length(ppq)
	for i := uint64(0); i < length; i++ {
		beat, strength := r.Pulse(i, ppq)
		if strength != 0 {
			fmt.Printf("%v: %v, %v\n", i, beat, strength)
		}
	}
}

func Test44Pulse(t *testing.T) {
	r := &Rhythm{
		Meter: &Meter{Beats: 4, Value: 4},
	}
	ppq := 64

	// first quarter note
	expectPulse(r, t, 0, ppq, 1, 1)
	expectPulse(r, t, 1, ppq, 1, 0)
	expectPulse(r, t, 2, ppq, 1, 0)
	expectPulse(r, t, 7, ppq, 1, 0)

	// we come to a 32nd note first...
	expectPulse(r, t, 8, ppq, 1, 32)
	expectPulse(r, t, 9, ppq, 1, 0)
	expectPulse(r, t, 15, ppq, 1, 0)

	// sixteenth...
	expectPulse(r, t, 16, ppq, 1, 16)
	expectPulse(r, t, 23, ppq, 1, 0)

	// another 32nd note...
	expectPulse(r, t, 24, ppq, 1, 32)
	expectPulse(r, t, 25, ppq, 1, 0)

	// an 8th note...
	expectPulse(r, t, 32, ppq, 1, 8)

	// second quarter note
	expectPulse(r, t, 64, ppq, 2, 4)

	// third quarter note
	expectPulse(r, t, 128, ppq, 3, 2)

	// fourth quarter note
	expectPulse(r, t, 192, ppq, 4, 4)
	expectPulse(r, t, 193, ppq, 4, 0)

	expectPulse(r, t, 200, ppq, 4, 32)

	if t.Failed() {
		dumpPulse(r, ppq)
	}
}

func Test34Pulse(t *testing.T) {
	r := &Rhythm{
		Meter: &Meter{Beats: 3, Value: 4},
	}
	ppq := 4
	expectPulse(r, t, 0, ppq, 1, 1)
	expectPulse(r, t, 1, ppq, 1, 16)
	expectPulse(r, t, 2, ppq, 1, 8)
	expectPulse(r, t, 3, ppq, 1, 16)

	expectPulse(r, t, 4, ppq, 2, 4)
	expectPulse(r, t, 5, ppq, 2, 16)
	expectPulse(r, t, 6, ppq, 2, 8)
	expectPulse(r, t, 7, ppq, 2, 16)

	expectPulse(r, t, 8, ppq, 3, 4)
	expectPulse(r, t, 9, ppq, 3, 16)
	expectPulse(r, t, 10, ppq, 3, 8)
	expectPulse(r, t, 11, ppq, 3, 16)

	ppq = 64
	// first quarter note
	expectPulse(r, t, 0, ppq, 1, 1)
	expectPulse(r, t, 1, ppq, 1, 0)
	expectPulse(r, t, 2, ppq, 1, 0)
	expectPulse(r, t, 7, ppq, 1, 0)

	expectPulse(r, t, 8, ppq, 1, 32)
	expectPulse(r, t, 9, ppq, 1, 0)
	expectPulse(r, t, 15, ppq, 1, 0)

	expectPulse(r, t, 16, ppq, 1, 16)
	expectPulse(r, t, 23, ppq, 1, 0)

	expectPulse(r, t, 24, ppq, 1, 32)
	expectPulse(r, t, 25, ppq, 1, 0)

	expectPulse(r, t, 32, ppq, 1, 8)

	// second quarter note
	expectPulse(r, t, 64, ppq, 2, 4)

	// third quarter note
	expectPulse(r, t, 128, ppq, 3, 4) // this is the difference; there's no half note

	if t.Failed() {
		dumpPulse(r, ppq)
	}
}

func Test68Pulse(t *testing.T) {
	r := &Rhythm{
		Meter: &Meter{Beats: 6, Value: 8},
	}
	ppq := 8
	expectPulse(r, t, 0, ppq, 1, 1)
	expectPulse(r, t, 1, ppq, 1, 16)
	expectPulse(r, t, 2, ppq, 1, 8)
	expectPulse(r, t, 4, ppq, 2, 4)
	expectPulse(r, t, 8, ppq, 3, 4)
	expectPulse(r, t, 12, ppq, 4, 2)
	expectPulse(r, t, 16, ppq, 5, 4)
	expectPulse(r, t, 20, ppq, 6, 4)

	if t.Failed() {
		dumpPulse(r, ppq)
	}
}

func Test54Pulse(t *testing.T) {
	r := &Rhythm{
		Meter: &Meter{Beats: 5, Value: 4},
	}
	ppq := 4
	expectPulse(r, t, 0, ppq, 1, 1)
	expectPulse(r, t, 1, ppq, 1, 16)
	expectPulse(r, t, 2, ppq, 1, 8)
	expectPulse(r, t, 3, ppq, 1, 16)

	expectPulse(r, t, 4, ppq, 2, 4)
	expectPulse(r, t, 5, ppq, 2, 16)
	expectPulse(r, t, 6, ppq, 2, 8)
	expectPulse(r, t, 7, ppq, 2, 16)

	expectPulse(r, t, 8, ppq, 3, 4)
	expectPulse(r, t, 9, ppq, 3, 16)
	expectPulse(r, t, 10, ppq, 3, 8)
	expectPulse(r, t, 11, ppq, 3, 16)

	expectPulse(r, t, 12, ppq, 4, 4)
	expectPulse(r, t, 13, ppq, 4, 16)
	expectPulse(r, t, 14, ppq, 4, 8)
	expectPulse(r, t, 15, ppq, 4, 16)

	expectPulse(r, t, 16, ppq, 5, 4)
	expectPulse(r, t, 17, ppq, 5, 16)
	expectPulse(r, t, 18, ppq, 5, 8)
	expectPulse(r, t, 19, ppq, 5, 16)
}
