package types

import (
	"fmt"
	"math/rand"
)

// Prob is a probabalistic rhythm/expression. The chance of playing a note is derived from the pulse.
type Prob struct {
	// TODO: An "invert" option for syncopation.
	beat     int
	strength int
	percent  int // Probability of playing a note, 0-100.
}

func NewProb(beat, strength, percent int) *Prob {
	return &Prob{
		beat:     beat,
		strength: strength,
		percent:  percent,
	}
}

func (p *Prob) Play(notesOut []Note, h *Harmony, r *Rhythm, counter uint64, step uint64, length uint64, ppq int) {
	beat, strength := r.Pulse(step, ppq)
	if p.beat == 0 || beat == p.beat {
		if p.strength == 0 || strength == p.strength {
			chance := rand.Intn(100)
			//fmt.Printf("chance: %v\n", chance)
			if chance <= p.percent {
				if h.Chord.HasValue() {
					/*
					   fmt.Printf("prob.Play, beat: %v, str: %v\n", beat, strength)
					   fmt.Printf("p.beat: %v, p.strength: %v\n", p.beat, p.strength)
					*/
					h.Chord.Play(notesOut, h) // Block chord will sound.
				} else {
					notesOut[0] = h.Pitch.At(h.Octave) // Just a single note.
				}
			}
		}
	}
}

func (p *Prob) String() string {
	return fmt.Sprintf("prob(%v, %v, %v)", p.beat, p.strength, p.percent)
}

func (p *Prob) HasValue() bool {
	return p != nil
}
