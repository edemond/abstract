// Package ast defines the types that can appear in the abstract syntax tree of an Abstract program.
// It uses marker interfaces to get a bit more type safety similar to sum/union types.
package ast

import (
	"fmt"
	"strings"
)

// Every AST node is either a statement or an expression.

// Declarations or imperative actions (i.e. to play something.)
type Statement interface {
	isStatement()
	String() string
}

// Abstraction over either a single or compound expression.
type Expression interface {
	isExpression()
	String() string
}

// Statements ---------------

// A binding of an expression to a name.
type LetStatement struct {
	Name string
	Expr Expression
	Line int
}

func (let *LetStatement) String() string {
	paramExpr, ok := let.Expr.(Parameterized)
	if ok && paramExpr.HasParameters() {
		params := []string{}
		for _, p := range paramExpr.Parameters() {
			params = append(params, string(p))
		}
		return fmt.Sprintf("let %v(%v) = %v", let.Name, strings.Join(params, ", "), let.Expr)
	}
	return fmt.Sprintf("let %v = %v", let.Name, let.Expr)
}

func NewLetStatement(name string, expr Expression, line int) *LetStatement {
	return &LetStatement{
		Name: name,
		Expr: expr,
		Line: line,
	}
}

// A statement that sets a default musical context for expressions that don't specify everything.
type DefaultStatement struct {
	Expr *SimpleExpr
	Line int
}

func (def *DefaultStatement) String() string {
	return fmt.Sprintf("default %v", def.Expr)
}

// An expression to be played.
type PlayStatement struct {
	Expr Expression
	Line int
}

func (p *PlayStatement) String() string {
	return p.Expr.String()
}

// Sets the BPM.
type BPMStatement struct {
	BPM  int
	Line int
}

func (b *BPMStatement) String() string {
	return fmt.Sprintf("bpm %v", b.BPM)
}

// Sets the PPQ.
type PPQStatement struct {
	PPQ  int
	Line int
}

func (p *PPQStatement) String() string {
	return fmt.Sprintf("ppq %v", p.PPQ)
}

// Expressions ---------------

// Identifiers: Variable references or chord symbols. e.g. C, dorian, bass, iii7
type IdentExpr string

func (i IdentExpr) String() string {
	return string(i)
}

type StringExpr string

func (s StringExpr) String() string {
	return string(s)
}

type NumberExpr struct {
	Value  uint64
	Digits int // Significant when we're treating a number like a bit pattern.
	Line   int
}

func (n *NumberExpr) String() string {
	return fmt.Sprintf("%v", n.Value) // TODO: Digits
}

// "Functions", essentially. arp(32), rhythm(5439, 234783), over(dorian), blah("whatever", 12)
type ParamExpr struct {
	Name   string
	Params []Expression
	Line   int
}

func (p *ParamExpr) String() string {
	params := []string{}
	for _, param := range p.Params {
		params = append(params, param.String())
	}
	return fmt.Sprintf("%v(%v)", p.Name, strings.Join(params, ", "))
}

// Syntax sugar for meter (e.g. 6/8 for meter(6,8)).
type MeterExpr struct {
	// TODO: It'd be interesting to allow any expression here, that way you could do:
	// let beats = 2
	// let meter = beats/4
	Beats *NumberExpr
	Value *NumberExpr
	Line  int
}

func (m *MeterExpr) String() string {
	return fmt.Sprintf("%v/%v", m.Beats, m.Value)
}

// Any kind of expression that can be parameterized.
type Parameterized interface {
	Expression
	Parameters() []IdentExpr
	HasParameters() bool
	AddParameter(IdentExpr)
}

type SimpleExpr struct {
	ValueExprs []Expression
	params     []IdentExpr
	Line       int
}

type CompoundExpr struct {
	SimpleExprs []*SimpleExpr
	params      []IdentExpr
	Line        int
}

type BlockExpr struct {
	Statements []Statement
	params     []IdentExpr
	Line       int
}

type SeqExpr struct {
	ValueExprs []Expression
	Line       int
}

func (b *BlockExpr) String() string {
	stmts := []string{}
	for _, s := range b.Statements {
		stmts = append(stmts, "  "+s.String())
	}
	return fmt.Sprintf("{\n%v\n}", strings.Join(stmts, "\n"))
}

func (e *SimpleExpr) AddParameter(param IdentExpr) {
	e.params = append(e.params, param)
}
func (e *CompoundExpr) AddParameter(param IdentExpr) {
	e.params = append(e.params, param)
}
func (e *BlockExpr) AddParameter(param IdentExpr) {
	e.params = append(e.params, param)
}

func (e *SimpleExpr) HasParameters() bool {
	return e.params != nil && len(e.params) > 0
}
func (e *CompoundExpr) HasParameters() bool {
	return e.params != nil && len(e.params) > 0
}
func (e *BlockExpr) HasParameters() bool {
	return e.params != nil && len(e.params) > 0
}

func (e *SimpleExpr) Parameters() []IdentExpr {
	return e.params
}
func (e *CompoundExpr) Parameters() []IdentExpr {
	return e.params
}
func (e *BlockExpr) Parameters() []IdentExpr {
	return e.params
}

func (e *SimpleExpr) String() string {
	exprs := []string{}
	for _, e := range e.ValueExprs {
		exprs = append(exprs, e.String())
	}
	return strings.Join(exprs, " ")
}

func (e *CompoundExpr) String() string {
	exprs := []string{}
	for _, e := range e.SimpleExprs {
		exprs = append(exprs, e.String())
	}
	return strings.Join(exprs, " | ")
}

func (e *SeqExpr) String() string {
	exprs := []string{}
	for _, e := range e.ValueExprs {
		exprs = append(exprs, e.String())
	}
	return fmt.Sprintf("[%v]", strings.Join(exprs, " "))
}

func (s *LetStatement) isStatement()     {}
func (s *DefaultStatement) isStatement() {}
func (s *PlayStatement) isStatement()    {}
func (s *BPMStatement) isStatement()     {}
func (s *PPQStatement) isStatement()     {}

func (e *SimpleExpr) isExpression()   {}
func (e *CompoundExpr) isExpression() {}
func (e *BlockExpr) isExpression()    {}
func (e *SeqExpr) isExpression()      {}
func (e IdentExpr) isExpression()     {}
func (e *ParamExpr) isExpression()    {}
func (s StringExpr) isExpression()    {}
func (n *NumberExpr) isExpression()   {}
func (m *MeterExpr) isExpression()    {}
