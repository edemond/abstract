package types

import (
	"edemond/abstract/msg"
)

// A list of parts to play concurrently.
type CompoundPart struct {
	parts  []Part
	length uint64
	scale  int
	// TODO: the combining mode (LCM, polyrhythm, straight-through (one-shot), or limited)
}

func NewCompoundPart() *CompoundPart {
	return &CompoundPart{
		parts:  make([]Part, 0),
		length: 0,
	}
}

func (c *CompoundPart) Add(p Part) {
	c.parts = append(c.parts, p)
}

func (c *CompoundPart) Play(buf msg.Buffer, ppq int, step uint64) {
	for _, part := range c.parts {
		// TODO: this currently only implements a limited looping play. need one-shot play, polymeter, polyrhythm
		s := step % part.Length(ppq)
		part.Play(buf, ppq, s)
	}
}

// Length returns the length of the part in steps. Memoized.
// For compound parts, this is the length of its longest constituent part.
func (c *CompoundPart) Length(ppq int) uint64 {
	if c.length <= 0 {
		for _, part := range c.parts {
			length := part.Length(ppq)
			if length > c.length {
				c.length = length / c.scale
			}
		}
	}
	return c.length
}

func (c *CompoundPart) String() string {
	return "compoundpart()"
}

func (c *CompoundPart) HasValue() bool {
	return c != nil
}
