package day2

// Day2 implements the Day interface for day 2
type Day2 struct{}

// Part1 implements the Day interface
func (d Day2) Part1(input string) interface{} {
	return SolveDay2Part1(input)
}

// Part2 implements the Day interface
func (d Day2) Part2(input string) (interface{}, bool) {
	return SolveDay2Part2(input), true
}
