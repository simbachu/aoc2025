package day7

import (
	"strings"
)

type DiagramRow struct {
	exits map[int]bool
}

type Diagram struct {
	rows [][]rune
}

type Position struct {
	row int
	col int
}

// naive solution took too much time and space.
// had to add memoization to make it work.
func (d *Diagram) CountPathsFromS() int {
	if len(d.rows) == 0 {
		return 0
	}

	// Find S in the first row
	var sCol int = -1
	for i, char := range d.rows[0] {
		if char == 'S' {
			sCol = i
			break
		}
	}
	if sCol == -1 {
		return 0
	}

	// cache results for each position
	memo := make(map[Position]int)

	// recursive function with memoization
	var countPaths func(row, col int) int
	countPaths = func(row, col int) int {
		// check if we've computed this before
		pos := Position{row: row, col: col}
		if val, ok := memo[pos]; ok {
			return val
		}

		// base case: reached end of diagram
		if row >= len(d.rows) {
			memo[pos] = 1
			return 1
		}

		// check if we've gone off the sides
		if col < 0 || col >= len(d.rows[row]) {
			memo[pos] = 0
			return 0
		}

		char := d.rows[row][col]

		// if we hit a splitter, sum paths from left and right branches
		if char == '^' {
			nextRow := row + 1

			// check if we've reached the end
			if nextRow >= len(d.rows) {
				memo[pos] = 1
				return 1
			}

			count := 0
			// left branch (col - 1)
			if col-1 >= 0 && col-1 < len(d.rows[nextRow]) {
				count += countPaths(nextRow, col-1)
			}
			// right branch (col + 1)
			if col+1 >= 0 && col+1 < len(d.rows[nextRow]) {
				count += countPaths(nextRow, col+1)
			}

			memo[pos] = count
			return count
		}

		// if we hit something that's not '.' or 'S', the path is blocked
		if char != '.' && char != 'S' {
			memo[pos] = 0
			return 0
		}

		// continue straight down (raymarch)
		nextRow := row + 1
		result := countPaths(nextRow, col)
		memo[pos] = result
		return result
	}

	return countPaths(0, sCol)
}

func (r *DiagramRow) AddExit(exitPos int) bool {
	if r.exits == nil {
		r.exits = make(map[int]bool)
	}
	if _, exists := r.exits[exitPos]; exists {
		return false
	}
	r.exits[exitPos] = true
	return true
}

func NewDiagramRow(
	gridRow []rune,
	rowIndex int,
	previousRow *DiagramRow,
	splitCallback func(int),
) *DiagramRow {
	row := &DiagramRow{}

	if previousRow == nil {
		for i, char := range gridRow {
			if char == 'S' {
				_ = row.AddExit(i)
			}
		}
		return row
	}

	for exitPos := range previousRow.exits {
		if exitPos < 0 || exitPos >= len(gridRow) {
			continue
		}
		char := gridRow[exitPos]

		switch char {
		case '.':
			_ = row.AddExit(exitPos)
		case '^':
			splitCallback(exitPos)
			if exitPos > 0 && gridRow[exitPos-1] == '.' {
				_ = row.AddExit(exitPos - 1)
			}
			if exitPos < len(gridRow)-1 && gridRow[exitPos+1] == '.' {
				_ = row.AddExit(exitPos + 1)
			}
		default:
			continue
		}
	}
	return row
}

func SolveDay7Part1(input string) int {
	lines := strings.Split(input, "\n")
	linesAsRunes := make([][]rune, len(lines))
	for i, line := range lines {
		linesAsRunes[i] = []rune(line)
	}
	count := 0
	countedSplitters := make(map[int]map[int]bool)
	var previousRow *DiagramRow
	for rowIndex := range linesAsRunes {
		previousRow = NewDiagramRow(
			linesAsRunes[rowIndex],
			rowIndex,
			previousRow,
			func(splitterPos int) {
				if countedSplitters[rowIndex] == nil {
					countedSplitters[rowIndex] = make(map[int]bool)
				}
				if !countedSplitters[rowIndex][splitterPos] {
					countedSplitters[rowIndex][splitterPos] = true
					count++
				}
			})
	}
	return count
}

func SolveDay7Part2(input string) int {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	nonEmptyLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	linesAsRunes := make([][]rune, len(nonEmptyLines))
	for i, line := range nonEmptyLines {
		linesAsRunes[i] = []rune(line)
	}
	return (&Diagram{rows: linesAsRunes}).CountPathsFromS()
}
