package types

import (
	"fmt"
)

type Octave uint

const noOctave Octave = 9001 // it's over 9000

func (o Octave) String() string {
	return fmt.Sprintf("octave(%v)", uint(o))
}

func (o Octave) HasValue() bool {
	return o != noOctave
}

func NoOctave() Octave {
	return noOctave
}

func DefaultOctave() Octave {
	return 4
}

func NewOctave(num uint) Octave {
	return Octave(num)
}

// IsOctave determines if a string is an octave expression.
func IsOctave(text string) bool {
	if len(text) != 2 {
		return false
	}
	if text[0] != 'O' {
		return false
	}
	last := text[1]
	return last >= '0' && last <= '9'
}
