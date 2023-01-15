package types

import (
	"edemond/abstract/msg"
	"edemond/abstract/util"
	"fmt"
	"strings"
)

type Seq struct {
	parts      []Part
	partRanges []partRange // [start, end] pairs for each part, in order.
	printed    bool
	length     uint64
	scale      int // Scaling factor.
}

func NewSeqPart() *Seq {
	s := &Seq{
		parts:      []Part{},
		parent:     nil,
		partRanges: nil,
		scale:      1,
	}
	return s
}

func (s *Seq) SetScale(scale int) {
	s.scale = scale
}

func (s *Seq) SetParts(parts []Part) {
	s.parts = parts
}

func (s *Seq) SetParent(parent Part) {
	s.parent = parent
}

func (s *Seq) NumParts() int {
	return len(s.parts)
}

// TODO: Can we get rid of this? It shouldn't be dependent on having ppq
// in order to build this either.
type partRange []int

func (p partRange) start() int { return p[0] }
func (p partRange) end() int   { return p[1] }

// This hinges on a length based on a given ppq, and thus
// cannot be used until after semantic analysis.
func getPartRanges(parts []Part, length int) []partRange {
	fmt.Printf("getPartRanges: length %v\n", length)
	// Distribute the parts evenly over the length using Bjorklund,
	// otherwise we'll end up with all the parts stacked at the
	// start and with a gap somewhere.
	// This assumes all parts are the same length, which is true by
	// the definition of a seq part.
	pattern := util.Bjorklund(len(parts), length)

	start := -1
	ranges := []partRange{}
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == 1 {
			// okay, saw a 1. this can be either the start of the first part,
			// in which case we have nothing to add,
			// or the start of the second one, in which case we add the last one.

			// now how do we know if we've seen anything yet?
			if start != -1 {
				ranges = append(ranges, []int{start, i - 1})
			}

			// either way, something starts now.
			start = i
		}
	}
	if start != -1 {
		ranges = append(ranges, []int{start, length - 1})
	}

	// We may not end up with a 1 in the first position. Rotate the pattern so
	// that we always do.
	rotate := ranges[0].start()
	for _, r := range ranges {
		// TODO: Cheating! Clean this up; add setters.
		r[0] -= rotate
		r[1] -= rotate
	}
	// The last range should extend all the way to the end so the whole length is covered.
	ranges[len(ranges)-1][1] = length

	return ranges
}

// A Seq is a Value.
func (s *Seq) HasValue() bool {
	return s != nil
}

func (s *Seq) String() string {
	parts := make([]string, len(s.parts))
	for i, p := range s.parts {
		parts[i] = p.String()
	}
	return fmt.Sprintf("[%v]", strings.Join(parts, " "))
}

// at returns which part to play for the given (local!) step, plus the
// (local!) step at which it starts.
func (s *Seq) at(step uint64, ppq int) (Part, uint64) {
	// TODO: Bjorklund during playback isn't ideal.
	if s.partRanges == nil {
		ln := int(s.Length(ppq))
		fmt.Printf("length: %v, ppq: %v, parent: %v\n", ln, ppq, s.parent)
		s.partRanges = getPartRanges(s.parts, ln)
	}

	for index, rng := range s.partRanges {
		if step >= uint64(rng.start()) && step <= uint64(rng.end()) {
			return s.parts[index], uint64(rng.start())
		}
	}
	panic(fmt.Sprintf("There was no part in the map for step %v! (ranges %v)", step, s.partRanges))
}

// A Seq is a Part. TODO: rename, lol
func (s *Seq) Play(buf msg.Buffer, ppq int, step uint64) {
	length := s.Length(ppq)
	if length == 0 {
		return
	}

	// We need to stretch this out evenly over the meter. If there are four values,
	// they get applied one per quarter note in a meter of 4/4, etc.
	// this enables neat stuff like drag triplets: 4/4 [a b c]
	step = step % length

	// TODO: This is going to fail or at least behave weirdly for small enough ppq.
	ppq = ppq / len(s.parts) // Scale ppq down to compress the subparts into one measure.
	part, start := s.at(step, ppq)

	if !s.printed {
		fmt.Println("----------------------------------")
		for i, pr := range s.partRanges {
			part := s.parts[i]
			fmt.Printf("%v, start: %v, end: %v\n", part, pr[0], pr[1])
		}
		s.printed = true
	}

	fmt.Printf("%v: %v\n", step, part)
	part.Play(buf, ppq, step-start)
}

// A seq part is exactly as long as the simple part in which it was found.
// TODO: WHOA. We DON'T want seq parts to be "contained" in simple parts.
// An expression like "C lydian [bleh blah]" should evaluate to a seq part,
// not a simple part!
func (s *Seq) Length(ppq int) uint64 {
	if s.length == 0 {
		s.length = s.parent.Length(ppq) / s.scale
	}
	return s.length
}
