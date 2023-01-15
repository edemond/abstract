package main

import (
	"github.com/edemond/abstract/ast"
	"github.com/edemond/abstract/parser"
	"github.com/edemond/abstract/types"
	"testing"
)

func testParse(t *testing.T, text string) *ast.PlayStatement {
	parser, err := parser.FromBytes([]byte(text))
	main, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	return main
}

func TestBasicSeqPart(t *testing.T) {
	text := `[C F G C]
    `
	a := NewAnalyzer()
	stmt := testParse(t, text)
	_, err := a.Analyze(stmt)
	if err != nil {
		t.Fatal(err)
	}
	// TODO
}

func TestNestedSeqPart(t *testing.T) {
	// TODO
}

func TestDefaultPartWorks(t *testing.T) {
	a := NewAnalyzer()

	text := `default A
        @I
        `
	parsed := testParse(t, text)
	part, err := a.Analyze(parsed)
	if err != nil {
		t.Fatal(err)
	}
	simple, ok := part.(*types.SimplePart)
	if !ok {
		t.Fatalf("Expected block part to collapse to simple part, got %v", part)
	}
	pitch, err := types.LookUpPitch("A")
	if err != nil {
		t.Fatalf("Error in test: %v", err)
	}
	if simple.Harmony.Pitch != pitch {
		t.Fatalf("Expected pitch %v, got %v", pitch, simple.Harmony.Pitch)
	}
}

func TestSimplePartsCanCombine(t *testing.T) {
	// TODO: Test the internal API more directly, like this.
	/*
	   s1 := types.NewSimplePart()
	   s1.Rhythm.Dynamics = types.NewDynamics(60)

	   s2 := types.NewSimplePart()
	   s2.Harmony.Pitch = types.NewPitch(5)

	   s3 := // combine them
	*/
	a := NewAnalyzer()

	text := `let x = C# dorian
        let y = dynamics(31) 5/4
        x y
        `
	parsed := testParse(t, text)
	part, err := a.Analyze(parsed)
	if err != nil {
		t.Fatal(err)
	}
	simple, ok := part.(*types.SimplePart)
	if !ok {
		t.Fatalf("Expected block part to collapse to simple part, got %v", part)
	}
	pitch, err := types.LookUpPitch("C#")
	if err != nil {
		t.Fatalf("Error in test: %v", err)
	}
	if simple.Harmony.Pitch != pitch {
		t.Fatalf("expected C#, got %v", simple.Harmony.Pitch)
	}
	dorian := []int{2, 1, 2, 2, 2, 1, 2}
	if !simple.Harmony.Scale.HasSteps(dorian) {
		t.Fatalf("expected a dorian scale, got: %v", simple.Harmony.Scale)
	}
	center := 31
	if !simple.Rhythm.Dynamics.HasValue() {
		t.Fatalf("Expected dynamics to have value.")
	}
	if simple.Rhythm.Dynamics.Center != center {
		t.Fatalf("expected dynamics(31), got %v", simple.Rhythm.Dynamics)
	}
	if simple.Rhythm.Meter.Beats != 5 || simple.Rhythm.Meter.Value != 4 {
		t.Fatalf("expected 5/4 time, got %v", simple.Rhythm.Meter)
	}
}
