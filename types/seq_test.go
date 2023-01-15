package types

import (
	"testing"
)

func newSimplePartWithPitch(pitch uint64) *SimplePart {
	s := NewSimplePart()
	s.Harmony.Pitch = NewPitch(pitch)
	return s
}

func TestBasicSeqPart(t *testing.T) {
	parent := newSimplePartWithPitch(1)
	s1 := newSimplePartWithPitch(2)
	s2 := newSimplePartWithPitch(3)

	seq := NewSeqPart(
		[]Part{s1, s2},
		parent,
	)

	ppq := 16
	p1, _ := seq.at(uint64(ppq*0), ppq) // First beat
	p2, _ := seq.at(uint64(ppq*2), ppq) // Third beat
	if p1 != s1 {
		t.Fatalf("Expected first beat to be %v, got %v", s1, p1)
	}
	if p2 != s2 {
		t.Fatalf("Expected third beat to be %v, got %v", s2, p2)
	}
}

// TODO: Test fails because the code we really want is locked up in the Play method.
// Protocols should solve this so that we can test what playable object we get at each time step.
/*
func TestNestedSeqPart(t *testing.T) {
    parent := newSimplePartWithPitch(1)
    s1 := newSimplePartWithPitch(2)
    s2 := newSimplePartWithPitch(3)
    s3 := newSimplePartWithPitch(4)
    seq := NewSeqPart(
        []Part{
            s1,
            NewSeqPart(
                []Part{s2, s3},
                parent,
            ),
        },
        parent,
    )

    ppq := 16
    p1,_ := seq.at(uint64(ppq*0), ppq) // First beat
    p2,_ := seq.at(uint64(ppq*2), ppq) // Third beat
    p3,_ := seq.at(uint64(ppq*3), ppq) // Fourth beat
    if p1 != s1 {
        t.Fatalf("Expected first beat to be %v, got %v", s1, p1)
    }
    if p2 != s2 {
        t.Fatalf("Expected third beat to be %v, got %v", s2, p2)
    }
    if p3 != s3 {
        t.Fatalf("Expected fourth beat to be %v, got %v", s3, p3)
    }
}
*/
