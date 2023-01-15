package chord

import (
	"github.com/edemond/abstract/types"
	"fmt"
	"testing"
)

func getPitch(s string) types.Pitch {
	pitch, err := types.LookUpPitch(s)
	if err != nil {
		panic(fmt.Sprintf("'%v' is not a pitch, go fix the test", s))
	}
	return pitch
}

func failWithExpected(t *testing.T, expected []string, got []types.Pitch) {
	t.Fatalf("Expected the set %v, got %v\n", expected, got)
}

// Test an absolute chord.
func testAbs(t *testing.T, symbol string, expectedPitches []string) {
	test(t, symbol, types.NoPitch(), MAJOR, expectedPitches)
}

// Verify that a symbol does NOT parse as a chord.
func testBadAbs(t *testing.T, symbol string) {
	testBadChord(t, symbol, types.NoPitch(), MAJOR)
}

// Test a relative chord.
func testRel(t *testing.T, symbol string, pitch string, expectedPitches []string) {
	test(t, symbol, getPitch(pitch), MAJOR, expectedPitches)
}

// Test a diatonic chord.
func testDia(t *testing.T, symbol string, pitch string, scale *types.Scale, expectedPitches []string) {
	test(t, symbol, getPitch(pitch), scale, expectedPitches)
}

// Test how a chord is interpreted in a given scale and pitch.
func test(t *testing.T, symbol string, pitch types.Pitch, scale *types.Scale, expectedPitches []string) {
	chord, err := ParseAndAnalyze(symbol)
	if err != nil {
		t.Fatalf("'%v' didn't parse: %v", symbol, err)
	}

	set := make(map[types.Pitch]bool)
	for _, pitch := range expectedPitches {
		set[getPitch(pitch)] = false
	}

	actualPitches := chord.ResolveIn(pitch, scale) // Scale doesn't matter for absolute or relative chords

	// Check for pitches we got but didn't expect...
	for _, pitch := range actualPitches {
		_, ok := set[pitch]
		if !ok {
			failWithExpected(t, expectedPitches, actualPitches)
		}
		set[pitch] = true
	}

	// ...and ones we expected but didn't get.
	for _, found := range set {
		if !found {
			failWithExpected(t, expectedPitches, actualPitches)
		}
	}
}

func testBadChord(t *testing.T, symbol string, pitch types.Pitch, scale *types.Scale) {
	chord, err := ParseAndAnalyze(symbol)
	if err == nil {
		pitches := chord.ResolveIn(pitch, scale)
		t.Fatalf("'%v' was not expected to parse, but parsed as: %v", symbol, pitches)
	}
}

func TestAbsoluteMajorChord(t *testing.T) {
	testAbs(t, "Cmaj", []string{"C", "E", "G"})
	testAbs(t, "Dmaj", []string{"D", "F#", "A"})
	testAbs(t, "Bmaj", []string{"B", "F#", "D#"})
	testAbs(t, "CM", []string{"C", "E", "G"})
	testAbs(t, "DM", []string{"D", "F#", "A"})
	testAbs(t, "BM", []string{"B", "F#", "D#"})
	testAbs(t, "CΔ", []string{"C", "E", "G"})
	testAbs(t, "D∆", []string{"D", "F#", "A"})
	testAbs(t, "DΔ", []string{"D", "F#", "A"}) // NOT a duplicate of the above line, different delta char
	testAbs(t, "BΔ", []string{"B", "F#", "D#"})
}

func TestSus2Chord(t *testing.T) {
	testAbs(t, "Csus2", []string{"C", "D", "G"})
	testAbs(t, "C5add9", []string{"C", "G", "D"})
}

func TestSus4Chord(t *testing.T) {
	testAbs(t, "Csus", []string{"C", "F", "G"})
	testAbs(t, "Csus4", []string{"C", "F", "G"})
}

func TestAdd2Chord(t *testing.T) {
	testAbs(t, "Cadd2", []string{"C", "D", "E", "G"})
}

func TestAddStuffChords(t *testing.T) {
	testAbs(t, "Badd9", []string{"B", "D#", "F#", "C#"})
	testAbs(t, "C6add2", []string{"C", "D", "E", "G", "A"})
}

func TestAbsoluteMajor6Chord(t *testing.T) {
	testAbs(t, "C6", []string{"C", "E", "G", "A"})
	testAbs(t, "Cmaj6", []string{"C", "E", "G", "A"})
	testAbs(t, "BM6", []string{"B", "F#", "D#", "G#"})
	testAbs(t, "D∆6", []string{"D", "F#", "A", "B"})
	testAbs(t, "DΔ6", []string{"D", "F#", "A", "B"}) // NOT a duplicate of the above line, different delta char
	testAbs(t, "BΔ6", []string{"B", "F#", "D#", "G#"})
}

func TestAbsoluteMinorChord(t *testing.T) {
	testAbs(t, "Cmin", []string{"C", "Eb", "G"})
	testAbs(t, "Bmin", []string{"B", "D", "F#"})
	testAbs(t, "Cm", []string{"C", "Eb", "G"})
	testAbs(t, "Bm", []string{"B", "D", "F#"})
	testAbs(t, "C-", []string{"C", "Eb", "G"})
	testAbs(t, "B-", []string{"B", "D", "F#"})
	testAbs(t, "C–", []string{"C", "Eb", "G"})
	testAbs(t, "B–", []string{"B", "D", "F#"})
	testAbs(t, "C—", []string{"C", "Eb", "G"})
	testAbs(t, "B—", []string{"B", "D", "F#"})
}

func TestAbsoluteAugmentedChord(t *testing.T) {
	testAbs(t, "C+", []string{"C", "E", "G#"})
	testAbs(t, "Baug", []string{"B", "D#", "F##"})
}

func TestAbsoluteDiminishedChord(t *testing.T) {
	testAbs(t, "Cdim", []string{"C", "Eb", "Gb"})
	testAbs(t, "Co", []string{"C", "Eb", "Gb"})
	testAbs(t, "B°", []string{"B", "D", "F"})
	testAbs(t, "A°", []string{"A", "C", "Eb"})
}

func TestAbsoluteHalfDiminished7Chord(t *testing.T) {
	testAbs(t, "C∅", []string{"C", "Eb", "Gb", "Bb"})
	testAbs(t, "Bø", []string{"B", "D", "F", "A"})
	testAbs(t, "AØ", []string{"A", "C", "Eb", "G"})
}

func TestAbsoluteMinor7Chord(t *testing.T) {
	testAbs(t, "Cmin7", []string{"C", "Eb", "G", "Bb"})
	testAbs(t, "Bm7", []string{"B", "D", "F#", "A"})
	testAbs(t, "A-7", []string{"A", "C", "E", "G"})
	testAbs(t, "A–7", []string{"A", "C", "E", "G"}) // not a duplicate
	testAbs(t, "A—7", []string{"A", "C", "E", "G"}) // also not a duplicate
}

func TestAbsolutePowerChord(t *testing.T) {
	testAbs(t, "C5", []string{"C", "G"})
}

func TestAbsoluteNo(t *testing.T) {
	testAbs(t, "Cno5", []string{"C", "E"})
	testAbs(t, "C7no3", []string{"C", "G", "Bb"})
}

// Dominants
func TestAbsoluteDominant7Chord(t *testing.T) {
	testAbs(t, "C7", []string{"C", "E", "G", "Bb"})
	testAbs(t, "Adom7", []string{"A", "C#", "E", "G"})
}

func TestAbsoluteDominant9Chord(t *testing.T) {
	testAbs(t, "C9", []string{"C", "E", "G", "Bb", "D"})
}

func TestAbsoluteMajor7Chord(t *testing.T) {
	testAbs(t, "Cmaj7", []string{"C", "E", "G", "B"})
	testAbs(t, "BM7", []string{"B", "D#", "F#", "A#"})
	testAbs(t, "AΔ7", []string{"A", "C#", "E", "G#"})
	testAbs(t, "AΔ7", []string{"A", "C#", "E", "G#"}) // not a duplicate
}

func TestAbsoluteMajorNinthChord(t *testing.T) {
	testAbs(t, "Cmaj9", []string{"C", "E", "G", "B", "D"})
	testAbs(t, "Cmaj9add6", []string{"C", "E", "G", "A", "B", "D"})
}

func TestAbsoluteDominant7Flat9(t *testing.T) {
	testAbs(t, "C7b9", []string{"C", "E", "G", "Bb", "Db"})
}

func TestAbsoluteMinorNinth(t *testing.T) {
	testAbs(t, "Cmin9", []string{"C", "Eb", "G", "Bb", "D"})
	testAbs(t, "Cm9", []string{"C", "Eb", "G", "Bb", "D"})
}

func TestAbsoluteMinorMajorChord(t *testing.T) {
	testAbs(t, "Cminmaj7", []string{"C", "Eb", "G", "B"})
}

func TestAbsoluteDominant7Sharp9(t *testing.T) {
	testAbs(t, "C7#9", []string{"C", "E", "G", "Bb", "D#"})
}

func TestAbsoluteAugmented7Chord(t *testing.T) {
	testAbs(t, "C+7", []string{"C", "E", "G#", "Bb"})
}

func TestRelativeChord(t *testing.T) {
	testRel(t, "bVI", "C", []string{"Ab", "C", "Eb"})
	testRel(t, "#VI", "C", []string{"A#", "C##", "E#"})

	testRel(t, "I", "C", []string{"C", "E", "G"})
	testRel(t, "ii", "C", []string{"D", "F", "A"})
	testRel(t, "iii", "C", []string{"E", "G", "B"})
	testRel(t, "IV", "C", []string{"F", "A", "C"})
	testRel(t, "V", "C", []string{"G", "B", "D"})
	testRel(t, "vi", "C", []string{"A", "C", "E"})
	testRel(t, "viio", "C", []string{"B", "D", "F"})
	testRel(t, "I", "B", []string{"B", "D#", "F#"})
	testRel(t, "i", "C", []string{"C", "Eb", "G"})

	testRel(t, "I7", "C", []string{"C", "E", "G", "Bb"})
	testRel(t, "I9", "C", []string{"C", "E", "G", "Bb", "D"})
	testRel(t, "Imaj9", "C", []string{"C", "E", "G", "B", "D"})
}

func TestBadChords(t *testing.T) {
	testBadAbs(t, "Hmin")
	testBadAbs(t, "Amx")
	testBadAbs(t, "C") // This should parse as a pitch.
	testBadAbs(t, "maj")
	testBadAbs(t, "maj9")
	testBadAbs(t, "C0")
	testBadAbs(t, "C1")
	testBadAbs(t, "CA")
	testBadAbs(t, "7")
	testBadAbs(t, "Fadd")
	// Questionable:
	/*
		testBadAbs(t, "Cadd45")
		testBadAbs(t, "Gm4")
	*/
}

func TestDiatonicMajorScaleChords(t *testing.T) {
	testDia(t, "@I", "C", MAJOR, []string{"C", "E", "G"})
	testDia(t, "@II", "C", MAJOR, []string{"D", "F", "A"})
	testDia(t, "@III", "C", MAJOR, []string{"E", "G", "B"})
	testDia(t, "@IV", "C", MAJOR, []string{"F", "A", "C"})
	testDia(t, "@V", "C", MAJOR, []string{"G", "B", "D"})
	testDia(t, "@VI", "C", MAJOR, []string{"A", "C", "E"})
	testDia(t, "@VII", "C", MAJOR, []string{"B", "D", "F"})
}

func TestDiatonicMinorScaleChords(t *testing.T) {
	testDia(t, "@I", "C", MINOR, []string{"C", "Eb", "G"})
	testDia(t, "@II", "C", MINOR, []string{"D", "F", "Ab"})
	testDia(t, "@III", "C", MINOR, []string{"Eb", "G", "Bb"})
	testDia(t, "@IV", "C", MINOR, []string{"F", "Ab", "C"})
	testDia(t, "@V", "C", MINOR, []string{"G", "Bb", "D"})
	testDia(t, "@VI", "C", MINOR, []string{"Ab", "C", "Eb"})
	testDia(t, "@VII", "C", MINOR, []string{"Bb", "D", "F"})
}

func TestDiatonicMajorScaleAddedNoteChords(t *testing.T) {
	testDia(t, "@Iadd2", "C", MAJOR, []string{"C", "D", "E", "G"})
	testDia(t, "@IIadd6", "C", MAJOR, []string{"D", "F", "A", "B"})
	testDia(t, "@IIIadd4", "C", MAJOR, []string{"E", "G", "A", "B"})
	testDia(t, "@IVadd9", "C", MAJOR, []string{"F", "A", "C", "G"})
	testDia(t, "@Vadd7", "C", MAJOR, []string{"G", "B", "D", "F"})
	testDia(t, "@VIadd2", "C", MAJOR, []string{"A", "B", "C", "E"})
}

func TestDiatonicMajorScaleSusChords(t *testing.T) {
	testDia(t, "@Isus2", "C", MAJOR, []string{"C", "D", "G"})
	testDia(t, "@IIsus4", "C", MAJOR, []string{"D", "G", "A"})
	testDia(t, "@IIIsus", "C", MAJOR, []string{"E", "A", "B"})
}

func TestDiatonicMajorScaleNoChords(t *testing.T) {
	testDia(t, "@VIadd2no5", "C", MAJOR, []string{"A", "B", "C"})
}
