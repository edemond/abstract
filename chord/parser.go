// This is both a parser and a semantic analyzer for chords, because honestly, why not
package chord

import (
	"github.com/edemond/abstract/types"
	"fmt"
	"strconv"
)

const PARSER_TRACE = false

type Parser struct {
	lex *Lexer
	tok Token  // current token
	val string // value of current token
}

func (p *Parser) trace(s string, args ...interface{}) {
	if PARSER_TRACE {
		fmt.Println(s, args)
	}
}

// Advance to the next non-comment token.
func (p *Parser) next() {
	p.tok, p.val = p.lex.Scan()
	p.trace("advanced to:", p.tok, p.val)
}

// NewParserFromString creates a new parser for the given source string.
// TODO: This seems like it'd create a lot of extra objects, can we reuse one chord parser? (Is it worth it? Measure first.)
func NewParserFromString(src string) (*Parser, error) {
	lex, err := NewLexerFromString(src)
	if err != nil {
		return nil, err
	}
	p := &Parser{}
	p.lex = lex
	p.tok = INVALID // doesn't matter what we start on
	return p, nil
}

func ParseAndAnalyze(text string) (types.Chord, error) {
	p, err := NewParserFromString(text)
	if err != nil {
		return nil, err
	}
	return p.ParseAndAnalyze()
}

// Parse parses and analyzes the given chord symbol.
func (p *Parser) ParseAndAnalyze() (types.Chord, error) {
	p.next() // start by advancing one step
	chord, err := p.parseChord()
	if err != nil {
		return nil, err
	}
	return Analyze(chord)
}

// Format an "expected this, got that" error.
func (p *Parser) expected(expect, got Token) error {
	return fmt.Errorf("expected %v, got %v", expect, got)
}

func (p *Parser) parseChord() (chordExpr, error) {
	p.trace("Parsing a chord symbol. -------------------------")
	switch p.tok {
	case DIATONIC:
		return p.parseDiatonicChord()
	case ROOT, FLAT, SHARP:
		return p.parseRelativeChord()
	case PITCH:
		return p.parseAbsoluteChord()
	}
	return nil, fmt.Errorf("Expected @, root, or pitch, got %v", p.tok)
}

// We enforce capital letters for diatonic chords, otherwise it'd be
// confusing if @I and @i were equivalent, etc.
func convertDiatonicRoot(text string) (degree int, err error) {
	switch text {
	case "Ⅰ":
	case "I":
		return 1, nil
	case "Ⅱ":
	case "II":
		return 2, nil
	case "Ⅲ":
	case "III":
		return 3, nil
	case "Ⅳ":
	case "IV":
		return 4, nil
	case "Ⅴ":
	case "V":
		return 5, nil
	case "Ⅵ":
	case "VI":
		return 6, nil
	case "Ⅶ":
	case "VII":
		return 7, nil
	}
	return 0, fmt.Errorf("invalid root in diatonic chord symbol (must be a capital Roman numeral): '%v'", text)
}

func convertRelativeRoot(text string) (degree int, quality Token, err error) {
	switch text {
	case "Ⅰ":
	case "I":
		return 1, MAJ, nil
	case "Ⅱ":
	case "II":
		return 2, MAJ, nil
	case "Ⅲ":
	case "III":
		return 3, MAJ, nil
	case "Ⅳ":
	case "IV":
		return 4, MAJ, nil
	case "Ⅴ":
	case "V":
		return 5, MAJ, nil
	case "Ⅵ":
	case "VI":
		return 6, MAJ, nil
	case "Ⅶ":
	case "VII":
		return 7, MAJ, nil
	case "ⅰ":
	case "i":
		return 1, MIN, nil
	case "ⅱ":
	case "ii":
		return 2, MIN, nil
	case "ⅲ":
	case "iii":
		return 3, MIN, nil
	case "ⅳ":
	case "iv":
		return 4, MIN, nil
	case "ⅴ":
	case "v":
		return 5, MIN, nil
	case "ⅵ":
	case "vi":
		return 6, MIN, nil
	case "ⅶ":
	case "vii":
		return 7, MIN, nil
	}
	return 0, INVALID, fmt.Errorf("invalid root in chord symbol: '%v'", text)
}

func (p *Parser) parseDiatonicChord() (*diatonicChordExpr, error) {
	p.trace("Parsing a diatonic chord.")

	// We're on a DIATONIC token, so just advance it one.
	p.next()

	// We expect a root (TODO: for now! Maybe the Nashville Number System
	// has ideas about how to notate extended diatonic chords?)
	if p.tok != ROOT {
		return nil, p.expected(ROOT, p.tok)
	}

	rootDegree, err := convertDiatonicRoot(p.val)
	if err != nil {
		return nil, err
	}
	p.next()

	// We now expect a list of qualities.
	qualities, err := p.parseAdditionalQualities()
	if err != nil {
		return nil, err
	}

	return &diatonicChordExpr{
		rootScaleDegree: rootDegree,
		qualities:       qualities,
	}, nil
}

func (p *Parser) parseRelativeChord() (*relativeChordExpr, error) {
	p.trace("Parsing a relative chord.")

	// We may have an accidental before the root (e.g. bV or #iii)
	accidental := 0
	if p.tok == SHARP {
		accidental = 1
		p.next()
	} else if p.tok == FLAT {
		accidental = -1
		p.next()
	}

	// We're on a ROOT token, which tells us both the scale degree (1-7) and quality (major, minor).
	rootDegree, quality, err := convertRelativeRoot(p.val)
	if err != nil {
		return nil, err
	}
	p.next()

	// TODO: sus2 and sus4 are cases we missed. Relative chords may imply major or minor in their
	// capitalization, but nothing rules out doing a sus2 or sus4!

	// TODO: Some additional qualities kind of override the main one implied by
	// the capitalization of the Roman numeral, e.g. diminished. We could skirt
	// this by saying that diminished just affects the 5th and MAJ/MIN doesn't?
	// actually I think that's correct, that's the only way.

	// Then we expect zero or more additional qualities.
	// (This advances the parser for us, no need for p.next() after.)
	additional, err := p.parseAdditionalQualities()
	if err != nil {
		return nil, err
	}

	qualities := []*qualityExpr{}
	qualities = append(qualities, &qualityExpr{quality: quality, interval: 0, implied: true})
	qualities = append(qualities, additional...)

	return &relativeChordExpr{
		rootScaleDegree: rootDegree,
		accidental:      accidental,
		qualities:       qualities,
	}, nil
}

// Parse a list of qualities.
func (p *Parser) parseAdditionalQualities() ([]*qualityExpr, error) {
	p.trace("Parsing additional qualities.")
	qualities := []*qualityExpr{}
	for p.tok != EOF {
		p.trace("Parsing a quality, because token isn't EOF, it's:", p.tok)
		q, err := p.parseQuality()
		if err != nil {
			return nil, err
		}
		qualities = append(qualities, q)
		p.next()
	}
	return qualities, nil
}

func (p *Parser) parseAbsoluteChord() (*absoluteChordExpr, error) {
	p.trace("Parsing an absolute chord.")

	// We're on a PITCH token.
	rootPitch, err := types.LookUpPitch(p.val)
	if err != nil {
		return nil, err
	}
	p.next()

	// For now, we expect at least one quality.
	quality, err := p.parseQuality()
	if err != nil {
		return nil, err
	}
	p.next()

	// Then we expect zero or more additional qualities.
	// (This advances the parser for us, no need for p.next() after.)
	additional, err := p.parseAdditionalQualities()
	if err != nil {
		return nil, err
	}

	qualities := []*qualityExpr{}
	qualities = append(qualities, quality)
	qualities = append(qualities, additional...)

	p.trace("Parsed an absolute chord.")
	return &absoluteChordExpr{
		pitch:     rootPitch,
		qualities: qualities,
	}, nil
}

// Returns whether or not the quality must be followed by an intervallic parameter.
// e.g. add; there's no Cadd. you need something like Cadd2.
func mustHaveParameter(quality Token) bool {
	switch quality {
	case ADD, NO, SLASH, FLAT, SHARP, MINMAJ:
		return true
	}
	return false
}

// Returns whether or not the quality can have an intervallic parameter or not.
// e.g. MAJ can be just maj or maj9, but we treat aug as a complete quality affecting
// the fifth, for practical purposes (aug7 is parsed as aug, 7).
func canHaveParameter(quality Token) bool {
	switch quality {
	case AUG: // TODO: DIM is a tricky case...dim7 is a thing, but not generally like dim6?
		return false
	}
	return true
}

func (p *Parser) parseQuality() (*qualityExpr, error) {
	// We expect a quality.
	p.trace("Parsing a quality.")

	// Certain numbers here can be read as a chord quality (e.g. C7, D13, for dominant, E5 for a power chord)
	if p.tok == NUMBER {
		switch p.val {
		case "5":
			return &qualityExpr{quality: POWER, interval: 5}, nil
		case "6":
			return &qualityExpr{quality: MAJ, interval: 6}, nil
		case "7":
			return &qualityExpr{quality: DOMINANT, interval: 7}, nil
		case "9":
			return &qualityExpr{quality: DOMINANT, interval: 9}, nil
		case "11":
			return &qualityExpr{quality: DOMINANT, interval: 11}, nil
		case "13":
			return &qualityExpr{quality: DOMINANT, interval: 13}, nil
		default:
			return nil, fmt.Errorf("%v is not a chord quality", p.val)
		}
	}
	var quality Token
	switch p.tok {
	case MAJ, MIN, MINMAJ, AUG, DIM, HALF_DIM, DOMINANT, ADD, SUS, NO, SHARP, FLAT:
		quality = p.tok
		p.trace("got quality:", quality)
	default:
		p.trace("found no quality :(")
		return nil, fmt.Errorf("expected a chord quality")
	}

	// Now, is this a parameterized quality?
	if canHaveParameter(quality) {
		p.next()
		if p.tok == NUMBER {
			num, err := strconv.ParseInt(p.val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad number format: '%v'", p.val)
			}
			p.trace("quality has a number; got", p.tok, p.val)
			return &qualityExpr{quality: quality, interval: int(num)}, nil
		} else {
			p.trace("next token wasn't a number; got", p.tok, p.val)
			if mustHaveParameter(quality) {
				return nil, fmt.Errorf("Quality '%v' must have an interval specified.", quality)
			}
			return &qualityExpr{quality: quality, interval: 0}, nil
		}
	}
	return &qualityExpr{quality: quality, interval: 0}, nil
}
