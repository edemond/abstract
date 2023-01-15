package types

import (
	"fmt"
	"strings"
)

// Pitch is a pitch class.
type Pitch int

const noPitch Pitch = -1

func (p Pitch) HasValue() bool {
	return p != noPitch
}

func NoPitch() Pitch {
	return noPitch
}

func DefaultPitch() Pitch {
	return Pitch(0) // C
}

func NewPitch(num uint64) Pitch {
	return Pitch(num % 12)
}

func (p Pitch) String() string {
	if p == noPitch {
		return "(no pitch)"
	}
	name, ok := pitchNames[p]
	if !ok {
		panic(fmt.Sprintf("Internal error: unaccounted-for pitch in pitchNames: %v", int(p)))
	}
	return name
}

// Pitch combined with Octave yields a specific Note.
func (p Pitch) At(o Octave) Note {
	note, err := NewNote(uint64(o*12) + uint64(p))
	if err != nil {
		panic(fmt.Sprintf("Internal error in Pitch.At: %v", err))
	}
	return note
}

// Add applies an interval, in positive or negative half-steps, to a root Pitch, resulting in a new Pitch.
func (p Pitch) Add(interval int) Pitch {
	return NewPitch(uint64(int(p) + interval))
}

// Calculate the distance in half-steps from Pitch p1 to Pitch p2, constrained to the given Octave.
// Not the absolute value; it could return something negative.
func (p1 Pitch) StepsTo(p2 Pitch, o Octave) int {
	return int(p1.At(o)) - int(p2.At(o))
}

var pitchNames = map[Pitch]string{
	NewPitch(0):  "C",
	NewPitch(1):  "C#/Db",
	NewPitch(2):  "D",
	NewPitch(3):  "D#/Eb",
	NewPitch(4):  "E",
	NewPitch(5):  "F",
	NewPitch(6):  "F#/Gb",
	NewPitch(7):  "G",
	NewPitch(8):  "G#/Ab",
	NewPitch(9):  "A",
	NewPitch(10): "A#/Bb",
	NewPitch(11): "B",
}

// Built in pitch literals.
var pitches = map[string]Pitch{
	"C♭♭": NewPitch(10),
	"Cbb": NewPitch(10),
	"C♭":  NewPitch(11),
	"Cb":  NewPitch(11),
	"C":   NewPitch(0),
	"C♮":  NewPitch(0),
	"C#":  NewPitch(1),
	"C♯":  NewPitch(1),
	"C##": NewPitch(2),
	"C♯♯": NewPitch(2),

	"D♭♭": NewPitch(0),
	"Dbb": NewPitch(0),
	"D♭":  NewPitch(1),
	"Db":  NewPitch(1),
	"D":   NewPitch(2),
	"D♮":  NewPitch(2),
	"D#":  NewPitch(3),
	"D♯":  NewPitch(3),
	"D##": NewPitch(4),
	"D♯♯": NewPitch(4),

	"E♭♭": NewPitch(2),
	"Ebb": NewPitch(2),
	"E♭":  NewPitch(3),
	"Eb":  NewPitch(3),
	"E":   NewPitch(4),
	"E♮":  NewPitch(4),
	"E#":  NewPitch(5),
	"E♯":  NewPitch(5),
	"E##": NewPitch(6),
	"E♯♯": NewPitch(6),

	"F♭♭": NewPitch(3),
	"Fbb": NewPitch(3),
	"F♭":  NewPitch(4),
	"Fb":  NewPitch(4),
	"F":   NewPitch(5),
	"F♮":  NewPitch(5),
	"F#":  NewPitch(6),
	"F♯":  NewPitch(6),
	"F##": NewPitch(7),
	"F♯♯": NewPitch(7),

	"G♭♭": NewPitch(5),
	"Gbb": NewPitch(5),
	"G♭":  NewPitch(6),
	"Gb":  NewPitch(6),
	"G":   NewPitch(7),
	"G♮":  NewPitch(7),
	"G#":  NewPitch(8),
	"G♯":  NewPitch(8),
	"G##": NewPitch(9),
	"G♯♯": NewPitch(9),

	"A♭♭": NewPitch(7),
	"Abb": NewPitch(7),
	"A♭":  NewPitch(8),
	"Ab":  NewPitch(8),
	"A":   NewPitch(9),
	"A♮":  NewPitch(9),
	"A#":  NewPitch(10),
	"A♯":  NewPitch(10),
	"A##": NewPitch(11),
	"A♯♯": NewPitch(11),

	"B♭♭": NewPitch(9),
	"Bbb": NewPitch(9),
	"B♭":  NewPitch(10),
	"Bb":  NewPitch(10),
	"B":   NewPitch(11),
	"B♮":  NewPitch(11),
	"B#":  NewPitch(0),
	"B♯":  NewPitch(0),
	"B##": NewPitch(1),
	"B♯♯": NewPitch(1),
}

func IsPitch(text string) bool {
	_, ok := pitches[text]
	return ok
}

// LookUpPitch converts a pitch notation string (e.g. "D", "C##", "B♭") to a Pitch.
func LookUpPitch(text string) (Pitch, error) {
	pitch, ok := pitches[text]
	if ok {
		return pitch, nil
	}
	return NoPitch(), fmt.Errorf("%v is not a pitch")
}

// strings.Join, but for types.Pitch
func JoinPitches(input []Pitch, sep string) string {
	strs := make([]string, len(input))
	for i, s := range input {
		strs[i] = s.String()
	}
	return strings.Join(strs, sep)
}
