package types

import (
	"fmt"
	"math"
	"math/rand"
)

// Humanized time.
type Humanize struct {
	Time int32 // +/- range of randomized time offset for each note
}

func NewHumanize(time uint64) (*Humanize, error) {
	if time > math.MaxInt32 {
		return nil, fmt.Errorf("Humanize must be 0-%v", math.MaxInt32)
	}

	return &Humanize{
		Time: int32(time),
	}, nil
}

func (h *Humanize) HasValue() bool {
	return h != nil
}

func (h *Humanize) String() string {
	return fmt.Sprintf("humanize(%v)", h.Time)
}

func NoHumanize() *Humanize {
	return nil
}

func DefaultHumanize() *Humanize {
	return NoHumanize()
}

func (h *Humanize) TimeOffset() int {
	if h.Time == 0 {
		return 0
	}
	offset := rand.Int31n(h.Time * 2)
	return int(offset - h.Time)
}
