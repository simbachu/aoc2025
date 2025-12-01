package day1

import (
	"testing"
)

const sampleInput = `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`

func TestDialTurn(t *testing.T) {
	tests := []struct {
		start     DialValue
		direction Direction
		steps     int
		expected  DialValue
	}{
		{50, Clockwise, 10, 60},
		{50, CounterClockwise, 20, 30},
		{0, CounterClockwise, 1, 99},
		{99, Clockwise, 1, 0},
		{95, Clockwise, 10, 5},
		{5, CounterClockwise, 10, 95},
	}
	for _, test := range tests {
		dial := test.start
		newValue := dial.Turn(test.direction, test.steps)
		if newValue != test.expected {
			t.Errorf("Turn(%d, %d) from %d: expected %d, got %d", test.direction, test.steps, test.start, test.expected, newValue)
		}
	}
}

func TestSolveDay1Part1(t *testing.T) {

	const expectedSampleOutput = 3
	t.Run("Sample Input", func(t *testing.T) {
		output := SolveDay1Part1(sampleInput)
		if output != expectedSampleOutput {
			t.Errorf("Expected %d, got %d", expectedSampleOutput, output)
		}
	})
}

func TestTurnAndCountZeros(t *testing.T) {
	tests := []struct {
		start          DialValue
		direction      Direction
		steps          int
		expectedValue  DialValue
		expectedPasses int
	}{
		{50, Clockwise, 150, 0, 2},         // should be 50 -> 0 -> 50 -> 0
		{50, CounterClockwise, 200, 50, 2}, // should be 50 -> 0 -> 50 -> 0 -> 50
		{0, Clockwise, 100, 0, 1},          // should be 0 -> 0
		{99, CounterClockwise, 101, 98, 1}, // should be 99 -> 98 (pass 0 once)
	}
	for _, test := range tests {
		dial := test.start
		newValue, passes := dial.TurnAndCountZeros(test.direction, test.steps)
		if newValue != test.expectedValue || passes != test.expectedPasses {
			t.Errorf("TurnAndCountZeros(%d, %d) from %d: expected (%d, %d), got (%d, %d)",
				test.direction, test.steps, test.start, test.expectedValue, test.expectedPasses, newValue, passes)
		}
	}
}

func TestSolveDay1Part2(t *testing.T) {
	// use sampleInput for part 2 as well
	// but expected is 6
	const expectedPart2Output = 6

	t.Run("Sample Input Part 2", func(t *testing.T) {
		output := SolveDay1Part2(sampleInput)
		if output != expectedPart2Output {
			t.Errorf("Expected %d, got %d", expectedPart2Output, output)
		}
	})
}
