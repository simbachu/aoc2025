package day4

// Day4 implements the Day interface for day 4
type Day4 struct{}

// Part1 implements the Day interface
func (d Day4) Part1(input string) interface{} {
	return SolveDay4Part1(input)
}

// Part2 implements the Day interface
func (d Day4) Part2(input string) (interface{}, bool) {
	return SolveDay4Part2(input), true
}
