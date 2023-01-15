package types

// Harmonic context.
type Harmony struct {
	Chord       Chord
	Octave      Octave
	Pitch       Pitch
	Scale       *Scale
	Voicing     Voicing
	defaultsSet bool
}

// Initialize the harmonic context with default values wherever
// something hasn't been provided by the user. This should
// only be called just before playing the part.
func (h *Harmony) SetDefaults() {
	if !h.defaultsSet {
		if !h.Chord.HasValue() {
			h.Chord = DefaultChord()
		}
		if !h.Octave.HasValue() {
			h.Octave = DefaultOctave()
		}
		if !h.Pitch.HasValue() {
			h.Pitch = DefaultPitch()
		}
		if !h.Scale.HasValue() {
			h.Scale = DefaultScale()
		}
		if !h.Voicing.HasValue() {
			h.Voicing = DefaultVoicing()
		}
		h.defaultsSet = true
	}
}

// Get the root pitch of the harmonic context, at least our best guess.
func (h *Harmony) Root() Pitch {
	if !h.Chord.HasValue() {
		return h.Pitch // There's always one of these.
	}
	root := h.Chord.Root()
	if !root.HasValue() {
		return h.Pitch
	}
	return root
}

func (h *Harmony) Bass() Pitch {
	return h.Root() // TODO: Inversions!
}
