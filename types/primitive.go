package types

// A Value wrapper for strings so they can be first-class values in the analyzer.
type String string

// A Value wrapper for numbers so they can be first-class values in the analyzer.
type Number struct {
	Value  uint64
	Digits int // Significant when we're treating a number like a bit pattern.
}

func (s String) String() string { // lol
	return string(s)
}

func (s String) HasValue() bool {
	return true // strings, uh, always have values, I guess
}

func (n *Number) String() string {
	return string(n.Value) // TODO: Digits
}

func (n *Number) HasValue() bool {
	return n != nil
}
