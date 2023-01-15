package types

import (
	"edemond/abstract/util"
	"fmt"
)

// Chord represents a set of intervals which can be applied to an unspecified pitch.
// There are several ways to specify a chord.
type Chord interface {
	Value
	Root() Pitch // Returns a pitch with .HasValue() == false if not an absolute chord.
	ResolveIn(key Pitch, scale *Scale) []Pitch
	Play(notesOut []Note, h *Harmony) // Used by Interpretation. Sound a block chord.
}

// Neither major/minor qualities nor root pitch is specified, just a set of scale degrees.
// To play a diatonicChord, you have to specify what scale to play it in.
// This is used to transform musical material to, for example, its parallel major or minor, or different modes.
// e.g. @III in C major produces E minor.
type diatonicChord struct {
	// TODO: How to indicate a diatonic 7th chord and beyond?
	// Triad vs. extended chords...should be handled in the parser.
	scaleDegrees []int // e.g. [1,3,5] for a triad on the root.
}

// relativeChord is a collection of intervals (in half-steps) from some unspecified root pitch.
// i.e. Major/minor quality is specified, but a relative root pitch instead of an absolute one.
// e.g. iiimin7.
type relativeChord struct {
	intervalsInHalfSteps []int // These can be offset by an accidental (even negative if the chord is on a flattened scale degree.)
	rootScaleDegree      int   // Typically 1-7, but could be higher for weird octatonic scales, etc.
}

// absoluteChord is a collection of pitches.
// e.g. emin7.
type absoluteChord struct {
	pitches    []Pitch
	avoidNotes []Pitch // TODO: Notes which are traditionally avoided. (This isn't used yet, just an idea.)
}

// Diatonic chords don't have a specified root pitch.
func (c *diatonicChord) Root() Pitch {
	return NoPitch()
}

// Relative chords don't have a specified root pitch.
func (c *relativeChord) Root() Pitch {
	return NoPitch()
}

func (c *absoluteChord) Root() Pitch {
	return c.pitches[0] // Treat the first note specified as the root.
}

// Create a chord from a set of absolute pitches, e.g. chord(C, Eb, G, Bb)
func NewAbsoluteChordFromPitches(pitches []Pitch) Chord {
	c := &absoluteChord{}
	dup := make(map[Pitch]bool)
	c.pitches = make([]Pitch, len(pitches))
	for index, pitch := range pitches {
		_, ok := dup[pitch]
		if !ok {
			c.pitches[index] = pitch
			dup[pitch] = true
		}
	}
	return c
}

// Create a chord from a root pitch plus a set of intervals in halfsteps.
// Used for resolving chord notation like Cmin7.
func NewAbsoluteChord(root Pitch, intervalsInHalfSteps []int) Chord {
	c := &absoluteChord{}
	dup := make(map[Pitch]bool)
	c.pitches = make([]Pitch, len(intervalsInHalfSteps))
	for index, interval := range intervalsInHalfSteps {
		pitch := root.Add(interval)
		_, ok := dup[pitch]
		if !ok {
			c.pitches[index] = pitch
			dup[pitch] = true
		}
	}
	return c
}

// Create a chord from a root scale degree plus a set of intervals in halfsteps.
// Used for resolving chord notation like IVmin7.
func NewRelativeChord(rootScaleDegree int, intervalsInHalfSteps []int) Chord {
	c := &relativeChord{rootScaleDegree: rootScaleDegree}
	c.intervalsInHalfSteps = make([]int, len(intervalsInHalfSteps))
	for index, interval := range intervalsInHalfSteps {
		c.intervalsInHalfSteps[index] = interval
	}
	return c
}

// Create a chord from set of scale degrees.
// Used for resolving chord notation like IVmin7.
func NewDiatonicChord(scaleDegrees []int) Chord {
	c := &diatonicChord{}
	c.scaleDegrees = make([]int, len(scaleDegrees))
	for index, interval := range scaleDegrees {
		c.scaleDegrees[index] = interval
	}
	return c
}

// NoChord creates a null chord.
func NoChord() *absoluteChord {
	return nil
}

func DefaultChord() Chord {
	return NoChord()
}

// In calculates a relative chord for a given root note, resulting in a set of Pitches.
// e.g. chord.In(C)
func (c *relativeChord) in(root Pitch) []Pitch {
	pitches := make([]Pitch, 0)
	for _, halfSteps := range c.intervalsInHalfSteps {
		pitches = append(pitches, root.Add(halfSteps))
	}
	return pitches
}

func (c *absoluteChord) String() string {
	return fmt.Sprintf("absolute chord(%v)", JoinPitches(c.pitches, ", "))
}

func (c *diatonicChord) String() string {
	return fmt.Sprintf("diatonic chord(%v)", util.JoinInts(c.scaleDegrees, ", "))
}

func (c *relativeChord) String() string {
	return fmt.Sprintf("relative chord(%v) on %v", util.JoinInts(c.intervalsInHalfSteps, ", "), c.rootScaleDegree)
}

func (c *absoluteChord) HasValue() bool {
	return c != nil
}

func (c *diatonicChord) HasValue() bool {
	return c != nil
}

func (c *relativeChord) HasValue() bool {
	return c != nil
}

func (c *absoluteChord) ResolveIn(key Pitch, scale *Scale) []Pitch {
	pitches := make([]Pitch, len(c.pitches))
	for i, p := range c.pitches {
		pitches[i] = p
	}
	return pitches
}

func (c *relativeChord) ResolveIn(key Pitch, scale *Scale) []Pitch {
	// Get the chord's root degree in the scale in the given key, that's a pitch
	rootPitch := key.Add(scale.StepsAtDegree(c.rootScaleDegree - 1))
	pitches := make([]Pitch, len(c.intervalsInHalfSteps))
	for index, interval := range c.intervalsInHalfSteps {
		pitches[index] = rootPitch.Add(interval)
	}
	return pitches
}

func (c *diatonicChord) ResolveIn(key Pitch, scale *Scale) []Pitch {
	pitches := make([]Pitch, len(c.scaleDegrees))
	for index, degree := range c.scaleDegrees {
		pitches[index] = key.Add(scale.StepsAtDegree(degree - 1))
	}
	return pitches
}

// TODO: The Play methods all have to be refactored into methods that get chordal info.
// Chord should no longer be an Interpretation, but a source of harmonic information for Interpretations.

// TODO: Should Chord really know about Voicing? I don't think so; that ought to be taken care of
// somewhere else, right?

func (c *absoluteChord) Play(notesOut []Note, h *Harmony) {
	if h.Voicing.HasValue() {
		voiced := h.Voicing.Voice(c.pitches)
		for i, note := range voiced {
			notesOut[i] = note
		}
	} else {
		// Default voicing.....hmmmmmmmm. Use the pitch and octave and just play it straight.
		for i, pitch := range c.pitches {
			notesOut[i] = pitch.At(h.Octave)
		}
	}
}

func (c *relativeChord) Play(notesOut []Note, h *Harmony) {
	// TODO: ARGH! more allocation during playback. can we allocate one buffer per thread or something?
	// or just make Voicing handle the conversion in its loop?
	pitches := c.ResolveIn(h.Pitch, h.Scale)

	if h.Voicing.HasValue() {
		// The voicing takes care of the octave of each pitch, so ignore context.Octave.
		// TODO: until...octave-relative voicings? someday? that oughtta be possible/useful.
		voiced := h.Voicing.Voice(pitches)
		for i, note := range voiced {
			notesOut[i] = note
		}
	} else {
		// Default voicing.....hmmmmmmmm. Use the pitch, octave, and scale, and just play it straight.
		// This implementation depends on the fact we have a default scale.
		for i, pitch := range pitches {
			// TODO: dedupe notes
			notesOut[i] = pitch.At(h.Octave)
		}
	}
}

func (c *diatonicChord) Play(notesOut []Note, h *Harmony) {
	pitches := c.ResolveIn(h.Pitch, h.Scale)

	if h.Voicing.HasValue() {
		// The voicing takes care of the octave of each pitch, so ignore context.Octave.
		// TODO: until...octave-relative voicings? someday? that oughtta be possible/useful.
		voiced := h.Voicing.Voice(pitches)
		for i, note := range voiced {
			notesOut[i] = note
		}
	} else {
		// Default voicing.....hmmmmmmmm. Use the pitch, octave, and scale, and just play it straight.
		// This implementation depends on the fact we have a default scale.
		for i, pitch := range pitches {
			// TODO: dedupe notes
			notesOut[i] = pitch.At(h.Octave)
		}
	}
}
