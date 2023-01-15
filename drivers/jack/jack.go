// Package jack defines an Abstract driver for the JACK Audio Connection Kit (http://jackaudio.org).
package jack

// #cgo LDFLAGS: -ljack
// #include <stdlib.h>
// #include "jack.h"
/*
// jack_client_open is variadic for some reason. Go can't directly call those.
static jack_client_t* open_jack_client(const char* name) {
    // TODO: If this fails, check the jack_status_t. There's more info to be had.
    return jack_client_open(name, JackNullOption, NULL);
}
static jack_port_t* open_jack_port(jack_client_t* client, const char* name) {
    return jack_port_register(
        client,
        name,
        JACK_DEFAULT_MIDI_TYPE,
        JackPortIsOutput | JackPortIsTerminal,
        0 // buffer_size, ignored since we're using JACK_DEFAULT_MIDI_TYPE
    );
}
*/
import "C"

import (
	"github.com/edemond/abstract/drivers"
	"github.com/edemond/abstract/msg"
	"github.com/edemond/abstract/types"
	"fmt"
	"unsafe"
)

const (
	// If you change these, change the corresponding values in jack.h!
	JACK_OK                  = 0 // No error.
	JACK_SET_CALLBACK_FAILED = 1 // Failed to set the process callback.
	JACK_ACTIVATE_FAILED     = 2 // Failed to activate JACK client.
	JACK_DEACTIVATE_FAILED   = 3
)

type jackDriver struct {
	client *C.jack_client_t
	ports  map[int]*C.jack_port_t // Instrument ID -> jack_port_t*

	// This is a bit ugly, but during playback, these fields are set so the callback
	// can get at them:
	buffers map[int]unsafe.Pointer // Instrument ID -> void* (output port buffer)
	ppq     int
	part    types.Part
	buf     msg.Buffer // Main note buffer that the piece's Parts dump notes into.
	length  uint64     // Total length of the piece in steps.
	loop    bool       // Whether or not to loop.
}

// Unique Instrument ID to be incremented each time we assign one.
var instrumentID int

// Prefer to store this in a global during playback so we don't have to complicate
// the interface between Go and C.
var _driver *jackDriver

// Create and initialize the JACK driver.
func NewJACKDriver() (drivers.Driver, error) {
	name := C.CString("abstract")
	defer C.free(unsafe.Pointer(name))

	client := C.open_jack_client(name)
	if client == nil {
		return nil, fmt.Errorf("Couldn't open JACK client.")
	}

	return &jackDriver{
		client:  client,
		ports:   make(map[int]*C.jack_port_t),
		buffers: make(map[int]unsafe.Pointer),
	}, nil
}

// Get the humanized time offset, clamped to within the number of frames.
func calculateHumanizedOffset(offset int, humanize int, nframes C.jack_nframes_t) int {
	//fmt.Printf("humanizing: %v %v %v -> ", offset, humanize, nframes)
	t := offset + humanize
	if t < 0 {
		return 0
	} else if t >= int(nframes) {
		return int(nframes) - 1
	}
	//fmt.Printf("%v\n", t)
	return t
}

func (j *jackDriver) donePlaying(step uint64) bool {
	return !_driver.loop && (step >= _driver.length)
}

// getErrorMessage translates a numeric error code (returned from the driver C code) into an error message.
func getErrorMessage(errno int) string {
	switch errno {
	case JACK_OK:
		return "Success."
	case JACK_ACTIVATE_FAILED:
		return "Failed to activate JACK client."
	case JACK_DEACTIVATE_FAILED:
		return "Failed to deactivate JACK client."
	case JACK_SET_CALLBACK_FAILED:
		return "Failed to set JACK process callback."
	}
	return fmt.Sprintf("Internal error: Unknown JACK driver error code: %v", errno)
}

//export PrepareBuffers
func PrepareBuffers(nframes C.jack_nframes_t) {
	// Pre-fetch and clear all the port buffers.
	for id, port := range _driver.ports {
		portBuffer := C.jack_port_get_buffer(port, nframes)
		_driver.buffers[id] = portBuffer
		C.jack_midi_clear_buffer(portBuffer)
	}
}

// Returns 1 to keep playing, 0 for done.
//export StepSong
func StepSong(step uint64, offset int, nframes C.jack_nframes_t) int {

	// Signal that we're done if we're past length and we're not looping.
	if _driver.donePlaying(step) {
		return 0
	}

	_driver.part.Play(_driver.buf, _driver.ppq, step) // Fill buf with notes to process.

	if _driver.buf.Any() {

		noteOffOffset := offset
		last := _driver.buf.Last()
		next := _driver.buf.Next()

		// JACK doesn't support writing events out of order, and we might have scrambled
		// them with the "humanize" feature, so we need to sort them. This shouldn't
		// take long because hopefully we're not sounding THAT many notes per step.
		if _driver.buf.NextLength() > 0 {
			_driver.buf.Sort()

			// The note offs all need to go before or at the same time as the earliest note on,
			// which could be slightly randomized (humanized) offset-wise. Now that we've just
			// sorted these, it's the first offset.
			noteOffOffset = calculateHumanizedOffset(offset, next[0].HumanizeTime, nframes)
		}

		// TODO: A human player would take a bit of extra time between lifting off the last
		// note and hitting the next note! This could be an opportunity to inject more feel.

		for i := 0; i < _driver.buf.LastLength(); i++ {
			m := last[i]
			buffer := _driver.buffers[m.Instrument]
			note := m.MidiMessage
			C.write_midi_event(
				buffer,
				C.int(noteOffOffset),
				C.uchar(0x8),
				C.uchar(note.Channel),
				C.uchar(note.Data1),
				C.uchar(note.Data2),
			)
		}

		for i := 0; i < _driver.buf.NextLength(); i++ {
			m := next[i]
			buffer := _driver.buffers[m.Instrument]
			note := m.MidiMessage
			result := C.write_midi_event(
				buffer,
				C.int(calculateHumanizedOffset(offset, m.HumanizeTime, nframes)),
				C.uchar(note.Command),
				C.uchar(note.Channel),
				C.uchar(note.Data1),
				C.uchar(note.Data2),
			)

			if result != 0 {
				fmt.Printf("Failed to write MIDI note!\n")
			}
		}
		_driver.buf.Flip()
	}
	return 1 // Keep looping.
}

func (j *jackDriver) OpenInstrument(name string) (int, error) {
	// TODO: This should also check that the name doesn't exceed the max length.
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	port := C.open_jack_port(
		j.client,
		cname,
	)

	if port == nil {
		return -1, fmt.Errorf("Error registering JACK port.")
	}

	id := instrumentID
	j.ports[id] = port
	fmt.Printf("Registered JACK port '%v' with ID %v.\n", name, id)
	instrumentID += 1
	return id, nil
}

func (j *jackDriver) CloseInstrument(id int) error {
	port, ok := j.ports[id]
	if !ok {
		panic(fmt.Sprintf("Internal error: Couldn't close instrument; no instrument with ID %v is open", id))
	}
	result := C.jack_port_unregister(j.client, port)
	if result != 0 {
		return fmt.Errorf("Error unregistering JACK port.")
	}
	delete(j.ports, id)
	return nil
}

func (j *jackDriver) Close() error {
	for id := range j.ports {
		err := j.CloseInstrument(id)
		if err != nil {
			fmt.Printf("Error closing instrument %v: %v\n", id, err)
		}
	}

	result := C.jack_client_close(j.client)
	if result != 0 {
		return fmt.Errorf("Error closing JACK client.")
	}

	j.client = nil
	return nil
}

func (j *jackDriver) Play(part types.Part, bpm int, ppq int, loop bool, polyphony int) error {
	buf, err := msg.NewBuffer(polyphony)
	if err != nil {
		return err
	}

	// TODO: Would be nice not to have to rely on globals, but doesn't make sense to
	// bounce this stuff off of C constantly, and we have to watch out for C storing
	// Go pointers, even temporarily.
	j.ppq = ppq
	j.part = part
	j.buf = buf
	j.length = part.Length(ppq)
	_driver = j

	result := C.run_jack_driver(j.client, C.int(bpm), C.int(ppq))
	if int(result) != JACK_OK {
		return fmt.Errorf("Error in JACK driver: %v", getErrorMessage(int(result)))
	}

	return nil
}
