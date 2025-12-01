package day1

import (
	"strconv"
	"strings"
)

type Direction int

const (
	Clockwise        Direction = 1
	CounterClockwise Direction = -1
)

type DialValue int

const (
	kMinDialValue   = 0
	kMaxDialValue   = 99
	kDialStartValue = 50
)

// Turn the dial in the given direction by the given number of steps
// and return the new dial value
// If we start at 50 and turn clockwise by 10 steps, we end at 60
// If we start at 50 and turn counter-clockwise by 20 steps, we end at 30
func (d *DialValue) Turn(direction Direction, steps int) DialValue {
	newValue := int(*d) + int(direction)*steps
	// wrap around the dial values
	for newValue < kMinDialValue {
		newValue += (kMaxDialValue - kMinDialValue + 1)
	}
	for newValue > kMaxDialValue {
		newValue -= (kMaxDialValue - kMinDialValue + 1)
	}
	*d = DialValue(newValue)
	return *d
}

// Turn the dial and count every time we pass 0, not just landing on it
// Returns the new dial value and the number of times we passed 0
func (d *DialValue) TurnAndCountZeros(direction Direction, steps int) (DialValue, int) {
	startValue := int(*d)

	// Count every position where we're at 0 during the turn
	passes := 0
	if direction == Clockwise {
		for i := 1; i <= steps; i++ {
			if (startValue+i)%(kMaxDialValue+1) == 0 {
				passes++
			}
		}
	} else {
		for i := 1; i <= steps; i++ {
			pos := (startValue - i) % (kMaxDialValue + 1)
			if pos < 0 {
				pos += (kMaxDialValue + 1)
			}
			if pos == 0 {
				passes++
			}
		}
	}

	newValue := d.Turn(direction, steps)
	return newValue, passes
}

// SolveDay1 starts with a value kDialStartValue on the dial
// Then take each line as string and convert to a DialRotation
// Each line is in the format "DS" where D is either "R" or "L" and S is the number of steps to turn
// After each rotation, if the dial lands on 0, we increment a counter
// Finally, return the counter
func SolveDay1Part1(input string) int {
	result := 0

	dial := DialValue(kDialStartValue)

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue // skip invalid lines
		}

		// First character is direction (L or R)
		dirChar := line[0]
		// Rest is the number of steps
		stepsStr := line[1:]

		steps, err := strconv.Atoi(stepsStr)
		if err != nil {
			continue // skip invalid lines
		}

		var dir Direction
		switch dirChar {
		case 'R':
			dir = Clockwise
		case 'L':
			dir = CounterClockwise
		default:
			continue // skip invalid lines
		}

		if dial.Turn(dir, steps) == 0 {
			result++
		}
	}
	return result
}

func SolveDay1Part2(input string) int {
	result := 0
	dial := DialValue(kDialStartValue)

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue // skip invalid lines
		}
		// First character is direction (L or R)
		dirChar := line[0]
		// Rest is the number of steps
		stepsStr := line[1:]
		steps, err := strconv.Atoi(stepsStr)
		if err != nil {
			continue // skip invalid lines
		}
		var dir Direction
		switch dirChar {
		case 'R':
			dir = Clockwise
		case 'L':
			dir = CounterClockwise
		default:
			continue // skip invalid lines
		}
		_, passes := dial.TurnAndCountZeros(dir, steps)
		result += passes
	}
	return result
}
