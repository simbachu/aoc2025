package day2

import "testing"

const kDay2SampleInput = `11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`
const kDay2SampleOutput = 1227775554

func TestToRange(t *testing.T) {
	tests := []struct {
		input    string
		expected Range
	}{
		{"11-22", Range{Start: 11, End: 22}},
		{"95-115", Range{Start: 95, End: 115}},
		{"998-1012", Range{Start: 998, End: 1012}},
		{"1188511880-1188511890", Range{Start: 1188511880, End: 1188511890}},
		{"222220-222224", Range{Start: 222220, End: 222224}},
		{"1698522-1698528", Range{Start: 1698522, End: 1698528}},
	}
	for _, test := range tests {
		result := ToRange(test.input)
		if result != test.expected {
			t.Errorf("ToRange(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestPart1IdOfConcern(t *testing.T) {
	tests := []struct {
		input       int
		expected    int
		expectedErr string
	}{
		{11, 11, ""},
		{1111, 1111, ""},
		{111111, 111111, ""},
		{1234, 0, "numbers do not match"},
		{12341234, 12341234, ""},
	}
	for _, test := range tests {
		result, err := Part1IdOfConcern(test.input)

		if test.expectedErr == "" {
			if err != nil {
				t.Errorf("Part1IdOfConcern(%d) unexpected error %v", test.input, err)
			}
		} else {
			if err == nil || err.Error() != test.expectedErr {
				t.Errorf("Part1IdOfConcern(%d) error = %v, expected %q", test.input, err, test.expectedErr)
			}
		}

		if result != test.expected {
			t.Errorf("Part1IdOfConcern(%d) result = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestPart2IdOfConcern(t *testing.T) {
	tests := []struct {
		input       int
		expected    int
		expectedErr string
	}{
		{11, 11, ""},
		{111, 111, ""},
		{1111, 1111, ""},
		{123123123, 123123123, ""},
		{1234, 0, "parts do not match"},
		{12341234, 12341234, ""},
		{111111, 111111, ""},
		{1234, 0, "parts do not match"},
		{12341234, 12341234, ""},
	}
	for _, test := range tests {
		result, err := Part2IdOfConcern(test.input)
		if test.expectedErr == "" {
			if err != nil {
				t.Errorf("Part2IdOfConcern(%d) unexpected error %v", test.input, err)
			}
		}
		if result != test.expected {
			t.Errorf("Part2IdOfConcern(%d) result = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestSolveDay2Part1(t *testing.T) {
	result := SolveDay2Part1(kDay2SampleInput)
	if result != kDay2SampleOutput {
		t.Errorf("SolveDay2Part1(%s) = %d, expected %d", kDay2SampleInput, result, kDay2SampleOutput)
	}
}
