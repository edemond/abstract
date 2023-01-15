package lexer

import (
	"testing"
)

func lex(t *testing.T, text string) (Token, string, error) {
	lexer := FromBytes([]byte(text))
	return lexer.Scan()
}

func expect(t *testing.T, tok Token, val string, err error, etok Token, eval string) {
	// TODO: We might want to test for errors later.
	if err != nil {
		t.Fatalf("lexer encountered an error: %v", err)
	}
	if tok != etok {
		t.Fatalf("expected %v, got %v", etok, tok)
	}
	if val != eval {
		t.Fatalf("expected '%v', got '%v'", eval, val)
	}
}

func scanAndExpect(t *testing.T, lexer *Lexer, token Token, value string) {
	tok, val, err := lexer.Scan()
	expect(t, tok, val, err, token, value)
}

func TestLeadingNewlineSkipped(t *testing.T) {
	lexer := FromBytes([]byte("\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, EOF, "")
}

func TestMultipleLeadingNewlinesSkipped(t *testing.T) {
	lexer := FromBytes([]byte("\n\n\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, EOF, "")
}

func TestMultipleLeadingNewlinesAndCommentsSkipped(t *testing.T) {
	lexer := FromBytes([]byte("\n/*blah */\n\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, EOF, "")
}

func TestLeadingCommentThenMultipleNewlinesSkipped(t *testing.T) {
	lexer := FromBytes([]byte("/*blah */\n\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, EOF, "")
}

func TestLeadingSingleLineCommentSkipped(t *testing.T) {
	lexer := FromBytes([]byte("// blah\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, EOF, "")
}

func TestLeadingBlockCommentSkipped(t *testing.T) {
	lexer := FromBytes([]byte("/*blah*/\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, EOF, "")
}

func TestExtraNewlineFollowingBlockCommentNotEmitted(t *testing.T) {
	lexer := FromBytes([]byte("statement\n/* here's a comment */\n"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, EOF, "")
}

func TestExtraNewlineFollowingMultilineBlockCommentNotEmitted(t *testing.T) {
	lexer := FromBytes([]byte("statement\n/* here's a comment\nthat goes for two lines */\n"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, EOF, "")
}

func TestSingleLineCommentSkipped(t *testing.T) {
	lexer := FromBytes([]byte("statement // here's a comment\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, IDENT, "statement")
}

func TestBlockCommentSkipped(t *testing.T) {
	lexer := FromBytes([]byte("statement /* here's a comment */\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, IDENT, "statement")
}

func TestMultilineBlockCommentSkipped(t *testing.T) {
	lexer := FromBytes([]byte("statement /* here's a comment\nit goes for two lines */\nstatement"))
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, IDENT, "statement")
}

func TestNewlinesIgnoredWhenNotStatementTerminator(t *testing.T) {
	lexer := FromBytes([]byte("{\nstatement\n}"))
	scanAndExpect(t, lexer, LBRACE, "{")
	// Newline gets ignored here.
	scanAndExpect(t, lexer, IDENT, "statement")
	// This one follows a statement, so it's a statement terminator.
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, RBRACE, "}")
}

func TestMultipleNewlinesIgnoredWhenNotStatementTerminator(t *testing.T) {
	lexer := FromBytes([]byte("{\n\n\nstatement\n\n\n}"))
	scanAndExpect(t, lexer, LBRACE, "{")
	// Many newlines ignored here.
	scanAndExpect(t, lexer, IDENT, "statement")
	// Many newlines, but one of them is a statement terminator and gets yielded.
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, RBRACE, "}")
}

func TestAllKindsOfWhitespaceIgnoredAfterNewlines(t *testing.T) {
	lexer := FromBytes([]byte("{\n\t\n\r\t\nstatement\n\t\t\r\n\t\n\r}"))
	scanAndExpect(t, lexer, LBRACE, "{")
	// mucho whitespace ignored here
	scanAndExpect(t, lexer, IDENT, "statement")
	// also mucho whitespace, but one of them is a statement terminator and gets yielded.
	scanAndExpect(t, lexer, NEWLINE, "\n")
	scanAndExpect(t, lexer, RBRACE, "}")
}

func TestNoWhitespaceInBrackets(t *testing.T) {
	lexer := FromBytes([]byte("{statement}"))
	scanAndExpect(t, lexer, LBRACE, "{")
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, RBRACE, "}")
}

func TestOtherWhitespaceInBracketsIsSkipped(t *testing.T) {
	lexer := FromBytes([]byte("{ statement }"))
	scanAndExpect(t, lexer, LBRACE, "{")
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, RBRACE, "}")
}

func TestCrazyNonNewlineWhitespaceInBracketsIsSkipped(t *testing.T) {
	lexer := FromBytes([]byte("{\t\r  \t  \r statement\r \t }"))
	scanAndExpect(t, lexer, LBRACE, "{")
	scanAndExpect(t, lexer, IDENT, "statement")
	scanAndExpect(t, lexer, RBRACE, "}")
}

func TestOctave(t *testing.T) {
	tok, val, err := lex(t, `O3`)
	expect(t, tok, val, err, IDENT, "O3")
}

func TestLet(t *testing.T) {
	tok, val, err := lex(t, `let`)
	expect(t, tok, val, err, LET, "let")
}

func TestIdent(t *testing.T) {
	tok, val, err := lex(t, `mixolydian`)
	expect(t, tok, val, err, IDENT, "mixolydian")
}

func TestSharps(t *testing.T) {
	tok, val, err := lex(t, `C#`)
	expect(t, tok, val, err, IDENT, "C#")
}

func Test1DigitDecimalNumber(t *testing.T) {
	tok, val, err := lex(t, `1`)
	expect(t, tok, val, err, NUMBER, "1")
}

func Test5DigitDecimalNumber(t *testing.T) {
	tok, val, err := lex(t, `12345`)
	expect(t, tok, val, err, NUMBER, "12345")
}

func TestString(t *testing.T) {
	tok, val, err := lex(t, `"a string, dude"`)
	expect(t, tok, val, err, STRING, "a string, dude")
}

func TestRest(t *testing.T) {
	tok, val, err := lex(t, `_`)
	expect(t, tok, val, err, IDENT, "_")
}
