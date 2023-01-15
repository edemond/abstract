package types

import (
	"fmt"
	"math/rand"
)

// e.g. dynamics(80), dynamics(80, 12) // 80 +-12
type Dynamics struct {
	Center int
	Human  int // +- bound for humanizing. 0 value means no humanization.
}

func (d *Dynamics) String() string {
	return fmt.Sprintf("dynamics(%v, %v)", d.Center, d.Human)
}

func (d *Dynamics) HasValue() bool {
	return d != nil
}

// NewDynamics creates a Dynamics centered around a given velocity.
func NewDynamics(velocity int) *Dynamics {
	return &Dynamics{Center: velocity}
}

func (d *Dynamics) SetHumanize(human int) {
	d.Human = human
}

func NoDynamics() *Dynamics {
	return nil
}

func DefaultDynamics() *Dynamics {
	return &Dynamics{Center: 127, Human: 0} // full volume
}

func (d *Dynamics) Humanize(velocity uint64) uint64 {
	if d.Human == 0 || velocity == 0 {
		return velocity
	}
	h := uint64(rand.Int31n(int32(d.Human)))
	return (velocity - h) + (h * 2)
}
