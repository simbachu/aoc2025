package day3

import (
	"errors"
	"regexp"
	"strings"
)

// each line is a string of digits. find the largest number made of a pair of digits, in order
// "987654321111111" -> "98"

type Bank []int // one line of input, as array of digits

func (b Bank) MakeBank(in string) (Bank, error) {
	// reject non-digits
	if regexp.MustCompile(`\D`).MatchString(in) {
		return nil, errors.New("input contains non-digit characters")
	}
	bank := make(Bank, len(in))
	for i, ch := range in {
		bank[i] = int(ch - '0')
	}
	return bank, nil
}

// in c++ i would use std::max_element on a range with len-1 and find the pivot,
// then do the same with rest of the array starting at pivot+1
func (b Bank) MaxPair() int {
	if len(b) < 2 {
		return 0
	}

	// Find the index of the max element in range [0, len-1)
	maxIdx := 0
	for i := 1; i < len(b)-1; i++ {
		if b[i] > b[maxIdx] {
			maxIdx = i
		}
	}

	// Find the max element in range [maxIdx+1, len)
	maxSecondIdx := maxIdx + 1
	for i := maxIdx + 2; i < len(b); i++ {
		if b[i] > b[maxSecondIdx] {
			maxSecondIdx = i
		}
	}

	return b[maxIdx]*10 + b[maxSecondIdx]
}

func (b Bank) Max_N(n int) int {
	// returns the n largest digits in order of appearance, as a number
	if n <= 0 || n > len(b) {
		return 0
	}

	// Recursive approach: for each possible first digit position in range [0, len-n],
	// find the max digit, then recursively solve for remaining (n-1) digits after it
	return b.maxNHelper(n, 0)
}

func (b Bank) maxNHelper(n int, startIdx int) int {
	if n == 0 {
		return 0
	}
	if startIdx >= len(b) {
		return 0
	}

	// Need n digits, so first digit can be at most at position len(b)-n
	endSearchPos := len(b) - n
	if endSearchPos < startIdx {
		return 0
	}

	// Find the maximum digit in the valid range
	maxDigit := -1
	maxPos := -1
	for i := startIdx; i <= endSearchPos; i++ {
		if b[i] > maxDigit {
			maxDigit = b[i]
			maxPos = i
		}
	}

	if maxPos == -1 {
		return 0
	}

	// Use this digit and recursively find the best (n-1) digits after it
	if n == 1 {
		return maxDigit
	}

	remaining := b.maxNHelper(n-1, maxPos+1)
	// Compute 10^(n-1) to properly position maxDigit
	multiplier := 1
	for i := 0; i < n-1; i++ {
		multiplier *= 10
	}

	return maxDigit*multiplier + remaining
}

func SolveDay3Part1(input string) interface{} {
	result := 0
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// Trim whitespace including carriage returns
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bank, err := Bank(nil).MakeBank(line)
		if err != nil {
			continue
		}
		result += bank.MaxPair()
	}
	return result
}

func SolveDay3Part2(input string) interface{} {
	result := 0
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// Trim whitespace including carriage returns
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bank, err := Bank(nil).MakeBank(line)
		if err != nil {
			continue
		}
		result += bank.Max_N(12)
	}
	return result
}
