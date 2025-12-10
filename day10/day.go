package day10

type Day10 struct{}

func (d Day10) Part1(input string) interface{} {
	return SolveDay10Part1(input)
}

func (d Day10) Part2(input string) (interface{}, bool) {
	return SolveDay10Part2(input), false
}
