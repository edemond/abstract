package types

import (
	"fmt"
)

// Instrument defines a instrument for output (ALSA devices, JACK output ports, etc.)
// Its "Name" field is essentially a connection string that a driver knows how to interpret.
type Instrument struct {
	ID      int    // Index of the instrument in the list of open ones.
	Name    string // Name of the synth (meaningful to the driver for opening the instrument.)
	Channel byte   // MIDI channel
	Voices  int    // number of polyphonic voices
	// TODO: note stealing modes! (none, steal first, steal random, etc...)
}

func (i *Instrument) String() string {
	return fmt.Sprintf("instrument(\"%v\", channel %v, %v voices)", i.Name, i.Channel, i.Voices)
}

func (i *Instrument) HasValue() bool {
	return i != nil
}

func NoInstrument() *Instrument {
	return nil
}

func NewInstrument(name string, channel byte, voices int) *Instrument {
	return &Instrument{
		Channel: byte(channel),
		Voices:  int(voices),
		Name:    name,
	}
}

func defaultInstrument() *Instrument {
	return &Instrument{
		Channel: 1,
		Voices:  0,
		Name:    "(no instrument)",
	}
}

func TotalVoices(insts []*Instrument) int {
	voices := 0
	for _, i := range insts {
		voices += i.Voices
	}
	return voices
}
