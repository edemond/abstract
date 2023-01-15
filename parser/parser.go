//line parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"github.com/edemond/abstract/ast"
	"github.com/edemond/abstract/lexer"
	"fmt"
	"strconv"
	"strings"
)

const PARSER_TRACE = false

//line parser.y:16
type abSymType struct {
	yys          int
	val          string
	statement    ast.Statement
	expr         ast.Expression
	simpleexpr   *ast.SimpleExpr
	compoundexpr *ast.CompoundExpr
	blockexpr    *ast.BlockExpr
	paramexpr    *ast.ParamExpr
	exprlist     []ast.Expression
}

const IDENT = 2
const NUMBER = 3
const STRING = 4
const LET = 14
const DEFAULT = 15
const BPM = 16
const PPQ = 17
const LBRACKET = 57353
const RBRACKET = 57354

var abToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"NUMBER",
	"STRING",
	"'='",
	"'('",
	"')'",
	"'{'",
	"'}'",
	"'|'",
	"','",
	"'\\n'",
	"'/'",
	"LET",
	"DEFAULT",
	"BPM",
	"PPQ",
	"LBRACKET",
	"'['",
	"RBRACKET",
	"']'",
}
var abStatenames = [...]string{}

const abEofCode = 1
const abErrCode = 2
const abInitialStackSize = 16

//line parser.y:407

// Wrap a lexer.Lexer in a struct that implements abLexer.
// All of lexer.Lexer's methods are forwarded here.
type abLexerImpl struct {
	*lexer.Lexer
	parseResult *ast.PlayStatement // The root of the parsed AST is stored here after parsing.
}

func (lex *abLexerImpl) Lex(yylval *abSymType) int {
	tok, val, err := lex.Scan()
	// TODO: this is total garbage
	if err != nil {
		panic(err)
	}
	yylval.val = val
	return int(tok)
}

func (lex *abLexerImpl) Error(e string) {
	fmt.Printf("line %v: %v\n", lex.Line(), e)
}

type generatedParser struct {
	lex *lexer.Lexer
}

func trace(format string, args ...interface{}) {
	if PARSER_TRACE {
		fmt.Printf(format, args...)
	}
}

func (p *generatedParser) Parse() (*ast.PlayStatement, error) {
	// Call the entry point of the yacc-generated parser.
	lex := &abLexerImpl{Lexer: p.lex}
	_ = abParse(lex) // TODO: Handle this return value?
	if lex.parseResult == nil {
		return nil, fmt.Errorf("Couldn't parse.") // TODO: actual error message here? filename?
	}
	trace("completed parse: %v\n", lex.parseResult)
	return lex.parseResult, nil
}

func convertNumber(text string) (num uint64, digits int, err error) {
	if strings.HasPrefix(text, "0x") {
		// Hex number/bit pattern.
		value := text[2:] // strip the 0x
		num, err := strconv.ParseUint(value, 16, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("bad hex number format: '%v'", value)
		}
		return num, len(value), nil
	} else {
		// Decimal number.
		num, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("bad decimal number format: '%v'", text)
		}
		return num, len(text), nil
	}
}

//line yacctab:1
var abExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 15,
	12, 36,
	-2, 19,
	-1, 17,
	12, 37,
	-2, 21,
}

const abNprod = 44
const abPrivate = 57344

var abTokenNames []string
var abStates []string

const abLast = 123

var abAct = [...]int{

	31, 13, 41, 51, 3, 17, 15, 24, 33, 18,
	2, 39, 69, 19, 22, 20, 70, 27, 28, 37,
	32, 36, 40, 38, 72, 34, 44, 45, 46, 43,
	23, 50, 61, 38, 36, 19, 22, 20, 60, 24,
	26, 14, 59, 56, 55, 56, 55, 54, 62, 57,
	63, 53, 23, 19, 22, 20, 25, 48, 49, 14,
	52, 36, 38, 73, 68, 12, 11, 9, 10, 71,
	23, 19, 22, 20, 74, 75, 35, 21, 19, 22,
	20, 47, 66, 30, 14, 29, 67, 65, 23, 64,
	12, 11, 9, 10, 58, 23, 19, 22, 20, 19,
	22, 20, 14, 33, 42, 19, 22, 20, 12, 11,
	9, 10, 16, 23, 8, 32, 23, 7, 6, 5,
	4, 1, 23,
}
var abPact = [...]int{

	92, -1000, 92, -1000, -1000, -1000, -1000, -1000, -1000, 51,
	35, 95, 81, 6, 74, 95, 7, 95, -1, 14,
	-1000, -1000, -13, 95, -1000, 6, 6, 101, 67, 50,
	6, -1000, -1000, -11, 49, -1000, -1000, 95, -1000, 95,
	31, 33, 9, -1000, -1000, -1000, -1000, -1000, 31, 83,
	-1000, -1000, -1000, -1000, -1000, 95, 95, -1000, 73, -1000,
	-1000, -1000, -1000, 6, 3, -1000, -1000, 31, -1000, 17,
	59, -1000, 31, -1000, 6, -1000,
}
var abPgo = [...]int{

	0, 121, 10, 4, 120, 119, 118, 117, 114, 6,
	112, 5, 1, 9, 104, 94, 89, 77, 0,
}
var abR1 = [...]int{

	0, 1, 2, 2, 3, 3, 3, 3, 3, 18,
	18, 8, 4, 5, 6, 6, 12, 12, 12, 12,
	12, 12, 7, 7, 7, 11, 11, 11, 11, 11,
	11, 14, 14, 17, 15, 15, 13, 13, 10, 10,
	9, 9, 16, 16,
}
var abR2 = [...]int{

	0, 1, 1, 2, 1, 1, 1, 1, 1, 1,
	2, 2, 3, 3, 3, 3, 3, 3, 2, 1,
	1, 1, 5, 8, 3, 1, 1, 1, 3, 1,
	3, 1, 2, 4, 1, 3, 1, 1, 3, 3,
	2, 2, 1, 3,
}
var abChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, 18,
	19, 17, 16, -12, 10, -9, -10, -11, -13, 4,
	6, -17, 5, 21, -3, 5, 5, -11, -9, 4,
	2, -18, 14, 2, -2, 2, -11, 12, -11, 12,
	8, 15, -14, -11, -18, -18, -18, 14, 7, 8,
	-18, 14, 11, 2, -13, -9, -11, -13, -15, -12,
	5, 23, -11, -12, -16, 4, 9, 13, -18, 9,
	13, -12, 7, 4, -12, -18,
}
var abDef = [...]int{

	0, -2, 1, 2, 4, 5, 6, 7, 8, 0,
	0, 0, 0, 0, 0, -2, 20, -2, 0, 25,
	26, 27, 29, 0, 3, 0, 0, 0, 0, 0,
	0, 11, 9, 0, 0, 18, 41, 0, 40, 0,
	0, 0, 0, 31, 12, 13, 14, 15, 0, 0,
	24, 10, 16, 17, 39, 36, 37, 38, 0, 34,
	28, 30, 32, 0, 0, 42, 33, 0, 22, 0,
	0, 35, 0, 43, 0, 23,
}
var abTok1 = [...]int{

	1, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 21, 23,
}
var abTok2 = [...]int{

	2, 3, 0, 0, 0, 0, 0, 0, 0, 20,
	22,
}
var abTok3 = [...]int{
	0,
}

var abErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	abDebug        = 0
	abErrorVerbose = false
)

type abLexer interface {
	Lex(lval *abSymType) int
	Error(s string)
}

type abParser interface {
	Parse(abLexer) int
	Lookahead() int
}

type abParserImpl struct {
	lval  abSymType
	stack [abInitialStackSize]abSymType
	char  int
}

func (p *abParserImpl) Lookahead() int {
	return p.char
}

func abNewParser() abParser {
	return &abParserImpl{}
}

const abFlag = -1000

func abTokname(c int) string {
	if c >= 1 && c-1 < len(abToknames) {
		if abToknames[c-1] != "" {
			return abToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func abStatname(s int) string {
	if s >= 0 && s < len(abStatenames) {
		if abStatenames[s] != "" {
			return abStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func abErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !abErrorVerbose {
		return "syntax error"
	}

	for _, e := range abErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + abTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := abPact[state]
	for tok := TOKSTART; tok-1 < len(abToknames); tok++ {
		if n := base + tok; n >= 0 && n < abLast && abChk[abAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if abDef[state] == -2 {
		i := 0
		for abExca[i] != -1 || abExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; abExca[i] >= 0; i += 2 {
			tok := abExca[i]
			if tok < TOKSTART || abExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if abExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += abTokname(tok)
	}
	return res
}

func ablex1(lex abLexer, lval *abSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = abTok1[0]
		goto out
	}
	if char < len(abTok1) {
		token = abTok1[char]
		goto out
	}
	if char >= abPrivate {
		if char < abPrivate+len(abTok2) {
			token = abTok2[char-abPrivate]
			goto out
		}
	}
	for i := 0; i < len(abTok3); i += 2 {
		token = abTok3[i+0]
		if token == char {
			token = abTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = abTok2[1] /* unknown char */
	}
	if abDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", abTokname(token), uint(char))
	}
	return char, token
}

func abParse(ablex abLexer) int {
	return abNewParser().Parse(ablex)
}

func (abrcvr *abParserImpl) Parse(ablex abLexer) int {
	var abn int
	var abVAL abSymType
	var abDollar []abSymType
	_ = abDollar // silence set and not used
	abS := abrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	abstate := 0
	abrcvr.char = -1
	abtoken := -1 // abrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		abstate = -1
		abrcvr.char = -1
		abtoken = -1
	}()
	abp := -1
	goto abstack

ret0:
	return 0

ret1:
	return 1

abstack:
	/* put a state and value onto the stack */
	if abDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", abTokname(abtoken), abStatname(abstate))
	}

	abp++
	if abp >= len(abS) {
		nyys := make([]abSymType, len(abS)*2)
		copy(nyys, abS)
		abS = nyys
	}
	abS[abp] = abVAL
	abS[abp].yys = abstate

abnewstate:
	abn = abPact[abstate]
	if abn <= abFlag {
		goto abdefault /* simple state */
	}
	if abrcvr.char < 0 {
		abrcvr.char, abtoken = ablex1(ablex, &abrcvr.lval)
	}
	abn += abtoken
	if abn < 0 || abn >= abLast {
		goto abdefault
	}
	abn = abAct[abn]
	if abChk[abn] == abtoken { /* valid shift */
		abrcvr.char = -1
		abtoken = -1
		abVAL = abrcvr.lval
		abstate = abn
		if Errflag > 0 {
			Errflag--
		}
		goto abstack
	}

abdefault:
	/* default state action */
	abn = abDef[abstate]
	if abn == -2 {
		if abrcvr.char < 0 {
			abrcvr.char, abtoken = ablex1(ablex, &abrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if abExca[xi+0] == -1 && abExca[xi+1] == abstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			abn = abExca[xi+0]
			if abn < 0 || abn == abtoken {
				break
			}
		}
		abn = abExca[xi+1]
		if abn < 0 {
			goto ret0
		}
	}
	if abn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			ablex.Error(abErrorMessage(abstate, abtoken))
			Nerrs++
			if abDebug >= 1 {
				__yyfmt__.Printf("%s", abStatname(abstate))
				__yyfmt__.Printf(" saw %s\n", abTokname(abtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for abp >= 0 {
				abn = abPact[abS[abp].yys] + abErrCode
				if abn >= 0 && abn < abLast {
					abstate = abAct[abn] /* simulate a shift of "error" */
					if abChk[abstate] == abErrCode {
						goto abstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if abDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", abS[abp].yys)
				}
				abp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if abDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", abTokname(abtoken))
			}
			if abtoken == abEofCode {
				goto ret1
			}
			abrcvr.char = -1
			abtoken = -1
			goto abnewstate /* try again in the same state */
		}
	}

	/* reduction by production abn */
	if abDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", abn, abStatname(abstate))
	}

	abnt := abn
	abpt := abp
	_ = abpt // guard against "declared and not used"

	abp -= abR2[abn]
	// abp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if abp+1 >= len(abS) {
		nyys := make([]abSymType, len(abS)*2)
		copy(nyys, abS)
		abS = nyys
	}
	abVAL = abS[abp+1]

	/* consult goto table to find next state */
	abn = abR1[abn]
	abg := abPgo[abn]
	abj := abg + abS[abp].yys + 1

	if abj >= abLast {
		abstate = abAct[abg]
	} else {
		abstate = abAct[abj]
		if abChk[abstate] != -abn {
			abstate = abAct[abg]
		}
	}
	// dummy call; replaced with literal code
	switch abnt {

	case 1:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:77
		{
			trace("Parsed a piece.\n")
			stmt := &ast.PlayStatement{
				Line: ablex.(*abLexerImpl).Line(),
				Expr: abDollar[1].blockexpr,
			}
			abVAL.statement = stmt
			// Stash the finished piece where we can get it after parsing.
			ablex.(*abLexerImpl).parseResult = stmt
		}
	case 2:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:89
		{
			trace("Parsed a statement: %v\n", abDollar[1].statement)
			abVAL.blockexpr = &ast.BlockExpr{
				Line:       ablex.(*abLexerImpl).Line(),
				Statements: []ast.Statement{abDollar[1].statement},
			}
		}
	case 3:
		abDollar = abS[abpt-2 : abpt+1]
		//line parser.y:97
		{
			trace("Parsed a statement list (more): %v\n", abDollar[2].statement)
			abDollar[1].blockexpr.Statements = append(abDollar[1].blockexpr.Statements, abDollar[2].statement)
			abVAL.blockexpr = abDollar[1].blockexpr
		}
	case 10:
		abDollar = abS[abpt-2 : abpt+1]
		//line parser.y:112
		{
			ablex.Error("Expected end of statement.")
		}
	case 11:
		abDollar = abS[abpt-2 : abpt+1]
		//line parser.y:117
		{
			stmt := &ast.PlayStatement{
				Line: ablex.(*abLexerImpl).Line(),
				Expr: abDollar[1].expr,
			}
			abVAL.statement = stmt
			trace("Parsed a play statement: %v\n", stmt)
		}
	case 12:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:127
		{
			bpm, err := strconv.ParseUint(abDollar[2].val, 10, 64)
			if err != nil {
				ablex.Error(fmt.Sprintf("bad number format: %v", abDollar[2].val))
			}

			stmt := &ast.BPMStatement{
				Line: ablex.(*abLexerImpl).Line(),
				BPM:  int(bpm),
			}
			abVAL.statement = stmt
			trace("Parsed a BPM statement: %v\n", stmt)
		}
	case 13:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:142
		{
			ppq, err := strconv.ParseUint(abDollar[2].val, 10, 64)
			if err != nil {
				ablex.Error(fmt.Sprintf("bad number format: %v", abDollar[2].val))
			} else {
				stmt := &ast.PPQStatement{
					Line: ablex.(*abLexerImpl).Line(),
					PPQ:  int(ppq),
				}
				trace("Parsed a PPQ statement: %v\n", stmt)
				abVAL.statement = stmt
			}
		}
	case 14:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:157
		{
			simple := &ast.SimpleExpr{
				ValueExprs: []ast.Expression{abDollar[2].expr},
			}
			def := &ast.DefaultStatement{
				Line: ablex.(*abLexerImpl).Line(),
				Expr: simple,
			}
			abVAL.statement = def
			trace("Parsed a default statement: %v\n", def)
		}
	case 15:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:169
		{
			def := &ast.DefaultStatement{
				Line: ablex.(*abLexerImpl).Line(),
				Expr: abDollar[2].simpleexpr,
			}
			abVAL.statement = def
			trace("Parsed a default statement: %v\n", def)
		}
	case 16:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:179
		{
			trace("Parsed a block expression: %v\n", abDollar[2].blockexpr)
			abVAL.expr = abDollar[2].blockexpr
		}
	case 17:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:184
		{
			fmt.Printf("expected '}'\n")
		}
	case 18:
		abDollar = abS[abpt-2 : abpt+1]
		//line parser.y:188
		{
			fmt.Printf("expected statements after '%v'\n", abDollar[1].val)
		}
	case 19:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:192
		{
			abVAL.expr = abDollar[1].simpleexpr
		}
	case 20:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:196
		{
			abVAL.expr = abDollar[1].compoundexpr
		}
	case 21:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:200
		{
			abVAL.expr = abDollar[1].expr
		}
	case 22:
		abDollar = abS[abpt-5 : abpt+1]
		//line parser.y:205
		{
			stmt := ast.NewLetStatement(abDollar[2].val, abDollar[4].expr, ablex.(*abLexerImpl).Line())
			trace("Parsed a let statement: %v\n", stmt)
			abVAL.statement = stmt
		}
	case 23:
		abDollar = abS[abpt-8 : abpt+1]
		//line parser.y:211
		{
			params := abDollar[4].exprlist
			expr := abDollar[7].expr

			if params != nil && len(params) > 0 {
				// We only need to add this parameter list to a Simple, Compound, or Block expr.
				// Everything else can't have parameters.......or can it? TODO
				e, ok := expr.(ast.Parameterized)
				if !ok {
					ablex.Error("expected simple, compound, or block expression")
				} else {
					for _, prm := range params {
						param, ok := prm.(ast.IdentExpr)
						if !ok {
							ablex.Error("expected only identifiers in parameter list")
						} else {
							e.AddParameter(param)
						}
					}
					expr = e
				}
			}

			stmt := ast.NewLetStatement(abDollar[2].val, expr, ablex.(*abLexerImpl).Line())
			abVAL.statement = stmt
			trace("Parsed a let statement (with params): %v\n", stmt)
		}
	case 24:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:239
		{
			trace("error in let statement")
		}
	case 25:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:244
		{
			abVAL.expr = ast.IdentExpr(abDollar[1].val)
			trace("Parsed an ident value expression: %v\n", abDollar[1].val)
		}
	case 26:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:249
		{
			abVAL.expr = ast.StringExpr(abDollar[1].val)
			trace("Parsed a string value expression: %v\n", abDollar[1].val)
		}
	case 27:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:254
		{
			abVAL.expr = abDollar[1].paramexpr
			trace("Parsed a parameterized value expression: %v\n", abDollar[1].paramexpr)
		}
	case 28:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:259
		{
			line := ablex.(*abLexerImpl).Line()
			beats, bdigits, err := convertNumber(abDollar[1].val)
			if err != nil {
				ablex.Error(fmt.Sprintf("%v", err)) // TODO bleh
			} else {
				value, vdigits, err := convertNumber(abDollar[3].val)
				if err != nil {
					ablex.Error(fmt.Sprintf("%v", err)) // TODO bleh
				} else {
					expr := &ast.MeterExpr{
						Line: line,
						Beats: &ast.NumberExpr{
							Line:   line,
							Value:  beats,
							Digits: bdigits,
						},
						Value: &ast.NumberExpr{
							Line:   line,
							Value:  value,
							Digits: vdigits,
						},
					}
					abVAL.expr = expr
					trace("Parsed a meter expression: %v\n", expr)
				}
			}
		}
	case 29:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:288
		{
			num, digits, err := convertNumber(abDollar[1].val)
			if err != nil {
				ablex.Error(fmt.Sprintf("%v", err)) // TODO bleh 
			} else {
				abVAL.expr = &ast.NumberExpr{
					Line:   ablex.(*abLexerImpl).Line(),
					Value:  num,
					Digits: digits,
				}
			}
			trace("Parsed a number value expression: %v\n", abDollar[1].val)
		}
	case 30:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:302
		{
			expr := &ast.SeqExpr{
				Line:       ablex.(*abLexerImpl).Line(),
				ValueExprs: abDollar[2].exprlist,
			}
			abVAL.expr = expr
			trace("Parsed a sequence expression: %v\n", expr)
		}
	case 31:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:312
		{
			exprs := []ast.Expression{abDollar[1].expr}
			abVAL.exprlist = exprs
			trace("Parsed a value expression list: %v\n", exprs)
		}
	case 32:
		abDollar = abS[abpt-2 : abpt+1]
		//line parser.y:318
		{
			abDollar[1].exprlist = append(abDollar[1].exprlist, abDollar[2].expr)
			abVAL.exprlist = abDollar[1].exprlist
			trace("Parsed a value expression list (more): %v\n", abDollar[1].exprlist)
		}
	case 33:
		abDollar = abS[abpt-4 : abpt+1]
		//line parser.y:325
		{
			expr := &ast.ParamExpr{
				Line:   ablex.(*abLexerImpl).Line(),
				Name:   abDollar[1].val,
				Params: abDollar[3].exprlist,
			}
			abVAL.paramexpr = expr
			trace("Parsed a parameterized expression: %v\n", expr)
		}
	case 34:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:336
		{
			expr := []ast.Expression{abDollar[1].expr}
			abVAL.exprlist = expr
			trace("Parsed an expression list (start): %v\n", expr)
		}
	case 35:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:342
		{
			abDollar[1].exprlist = append(abDollar[1].exprlist, abDollar[3].expr)
			abVAL.exprlist = abDollar[1].exprlist
			trace("Parsed an expression list (more): %v\n", abDollar[1].exprlist)
		}
	case 36:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:349
		{
			abVAL.simpleexpr = abDollar[1].simpleexpr
		}
	case 37:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:353
		{
			expr := &ast.SimpleExpr{
				ValueExprs: []ast.Expression{abDollar[1].expr},
			}
			abVAL.simpleexpr = expr
			trace("Upgraded a value expr to a simple expression: %v\n", expr)
		}
	case 38:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:362
		{
			expr := &ast.CompoundExpr{
				Line:        ablex.(*abLexerImpl).Line(),
				SimpleExprs: []*ast.SimpleExpr{abDollar[1].simpleexpr, abDollar[3].simpleexpr},
			}
			abVAL.compoundexpr = expr
			trace("Parsed a compound expression: %v\n", expr)
		}
	case 39:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:371
		{
			abDollar[1].compoundexpr.SimpleExprs = append(abDollar[1].compoundexpr.SimpleExprs, abDollar[3].simpleexpr)
			abVAL.compoundexpr = abDollar[1].compoundexpr
			trace("Parsed a compound expression (more): %v\n", abDollar[1].compoundexpr)
		}
	case 40:
		abDollar = abS[abpt-2 : abpt+1]
		//line parser.y:378
		{
			expr := &ast.SimpleExpr{
				ValueExprs: []ast.Expression{abDollar[1].expr, abDollar[2].expr},
			}
			abVAL.simpleexpr = expr
			trace("Parsed a simple expression: %v\n", expr)
		}
	case 41:
		abDollar = abS[abpt-2 : abpt+1]
		//line parser.y:386
		{
			abDollar[1].simpleexpr.ValueExprs = append(abDollar[1].simpleexpr.ValueExprs, abDollar[2].expr)
			abVAL.simpleexpr = abDollar[1].simpleexpr
			trace("Parsed a simple expression (more): %v\n", abDollar[1].simpleexpr)
		}
	case 42:
		abDollar = abS[abpt-1 : abpt+1]
		//line parser.y:393
		{
			exprs := []ast.Expression{
				ast.IdentExpr(abDollar[1].val),
			}
			abVAL.exprlist = exprs
			trace("Parsed a formal parameter list: %v\n", abDollar[1].val)
		}
	case 43:
		abDollar = abS[abpt-3 : abpt+1]
		//line parser.y:401
		{
			abDollar[1].exprlist = append(abDollar[1].exprlist, ast.IdentExpr(abDollar[3].val))
			abVAL.exprlist = abDollar[1].exprlist
			trace("Parsed a formal parameter list (more): %v\n", abDollar[1].exprlist)
		}
	}
	goto abstack /* stack new state and value */
}
