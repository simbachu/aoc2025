package main

// Day represents a single day's solution for Advent of Code
type Day interface {
	// Part1 solves part 1 of the day's puzzle
	// Returns any printable value
	Part1(input string) interface{}
	// Part2 solves part 2 of the day's puzzle
	// Returns the result and true if implemented, nil and false if not yet unlocked
	Part2(input string) (interface{}, bool)
}
