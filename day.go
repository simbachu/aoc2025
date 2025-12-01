package main

// Day represents a single day's solution for Advent of Code
type Day interface {
	// Part1 solves part 1 of the day's puzzle
	Part1(input string) int
	// Part2 solves part 2 of the day's puzzle
	// Returns the result and true if implemented, 0 and false if not yet unlocked
	Part2(input string) (int, bool)
}
