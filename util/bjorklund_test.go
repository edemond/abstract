package util

import (
	"testing"
)

func expect(t *testing.T, actual []int, expected []int) {
	if len(actual) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	}
}

func TestRotate(t *testing.T) {
	expect(t, rotate(1, []int{0, 1, 0, 0, 1, 0, 0, 1}), []int{1, 0, 0, 1, 0, 0, 1, 0})
	expect(t, rotate(2, []int{1, 0, 1, 0, 0, 1, 0, 0}), []int{1, 0, 0, 1, 0, 0, 1, 0})
}

func TestBjorklund(t *testing.T) {
	expect(t, Bjorklund(3, 8), []int{1, 0, 0, 1, 0, 0, 1, 0})
}
