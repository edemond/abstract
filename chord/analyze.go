package chord

import (
	"github.com/edemond/abstract/types"
	"fmt"
)

const OCTAVE = 12 // for convenience
var MAJOR = types.NewScale([]int{2, 2, 1, 2, 2, 2, 1})
var MINOR = types.NewScale([]int{2, 1, 2, 2, 1, 2, 2})
var MINORMAJOR = types.NewScale([]int{2, 1, 2, 2, 1, 3, 1})

func Analyze(expr chordExpr) (types.Chord, error) {
	switch e := expr.(type) {
	case *absoluteChordExpr:
		return analyzeAbsoluteChordExpr(e)
	case *relativeChordExpr:
		return analyzeRelativeChordExpr(e)
	case *diatonicChordExpr:
		return analyzeDiatonicChordExpr(e)
	}
	panic("Internal error: unhandled chord type")
}

// Get the scale implied by seeing this quality in the first position (i.e. immediately
// after the pitch in an absolute chord.)
func getScaleImpliedByQuality(quality Token) *types.Scale {
	switch quality {
	// TODO: Consult with mad on these. This is really dumb and is missing a lot of
	// the finer points of what scales are implied by what chords.
	case MAJ, AUG, DOMINANT, SUS, ADD, POWER, NO: // TODO: POWER is really neither...
		return MAJOR
	case MIN, DIM, HALF_DIM:
		return MINOR
	case MINMAJ:
		return MINORMAJOR
	}
	panic(fmt.Sprintf("Internal error: unhandled chord quality (%v)", quality))
}

func addInterval(chord map[int]int, interval int, scale *types.Scale) {
	if interval != 0 {
		// TODO: This will allow for some wacky no-op stuff like maj4, min5, maj1. Do we want?
		chord[interval] = scale.StepsAtDegree(interval - 1)
	}
}

func addMajorExtendedIntervals(chord map[int]int, interval int, scale *types.Scale) {
	switch interval {
	case 7:
		chord[7] = 11
	case 9:
		chord[7] = 11
		chord[9] = OCTAVE + 2
	// TODO: Check these extended ones. We may need to work in the concept of "avoid" notes.
	// There are a lot of variations on these...
	case 11:
		chord[7] = 11
		chord[9] = OCTAVE + 2
		chord[11] = OCTAVE + 5
	case 13:
		chord[7] = 11
		chord[9] = OCTAVE + 2
		chord[11] = OCTAVE + 5
		chord[13] = OCTAVE + 9
	default:
		addInterval(chord, interval, scale)
	}
}

func addMinorExtendedIntervals(chord map[int]int, interval int, scale *types.Scale) {
	switch interval {
	case 7:
		chord[7] = 10
	case 9:
		chord[7] = 10
		chord[9] = OCTAVE + 2
	case 11: // TODO: Check these extended ones.
		chord[7] = 10
		chord[9] = OCTAVE + 2
		chord[11] = OCTAVE + 5
	case 13:
		chord[7] = 10
		chord[9] = OCTAVE + 2
		chord[11] = OCTAVE + 5
		chord[13] = OCTAVE + 9 // It's actually a major 6 in the minor 13th chord.
	default:
		addInterval(chord, interval, scale)
	}
}

func addDominantExtendedIntervals(chord map[int]int, interval int) error {
	// "stacking donuts" --dave stewart (of hatfield and the north)
	switch interval {
	case 7:
		chord[7] = 10
	case 9:
		chord[7] = 10
		chord[9] = OCTAVE + 2
	case 11:
		chord[7] = 10
		chord[9] = OCTAVE + 2
		chord[11] = OCTAVE + 5
	case 13:
		chord[7] = 10
		chord[9] = OCTAVE + 2
		chord[11] = OCTAVE + 5
		chord[13] = OCTAVE + 9
	default:
		return fmt.Errorf("only 7, 9, 11, and 13 dominant chords supported (got '%v')", interval)
	}
	return nil
}

func getIntervals(qualities []*qualityExpr) (map[int]int, error) {
	if len(qualities) <= 0 {
		return nil, fmt.Errorf("A chord must have at least one quality.")
	}

	main := qualities[0]
	additional := qualities[1:]

	chord := map[int]int{} // scale degree -> value in halfsteps
	scale := getScaleImpliedByQuality(main.quality)

	// TODO: Wonder if we can unify the "getXExtendedIntervals" functions
	// using the scale implied by the quality?

	switch main.quality {
	case MAJ:
		chord[1] = 0
		chord[3] = 4
		chord[5] = 7
		addMajorExtendedIntervals(chord, main.interval, scale)
	case MIN:
		chord[1] = 0
		chord[3] = 3
		chord[5] = 7
		addMinorExtendedIntervals(chord, main.interval, scale)
	case AUG:
		chord[1] = 0
		chord[3] = 4
		chord[5] = 8
		// Augmented doesn't imply any scale that we could base extended intervals upon.
	case DIM:
		chord[1] = 0
		chord[3] = 3
		chord[5] = 6
		// Fully diminished seventh is a thing.
		if main.interval == 7 {
			chord[7] = 9
		} else {
			addInterval(chord, main.interval, scale) // TODO: Does this make sense for dim?
		}
	case POWER:
		chord[1] = 0
		chord[5] = 7
	case ADD:
		// If we have add in this position, e.g. Cadd2, it implies major. (Relative chords are different.)
		chord[1] = 0
		chord[3] = 4
		chord[5] = 7
		addInterval(chord, main.interval, MAJOR)
	case NO:
		// If we have "no" in this position, e.g. Cno5, it implies major. (Relative chords are different.)
		chord[1] = 0
		chord[3] = 4
		chord[5] = 7
		delete(chord, main.interval) // TODO: What happens if you do Cadd2no9?
	case SUS:
		// Sus means replace the third.
		if main.interval == 2 {
			chord[1] = 0
			chord[2] = 2
			chord[5] = 7
		} else if main.interval == 4 || main.interval == 0 {
			// sus without an interval, like Csus, means sus4.
			chord[1] = 0
			chord[4] = 5
			chord[5] = 7
		} else {
			return nil, fmt.Errorf("only sus, sus2, and sus4 supported (got '%v')", main.interval)
		}
	case HALF_DIM:
		chord[1] = 0
		chord[3] = 3
		chord[5] = 6
		chord[7] = 10
	case MINMAJ:
		chord[1] = 0
		chord[3] = 3
		chord[5] = 7
		switch main.interval {
		case 7:
			chord[7] = 11
		default:
			return nil, fmt.Errorf("only minmaj7 supported right now")
		}
	case DOMINANT:
		chord[1] = 0
		chord[3] = 4
		chord[5] = 7
		err := addDominantExtendedIntervals(chord, main.interval)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("unhandled main quality: %v", main.quality))
	}

	for i, q := range additional {
		switch q.quality {
		case MAJ:
			addMajorExtendedIntervals(chord, q.interval, scale)
		case MIN:
			addMinorExtendedIntervals(chord, q.interval, scale)
		case SUS:
			// sus, as the first explicit quality, overrides the implied main quality of a relative chord (e.g. Vsus4 is sus4, not major.)
			if main.implied && i == 0 {
				delete(chord, 3)
				if q.interval == 2 {
					chord[2] = 2
				} else if q.interval == 4 || q.interval == 0 {
					// sus without an interval, like Csus, means sus4.
					chord[4] = 5
				} else {
					return nil, fmt.Errorf("only sus, sus2, and sus4 supported (got '%v')", main.interval)
				}
			} else {
				return nil, fmt.Errorf("sus is only supported as the first chord quality.")
			}
		case ADD:
			// Add the note, flatted if implied by the main Quality.
			chord[q.interval] = scale.StepsAtDegree(q.interval - 1)
		case FLAT:
			chord[q.interval] = scale.StepsAtDegree(q.interval-1) - 1
		case SHARP:
			// Raise or lower the scale degree at the given interval, relative to the kind of
			// scale implied by the main Quality.
			chord[q.interval] = scale.StepsAtDegree(q.interval-1) + 1
		case NO:
			delete(chord, q.interval) // TODO: What about, say, Cadd2no9?
		case AUG:
			// Possible weirdness: This will allow augmented minor chords, like C Eb G#.
			_, ok := chord[5]
			if ok {
				chord[5] = scale.StepsAtDegree(5-1) + 1
			}
		case DIM:
			// Possible weirdness: This will allow diminished major chords, like C E Gb.
			_, ok := chord[5]
			if ok {
				chord[5] = scale.StepsAtDegree(5-1) - 1
			}
		case DOMINANT:
			err := addDominantExtendedIntervals(chord, q.interval)
			if err != nil {
				return nil, err
			}
		case POWER:
			chord[1] = 0
			delete(chord, 2)
			delete(chord, 3)
			delete(chord, 4)
			chord[5] = 7
			delete(chord, 6)
			delete(chord, 7)
		default:
			return nil, fmt.Errorf("unsupported additional quality: '%v'", q.quality) // TODO: better error message
		}
	}

	return chord, nil
}

func analyzeAbsoluteChordExpr(expr *absoluteChordExpr) (types.Chord, error) {
	chord, err := getIntervals(expr.qualities)
	if err != nil {
		return types.NoChord(), err
	}

	// Cobble the chord together into a list of intervals (in half steps).
	intervals := []int{}
	for _, halfSteps := range chord {
		intervals = append(intervals, halfSteps)
	}

	return types.NewAbsoluteChord(expr.pitch, intervals), nil
}

func analyzeRelativeChordExpr(expr *relativeChordExpr) (types.Chord, error) {
	chord, err := getIntervals(expr.qualities)
	if err != nil {
		return types.NoChord(), err
	}

	// Cobble the chord together into a list of intervals (in half steps).
	intervals := []int{}
	for _, halfSteps := range chord {
		intervals = append(intervals, halfSteps+expr.accidental)
	}

	return types.NewRelativeChord(expr.rootScaleDegree, intervals), nil
}

func analyzeDiatonicChordExpr(expr *diatonicChordExpr) (types.Chord, error) {
	// A set of scale degrees (Scale degree -> included or not)
	chord := map[int]bool{}

	// Start with a basic triad.
	// TODO: Extended chords. How do you do, say, a diatonic 7th chord? What's the notation?
	// The traditional seventh chord notation implies a quality; if you see @V7 you're going
	// to assume it's a dominant 7th chord, but if you're in a minor scale, you'd get Vmin7!
	chord[1] = true
	chord[3] = true
	chord[5] = true

	// Add additional qualities.
	for _, q := range expr.qualities {
		switch q.quality {
		// TODO: There's got to be more qualities these could support, right?
		case ADD:
			// TODO: Voicing hints, like add2 vs. add9
			chord[q.interval] = true
		case SUS:
			if q.interval == 2 {
				chord[2] = true
				chord[3] = false
			} else if q.interval == 4 || q.interval == 0 {
				// sus without an interval, like Csus, means sus4.
				chord[3] = false
				chord[4] = true
			} else {
				return nil, fmt.Errorf("only sus, sus2, and sus4 supported (got '%v')", q.interval)
			}
		case NO:
			chord[q.interval] = false
		case POWER:
			chord[2] = false
			chord[3] = false
			chord[4] = false
			chord[5] = true
		default:
			return nil, fmt.Errorf("Unsupported quality in diatonic chord: %v", q.quality)
		}
	}

	// Cobble the chord together into a list of scale degrees.
	degrees := []int{}
	for degree, included := range chord {
		if included {
			degrees = append(degrees, (degree + expr.rootScaleDegree - 1))
		}
	}

	return types.NewDiatonicChord(degrees), nil
}
