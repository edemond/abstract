package types

import (
	"edemond/abstract/msg"
)

// A Part that sends a message (e.g. MIDI, OSC) immediately.
type MessagePart interface {
	Part
}

type MIDIMessagePart struct {
	// TODO: Make this support more than one message?
	// TODO: Is this struct redundant?
	Command byte
	Channel byte
	Data1   byte
	Data2   byte
}

func (p *MIDIMessagePart) Play(buf msg.Buffer, ppq int, step uint64) {
	var m msg.Message
	m.MidiMessage.Command = p.Command
	m.MidiMessage.Channel = p.Channel
	m.MidiMessage.Data1 = p.Data1
	m.MidiMessage.Data2 = p.Data2
	//msg.Instrument = s.Instrument.ID // TODO: Do we need this?! What is this used for?
	buf.Add(&m)
}

func (p *MIDIMessagePart) HasValue() bool {
	return p != nil
}
