package day6

import (
	"testing"
)

const kDay6SampleInput = `
123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  `

var kDay6Part1ExpectedProblems = []struct {
	problem mathProblem
	result  int
}{
	{problem: mathProblem{operands: []int{123, 45, 6}, operator: '*'}, result: 33210},
	{problem: mathProblem{operands: []int{328, 64, 98}, operator: '+'}, result: 490},
	{problem: mathProblem{operands: []int{51, 387, 215}, operator: '*'}, result: 4243455},
	{problem: mathProblem{operands: []int{64, 23, 314}, operator: '+'}, result: 401},
}

var kDay6Part2ExpectedProblems = []struct {
	problem mathProblem
	result  int
}{
	{problem: mathProblem{operands: []int{4, 431, 623}, operator: '+'}, result: 1058},
	{problem: mathProblem{operands: []int{175, 581, 32}, operator: '*'}, result: 3253600},
	{problem: mathProblem{operands: []int{8, 248, 369}, operator: '+'}, result: 625},
	{problem: mathProblem{operands: []int{356, 24, 1}, operator: '*'}, result: 8544},
}

const kDay6SampleOutputPart1 = 4277556
const kDay6SampleOutputPart2 = 3263827

func TestReadInput(t *testing.T) {
	t.Run("Part 1", func(t *testing.T) {
		problems := ReadInput(kDay6SampleInput, false)
		if len(problems) != len(kDay6Part1ExpectedProblems) {
			t.Errorf("Expected %d problems, got %d", len(kDay6Part1ExpectedProblems), len(problems))
			return
		}
		for i, expected := range kDay6Part1ExpectedProblems {
			actual := problems[i]
			if actual.Solve() != expected.result {
				t.Errorf("Problem %d: expected result %d, got %d", i, expected.result, actual.Solve())
			}
		}
	})
	t.Run("Part 2", func(t *testing.T) {
		problems := ReadInput(kDay6SampleInput, true)
		if len(problems) != len(kDay6Part2ExpectedProblems) {
			t.Errorf("Expected %d problems, got %d", len(kDay6Part2ExpectedProblems), len(problems))
			return
		}
		for i, expected := range kDay6Part2ExpectedProblems {
			actual := problems[i]
			if actual.Solve() != expected.result {
				t.Errorf("Problem %d: expected result %d, got %d", i, expected.result, actual.Solve())
			}
		}
	})
}

func TestSolvePart1(t *testing.T) {
	result := SolveDay6Part1(kDay6SampleInput)
	if result != kDay6SampleOutputPart1 {
		t.Errorf("Expected %d, got %d", kDay6SampleOutputPart1, result)
	}
}

func TestSolvePart2(t *testing.T) {
	result := SolveDay6Part2(kDay6SampleInput)
	if result != kDay6SampleOutputPart2 {
		t.Errorf("Expected %d, got %d", kDay6SampleOutputPart2, result)
	}
}
