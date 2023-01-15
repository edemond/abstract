package types

import (
	"testing"
)

func testScale(t *testing.T, scale *Scale, tests [][2]int) {
	for n, test := range tests {
		expected := test[1]
		if actual := scale.StepsAtDegree(test[0]); actual != expected { // root
			t.Fatalf("expected %v, got %v (%vth test case)", expected, actual, n)
		}
	}
}

func TestMajorScale(t *testing.T) {
	major := NewScale([]int{2, 2, 1, 2, 2, 2, 1})
	tests := [][2]int{
		{0, 0}, {1, 2}, {2, 4}, {3, 5}, {4, 7}, {5, 9}, {6, 11}, {7, 12}, {8, 14}, {14, 24},
		{-1, -1}, {-2, -3}, {-3, -5}, {-7, -12}, {-8, -13}, {-14, -24},
	}
	testScale(t, major, tests)
}

func TestMinorScale(t *testing.T) {
	minor := NewScale([]int{2, 1, 2, 2, 1, 2, 2})
	tests := [][2]int{
		{0, 0}, {1, 2}, {2, 3}, {3, 5}, {4, 7}, {5, 8}, {6, 10}, {7, 12}, {8, 14}, {14, 24},
		{-1, -2}, {-2, -4}, {-3, -5}, {-7, -12}, {-8, -14}, {-14, -24},
	}
	testScale(t, minor, tests)
}

func TestWholeToneScale(t *testing.T) {
	wt := NewScale([]int{2, 2, 2, 2, 2, 2})
	tests := [][2]int{
		{0, 0}, {1, 2}, {2, 4}, {3, 6}, {4, 8}, {5, 10}, {6, 12}, {7, 14}, {8, 16}, {14, 28},
		{-1, -2}, {-2, -4}, {-3, -6}, {-7, -14}, {-8, -16}, {-14, -28},
	}
	testScale(t, wt, tests)
}

func TestChromaticScale(t *testing.T) {
	chromatic := NewScale([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	tests := [][2]int{
		{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {14, 14},
		{-1, -1}, {-2, -2}, {-3, -3}, {-7, -7}, {-8, -8}, {-14, -14},
	}
	testScale(t, chromatic, tests)
}

func TestShorterChromaticScale(t *testing.T) {
	chromatic := NewScale([]int{1, 1, 1})
	tests := [][2]int{
		{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {14, 14},
		{-1, -1}, {-2, -2}, {-3, -3}, {-7, -7}, {-8, -8}, {-14, -14},
	}
	testScale(t, chromatic, tests)
}

func TestNegativeScale(t *testing.T) {
	negative := NewScale([]int{-2, -2, -1, -2, -2, -2, -1})
	tests := [][2]int{
		{0, 0}, {1, -2}, {2, -4}, {3, -5}, {4, -7}, {5, -9}, {6, -11}, {7, -12}, {8, -14}, {14, -24},
		{-1, 1}, {-2, 3}, {-3, 5}, {-7, 12}, {-8, 13}, {-14, 24},
	}
	testScale(t, negative, tests)
}

func TestSelfDestructingScale(t *testing.T) {
	weird := NewScale([]int{2, -2})
	tests := [][2]int{
		{0, 0}, {1, 2}, {2, 0}, {3, 2}, {4, 0}, {5, 2}, {6, 0}, {7, 2}, {8, 0}, {14, 0},
		{-1, 2}, {-2, 0}, {-3, 2}, {-7, 2}, {-8, 0}, {-14, 0},
	}
	testScale(t, weird, tests)
}
