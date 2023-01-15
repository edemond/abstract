package types

import (
	"fmt"
)

type Note byte

const noNote = 128

func (n Note) String() string {
	return fmt.Sprintf("note(%v)", byte(n))
}

func (n Note) HasValue() bool {
	return n != noNote
}

func IsValidNote(num uint64) bool {
	return num <= 127 && num >= 0
}

func NoNote() Note {
	return noNote
}

func NewNote(num uint64) (Note, error) {
	if !IsValidNote(num) {
		return NoNote(), fmt.Errorf("note must be from 0-127")
	}
	return Note(num), nil
}

func (n Note) Adjust(offset int) Note {
	result := int(n) + offset
	note, err := NewNote(uint64(result % 127))
	if err != nil {
		panic(fmt.Sprintf("Internal error in Note.Adjust: %v", err))
	}
	return note
}
