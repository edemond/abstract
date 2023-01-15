package types

import (
	"testing"
)

func TestSimplePartLength44Time(t *testing.T) {
	part := NewSimplePart()
	part.Rhythm.Meter = &Meter{Beats: 4, Value: 4}
	length := part.Length(64)
	if length != 256 {
		t.Fatalf("expected 64 PPQ * 4/4 time == 256 ticks, got %v", length)
	}
}

func TestSimplePartLength22Time(t *testing.T) {
	part := NewSimplePart()
	part.Rhythm.Meter = &Meter{Beats: 2, Value: 2}
	length := part.Length(64)
	if length != 256 {
		t.Fatalf("expected 64 PPQ * 2/2 time == 256 ticks, got %v", length)
	}
}

func TestSimplePartLength54Time(t *testing.T) {
	part := NewSimplePart()
	part.Rhythm.Meter = &Meter{Beats: 5, Value: 4}
	length := part.Length(64)
	if length != (64 * 5) {
		t.Fatalf("expected 64 PPQ * 2/2 time == 256 ticks, got %v", length)
	}
}

func TestSimplePartLength58Time(t *testing.T) {
	part := NewSimplePart()
	part.Rhythm.Meter = &Meter{Beats: 5, Value: 8}
	length := part.Length(64)
	if length != (32 * 5) {
		t.Fatalf("expected 64 PPQ * 5/4 time == 160 ticks, got %v", length)
	}
}
