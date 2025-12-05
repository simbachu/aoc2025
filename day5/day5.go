package day5

import (
	"sort"
	"strconv"
	"strings"
)

type IdRange struct {
	Start int
	End   int
}

func ToIdRange(s string) IdRange {
	parts := strings.Split(s, "-")
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return IdRange{Start: start, End: end}
}

func IsIdInRange(id int, r IdRange) bool {
	return id >= r.Start && id <= r.End
}

func WhenIdIsInRanges(id int, ranges []IdRange, f func(int)) {
	for _, r := range ranges {
		if IsIdInRange(id, r) {
			f(id)
		}
	}
}

// take a list of ranges and give back a list of non-overlapping ranges,
// each with start and end encompassing all overlapping ranges
func FlattenRanges(ranges []IdRange) []IdRange {
	if len(ranges) == 0 {
		return ranges
	}

	// sort ranges by start value
	sortedRanges := make([]IdRange, len(ranges))
	copy(sortedRanges, ranges)
	sort.Slice(sortedRanges, func(i, j int) bool {
		return sortedRanges[i].Start < sortedRanges[j].Start
	})

	// merge overlapping ranges
	newRanges := make([]IdRange, 0)
	newRanges = append(newRanges, sortedRanges[0])

	for i := 1; i < len(sortedRanges); i++ {
		lastRange := &newRanges[len(newRanges)-1]
		currentRange := sortedRanges[i]

		// check if current range overlaps with the last merged range
		if currentRange.Start <= lastRange.End {
			// do the merge
			if currentRange.End > lastRange.End {
				lastRange.End = currentRange.End
			}
		} else {
			// add the current range as a new range
			newRanges = append(newRanges, currentRange)
		}
	}

	return newRanges
}

func CountTotalRangeSpan(ranges []IdRange) int {
	total := 0
	for _, r := range ranges {
		total += r.End - r.Start + 1
	}
	return total
}

// input is first ranges, blank line, then ids, all newline separated
func ReadInput(input string) ([]IdRange, []int) {
	lines := strings.Split(input, "\n")
	ranges := make([]IdRange, 0)
	ids := make([]int, 0)
	start_id_index := 0

	// read up until the first blank line, putting the ranges in the ranges array
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			start_id_index = i + 1
			break
		}
		ranges = append(ranges, ToIdRange(lines[i]))
	}
	// read the ids
	for i := start_id_index; i < len(lines); i++ {
		id, err := strconv.Atoi(lines[i])
		if err != nil {
			break // stop reading ids if we can't convert the line to an int
		}
		ids = append(ids, id)
	}
	return ranges, ids
}

func IsIdInAnyRange(id int, ranges []IdRange) bool {
	for _, r := range ranges {
		if IsIdInRange(id, r) {
			return true
		}
	}
	return false
}

func SolveDay5Part1(input string) int {
	ranges, ids := ReadInput(input)
	result := 0
	for _, id := range ids {
		if IsIdInAnyRange(id, ranges) {
			result++
		}
	}
	return result
}

func SolveDay5Part2(input string) int {
	ranges, _ := ReadInput(input)               // we only care about the ranges
	flattenedRanges := FlattenRanges(ranges)    // we only want the total span of the ranges
	return CountTotalRangeSpan(flattenedRanges) // answer is the total number included in the ranges
}
