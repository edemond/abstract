package chord

import (
	"github.com/edemond/abstract/types"
)

type chordExpr interface {
	isChordExpr()
}

type absoluteChordExpr struct {
	pitch     types.Pitch
	qualities []*qualityExpr
}

type relativeChordExpr struct {
	rootScaleDegree int
	accidental      int // -1 for flat, 1 for sharp, 0 for natural. We don't support double sharps or flats here.
	qualities       []*qualityExpr
}

type diatonicChordExpr struct {
	rootScaleDegree int
	// Some qualities are OK in a diatonic context (sus4, power chord), but not others (major, minor)!
	// The analyzer will enforce this.
	qualities []*qualityExpr
}

type rootExpr string // I, ii, @iii

type qualityExpr struct {
	quality  Token // maj, min, aug, dim, ∅, whatever. TODO: Token is doing double-duty here, is that wise?
	interval int   // optional, defaults to 0. for maj7, min11, etc.
	implied  bool  // A quality that was not explicitly specified, but inferred from, say, the capitalization of the root symbol (II vs. ii).
}

func (a *absoluteChordExpr) isChordExpr() {}
func (a *relativeChordExpr) isChordExpr() {}
func (a *diatonicChordExpr) isChordExpr() {}
