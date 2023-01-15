package types

import (
	"fmt"
)

// how do we want to represent this?
// should we have a type that's relative to an octave??
// 		well that's a hint right there...maybe it consists of an octave and a map of notes to a list of offsets
// wait is it even possible to talk about a voicing without a chord? or at least a number of notes??
// cause um. you need to know how many notes there are in order to map them to the octaves they're active in
// AND...should it be sensitive to which notes are just in there for color? like cminadd2, what if you get like eight Ds?
// should it be probabalistic?
// maybe we need to get voicing behind an interface first of all
type Voicing interface {
	Value
	//Voice(Pitch, *Chord) []Note
	Voice([]Pitch) []Note
	IsSounding(octave uint, note int) bool
}

// Each byte is an octave; each bit is if that note in that octave.
// Thus, supports up to 8-pitch chords, voiced in 8 octaves (64-bit int.)
type BitmapVoicing uint64

const noVoicing BitmapVoicing = 0 // TODO: might bite us later. do we always want there to be a voicing?

func (v BitmapVoicing) String() string {
	return fmt.Sprintf("voicing(%v)", uint64(v))
}

func (v BitmapVoicing) HasValue() bool {
	return v != noVoicing
}

func NewVoicing(num uint64) Voicing {
	return BitmapVoicing(num)
}

func NoVoicing() Voicing {
	return noVoicing
}

func DefaultVoicing() Voicing {
	return NoVoicing()
}

// IsSounding checks to see if the given note (a numeric index into a Chord) is sounding in the given Octave in this Voicing.
func (v BitmapVoicing) IsSounding(octave uint, note int) bool {
	if note < 0 {
		panic("note can't be negative")
	}
	return (uint64(v) & ((1 << uint(note)) << (octave * 8))) > 0
}

// Voice applies a Voicing to a set of pitches to produce a concrete set of notes.
//func (v BitmapVoicing) Voice(key Pitch, chord *Chord) []Note {
func (v BitmapVoicing) Voice(pitches []Pitch) []Note {
	notes := make([]Note, 0) // TODO: Agh. Allocation during playback. Fix soon.
	//pitches := chord.In(key)

	var octave uint
	for i, p := range pitches {
		for octave = 0; octave < 8; octave++ {
			// Is pitch i sounding in octave?
			if v.IsSounding(octave, i) {
				notes = append(notes, p.At(NewOctave(octave)))
			}
		}
	}

	return notes
}
