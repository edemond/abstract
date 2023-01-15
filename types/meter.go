package types

import (
	"fmt"
)

type Meter struct {
	Beats     int
	Value     int // lower number (value) of the meter, not a language Value.
	multiples []int
	divisions []int
}

func (m *Meter) String() string {
	return fmt.Sprintf("meter(%v, %v)", m.Beats, m.Value)
}

func (m *Meter) HasValue() bool {
	return m != nil
}

func NoMeter() *Meter {
	return nil
}

func DefaultMeter() *Meter {
	return &Meter{Beats: 4, Value: 4} // 4/4 time
}

// Length gets the number of steps encompassed by this meter at the given ppq.
func (m *Meter) Length(ppq int) uint64 {
	// TODO: Theoretically, this can overflow, buuuuuuut...
	// TODO: Memoize this, maybe? Meter is immutable.
	return uint64(m.Beats) * uint64((ppq*4)/m.Value) // times 4, because ppq is Pulses Per QUARTER (Note)
}

func (m *Meter) IsDuple() bool {
	return m.Beats%2 == 0
}
