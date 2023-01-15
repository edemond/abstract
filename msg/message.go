package msg

import (
	"github.com/edemond/midi"
)

// When played, parts emit a series of Messages.
type Message struct {
	// TODO: Decouple this from MIDI.
	MidiMessage      midi.Message // MIDI message to send.
	Instrument       int          // Instrument ID to send the MIDI message to.
	HumanizeTime     int          // TODO: This is currently JACK frames, but ought to be independent.
	HumanizeVelocity int
}
