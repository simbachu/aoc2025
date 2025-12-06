package day6

type Day6 struct{}

func (d Day6) Part1(input string) interface{} {
	return SolveDay6Part1(input)
}

func (d Day6) Part2(input string) (interface{}, bool) {
	return SolveDay6Part2(input), true
}
