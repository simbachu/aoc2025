package day6

import (
	"testing"
)

const kDay6SampleInput = `
123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  `

const kDay6SampleOutputPart1 = 4277556

func TestReadInput(t *testing.T) {
	problems := ReadInput(kDay6SampleInput)
	if len(problems) != 4 {
		t.Errorf("Expected 4 problems, got %d", len(problems))
	}
}

func TestSolvePart1(t *testing.T) {
	result := SolveDay6Part1(kDay6SampleInput)
	if result != kDay6SampleOutputPart1 {
		t.Errorf("Expected %d, got %d", kDay6SampleOutputPart1, result)
	}
}
