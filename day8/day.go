package day8

// Day8 implements the Day interface for day 8
type Day8 struct{}

// Part1 implements the Day interface
func (d Day8) Part1(input string) interface{} {
	return SolveDay8Part1(input)
}

// Part2 implements the Day interface
func (d Day8) Part2(input string) (interface{}, bool) {
	return SolveDay8Part2(input), true
}
