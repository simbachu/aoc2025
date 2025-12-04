package day4

import "strings"

type Roll bool // location on the grid, true has roll, false is empty

type Grid [][]Roll // 2D grid of rolls

type Callback func(x, y int) // callback function for grid operations

// construct a Grid from input strings, taking list of strings and which rune is a roll
func MakeGrid(input []string, mark rune) Grid {
	grid := make(Grid, len(input))
	for i, line := range input {
		grid[i] = make([]Roll, len(line))
		for j, char := range line {
			grid[i][j] = Roll(char == mark)
		}
	}
	return grid
}

func (g Grid) CountAdjacent(x, y int) int {
	// counts number of adjacent rolls (true) around position (x, y)
	// where x is column (0 = leftmost) and y is row (0 = topmost)
	// returns 0 if position is out of bounds
	if y < 0 || y >= len(g) || x < 0 || x >= len(g[0]) {
		return 0
	}
	count := 0
	directions := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for _, dir := range directions {
		newX, newY := x+dir[0], y+dir[1]
		// check bounds
		if newY >= 0 && newY < len(g) && newX >= 0 && newX < len(g[0]) {
			if g[newY][newX] {
				count++
			}
		}
	}
	return count
}

func (g Grid) pop(x, y int) {
	// sets position (x, y) to false
	if y >= 0 && y < len(g) && x >= 0 && x < len(g[0]) {
		g[y][x] = false
	}
}

func (g Grid) popMany(coordinates [][2]int) {
	//takes a list of coordinates and sets those positions to false
	for _, coord := range coordinates {
		x, y := coord[0], coord[1]
		g[y][x] = false
	}
}

// take a function and call it when a roll has less than 4 adjacent rolls
func (g Grid) ApplyWhenSparse(callback Callback) {
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[0]); x++ {
			if g[y][x] && g.CountAdjacent(x, y) < 4 {
				callback(x, y)
			}
		}
	}
}

func ParseGrid(input string, mark rune) Grid {
	// parse input into grid
	lines := []string{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		lines = append(lines, line)
	}
	return MakeGrid(lines, mark)
}

func SolveDay4Part1(input string) int {
	grid := ParseGrid(input, '@')
	count := 0
	grid.ApplyWhenSparse(func(x, y int) {
		count++
	})
	return count
}

// apply the sparse rule, counting and adding the rolls to a list
// pop the counted rolls, then repeat until no more rolls can be removed
// and count the total number of rolls removed
func SolveDay4Part2(input string) int {
	grid := ParseGrid(input, '@')
	count := 0
	changed := true
	for changed {
		changed = false
		toPop := [][2]int{}
		grid.ApplyWhenSparse(func(x, y int) {
			count++
			toPop = append(toPop, [2]int{x, y})
			changed = true
		})
		grid.popMany(toPop)
	}
	return count
}
