package types

type Rest int

func (r Rest) HasValue() bool {
	return true
}

func (r Rest) String() string {
	return "_"
}

func IsRest(value string) bool {
	return value == "_"
}

func NewRest(value string) Rest {
	return Rest(0)
}

// A Rest doesn't sound, by definition.
func (r Rest) Play(notes []Note, h *Harmony, rh *Rhythm, counter uint64, step uint64, length uint64, ppq int) {
	// TODO: Turn off previous notes here.
}
