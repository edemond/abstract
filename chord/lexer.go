// A lexer for chords, heavily based on Go's "go/scanner".
package chord

import (
	"fmt"
	"unicode/utf8"
)

type Token int

const (
	INVALID Token = iota
	EOF           // well, not really EOF, more like EOC

	NUMBER
	PITCH // C D E F G A B
	ROOT  // I II III IV V VI VII i ii iii iv v vi vii Ⅰ Ⅱ Ⅲ Ⅳ Ⅴ Ⅵ Ⅶ ⅰ ⅱ ⅲ ⅳ ⅴ ⅵ ⅶ - roots

	// symbols/operators
	DIATONIC // @ - indicates a diatonic chord
	FLAT     // b ♭
	SHARP    // # ♯
	NATURAL  // ♮
	SLASH    // /
	NO       // no

	// TODO: There are groups of these that are mutually exclusive. For example, we can't
	// have both "min" and "maj" a chord, sharp and flat, or "dom" and "min". Need to figure
	// out how to represent these groups and verify them.

	// qualities
	MAJ      // maj
	MIN      // min m
	MINMAJ   // minmaj
	ADD      // add
	SUS      // sus
	DOMINANT // dom
	AUG      // aug +
	DIM      // dim o °
	HALF_DIM // Ø ø ∅ m7b5 m7♭5
	POWER    // TODO: there's no symbol for this one, but we need a Token type for it because Token is pulling double-duty as a chord quality enum
)

var lookup = map[string]Token{

	"@": DIATONIC,

	// Sometimes lowercase letters are used to indicate minor chords, but we disallow it
	// because of ambiguity with "b" (is it a flat or B minor?)
	"C": PITCH,
	"D": PITCH,
	"E": PITCH,
	"F": PITCH,
	"G": PITCH,
	"A": PITCH,
	"B": PITCH,

	"I":   ROOT,
	"II":  ROOT,
	"III": ROOT,
	"IV":  ROOT,
	"V":   ROOT,
	"VI":  ROOT,
	"VII": ROOT,

	"i":   ROOT,
	"ii":  ROOT,
	"iii": ROOT,
	"iv":  ROOT,
	"v":   ROOT,
	"vi":  ROOT,
	"vii": ROOT,

	"Ⅰ": ROOT,
	"Ⅱ": ROOT,
	"Ⅲ": ROOT,
	"Ⅳ": ROOT,
	"Ⅴ": ROOT,
	"Ⅵ": ROOT,
	"Ⅶ": ROOT,

	"ⅰ": ROOT,
	"ⅱ": ROOT,
	"ⅲ": ROOT,
	"ⅳ": ROOT,
	"ⅴ": ROOT,
	"ⅵ": ROOT,
	"ⅶ": ROOT,

	"b": FLAT,
	"♭": FLAT,
	"#": SHARP,
	"♯": SHARP,
	"♮": NATURAL,

	// Qualities -----------

	"maj": MAJ,
	"M":   MAJ,
	"∆":   MAJ, // delta
	"Δ":   MAJ, // a more different delta

	"min": MIN,
	"m":   MIN,
	"-":   MIN, // minus
	"–":   MIN, // en dash
	"—":   MIN, // em dash

	"aug": AUG,
	"+":   AUG,

	"dim": DIM,
	"o":   DIM,
	"°":   DIM,

	// half-diminished 7th
	"∅":    HALF_DIM,
	"ø":    HALF_DIM,
	"Ø":    HALF_DIM,
	"m7b5": HALF_DIM,
	"m7♭5": HALF_DIM,

	// Parameterized qualities --------
	// stuff that needs an interval
	"add":    ADD,
	"sus":    SUS,
	"minmaj": MINMAJ, // minor-major, needs to be filled out with an interval
	"dom":    DOMINANT,
	"no":     NO,
	// stuff that takes a pitch
	"/": SLASH, // e.g. F/Ab
}

func (t Token) String() string {
	switch t {
	case INVALID:
		return "(invalid chord symbol token)"
	case EOF:
		return "(end of chord)"
	case ROOT:
		return "root"
	case PITCH:
		return "pitch"
	case NUMBER:
		return "number"
	case DIATONIC:
		return "diatonic"
	case FLAT:
		return "flat"
	case SHARP:
		return "sharp"
	case NATURAL:
		return "natural"
	case AUG:
		return "augmented"
	case DIM:
		return "diminished"
	case HALF_DIM:
		return "half-diminished 7th"
	case SLASH:
		return "slash"
	case MAJ:
		return "major"
	case MIN:
		return "minor"
	case MINMAJ: // minmaj
		return "minor/major"
	case ADD:
		return "add"
	case SUS:
		return "sus"
	case DOMINANT:
		return "dom"
	default:
		panic(fmt.Sprintf("unknown chord token type: %v", int(t)))
	}
}

type Lexer struct {
	source []byte
	start  int // start position
	end    int // end position
	char   rune
}

func NewLexerFromString(src string) (*Lexer, error) {
	lex := &Lexer{
		source: []byte(src),
		start:  0,
		end:    0,
		char:   ' ',
	}
	lex.next()
	return lex, nil
}

// next advances the lexer one rune. As in go/scanner, lex.char == -1 is EOF.
func (lex *Lexer) next() {
	lex.start = lex.end
	r, size := utf8.DecodeRune(lex.source[lex.start:])
	if size == 0 {
		lex.char = -1 // EOF.
	} else if size == 1 && r == utf8.RuneError {
		// TODO: handle errors
		panic("invalid UTF-8 character in chord symbol")
	} else {
		lex.end += size
		lex.char = r
	}
}

func (lex *Lexer) scanASCIIRoot() string {
	start := lex.start
	for isASCIIRoot(lex.char) {
		lex.next()
	}
	return string(lex.source[start:lex.start])
}

// Alphabetical symbols like "add", "sus", "dim", etc.
func (lex *Lexer) scanAlpha() string {
	start := lex.start
	for isLetter(lex.char) {
		lex.next()
	}
	return string(lex.source[start:lex.start])
}

func (lex *Lexer) scanNumber() string {
	start := lex.start
	for isDigit(lex.char) {
		lex.next()
	}
	return string(lex.source[start:lex.start])
}

func (lex *Lexer) scanPitch() string {
	start := lex.start
	for isPitch(lex.char) {
		lex.next()
	}
	return string(lex.source[start:lex.start])
}

// ASCII characters for Roman numeral chord root symbols.
func isASCIIRoot(r rune) bool {
	switch r {
	case 'i', 'I', 'v', 'V':
		return true
	}
	return false
}

// Single-rune symbols that we recognize.
func IsChordNotationSymbol(r rune) bool {
	switch r {
	case 'C', 'D', 'E', 'F', 'G', 'A', 'B':
		return true
	case 'Ⅰ', 'Ⅱ', 'Ⅲ', 'Ⅳ', 'Ⅴ', 'Ⅵ', 'Ⅶ', 'ⅰ', 'ⅱ', 'ⅲ', 'ⅳ', 'ⅴ', 'ⅵ', 'ⅶ':
		return true
	case '@', 'b', '♭', '#', '♯', '♮', '+', 'o', '°', '∘', '○':
		return true
	case '〇', 'Ø', 'ø', '∅', '/', '∆', 'Δ', '-', '–', '—':
		return true
	}
	return false
}

func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isPitchLetter(r rune) bool {
	return (r >= 'A' && r <= 'G')
}

func isPitch(r rune) bool {
	return isPitchLetter(r) || r == 'b' || r == '♭' || r == '#' || r == '♯' || r == '♮'
}

// getTokenType looks up a token's type and returns the (type, value) of the token.
func getTokenType(s string) (Token, string) {
	tok, ok := lookup[s]
	if !ok {
		return INVALID, s
	}
	return tok, s
}

// Returns the next token as (token type, value).
func (lex *Lexer) Scan() (Token, string) {
	switch ch := lex.char; {

	case isDigit(ch):
		return NUMBER, lex.scanNumber()

	// Roman numeral root symbols spelled out with [IVvi].
	case isASCIIRoot(ch):
		return getTokenType(lex.scanASCIIRoot())

	// Pitch literals like C#, Eb, F, G♭♭.
	case isPitchLetter(ch):
		return PITCH, lex.scanPitch()

	// This case needs to come before isLetter to correctly handle "b" and "o" (flat and diminished).
	case IsChordNotationSymbol(ch):
		lex.next()
		return getTokenType(string(ch))

	// This case needs to be after isSymbol and isRoot to correctly handle [boIVvi] (flat, diminished, and Roman numerals).
	// That means none of the chord quality symbols can start with those letters, but I don't know of any that do.
	case isLetter(ch):
		val := lex.scanAlpha()
		tok, ok := lookup[val]
		if ok {
			return tok, val
		}
		return INVALID, val

	case ch == -1:
		return EOF, ""

	default:
		lex.next()
		return INVALID, string(ch)
	}
}
