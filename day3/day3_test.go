package day3

import "testing"

const kDay3SampleInput = `987654321111111
811111111111119
234234234234278
818181911112111`

func TestMakeBank(t *testing.T) {
	input := "1234567890"
	expected := Bank{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	bank, err := Bank(nil).MakeBank(input)
	if err != nil {
		t.Errorf("MakeBank(%s) returned error: %v", input, err)
	}
	for i, v := range bank {
		if v != expected[i] {
			t.Errorf("MakeBank(%s)[%d] = %d; want %d", input, i, v, expected[i])
		}
	}
}

func TestMaxPair(t *testing.T) {
	tests := []struct {
		input    Bank
		expected int
	}{
		{Bank{9, 8, 7, 6, 5, 4, 3, 2, 1, 1, 1, 1, 1, 1}, 98},
		{Bank{8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9}, 89},
		{Bank{2, 3, 4, 2, 3, 4, 2, 3, 4, 2, 7, 8}, 78},
		{Bank{8, 1, 8, 1, 8, 1, 9, 1, 1, 1, 2, 1, 1}, 92},
	}
	for _, test := range tests {
		result := test.input.MaxPair()
		if result != test.expected {
			t.Errorf("MaxPair(%v) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestMaxN(t *testing.T) {
	tests := []struct {
		input    Bank
		n        int
		expected int
	}{
		{Bank{9, 8, 7, 6, 5, 4, 3, 2, 1, 1, 1, 1, 1, 1}, 3, 987},
		{Bank{8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9}, 2, 89},
		{Bank{2, 3, 4, 2, 3, 4, 2, 3, 4, 2, 7, 8}, 4, 4348},
		{Bank{8, 1, 8, 1, 8, 1, 9, 1, 1, 1, 2, 1, 1}, 5, 88892},
	}
	for _, test := range tests {
		result := test.input.Max_N(test.n)
		if result != test.expected {
			t.Errorf("Max_N(%v, %d) = %d; want %d", test.input, test.n, result, test.expected)
		}
	}
}

func TestSolveDay3Part1(t *testing.T) {
	expected := 357
	result := SolveDay3Part1(kDay3SampleInput)
	if result != expected {
		t.Errorf("SolveDay3Part1() = %v; want %v", result, expected)
	}
}

func TestSolveDay3Part2(t *testing.T) {
	expected := 3121910778619
	result := SolveDay3Part2(kDay3SampleInput)
	if result != expected {
		t.Errorf("SolveDay3Part2() = %v; want %v", result, expected)
	}
}
