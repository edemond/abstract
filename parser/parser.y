%{
package parser

import(
	"edemond/abstract/lexer"
	"edemond/abstract/ast"
	//"fmt"
	//"strconv"
	//"strings"
)

const PARSER_TRACE = false
%}

/* Reminder: this declares the contents of the generated "yylval" (abSymType). */
%union {
    val string 
    statement ast.Statement
    expr ast.Expression
    simpleexpr *ast.SimpleExpr
    compoundexpr *ast.CompoundExpr
    blockexpr *ast.BlockExpr
    paramexpr *ast.ParamExpr
    exprlist []ast.Expression
}

/*
    If these change, have to change the order in lexer/lex.go.
    TODO: Don't like this duplication. Can we have the lexer just use 
    stuff from the parser? Probably not because that'd be a 
    circular dependency?
*/
/*
    TODO: Why do some of these need to have types/fields in the union?
*/
%token <val> IDENT 2
%token <val> NUMBER 3
%token <val> STRING 4
%token '=' 5
%token '(' 6
%token ')' 7
%token <val> '{' 8
%token '}' 9
%token '|' 10
%token ',' 11
%token '\n' 12
%token '/' 13
%token LET 14
%token DEFAULT 15
%token BPM 16
%token PPQ 17
%token LBRACKET '[' 18
%token RBRACKET ']' 19

%type <val> error
%type <statement> piece 
%type <blockexpr> statementlist 
%type <statement> statement
%type <statement> bpmstatement
%type <statement> ppqstatement
%type <statement> defaultstatement
%type <statement> letstatement
%type <statement> playstatement
%type <simpleexpr> simpleexpr
%type <compoundexpr> compoundexpr
%type <expr> valueexpr
%type <expr> expr
%type <simpleexpr> simpleorvalueexpr
%type <exprlist> valueexprlist
%type <exprlist> exprlist
%type <exprlist> formalparameterlist
%type <paramexpr> paramexpr

%%

piece : statementlist
{
    trace("Parsed a piece.\n")
    stmt := &ast.PlayStatement{
        Line: ablex.(*abLexerImpl).Line(),
        Expr: $1,
    }
    $$ = stmt
    // Stash the finished piece where we can get it after parsing.
    ablex.(*abLexerImpl).parseResult = stmt
}

statementlist : statement
{
    trace("Parsed a statement: %v\n", $1)
    $$ = &ast.BlockExpr{
        Line: ablex.(*abLexerImpl).Line(),
        Statements: []ast.Statement{$1},
    }
}
    | statementlist statement
{
    trace("Parsed a statement list (more): %v\n", $2)
    $1.Statements = append($1.Statements, $2)
    $$ = $1
}

statement : bpmstatement
    | ppqstatement
    | defaultstatement
    | letstatement
    | playstatement
    ;

terminator : '\n'
    | error '\n'
{
    ablex.Error("Expected end of statement.")
}

playstatement : expr terminator
{
    stmt := &ast.PlayStatement{
        Line: ablex.(*abLexerImpl).Line(),
        Expr: $1,
    }
    $$ = stmt
    trace("Parsed a play statement: %v\n", stmt)
}

bpmstatement : BPM NUMBER terminator
{
	bpm, err := strconv.ParseUint($2, 10, 64)
	if err != nil {
        ablex.Error(fmt.Sprintf("bad number format: %v", $2))
	}

    stmt := &ast.BPMStatement{
        Line: ablex.(*abLexerImpl).Line(),
        BPM: int(bpm),
    }
    $$ = stmt
    trace("Parsed a BPM statement: %v\n", stmt)
}

ppqstatement : PPQ NUMBER terminator
{
	ppq, err := strconv.ParseUint($2, 10, 64)
	if err != nil {
        ablex.Error(fmt.Sprintf("bad number format: %v", $2))
	} else {
        stmt := &ast.PPQStatement{
            Line: ablex.(*abLexerImpl).Line(),
            PPQ: int(ppq),
        }
        trace("Parsed a PPQ statement: %v\n", stmt)
        $$ = stmt
    }
}

defaultstatement : DEFAULT valueexpr terminator
{
    simple := &ast.SimpleExpr{
        ValueExprs: []ast.Expression{$2},
    }
    def := &ast.DefaultStatement{
        Line: ablex.(*abLexerImpl).Line(),
        Expr: simple,
    }
    $$ = def
    trace("Parsed a default statement: %v\n", def)
}
    | DEFAULT simpleexpr '\n'
{
    def := &ast.DefaultStatement{
        Line: ablex.(*abLexerImpl).Line(),
        Expr: $2,
    }
    $$ = def
    trace("Parsed a default statement: %v\n", def)
}

expr : '{' statementlist '}'
{
    trace("Parsed a block expression: %v\n", $2)
    $$ = $2
}
    | '{' statementlist error
{
    fmt.Printf("expected '}'\n")
}
    | '{' error
{
    fmt.Printf("expected statements after '%v'\n", $1)
}
    | simpleexpr
{
    $$ = $1
}
    | compoundexpr
{
    $$ = $1
}
    | valueexpr
{
    $$ = $1
}

letstatement : LET IDENT '=' expr terminator
{
    stmt := ast.NewLetStatement($2, $4, ablex.(*abLexerImpl).Line())
    trace("Parsed a let statement: %v\n", stmt)
    $$ = stmt
}
    | LET IDENT '(' formalparameterlist ')' '=' expr terminator
{
    params := $4
    expr := $7

	if params != nil && len(params) > 0 {
		// We only need to add this parameter list to a Simple, Compound, or Block expr. 
		// Everything else can't have parameters.......or can it? TODO
		e, ok := expr.(ast.Parameterized)
		if !ok {
			ablex.Error("expected simple, compound, or block expression")
		} else {
            for _,prm := range params {
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

    stmt := ast.NewLetStatement($2, expr, ablex.(*abLexerImpl).Line()) 
    $$ = stmt
    trace("Parsed a let statement (with params): %v\n", stmt)
}
    | LET error terminator
{
    trace("error in let statement")
}

valueexpr : IDENT
{
    $$ = ast.IdentExpr($1)
    trace("Parsed an ident value expression: %v\n", $1)
}
    | STRING
{
    $$ = ast.StringExpr($1)
    trace("Parsed a string value expression: %v\n", $1)
}
    | paramexpr
{
    $$ = $1
    trace("Parsed a parameterized value expression: %v\n", $1)
}
    | NUMBER '/' NUMBER
{
    line := ablex.(*abLexerImpl).Line()
    beats, bdigits, err := convertNumber($1)
    if err != nil {
        ablex.Error(fmt.Sprintf("%v", err)) // TODO bleh 
    } else {
        value, vdigits, err := convertNumber($3)
        if err != nil {
            ablex.Error(fmt.Sprintf("%v", err)) // TODO bleh
        } else {
            expr := &ast.MeterExpr{
                Line: line,
                Beats: &ast.NumberExpr{
                    Line: line,
                    Value: beats,
                    Digits: bdigits,
                },
                Value: &ast.NumberExpr{
                    Line: line,
                    Value: value,
                    Digits: vdigits,
                },
            }
            $$ = expr
            trace("Parsed a meter expression: %v\n", expr)
        }
    }
}
    | NUMBER
{
    num, digits, err := convertNumber($1)
    if err != nil {
        ablex.Error(fmt.Sprintf("%v", err)) // TODO bleh
    } else {
        $$ = &ast.NumberExpr{
            Line: ablex.(*abLexerImpl).Line(),
            Value: num, 
            Digits: digits,
        }
    }
    trace("Parsed a number value expression: %v\n", $1)
}
    | '[' valueexprlist ']'
{
    expr := &ast.SeqExpr{
        Line: ablex.(*abLexerImpl).Line(),
        ValueExprs: $2,
    }
    $$ = expr
    trace("Parsed a sequence expression: %v\n", expr)
}

valueexprlist : valueexpr
{
    exprs := []ast.Expression{$1}
    $$ = exprs
    trace("Parsed a value expression list: %v\n", exprs)
}
    | valueexprlist valueexpr
{
    $1 = append($1, $2)
    $$ = $1
    trace("Parsed a value expression list (more): %v\n", $1)
}

paramexpr : IDENT '(' exprlist ')'
{
    expr := &ast.ParamExpr{
        Line: ablex.(*abLexerImpl).Line(),
        Name: $1,
        Params: $3,
    }
    $$ = expr
    trace("Parsed a parameterized expression: %v\n", expr)
}

exprlist : expr
{
    expr := []ast.Expression{$1}
    $$ = expr
    trace("Parsed an expression list (start): %v\n", expr)
}
    | exprlist ',' expr
{
    $1 = append($1, $3)
    $$ = $1
    trace("Parsed an expression list (more): %v\n", $1)
}

simpleorvalueexpr : simpleexpr
{
    $$ = $1
}
    | valueexpr
{
	expr := &ast.SimpleExpr{
        ValueExprs: []ast.Expression{$1},
    }
    $$ = expr
    trace("Upgraded a value expr to a simple expression: %v\n", expr)
}

compoundexpr : simpleorvalueexpr '|' simpleorvalueexpr
{
    expr := &ast.CompoundExpr{
        Line: ablex.(*abLexerImpl).Line(),
        SimpleExprs: []*ast.SimpleExpr{$1, $3},
    }
    $$ = expr
    trace("Parsed a compound expression: %v\n", expr)
}
    | compoundexpr '|' simpleorvalueexpr
{
    $1.SimpleExprs = append($1.SimpleExprs, $3)
    $$ = $1
    trace("Parsed a compound expression (more): %v\n", $1)
}

simpleexpr : valueexpr valueexpr
{
	expr := &ast.SimpleExpr{
        ValueExprs: []ast.Expression{$1, $2},
    }
    $$ = expr
    trace("Parsed a simple expression: %v\n", expr)
}
    | simpleexpr valueexpr
{
    $1.ValueExprs = append($1.ValueExprs, $2)
    $$ = $1
    trace("Parsed a simple expression (more): %v\n", $1)
}

formalparameterlist : IDENT
{
    exprs := []ast.Expression{
        ast.IdentExpr($1),
    }
    $$ = exprs
    trace("Parsed a formal parameter list: %v\n", $1)
}
    | formalparameterlist ',' IDENT
{
    $1 = append($1, ast.IdentExpr($3))
    $$ = $1
    trace("Parsed a formal parameter list (more): %v\n", $1)
}

%%

// Wrap a lexer.Lexer in a struct that implements abLexer.
// All of lexer.Lexer's methods are forwarded here.
type abLexerImpl struct {
    *lexer.Lexer 
    parseResult *ast.PlayStatement // The root of the parsed AST is stored here after parsing.
}

func (lex *abLexerImpl) Lex(yylval *abSymType) int {
	tok, val, err := lex.Scan()
    if err != nil {
        panic(err) // TODO: this is total garbage
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
