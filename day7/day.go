package day7

// Day7 implements the Day interface for day 7
type Day7 struct{}

// Part1 implements the Day interface
func (d Day7) Part1(input string) interface{} {
	return SolveDay7Part1(input)
}

// Part2 implements the Day interface
func (d Day7) Part2(input string) (interface{}, bool) {
	return SolveDay7Part2(input), true
}
