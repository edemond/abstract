//go:generate -command yacc go tool yacc
//go:generate yacc -o parser.go -p "ab" parser.y

package parser

import (
	"github.com/edemond/abstract/ast"
	"github.com/edemond/abstract/lexer"
	//"fmt"
	//"strconv"
	//"strings"
)

type Parser interface {
	Parse() (*ast.PlayStatement, error)
}

// FromFile creates a new parser for the given file.
func FromFile(filename string) (Parser, error) {
	lex, err := lexer.FromFile(filename)
	if err != nil {
		return nil, err
	}
	p := &generatedParser{}
	p.lex = lex
	return p, nil
}

// FromBytes creates a new parser for the given source string.
func FromBytes(src []byte) (Parser, error) {
	p := &generatedParser{}
	p.lex = lexer.FromBytes(src)
	return p, nil
}

/*
// handwrittenParser implements an ad-hoc recursive-descent parser for the Abstract language.
// TODO: This is the original parser and can be deprecated once we've got the generated
// parser working. It won't be maintained thereafter.
type handwrittenParser struct {
	lex *lexer.Lexer
	tok lexer.Token // current token
	val string // value of current token
	scope int // number of scopes we're indented
}
*/

/*
func (p *handwrittenParser) trace(s string, args ...interface{}) {
	if PARSER_TRACE {
		fmt.Println(s, args)
	}
}

func (p *handwrittenParser) openScope() {
	p.scope++
}

func (p *handwrittenParser) closeScope() {
	p.scope--
}

// Parse parses the given file into an AST representation.
func (p *handwrittenParser) Parse() (*ast.PlayStatement, error) {
	block, err := p.parseBlockContents()
	if err != nil {
		return nil, err
	}
	// The whole program is one big play statement.
	play := &ast.PlayStatement{Line: p.lex.Line(), Expr: block}
	return play, nil
}

// Format an error with the current line number.
func (p *handwrittenParser) errorf(err string, args ...interface{}) error {
	return fmt.Errorf("line %v: %v", p.lex.Line(), fmt.Sprintf(err, args...))
}

// Format an "expected this, got that" error.
func (p *handwrittenParser) expected(expect, got lexer.Token) error {
	return p.errorf("expected %v, got %v", expect, got)
}

func (p *handwrittenParser) parseBlockContents() (*ast.BlockExpr, error) {
	block := &ast.BlockExpr{
        Line: p.lex.Line(),
        Statements: []ast.Statement{},
    }

	stmt, err := p.parseStatement()
	for stmt != nil {
		block.Statements = append(block.Statements, stmt)
		stmt, err = p.parseStatement()
	}

	if err != nil {
		return nil, err
	}

	return block, nil
}

func (p *handwrittenParser) parseStatement() (ast.Statement, error) {
	// Skip newlines until we get to the start of a statement.
	for p.tok == lexer.NEWLINE {
		p.next()
	}

	// This is the token we see at the START of the statement.
	switch p.tok {
	case lexer.INVALID:
		return nil, p.errorf("invalid token: '%v'", p.val)
	case lexer.EOF:
		if p.scope > 0 {
			return nil, p.expected(lexer.RBRACE, lexer.EOF)
		}
		return nil, nil // no more statements
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.DEFAULT:
		return p.parseDefaultStatement()
	case lexer.BPM:
		return p.parseBPMStatement()
	case lexer.PPQ:
		return p.parsePPQStatement()
	case lexer.IDENT, lexer.NUMBER:
		// A simple or compound expression to be played.
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.PlayStatement{Line: p.lex.Line(), Expr: expr}, nil
	case lexer.LBRACE:
		// A block to be played.
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.PlayStatement{Line: p.lex.Line(), Expr: expr}, nil
	case lexer.RBRACE:
		if p.scope <= 0 {
			return nil, p.errorf("unexpected '%v'", p.tok)
		}
		// Close the current expression. No more statements.
		p.closeScope()
		return nil, nil
	default:
		return nil, p.errorf("unexpected '%v'", p.tok)
	}
}

func (p *handwrittenParser) parseLetStatement() (ast.Statement, error) {
	// "let" is current token, expecting [ident] = [expr] remaining
	p.next()
	if p.tok != lexer.IDENT {
		return nil, p.expected(lexer.IDENT, p.tok)
	}
	ident := p.val // Save the identifier for later.
	line := p.lex.Line()

	p.next()

	// Here we might see a list of parameters, or we might not.
	var params *ast.ParamExpr
	var err error
	if p.tok == lexer.LPAREN {
		p.next()
		params, err = p.parseFormalParameterList(ident)
		if err != nil {
			return nil, err
		}
		// Save these values for later and cram them into the expression that follows.
	}

	if p.tok != lexer.ASSIGN {
		return nil, p.expected(lexer.ASSIGN, p.tok)
	}

	p.next()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if params != nil && len(params.Params) > 0 {
		// We only need to add this parameter list to a Simple, Compound, or Block expr.
		// Everything else can't have parameters.......or can it? TODO
		e, ok := expr.(ast.Parameterized)
		if !ok {
			return nil, p.errorf("expected simple, compound, or block expression")
		}

		for _,prm := range params.Params {
			param, ok := prm.(ast.IdentExpr)
			if !ok {
				return nil, p.errorf("expected all identifiers in parameter list")
			}
			e.AddParameter(param)
		}
	}

	p.trace("Parsed a let statement.")
	return ast.NewLetStatement(ident, expr, line), nil
}

func (p *handwrittenParser) parseDefaultStatement() (ast.Statement, error) {
	// "default" is current token, expecting [simpleExpr]
	p.next()

	expr, err := p.parseSimpleExpr()
	if err != nil {
		return nil, err
	}

	p.trace("Parsed a default statement.")
	return &ast.DefaultStatement{
        Line: p.lex.Line(),
        Expr: expr,
    }, nil
}

func (p *handwrittenParser) parseBPMStatement() (ast.Statement, error) {
	// "bpm" is current token, expecting number
	p.next()
	if (p.tok != lexer.NUMBER) {
		return nil, p.expected(lexer.NUMBER, p.tok)
	}
	p.trace("parsing bpm value of", p.val)
	bpm, err := strconv.ParseUint(p.val, 10, 64)
	if err != nil {
		return nil, p.errorf("bad number format: %v", p.val)
	}
	p.next()
	if (p.tok != lexer.NEWLINE) {
		return nil, p.expected(lexer.NEWLINE, p.tok)
	}
	return &ast.BPMStatement{Line: p.lex.Line(), BPM: int(bpm)}, nil
}

func (p *handwrittenParser) parsePPQStatement() (ast.Statement, error) {
	// "ppq" is current token, expecting number
	p.next()
	if (p.tok != lexer.NUMBER) {
		return nil, p.expected(lexer.NUMBER, p.tok)
	}

	ppq, err := strconv.ParseUint(p.val, 10, 64)
	if err != nil {
		return nil, p.errorf("bad number format: %v", p.val)
	}
	p.next()
	if (p.tok != lexer.NEWLINE) {
		return nil, p.expected(lexer.NEWLINE, p.tok)
	}
	return &ast.PPQStatement{Line: p.lex.Line(), PPQ: int(ppq)}, nil
}

func (p *handwrittenParser) parseExpr() (ast.Expression, error) {

	// First, are we in a block expression?
	if p.tok == lexer.LBRACE {
		p.openScope()
		return p.parseBlockExpr()
	}

	// Expect a single value expression or a simple expr.
	val, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}

	// If this is the end of the line, it's a loose value expr. Otherwise go on to make
	// it a simple or compound expr.
	if p.tok == lexer.NEWLINE || p.tok == lexer.EOF || p.tok == lexer.RBRACE ||
		p.tok == lexer.COMMA || p.tok == lexer.RPAREN {
		p.trace("Parsed a value expression.")
		return val, nil
	}

	// This is a simple expr, possibly part of a compound expr.
	p.trace("Upgrading to simple expression.")
	simple := &ast.SimpleExpr{Line: p.lex.Line(), ValueExprs: []ast.Expression{val}}

	// If it's not a pipe right away, there's more to the simple expr.
	if p.tok != lexer.PIPE {
        p.trace("token after parsing the thing is: %v", p.tok)
		// Expect a simple expression, possibly the first of a compound expr.
		simple, err = p.parseSimpleExpr()
		if err != nil {
			return nil, err
		}
		// Put back in the value we already consumed.
		simple.ValueExprs = append(simple.ValueExprs, val)
	}

	// Now, if we hit a pipe, it's a compound. If we hit a newline, we're done.
	if p.tok != lexer.PIPE {
		p.trace("Parsed a simple expression.")
		return simple, nil
	}

	p.trace("Upgrading to compound expression.")
	compound := &ast.CompoundExpr{Line: p.lex.Line(), SimpleExprs: []*ast.SimpleExpr{simple}}
	for p.tok == lexer.PIPE {
		p.next()
		simple, err = p.parseSimpleExpr()
		if err != nil {
			return nil, err
		}
		compound.SimpleExprs = append(compound.SimpleExprs, simple)
	}

	p.trace("Parsed a compound expression.")
	return compound, nil
}

// Parse a block expression surrounded by left and right braces.
func (p *handwrittenParser) parseBlockExpr() (*ast.BlockExpr, error) {
	// Opening { is current token, start a new scope, end it with closing }.

	p.next() // Consume the '{'.
	block, err := p.parseBlockContents()
	if err != nil {
		return nil, err
	}

	if p.tok != lexer.RBRACE {
		return nil, p.expected(lexer.RBRACE, p.tok)
	}
	p.next() // Consume the '}'.

	if len(block.Statements) <= 0 {
		p.trace("Parsed an empty block.")
		return block, nil
	}

	p.trace("Parsed a block.")
	return block, nil
}

func (p *handwrittenParser) parseSimpleExpr() (*ast.SimpleExpr, error) {
	expr := &ast.SimpleExpr{ValueExprs: []ast.Expression{}}

	for p.tok == lexer.IDENT || p.tok == lexer.NUMBER {
		val, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		expr.ValueExprs = append(expr.ValueExprs, val)
		p.trace("parsed a value")
	}

	// If we hit something other than this, we're done here.
	p.trace("Parsed a simple expression.")
	return expr, nil
}

// Parse a syntax-sugar meter expression (e.g. 4/4).
func (p *handwrittenParser) parseMeterExpr(beats string) (*ast.MeterExpr, error) {
    // The opening number (beats) has already been consumed.
    if p.tok != lexer.SLASH {
        return nil, p.expected(lexer.SLASH, p.tok)
    }
    p.next()

    if p.tok != lexer.NUMBER {
        return nil, p.expected(lexer.NUMBER, p.tok)
    }

    b, err := p.convertNumber(beats)
    if err != nil {
        return nil, err
    }
    v, err := p.convertNumber(p.val)
    if err != nil {
        return nil, err
    }
    p.next()

	p.trace("Parsed a meter syntax sugar expression.")
    return &ast.MeterExpr{
        Line: p.lex.Line(),
        Beats: b,
        Value: v,
    }, nil
}

// Parse a "value" expression (meaning we expect an identifier, number, or string in this position)
func (p *handwrittenParser) parseValueExpr() (ast.Expression, error) {
	if (p.tok != lexer.IDENT) && (p.tok != lexer.NUMBER) && (p.tok != lexer.STRING) {
		// TODO: parameterized exprs should be able to go here too
		return nil, p.errorf("expected identifier, number, or string, got %v", p.tok)
	}
	val := p.val

	// Check the next token and see if we need to upgrade to a parameterized value.
	lastTok := p.tok
	p.next()
	if lastTok == lexer.IDENT && p.tok == lexer.LPAREN {
		p.next()
		return p.parseParameterizedExpr(val)
	}

	switch lastTok {
	case lexer.IDENT:
		return ast.IdentExpr(val), nil
	case lexer.NUMBER:
        if p.tok == lexer.SLASH {
            return p.parseMeterExpr(val)
        }
		return p.convertNumber(val)
	case lexer.STRING:
		return ast.StringExpr(val), nil
	default:
		panic("error in parseValueExpr") // should never get here
	}
}

func (p *handwrittenParser) convertNumber(text string) (*ast.NumberExpr, error) {
	p.trace("Parsing a number expression.")
	if strings.HasPrefix(text, "0x") {
		// Hex number/bit pattern.
		value := text[2:] // strip the 0x
		num, err := strconv.ParseUint(value, 16, 64)
		if err != nil {
			return nil, p.errorf("bad hex number format: '%v'", value)
		}
		return &ast.NumberExpr{Line: p.lex.Line(), Value: num, Digits: len(value)}, nil
	} else {
		// Decimal number.
		num, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return nil, p.errorf("bad decimal number format: '%v'", text)
		}
		return &ast.NumberExpr{Line: p.lex.Line(), Value: num, Digits: len(text)}, nil
	}
}

// Parse a parameter list in a declaration (i.e. let statement) and returns it as an ast.ParamExpr.
// TODO: Does it make sense to represent  a list of formal parameters and a parameterized expression
// (i.e. function call) as the same type?
func (p *handwrittenParser) parseFormalParameterList(name string) (*ast.ParamExpr, error) {
	// ident and opening ( already consumed.
	param := &ast.ParamExpr{
        Line: p.lex.Line(),
		Name: name,
		Params: []ast.Expression{},
	}

	for {
		// A param can be any type of expression. It's got to be followed up by
		// a comma and another expression, or a right paren.

		if p.tok != lexer.IDENT {
			return nil, p.errorf("expected identifier, got %v", p.tok)
		}
		param.Params = append(param.Params, ast.IdentExpr(p.val))

		p.next()
		if p.tok == lexer.COMMA {
			p.next()
		} else if p.tok == lexer.RPAREN {
			p.next()
			p.trace("Parsed a parameter list.")
			return param, nil
		} else {
			// TODO: Better error here?
			return nil, p.errorf("unexpected %v ('%v')", p.tok, p.val)
		}
	}

	p.trace("Parsed a parameter list.")
	return param, nil
}

func (p *handwrittenParser) parseParameterizedExpr(name string) (*ast.ParamExpr, error) {
	p.trace("Parsing a parameterized expression.")
	// ident and opening ( already consumed.
	param := &ast.ParamExpr{
        Line: p.lex.Line(),
		Name: name,
		Params: []ast.Expression{},
	}

	// We expect an expression, followed by either an RPAREN, or a comma and another expression.
	for {
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		param.Params = append(param.Params, val)

		if p.tok == lexer.COMMA {
			p.next()
		} else if p.tok == lexer.RPAREN {
			p.next()
			p.trace("Parsed a parameterized value.")
			return param, nil
		} else {
			// TODO: Better error here?
			return nil, p.errorf("unexpected %v ('%v')", p.tok, p.val)
		}
	}

	p.trace("Parsed a parameterized value.")
	return param, nil
}

// Advance to the next non-comment token.
func (p *handwrittenParser) next() {
    var err error
	p.tok, p.val, err = p.lex.Scan()
    // TODO: this is total garbage
    if err != nil {
        panic(err)
    }
}
*/
