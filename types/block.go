package types

// A "block" expression is the default way to play a chord or a note (Interpretation).
// It's a direction to play the entire block chord or single note right on the downbeat.
// Not to be confused with BlockPart.
type Block struct{}

func NewBlock() *Block {
	return &Block{}
}

func (b *Block) Play(notesOut []Note, h *Harmony, r *Rhythm, counter uint64, step uint64, length uint64, ppq int) {
	beat, strength := r.Pulse(step, ppq)
	// Block chords only play on the downbeat (1,1).
	if beat == 1 && strength == 1 {
		if h.Chord.HasValue() {
			h.Chord.Play(notesOut, h) // Block chord will sound.
		} else {
			notesOut[0] = h.Pitch.At(h.Octave) // Just a single note.
		}
	}
}

func (b *Block) String() string {
	return "block chord"
}

func (b *Block) HasValue() bool {
	// A Block is the "default" and hence null object Interpretation. There's no syntax
	// to specify a Block, it's just the default thing that happens, so you can override it.
	return false
}
