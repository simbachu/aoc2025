package day6

import "strings"

type mathProblem struct {
	operands []int
	operator rune
}

// make a math problem from a 2D grid of runes
func makeMathProblem(grid [][]rune) mathProblem {
	operands := make([]int, 0)
	operator := ' '
	for _, row := range grid {
		hasOperator := false
		for _, cell := range row {
			if cell == '+' || cell == '*' {
				if operator != ' ' {
					panic("multiple operators in a row")
				}
				operator = cell
				hasOperator = true
				break
			}
		}
		if hasOperator {
			continue
		}
		currentNumber := 0
		hasDigits := false
		for _, cell := range row {
			if cell >= '0' && cell <= '9' {
				currentNumber = currentNumber*10 + int(cell-'0')
				hasDigits = true
			} else if cell != ' ' {
				panic("invalid character in problem")
			}
		}
		if hasDigits {
			operands = append(operands, currentNumber)
		}
	}
	if operator == ' ' {
		panic("no operator in problem")
	}
	if len(operands) < 2 {
		panic("invalid number of operands in problem")
	}
	return mathProblem{
		operands: operands,
		operator: operator,
	}
}

// read the grid right-to-left columnar, so
//
// 123
//  45
//   6
// *
//
// should read as 356 * 24 * 1
func makeMathProblemColumnar(grid [][]rune) mathProblem {
	if len(grid) == 0 {
		panic("empty grid")
	}
	maxWidth := 0
	for _, row := range grid {
		if len(row) > maxWidth {
			maxWidth = len(row)
		}
	}

	operands := make([]int, 0)
	operator := ' '
	// iterate columns from right to left
	for col := maxWidth - 1; col >= 0; col-- {
		// collect digits from top to bottom in this column
		digits := make([]int, 0)
		foundOperator := false
		for row := 0; row < len(grid); row++ {
			if col >= len(grid[row]) {
				continue
			}
			cell := grid[row][col]
			if cell >= '0' && cell <= '9' {
				digits = append(digits, int(cell-'0'))
			} else if cell == '+' || cell == '*' {
				if operator != ' ' {
					panic("multiple operators in a column")
				}
				operator = cell
				foundOperator = true
				break
			} else if cell != ' ' {
				panic("invalid character in problem")
			}
		}
		if len(digits) > 0 {
			currentNumber := 0
			for _, digit := range digits {
				currentNumber = currentNumber*10 + digit
			}
			operands = append(operands, currentNumber)
		}
		if foundOperator {
			break
		}
	}
	if operator == ' ' {
		panic("no operator in problem")
	}
	if len(operands) < 2 {
		panic("invalid number of operands in problem")
	}
	return mathProblem{
		operands: operands,
		operator: operator,
	}
}

// solve a math problem, return the result
// this assumes that the problem is valid, and that the operands are valid,
// since we take in some type that we assume use a constructor for correctness
func (problem mathProblem) Solve() int {
	switch problem.operator {
	case '+':
		result := 0
		for _, operand := range problem.operands {
			result += operand
		}
		return result
	case '*':
		result := 1
		for _, operand := range problem.operands {
			result *= operand
		}
		return result
	default:
		panic("invalid operator")
	}
}

// a raw math problem is a 2D grid of runes, each rune is either a digit or an operator, or ' '
type mathProblemRaw struct {
	width  int      // number of columns that were surrounded by a full column of ' '
	height int      // number of lines in the input
	grid   [][]rune // 2D grid of runes, each rune is either a digit or an operator, or ' '
}

// make a raw math problem from a 2D grid of runes
func makeMathProblemRaw(input [][]rune) mathProblemRaw {
	return mathProblemRaw{
		width:  len(input[0]),
		height: len(input),
		grid:   input,
	}
}

func parseStringToRawProblems(input string) []mathProblemRaw {
	lineStrings := strings.Split(input, "\n")
	lines := make([][]rune, 0)
	maxWidth := 0

	for _, lineStr := range lineStrings {
		// only trim carriage returns, preserve all spaces for columnar reading
		lineStr = strings.TrimRight(lineStr, "\r")
		if len(lineStr) == 0 {
			continue
		}
		line := []rune(lineStr)
		lines = append(lines, line)
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	if len(lines) == 0 {
		return []mathProblemRaw{}
	}

	// find the spaces between problems (using maxWidth to check all possible columns)
	// a column is a boundary if it's all spaces in all rows that have that column
	isBoundaryColumn := make([]bool, maxWidth)
	for col := 0; col < maxWidth; col++ {
		allSpaces := true
		hasAnyContent := false
		for row := 0; row < len(lines); row++ {
			if col < len(lines[row]) {
				hasAnyContent = true
				if lines[row][col] != ' ' {
					allSpaces = false
					break
				}
			}
		}
		// only consider it a boundary if there's content and it's all spaces
		isBoundaryColumn[col] = hasAnyContent && allSpaces
	}

	// extract problems
	rawProblems := make([]mathProblemRaw, 0)
	startCol := -1

	for col := 0; col < maxWidth; col++ {
		if !isBoundaryColumn[col] {
			if startCol == -1 {
				startCol = col
			}
		} else {
			if startCol != -1 {
				problemWidth := col - startCol
				problemGrid := make([][]rune, len(lines))
				for row := 0; row < len(lines); row++ {
					problemGrid[row] = make([]rune, problemWidth)
					// copy what exists, pad with spaces if row is shorter
					for c := 0; c < problemWidth; c++ {
						srcCol := startCol + c
						if srcCol < len(lines[row]) {
							problemGrid[row][c] = lines[row][srcCol]
						} else {
							problemGrid[row][c] = ' '
						}
					}
				}
				rawProblems = append(rawProblems, makeMathProblemRaw(problemGrid))
				startCol = -1
			}
		}
	}

	if startCol != -1 {
		problemWidth := maxWidth - startCol
		problemGrid := make([][]rune, len(lines))
		for row := 0; row < len(lines); row++ {
			problemGrid[row] = make([]rune, problemWidth)
			// copy what exists, pad with spaces if row is shorter
			for c := 0; c < problemWidth; c++ {
				srcCol := startCol + c
				if srcCol < len(lines[row]) {
					problemGrid[row][c] = lines[row][srcCol]
				} else {
					problemGrid[row][c] = ' '
				}
			}
		}
		rawProblems = append(rawProblems, makeMathProblemRaw(problemGrid))
	}

	return rawProblems
}

// input is in this format:
// 123 328  51 64
//
//	45 64  387 23
//	 6 98  215 314
//
// *   +   *   +
// problems are arranged horizontally, separated by columns of spaces
// return a list of raw math problems; might make more sense to return constructed actual problems
func ReadInput(input string, columnar bool) []mathProblem {
	rawProblems := parseStringToRawProblems(input)
	problems := make([]mathProblem, 0)
	// if columnar, add in reverse order
	if columnar {
		for i := len(rawProblems) - 1; i >= 0; i-- {
			problems = append(problems, makeMathProblemColumnar(rawProblems[i].grid))
		}
	} else {
		for _, rawProblem := range rawProblems {
			problems = append(problems, makeMathProblem(rawProblem.grid))
		}
	}
	return problems
}

func SolveDay6Part1(input string) int {
	problems := ReadInput(input, false)
	total := 0
	for _, problem := range problems {
		answer := problem.Solve()
		total += answer
	}
	return total
}

func SolveDay6Part2(input string) int {
	problems := ReadInput(input, true)
	total := 0
	for _, problem := range problems {
		answer := problem.Solve()
		total += answer
	}
	return total
}
