package day9

// Day9 implements the Day interface for day 9
type Day9 struct{}

// Part1 implements the Day interface
func (d Day9) Part1(input string) interface{} {
	return SolveDay9Part1(input)
}

// Part2 implements the Day interface
func (d Day9) Part2(input string) (interface{}, bool) {
	return SolveDay9Part2(input), true
}
