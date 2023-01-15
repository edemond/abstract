package types

// Rhythmic/expressive context.
type Rhythm struct {
	Dynamics    *Dynamics
	Humanize    *Humanize
	Meter       *Meter
	defaultsSet bool
}

// Initialize the rhythmic context with default values wherever
// something hasn't been provided by the user. This should
// only be called just before playing the part.
func (r *Rhythm) SetDefaults() {
	if !r.defaultsSet {
		if !r.Dynamics.HasValue() {
			r.Dynamics = DefaultDynamics()
		}
		if !r.Humanize.HasValue() {
			r.Humanize = DefaultHumanize()
		}
		if !r.Meter.HasValue() {
			r.Meter = DefaultMeter()
		}
		r.defaultsSet = true
	}
}

// Get the pulse (beat number, beat strength) at the given step.
// The "beat number" counts which beat in the measure we're on, at the beat level (e.g. 1-4 in a measure of 4/4).
// The "beat strength" indicates the strength of the beat, ranging from 0 (no beat) to 32nd notes.
// (See the test cases of the Pulse function for examples.)
//
// e.g. throughout a bar of 4/4 this will return (omitting 32nd notes and 0s for space):
// number:   1  1  1  1  2  2  2  2  3  3  3  3  4  4  4  4
// strength: 1  16 8  16 4  16 8  16 2  16 8  16 4  16 8  16
//
func (r *Rhythm) Pulse(step uint64, ppq int) (int, int) {
	stepsPerBeat := uint64((ppq * 4) / r.Meter.Value)
	beat := int((step / stepsPerBeat) + 1)

	// Start at the beat level, because the measure is by definition divisible by that number of beats of that value
	// Then go down to division levels (eighth notes, sixteenth notes...)
	// Then go up to multiple levels (half notes, whole notes...)

	// Downbeat, the 1, is a special case.
	if step == 0 {
		return beat, 1
	}

	// Multiples (half notes)
	// We only have half note strengths if the meter is duple (i.e. upper number divisible by 2).
	// A measure of 3/4 would look like 1 4 4, because it can't divide in half.
	if r.Meter.IsDuple() && (step%(stepsPerBeat*uint64(r.Meter.Beats/2)) == 0) {
		return beat, 2
	}

	// Beat level and divisions, down to 32nd notes.
	// Start at the beat level, then go below the beat level by dividing by multiples of 2.
	// e.g. Divide by 1 (beat level), 2 (eigth notes), 4 (sixteenths), 8 (thirty-second notes).
	for div := uint64(1); div <= 8; div *= 2 {
		if stepsPerBeat%div == 0 {
			if step%(stepsPerBeat/div) == 0 {
				return beat, int(div * 4)
			}
		}
	}

	return beat, 0 // No recognized beat division.

	// TODO: Wait. The code below was a mistake for what I was trying to do. It works
	// in 4/4, but in 3/4, I think it actually calculates 4-against-3 polyrhythm!
	// This might be useful later on.
	/*
	   for i := 1; i <= 32; i *= 2 {
	       if (length-step) % (length / uint64(i)) == 0 {
	           strength = i
	           break
	       }
	   }
	*/
}
