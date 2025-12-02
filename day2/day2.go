package day2

import (
	"errors"
	"strconv"
	"strings"
)

type Range struct {
	Start int
	End   int
}

func ToRange(s string) Range {
	parts := strings.Split(s, "-")
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return Range{Start: start, End: end}
}

// take number and return string if the number is only consisting
// of two groups of equal digits, otherwise return empty string
func Part1IdOfConcern(n int) (int, error) {
	digits := strconv.Itoa(n)
	if len(digits)%2 != 0 {
		return 0, errors.New("number of digits is odd")
	}
	// split string in the middle, if strings match, retun string, else return empty string
	mid := len(digits) / 2
	if digits[:mid] == digits[mid:] {
		return n, nil
	}
	return 0, errors.New("numbers do not match")
}

// return int if the number contains only repeating parts, otherwise error
// 11, 222, 123123123
// checks if the number is composed of a repeating pattern
func Part2IdOfConcern(n int) (int, error) {
	digits := strconv.Itoa(n)
	length := len(digits)

	// try all possible segment lengths that divide evenly into the total length
	for segment_len := 1; segment_len <= length/2; segment_len++ {
		if length%segment_len != 0 {
			continue
		}
		segment := digits[:segment_len]
		repetitions := length / segment_len
		expected := strings.Repeat(segment, repetitions)
		if digits == expected {
			return n, nil
		}
	}
	return 0, errors.New("parts do not match")
}

// take csv string and return string with all ids of concern
func SolveDay2Part1(input string) int {
	values := strings.Split(input, ",")
	result := 0
	for _, value := range values {
		r := ToRange(value)
		for i := r.Start; i <= r.End; i++ {
			id, err := Part1IdOfConcern(i)
			if err == nil {
				result += id
			}
		}
	}
	return result
}

func SolveDay2Part2(input string) int {
	values := strings.Split(input, ",")
	result := 0
	for _, value := range values {
		r := ToRange(value)
		for i := r.Start; i <= r.End; i++ {
			id, err := Part2IdOfConcern(i)
			if err == nil {
				result += id
			}
		}
	}
	return result
}
