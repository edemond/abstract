package types

import (
	"edemond/abstract/msg"
	"fmt"
)

const BUFFER_SIZE = 128

var SIMPLE_PART_ID = 0

func makeNoteBuffer(size int) []Note {
	buf := make([]Note, size)
	for i := 0; i < size; i++ {
		buf[i] = NoNote()
	}
	return buf
}

type SimplePart struct {
	Harmony        *Harmony
	Rhythm         *Rhythm
	Instrument     *Instrument
	Interpretation Interpretation
	scale          int
	// bookkeeping
	id int
	// playback data
	// TODO: remove this from the parts. playback should be separate
	playing []Note
	length  uint64
	counter uint64
}

// NewSimplePart creates a new blank SimplePart.
func NewSimplePart() *SimplePart {
	SIMPLE_PART_ID += 1
	return &SimplePart{
		Rhythm: &Rhythm{
			Dynamics: NoDynamics(),
			Humanize: NoHumanize(),
			Meter:    NoMeter(),
		},
		Harmony: &Harmony{
			Chord:   NoChord(),
			Octave:  NoOctave(),
			Pitch:   NoPitch(),
			Scale:   NoScale(),
			Voicing: NoVoicing(),
		},
		Instrument:     NoInstrument(),
		Interpretation: NewBlock(), // TODO: Interpretation? lol pls. Voicing + seqs take care of this!
		playing:        makeNoteBuffer(BUFFER_SIZE),
		id:             SIMPLE_PART_ID,
		scale:          1,
	}
}

// Copy makes a memberwise copy of a simple part.
func (s *SimplePart) Copy() *SimplePart {
	SIMPLE_PART_ID += 1
	return &SimplePart{
		Rhythm: &Rhythm{
			Dynamics: s.Rhythm.Dynamics,
			Humanize: s.Rhythm.Humanize,
			Meter:    s.Rhythm.Meter,
		},
		Harmony: &Harmony{
			Chord:   s.Harmony.Chord,
			Octave:  s.Harmony.Octave,
			Pitch:   s.Harmony.Pitch,
			Scale:   s.Harmony.Scale,
			Voicing: s.Harmony.Voicing,
		},
		Instrument:     s.Instrument,
		Interpretation: s.Interpretation, // TODO: Interpretation? lol pls. Voicing + seqs take care of this!
		playing:        makeNoteBuffer(BUFFER_SIZE),
		id:             SIMPLE_PART_ID,
		scale:          s.scale,
	}
}

func (s *SimplePart) Play(buf msg.Buffer, ppq int, step uint64) {

	// TODO: This needs to be an exemplar of the open/closed principle. It's the
	// Harmonic and Rhythmic contexts together (and later Timbral!) that determine
	// what to play. SimplePart.Play's ONLY knowledge should be how to orchestrate
	// this (not in the musical sense.)

	s.Harmony.SetDefaults()
	s.Rhythm.SetDefaults()
	if s.Instrument == nil {
		s.Instrument = defaultInstrument()
	}

	length := s.Length(ppq)

	// Have the Interpretation update the note buffer.
	s.Interpretation.Play(s.playing, s.Harmony, s.Rhythm, s.counter, step, length, ppq)

	// Write all of the buffered notes out to the main message buffer.
	for i := 0; i < len(s.playing); i++ {
		note := s.playing[i]
		if note.HasValue() {
			human := 0
			if s.Rhythm.Humanize.HasValue() {
				human = s.Rhythm.Humanize.TimeOffset()
			}

			var m msg.Message
			m.MidiMessage.Command = 0x9 // note on
			m.MidiMessage.Channel = s.Instrument.Channel
			m.MidiMessage.Data1 = byte(note)
			m.MidiMessage.Data2 = byte(s.Rhythm.Dynamics.Center) // TODO: humanize
			m.Instrument = s.Instrument.ID
			m.HumanizeTime = human
			buf.Add(&m)
			s.playing[i] = NoNote()
		}
	}

	s.counter++
}

func (s *SimplePart) SetScale(scale int) {
	s.scale = scale
}

func (s *SimplePart) Length(ppq int) uint64 {
	// TODO: This has an obvious bug! It memoizes a length based on the first ppq it's given!
	// Ideally, ppq should not change, but.....
	// That's interesting. Maybe we remove the ppq statement from the language and make it
	// a flag that's passed in, with a sane default like 64? I can't think of a good use
	// case for letting ppq vary throughout the piece.
	if s.length == 0 {
		s.Rhythm.SetDefaults()
		// TODO: Aw, this memoization wrecks seq part lengths.
		// buuuuuut it doesn't exactly work without it either

		// TODO: okay, simple part needs to know about its parent now.
		// the length should be the parent length / s.scale, and only if there's
		// no parent, that's where we take meter into account.
		s.length = s.Rhythm.Meter.Length(ppq) / s.scale
	}
	return s.length
}

func (s *SimplePart) String() string {
	return fmt.Sprintf("simplepart(%v, %p)", s.id, s)
}

func (s *SimplePart) HasValue() bool {
	return s != nil
}
