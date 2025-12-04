package day4

import (
	"reflect"
	"testing"
)

const kDay4SampleInput = `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

const kDay4Part1Expected = 13
const kDay4Part2Expected = 43

func TestMakeGrid(t *testing.T) {
	input := []string{
		"@..@",
		".@@.",
		"....",
		"@.@.",
	}
	grid := MakeGrid(input, '@')
	expected := Grid{
		{true, false, false, true},
		{false, true, true, false},
		{false, false, false, false},
		{true, false, true, false},
	}
	if !reflect.DeepEqual(grid, expected) {
		t.Errorf("MakeGrid() = %v; want %v", grid, expected)
	}
}

func TestCountAdjacent(t *testing.T) {
	input := []string{
		"@.@",
		".@.",
		"...",
		"@.@",
	}
	grid := MakeGrid(input, '@')
	tests := []struct {
		x, y     int
		expected int
	}{
		{0, 0, 1},
		{2, 2, 2},
		{1, 3, 2},
		{3, 3, 0},
	}
	for _, test := range tests {
		if got := grid.CountAdjacent(test.x, test.y); got != test.expected {
			t.Errorf("CountAdjacent(%d, %d) = %d; want %d", test.x, test.y, got, test.expected)
		}
	}
}

func TestSolveDay4Part1(t *testing.T) {
	input := kDay4SampleInput
	expected := kDay4Part1Expected
	result := SolveDay4Part1(input)
	if result != expected {
		t.Errorf("SolveDay4Part1() = %v; want %v", result, expected)
	}
}

func TestSolveDay4Part2(t *testing.T) {
	input := kDay4SampleInput
	expected := kDay4Part2Expected
	result := SolveDay4Part2(input)
	if result != expected {
		t.Errorf("SolveDay4Part2() = %v; want %v", result, expected)
	}
}
