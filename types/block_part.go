package types

import (
	"github.com/edemond/abstract/msg"
)

// A list of parts to play sequentially.
type BlockPart struct {
	parts  []Part
	length uint64
}

func NewBlockPart() *BlockPart {
	return &BlockPart{
		parts:  make([]Part, 0),
		length: 0,
	}
}

func (b *BlockPart) Add(p Part) {
	b.parts = append(b.parts, p)
}

// TODO: These two functions are kind of a hack to support that optimization in
// the analyzer where we discard the outer part if it only contains one child part.
func (b *BlockPart) NumParts() int {
	return len(b.parts)
}

func (b *BlockPart) FirstPart() Part {
	return b.parts[0]
}

func (b *BlockPart) Play(buf msg.Buffer, ppq int, step uint64) {
	length := b.Length(ppq)
	if length == 0 {
		return
	}

	step = step % length
	steps := uint64(0)
	played := false

	for _, p := range b.parts {
		length := p.Length(ppq)
		if step < (steps + length) {
			// okay, here we need to give it what step it is in the child part.
			// that's (local step - start of child part in steps)
			p.Play(buf, ppq, step-steps)
			played = true
			// Zero-length parts don't count; play the next one immediately.
			if length != 0 {
				break
			}
		} else {
			steps += length
		}
	}
	if !played {
		panic("Internal error: Step counter exceeded block part index!")
	}
}

func (b *BlockPart) String() string {
	return "blockpart()"
}

func (b *BlockPart) HasValue() bool {
	return b != nil
}

// Length returns the length of the part in steps. Memoized.
// For block parts, this is the sum of the lengths of its constituent parts.
func (b *BlockPart) Length(ppq int) uint64 {
	if b.length <= 0 {
		for _, part := range b.parts {
			b.length += part.Length(ppq)
		}
	}
	return b.length
}
