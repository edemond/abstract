// Package lexer implements a lexer for Abstract. It's heavily based
// on Go's "go/scanner".
package lexer

import (
	"github.com/edemond/abstract/chord"
	"fmt"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

type Token int

const (
	/*
	   If these change, have to change the order in parser/parser.y.
	*/
	EOF Token = iota
	INVALID
	IDENT // any identifier
	NUMBER
	STRING  // a double-quoted string
	ASSIGN  // =
	LPAREN  // (
	RPAREN  // )
	LBRACE  // {
	RBRACE  // }
	PIPE    // |
	COMMA   // ,
	NEWLINE // "\n"
	SLASH   // /

	// keywords
	LET     // let
	DEFAULT // default
	BPM     // bpm
	PPQ     // ppq

	LBRACKET // [
	RBRACKET // ]
)

func (t Token) String() string {
	switch t {
	case INVALID:
		return "(invalid token)"
	case EOF:
		return "(eof)"
	case IDENT:
		return "identifier"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case ASSIGN:
		return "="
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case PIPE:
		return "|"
	case COMMA:
		return ","
	case SLASH:
		return "/"
	case NEWLINE:
		return "newline"
	case LET:
		return "let"
	case DEFAULT:
		return "default"
	case BPM:
		return "bpm"
	case PPQ:
		return "ppq"
	default:
		panic("unknown token type")
	}
}

// Lookup to tell if a token is a keyword.
var keywords = map[string]Token{
	"let":     LET,
	"default": DEFAULT,
	"bpm":     BPM,
	"ppq":     PPQ,
}

type Lexer struct {
	source []byte
	start  int  // start position
	end    int  // end position
	line   int  // current line
	last   rune // last non-whitespace char (TODO: It'd be easier if this were a token.)
	char   rune
}

func FromFile(filename string) (*Lexer, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return FromBytes(source), nil
}

func FromBytes(src []byte) *Lexer {
	lex := &Lexer{
		source: src,
		start:  0,
		end:    0,
		line:   1,    // lines start at 1
		char:   ' ',  // lexer starts by skipping whitespace
		last:   '\n', // TODO hack to cause initial lines to be skipped
	}
	return lex
}

func (lex *Lexer) Line() int {
	return lex.line
}

// Format an error with the current line number.
func (lex *Lexer) errorf(err string, args ...interface{}) error {
	return fmt.Errorf("line %v: %v", lex.Line(), fmt.Sprintf(err, args...))
}

// next advances the lexer one rune. As in go/scanner, lex.char == -1 is EOF.
func (lex *Lexer) next() error {
	lex.start = lex.end
	if !isWhitespace(lex.char) {
		lex.last = lex.char
	}
	r, size := utf8.DecodeRune(lex.source[lex.start:])
	if size == 0 {
		lex.char = -1 // EOF.
	} else if size == 1 && r == utf8.RuneError {
		return lex.errorf("invalid UTF-8 character")
	} else {
		lex.end += size
		lex.char = r
	}
	return nil
}

func (lex *Lexer) skipWhitespace() error {
	for isWhitespace(lex.char) {
		if err := lex.next(); err != nil {
			return err
		}
	}
	return nil
}

func (lex *Lexer) scanIdent() (string, error) {
	start := lex.start
	// Chord notation is treated as an identifier, and is lexed/parsed/evaluated in the analyzer.
	for isLetter(lex.char) || isDigit(lex.char) || lex.char == '#' || lex.char == '_' || chord.IsChordNotationSymbol(lex.char) {
		if err := lex.next(); err != nil {
			return "", err
		}
	}
	return string(lex.source[start:lex.start]), nil
}

func (lex *Lexer) scanString() (string, error) {
	start := lex.start
	for lex.char != '"' && lex.char != -1 {
		if lex.char == '\n' {
			return "", lex.errorf("unexpected newline in string")
		}
		if err := lex.next(); err != nil {
			return "", err
		}
	}
	val := string(lex.source[start:lex.start])
	// Consume the closing ".
	if err := lex.next(); err != nil {
		return "", err
	}
	return val, nil
}

func (lex *Lexer) scanNumber() (string, error) {
	start := lex.start
	for isDigit(lex.char) {
		if err := lex.next(); err != nil {
			return "", err
		}
	}
	val := string(lex.source[start:lex.start])
	if strings.ContainsRune(val, 'x') && (!strings.HasPrefix(val, "0x") || (strings.Count(val, "x") > 1)) {
		return val, lex.errorf("invalid number: '%v'", val)
	}
	if !strings.HasPrefix(val, "0x") && strings.ContainsAny(val, "abcdefABCDEF") {
		return val, lex.errorf("invalid number: '%v'", val)
	}
	return val, nil
}

// Scan and skip a single-line (C-style) comment starting with // until the end of the line.
func (lex *Lexer) skipComment() error {
	for lex.char != '\n' && lex.char != -1 {
		if err := lex.next(); err != nil {
			return err
		}
	}
	return nil
}

// Scan and skip a potentially multi-line comment starting with /* and ending with */.
// The opening /* is consumed, so we just look for the closing */.
func (lex *Lexer) skipBlockComment() error {
	for lex.char != -1 {
		for lex.char != '*' && lex.char != -1 {
			if lex.char == '\n' {
				lex.line++ // We don't need to emit this newline to the parser; we're in a comment.
			}
			if err := lex.next(); err != nil {
				return err
			}
		}
		if err := lex.next(); err != nil {
			return err
		}
		if lex.char == '/' {
			if err := lex.next(); err != nil {
				return err
			}
			break
		}
		if err := lex.next(); err != nil {
			return err
		}
	}
	return nil
}

func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

// Any character that can appear in a numeric literal.
func isDigit(r rune) bool {
	return (r >= '0' && r <= '9') || r == 'x' || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

// All kinds of whitespace except newlines, which are syntactically significant
// as our statement terminator in the right context.
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

// getTokenType determines if a token is a keyword or a regular identifier.
func getTextTokenType(s string) Token {
	tok, ok := keywords[s]
	if ok {
		//fmt.Println("lexed a keyword", s)
		return tok
	} else {
		return IDENT
	}
}

func (lex *Lexer) shouldIgnoreNewline() bool {
	return lex.last == '{' || lex.last == '\n'
}

func (lex *Lexer) Scan() (tok Token, val string, err error) {
	keepGoing := true

	for keepGoing {
		keepGoing = false

		lex.skipWhitespace()

		// So that we can use newlines as the statement terminator but as whitespace
		// everywhere else, we do sort of the inverse of what BCPL and Go do with semicolons
		// (these languages insert semicolons into the token stream based on context.)
		// Instead, we selectively ignore newlines:
		// Newlines act as a statement terminator in some contexts (in which case we yield
		// them to the parser), but are ignored as whitespace the rest of the time.
		if lex.char == '\n' {
			// Found a newline, do we yield it (as the statement terminator) or skip it?
			lex.line++

			if !lex.shouldIgnoreNewline() {
				// Statement terminator.
				val = string(lex.char)
				lex.next()
				return NEWLINE, val, nil
			}

			for (lex.shouldIgnoreNewline() && lex.char == '\n') || isWhitespace(lex.char) {
				if err = lex.next(); err != nil {
					return INVALID, string(lex.char), err
				}
			}
		}

		switch ch := lex.char; {
		case ch == '/':
			last := lex.last
			// For now, we expect // or the start of a block comment, /*
			if err = lex.next(); err != nil {
				return INVALID, "", err
			}
			if lex.char == '/' {
				if err = lex.skipComment(); err != nil {
					return INVALID, "", err
				}
				// TODO: This is an ugly hack. We patch up the lexer with an old "last"
				// character here to avoid emitting extra newlines after a comment.
				lex.last = last
				keepGoing = true
			} else if lex.char == '*' {
				if err = lex.skipBlockComment(); err != nil {
					return INVALID, "", err
				}
				// TODO: This is an ugly hack. We patch up the lexer with an old "last"
				// character here to avoid emitting extra newlines after a comment.
				lex.last = last
				keepGoing = true
			} else {
				return SLASH, string(ch), nil
			}
		case isLetter(ch) || chord.IsChordNotationSymbol(ch) || ch == '_':
			val, err = lex.scanIdent()
			tok = getTextTokenType(val)
			return tok, val, err
		case isDigit(ch):
			tok = NUMBER
			val, err = lex.scanNumber()
			if err != nil {
				return INVALID, val, err
			}
		case ch == '"':
			// Consume the opening ".
			if err = lex.next(); err != nil {
				return INVALID, val, err
			}
			tok = STRING
			val, err = lex.scanString()
		case ch == -1:
			return EOF, "", nil

		default:
			switch ch {
			case '=':
				tok = ASSIGN
			case '(':
				tok = LPAREN
			case ')':
				tok = RPAREN
			case '[':
				tok = LBRACKET
			case ']':
				tok = RBRACKET
			case '{':
				tok = LBRACE
			case '}':
				tok = RBRACE
			case ',':
				tok = COMMA
			case '|':
				tok = PIPE
			default:
				// TODO: Unsure about this case, why are we advancing the lexer here?
				// Why not just return INVALID?
				if err = lex.next(); err != nil {
					return INVALID, val, nil
				}
				fmt.Printf("lexed something invalid: %q", string(ch))
				return INVALID, string(ch), nil
			}
			val = string(ch)
			lex.next()
		}
	}

	return tok, val, err
}
