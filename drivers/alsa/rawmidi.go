// rawmidi.go implements an Abstract driver for the ALSA "rawmidi" API.
// It keeps time itself, using a Go time.Ticker.
package alsa

import (
	"github.com/edemond/abstract/drivers"
	"github.com/edemond/abstract/msg"
	"github.com/edemond/abstract/types"
	"github.com/edemond/midi"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

var instrumentID int = 0

type rawMidiDriver struct {
	devices     []midi.Device       // List of ALSA rawmidi devices we know about.
	openDevices map[int]midi.Device // Instrument ID -> open output device. The "instruments", basically.
}

func NewRawMidiDriver() (drivers.Driver, error) {
	devices, err := midi.GetDevices("alsa")
	if err != nil {
		return nil, err
	}

	return &rawMidiDriver{
		devices:     devices,
		openDevices: make(map[int]midi.Device),
	}, nil
}

func timer(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}

func swingTimer(d1 time.Duration, d2 time.Duration, cond *sync.Cond) {
	ticker1 := time.NewTicker(d1)
	ticker2 := time.NewTicker(d1 + d2)
	for {
		<-ticker1.C
		cond.Broadcast()
		<-ticker2.C
		cond.Broadcast()
	}
}

// Returns a time.Duration representing how long we should wait for each tick at the given BPM.
func bpmToDuration(bpm int, ticksPerBeat int) time.Duration {
	oneBeat := (time.Minute / time.Duration(bpm)) // duration of one beat
	return oneBeat / time.Duration(ticksPerBeat)  // duration of one tick
}

func stopAll(insts map[int]midi.Device) {
	// Send Controller Change 123 ("All notes off")
	// TODO: Also send All Sound Off (120)?
	for _, device := range insts {
		for i := byte(1); i <= 16; i++ {
			device.ControllerChange(i, 123, 0) // TODO: channel 0 ok?! or all channels?
		}
	}
}

func playNote(m *msg.Message, device midi.Device) {
	note := m.MidiMessage
	switch note.Command {
	case 0x8:
		device.NoteOff(note.Channel, note.Data1, note.Data2)
	case 0x9:
		device.NoteOn(note.Channel, note.Data1, note.Data2)
	case 0xC:
		device.ControllerChange(note.Channel, note.Data1, note.Data2)
	}
}

func (r *rawMidiDriver) OpenInstrument(name string) (int, error) {
	// TODO: this could be a map lookup, whatever
	for _, device := range r.devices {
		if device.Name() == name {
			// TODO: Right now we only support output devices.
			if !device.IsOutput() {
				return 0, fmt.Errorf("ALSA rawmidi device '%v' is an input device and cannot be used for output.", name)
			}
			err := device.OpenOutput()
			if err != nil {
				return 0, err
			}
			fmt.Printf("Opened ALSA rawmidi device '%v' for output.\n", name)
			id := instrumentID
			r.openDevices[id] = device

			instrumentID += 1
			return id, nil
		}
	}
	return 0, fmt.Errorf("ALSA rawmidi device '%v' not found.", name)
}

func (r *rawMidiDriver) CloseInstrument(id int) error {
	inst, ok := r.openDevices[id]
	if !ok {
		panic(fmt.Sprintf("Internal error: Couldn't close instrument '%v': no instrument found", id))
	}
	err := inst.Close()
	if err != nil {
		return err
	}
	delete(r.openDevices, id)
	return nil
}

func (r *rawMidiDriver) Close() error {
	return nil
}

func (r *rawMidiDriver) Play(part types.Part, bpm int, ppq int, loop bool, polyphony int) error {
	// Use a buffer big enough to handle all the notes that might be playing at once.
	buf, err := msg.NewBuffer(polyphony)
	if err != nil {
		return err
	}

	// Handle SIGINT and SIGKILL so we can cut off any notes that are still ringing.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	ticker := timer(bpmToDuration(bpm, ppq))
	defer stopAll(r.openDevices)

	var steps uint64
	for {
		length := part.Length(ppq)
		for step := uint64(0); step < length; step++ {
			part.Play(buf, ppq, step)
			select {
			case <-ticker.C:
				if buf.Any() {
					last := buf.Last()
					for i := 0; i < buf.LastLength(); i++ {
						m := last[i]
						note := m.MidiMessage
						device := r.openDevices[m.Instrument]
						device.NoteOff(note.Channel, note.Data1, note.Data2)
					}
					next := buf.Next()
					for i := 0; i < buf.NextLength(); i++ {
						m := next[i]
						playNote(m, r.openDevices[m.Instrument])
					}
				}
				buf.Flip()
			case <-signals:
				return nil
			}
		}
		if !loop {
			return nil
		}
		fmt.Println("Looping.")
		steps += part.Length(ppq)
	}

	return nil
}
