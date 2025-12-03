package main

import (
	"fmt"
	"os"
	"strconv"

	"aoc2025/day1"
	"aoc2025/day2"
	"aoc2025/day3"
)

var days = map[int]Day{
	1: day1.Day1{},
	2: day2.Day2{},
	3: day3.Day3{},
}

func main() {
	var daysToRun []int

	if len(os.Args) > 1 {
		dayNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid day number: %s\n", os.Args[1])
			os.Exit(1)
		}
		daysToRun = []int{dayNumber}
	} else {
		for dayNum := range days {
			daysToRun = append(daysToRun, dayNum)
		}
	}

	for _, dayNumber := range daysToRun {
		day, exists := days[dayNumber]
		if !exists {
			fmt.Fprintf(os.Stderr, "Day %d not implemented yet\n", dayNumber)
			continue
		}

		dataFile := fmt.Sprintf("day%d/data.txt", dayNumber)
		data, err := os.ReadFile(dataFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", dataFile, err)
			continue
		}

		fmt.Printf("Day %d:\n", dayNumber)
		part1 := day.Part1(string(data))
		fmt.Printf("  Part 1: %v\n", part1)

		if part2, unlocked := day.Part2(string(data)); unlocked {
			fmt.Printf("  Part 2: %v\n", part2)
		} else {
			fmt.Printf("  Part 2: Not yet unlocked\n")
		}
		fmt.Println()
	}
}
