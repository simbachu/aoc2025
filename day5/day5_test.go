package day5

import "testing"

const kDay5SampleInput = `3-5
10-14
16-20
12-18

1
5
8
11
17
32`

const kDay5SampleOutputPart1 = 3
const kDay5SampleOutputPart2 = 14

func TestReadInput(t *testing.T) {
	ranges, ids := ReadInput(kDay5SampleInput)
	if len(ranges) != 4 {
		t.Errorf("Expected 4 ranges, got %d", len(ranges))
	}
	if len(ids) != 6 {
		t.Errorf("Expected 6 ids, got %d", len(ids))
	}
}

func TestFlattenRanges(t *testing.T) {
	ranges := []IdRange{
		{Start: 3, End: 5},
		{Start: 10, End: 14},
		{Start: 16, End: 20},
		{Start: 12, End: 18},
	}
	expected := []IdRange{
		{Start: 3, End: 5},
		{Start: 10, End: 20},
	}
	result := FlattenRanges(ranges)
	if len(result) != len(expected) {
		t.Errorf("Expected %d ranges, got %d", len(expected), len(result))
	}
	for i, r := range result {
		if r != expected[i] {
			t.Errorf("Expected range %d to be %v, got %v", i, expected[i], r)
		}
	}
}

func TestSolveDay5Part1(t *testing.T) {
	result := SolveDay5Part1(kDay5SampleInput)
	if result != kDay5SampleOutputPart1 {
		t.Errorf("SolveDay5Part1(%s) = %d, expected %d", kDay5SampleInput, result, kDay5SampleOutputPart1)
	}
}

func TestSolveDay5Part2(t *testing.T) {
	result := SolveDay5Part2(kDay5SampleInput)
	if result != kDay5SampleOutputPart2 {
		t.Errorf("SolveDay5Part2(%s) = %d, expected %d", kDay5SampleInput, result, kDay5SampleOutputPart2)
	}
}
