package day5

// Day5 implements the Day interface for day 5
type Day5 struct{}

// Part1 implements the Day interface
func (d Day5) Part1(input string) interface{} {
	return SolveDay5Part1(input)
}

// Part2 implements the Day interface
func (d Day5) Part2(input string) (interface{}, bool) {
	return SolveDay5Part2(input), true
}
