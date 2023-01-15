package types

import (
	"edemond/abstract/msg"
	"fmt"
)

// A Part is a compiled ast.Expression. We have four types right now:
// - Simple part: The smallest playable unit of music.
// - Compound part: Two or more parts played simultaneously.
// - Block part: Two or more parts played sequentially, each part filling out one meter.
// - Seq part: Two or more parts played sequentially, with all parts condensed into one meter.
type Part interface {
	Value
	Play(buf msg.Buffer, ppq int, step uint64)
	Length(ppq int) uint64
}

// Interpretation is "how to play it", acting in a
// given rhythmic and harmonic context (Rhythm and Harmony).
// TODO: give this the axe; instead build interpretations in the language
// using seq parts and other primitives.
type Interpretation interface {
	Value
	// TODO: This method signature is waxing large
	Play(notes []Note, h *Harmony, r *Rhythm, counter uint64, step uint64, length uint64, ppq int)
}

// Any type of value in the language. Use Go type assertions to figure out what (sorry.)
type Value interface {
	HasValue() bool // All types may or may not have values.
	String() string
}

// MIDI program change message.
type ProgramChange struct {
	Channel uint8
	Program uint8
}

func (pc *ProgramChange) HasValue() bool {
	return pc != nil
}

func (pc *ProgramChange) String() string {
	return fmt.Sprintf("pc(%v, %v)", pc.Channel, pc.Program)
}

// MIDI controller change message.
type ControllerChange struct {
	Channel    uint8
	Controller uint8
	Value      uint8
}

func (cc *ControllerChange) HasValue() bool {
	return cc != nil
}

func (cc *ControllerChange) String() string {
	return fmt.Sprintf("cc(%v, %v, %v)", cc.Channel, cc.Controller, cc.Value)
}
