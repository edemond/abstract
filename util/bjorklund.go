// Implements Bjorklund's algorithm in order to generate "Euclidean" rhythms.
// References:
// "The Theory of Rep-Rate Pattern Generation in the SNS Timing System", E. Bjorklund
// "The Euclidean Algorithm Generates Traditional Musical Rhythms", Godfried Toussaint
package util

// Rotates an array so that the value at pos becomes the start.
func rotate(pos int, steps []int) []int {
	length := len(steps)
	rotated := make([]int, length)
	for i := 0; i < length; i++ {
		rotated[i] = steps[(i+pos)%length]
	}
	return rotated
}

// Returns the position of the first 1 in the list of steps.
// Returns -1 if none is found.
func findFirstPulse(steps []int) int {
	for i := 0; i < len(steps); i++ {
		if steps[i] == 1 {
			return i
		}
	}
	return -1
}

// this should be something like bjork(2, 5)
func Bjorklund(pulses, steps int) []int {
	if pulses > int(steps/2) {
		steps, pulses = pulses, steps
	}

	// TODO: We should add some kind of optimization for
	// pulses that divide evenly into the number of steps.
	// And benchmark that to prove it's actually doing anything.

	count := []int{}
	remainder := []int{pulses}
	divisor := steps - pulses
	level := 0

	for {
		count = append(count, int(divisor/remainder[level])) // truncates toward 0
		remainder = append(remainder, divisor%remainder[level])
		divisor = remainder[level]
		level++
		if remainder[level] <= 1 {
			break
		}
	}

	count = append(count, divisor)

	// Build the string.
	result := build(level, count, remainder)

	// Rotate it so that we always have a 1 in the first position.
	firstPulse := findFirstPulse(result)
	if firstPulse == -1 {
		// No pulses in this.
		return result
	}
	return rotate(firstPulse, result)
}

func build(level int, count []int, remainder []int) []int {
	if level == -2 {
		return []int{1}
	} else if level == -1 {
		return []int{0}
	} else {
		rest := []int{}
		for i := 0; i < count[level]; i++ {
			rest = append(rest, build(level-1, count, remainder)...)
		}
		if remainder[level] != 0 {
			rest = append(rest, build(level-2, count, remainder)...)
		}
		return rest
	}
}
