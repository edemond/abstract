package drivers

import (
	"github.com/edemond/abstract/types"
)

type Driver interface {
	// Play the piece from the given root part.
	// polyphony: The maximum number of voices that might be playing at once.
	Play(part types.Part, bpm int, ppq int, loop bool, polyphony int) error

	// Open an instrument for playback. Returns an instrument ID.
	OpenInstrument(name string) (int, error)

	// Close an instrument, given the ID returned from OpenInstrument.
	CloseInstrument(id int) error

	// Close the driver.
	Close() error
}
