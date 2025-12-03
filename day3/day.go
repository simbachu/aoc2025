package day3

type Day3 struct{}

func (d Day3) Part1(input string) interface{} {
	return SolveDay3Part1(input)
}

func (d Day3) Part2(input string) (interface{}, bool) {
	return SolveDay3Part2(input), true
}
