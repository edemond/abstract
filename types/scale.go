package types

import (
	"fmt"
	"math"
)

var majorScale = NewScale([]int{2, 2, 1, 2, 2, 2, 1})

// e.g. scale(2,2,1,2,2,2,1)
type Scale struct {
	steps     []int // [2, 2, 1, 2, etc.] // Relative distances, in half-steps, between each scale degree, stored for convenience.
	intervals []int // [0, 2, 4, 5, etc.] // Absolute distance, in half-steps, of each scale degree from the root.
	width     int   // total span of the scale in half-steps
}

func (s *Scale) String() string {
	return fmt.Sprintf("scale(%v)", s.steps)
}

func (s *Scale) HasValue() bool {
	return s != nil
}

func NoScale() *Scale {
	return nil
}

func DefaultScale() *Scale {
	return majorScale
}

func NewScale(steps []int) *Scale {
	s := &Scale{steps: steps, intervals: make([]int, len(steps))}
	interval := 0
	for i, step := range s.steps {
		s.intervals[i] = interval
		interval += step
	}
	s.width = interval
	return s
}

// StepsAtDegree gets the number of half-steps from root at the given scale degree.
// Scale degree arguments to this function are zero-based and can be negative.
func (s *Scale) StepsAtDegree(degree int) int {
	length := len(s.steps)
	octaves := int(math.Floor(float64(degree) / float64(length)))
	i := ((degree % length) + length) % length
	return s.intervals[i] + (s.width * octaves)
}

// Test if the scale has the given steps. Used for unit testing.
func (s *Scale) HasSteps(steps []int) bool {
	if len(steps) != len(s.steps) {
		fmt.Printf("length mismatch\n")
		return false
	}
	for i := 0; i < len(steps); i++ {
		if steps[i] != s.steps[i] {
			fmt.Printf("step mismatch: %v %v\n", steps[i], s.steps[i])
			return false
		}
	}
	return true
}
