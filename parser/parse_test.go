package parser

import (
	"github.com/edemond/abstract/ast"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	parser, _ := FromBytes([]byte(`let x = y`))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse()
	}
}

func contains(items []ast.Expression, item ast.Expression) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func coreTests(t *testing.T, text string) *ast.BlockExpr {
	parser, err := FromBytes([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	play, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected parsing error: %v", err)
	}
	if play == nil {
		t.Fatal("parser didn't return anything")
	}
	if play.Expr == nil {
		t.Fatal("expected block")
	}
	block, ok := play.Expr.(*ast.BlockExpr)
	if !ok {
		t.Fatal("expected *ast.BlockExpr")
	}
	return block
}

func expectParseError(t *testing.T, text string) {
	parser, err := FromBytes([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	_, err = parser.Parse()
	if err == nil {
		t.Fatalf("expected parsing error")
	}
}

func TestLetWithIdentRHS(t *testing.T) {
	text := "let x = y\n"
	block := coreTests(t, text)
	if len(block.Statements) != 1 {
		t.Fatalf("expected 1 statement in block, got %v", len(block.Statements))
	}
	stmt, ok := block.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatal("expected *ast.LetStatement")
	}
	if stmt.Name != "x" {
		t.Fatal("expected name 'x' on LHS")
	}
	ident, ok := stmt.Expr.(ast.IdentExpr)
	if !ok {
		t.Fatal("expected ast.IdentExpr")
	}
	if string(ident) != "y" {
		t.Fatalf("expected ident expr 'y', got '%v'", ident)
	}
}

func TestLetWith2IdentsRHS(t *testing.T) {
	text := "let x = C# lydian\n"
	block := coreTests(t, text)
	if len(block.Statements) != 1 {
		t.Fatalf("expected 1 statement in block, got %v", len(block.Statements))
	}
	stmt, ok := block.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatal("expected *ast.LetStatement")
	}
	if stmt.Name != "x" {
		t.Fatal("expected name 'x' on LHS")
	}
	expr, ok := stmt.Expr.(*ast.SimpleExpr)
	if !ok {
		t.Fatal("expected *ast.SimpleExpr on RHS")
	}
	if len(expr.ValueExprs) != 2 {
		t.Fatalf("expected 2 value exprs in simple expr, got %v", len(expr.ValueExprs))
	}
	if !contains(expr.ValueExprs, ast.IdentExpr("C#")) {
		t.Fatalf("expected ident expr 'C#' in RHS")
	}
	if !contains(expr.ValueExprs, ast.IdentExpr("lydian")) {
		t.Fatalf("expected ident expr 'lydian' in RHS")
	}
}

func TestLetWithBlockRHS(t *testing.T) {
	text := `let barf = { 
        C# lydian 
    }
    `
	block := coreTests(t, text)
	if len(block.Statements) != 1 {
		t.Fatalf("expected 1 statement in block, got %v", len(block.Statements))
	}
	let, ok := block.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatal("expected *ast.LetStatement")
	}
	if let.Name != "barf" {
		t.Fatal("expected name 'barf' on LHS")
	}
	expr, ok := let.Expr.(*ast.BlockExpr)
	if !ok {
		t.Fatal("expected *ast.BlockExpr on RHS")
	}
	if len(expr.Statements) != 1 {
		t.Fatalf("expected 1 statement in block expr, got %v", len(expr.Statements))
	}
	play, ok := expr.Statements[0].(*ast.PlayStatement)
	if !ok {
		t.Fatal("expected ast.PlayStatement")
	}
	simple, ok := play.Expr.(*ast.SimpleExpr)
	if !ok {
		t.Fatal("expected *ast.SimpleExpr in block")
	}
	if len(simple.ValueExprs) != 2 {
		t.Fatalf("expected 2 value exprs in simple expr, got %v", len(simple.ValueExprs))
	}
	if !contains(simple.ValueExprs, ast.IdentExpr("C#")) {
		t.Fatalf("expected ident expr 'C#' in RHS")
	}
	if !contains(simple.ValueExprs, ast.IdentExpr("lydian")) {
		t.Fatalf("expected ident expr 'lydian' in RHS")
	}
}

func TestLetWithOctaveRHS(t *testing.T) {
	text := "let x = O3\n"
	block := coreTests(t, text)
	if len(block.Statements) != 1 {
		t.Fatalf("expected 1 statement in block, got %v", len(block.Statements))
	}
	stmt, ok := block.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatal("expected *ast.LetStatement")
	}
	if stmt.Name != "x" {
		t.Fatal("expected name 'x' on LHS")
	}
	ident, ok := stmt.Expr.(ast.IdentExpr)
	if !ok {
		t.Fatal("expected ast.IdentExpr on RHS")
	}
	if string(ident) != "O3" {
		t.Fatalf("expected ident expr 'O3', got '%v'", ident)
	}
}

func TestParamExpr(t *testing.T) {
	text := "pc(1, 2)\n"
	block := coreTests(t, text)
	if len(block.Statements) != 1 {
		t.Fatalf("expected 1 statement in block, got %v", len(block.Statements))
	}
	stmt, ok := block.Statements[0].(*ast.PlayStatement)
	if !ok {
		t.Fatal("expected *ast.PlayStatement")
	}
	expr, ok := stmt.Expr.(*ast.ParamExpr)
	if !ok {
		t.Fatalf("expected *ast.ParamExpr, got %v", stmt.Expr)
	}
	if expr.Name != "pc" {
		t.Fatalf("expected param expr named 'pc', got: '%v'", expr.Name)
	}
}

func TestSeqExpr(t *testing.T) {
	text := "[a b _]\n"
	block := coreTests(t, text)
	if len(block.Statements) != 1 {
		t.Fatalf("expected 1 statement in block, got %v", len(block.Statements))
	}
	stmt, ok := block.Statements[0].(*ast.PlayStatement)
	if !ok {
		t.Fatal("expected *ast.PlayStatement")
	}
	seq, ok := stmt.Expr.(*ast.SeqExpr)
	if !ok {
		t.Fatalf("expected *ast.SeqExpr, got %v", stmt.Expr)
	}
	a, ok := seq.ValueExprs[0].(ast.IdentExpr)
	if !ok {
		t.Fatalf("expected ast.IdentExpr, got %v", seq.ValueExprs[0])
	}
	if a.String() != "a" {
		t.Fatalf("expected first value to be 'a', got '%v'", a)
	}
	b, ok := seq.ValueExprs[1].(ast.IdentExpr)
	if !ok {
		t.Fatalf("expected ast.IdentExpr, got %v", seq.ValueExprs[1])
	}
	if b.String() != "b" {
		t.Fatalf("expected second value to be 'b', got '%v'", b)
	}
	rest, ok := seq.ValueExprs[2].(ast.IdentExpr)
	if !ok {
		t.Fatalf("expected ast.IdentExpr, got %v", seq.ValueExprs[2])
	}
	if rest.String() != "_" {
		t.Fatalf("expected third value to be '_', got '%v'", rest)
	}
}
