package day1

// Day1 implements the Day interface for day 1
type Day1 struct{}

// Part1 implements the Day interface
func (d Day1) Part1(input string) interface{} {
	return SolveDay1Part1(input)
}

// Part2 implements the Day interface
func (d Day1) Part2(input string) (interface{}, bool) {
	return SolveDay1Part2(input), true
	// Return (nil, false) if Part 2 is not yet unlocked
}
