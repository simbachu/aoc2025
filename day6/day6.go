package day6

import "strings"

type mathProblem struct {
	operands []int
	operator rune
}

// make a math problem from a 2D grid of runes
func makeMathProblem(raw [][]rune) mathProblem {
	operands := make([]int, 0)
	operator := ' '
	for _, row := range raw {
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

// input is in this format:
// 123 328  51 64
//
//	45 64  387 23
//	 6 98  215 314
//
// *   +   *   +
// problems are arranged horizontally, separated by columns of spaces
// return a list of raw math problems; might make more sense to return constructed actual problems
func ReadInput(input string) []mathProblemRaw {
	lineStrings := strings.Split(input, "\n")
	lines := make([][]rune, 0)
	maxWidth := 0

	for _, lineStr := range lineStrings {
		lineStr = strings.TrimRight(lineStr, " \r")
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

	for i := range lines {
		for len(lines[i]) < maxWidth {
			lines[i] = append(lines[i], ' ')
		}
	}

	// find the spaces between problems
	isBoundaryColumn := make([]bool, maxWidth)
	for col := 0; col < maxWidth; col++ {
		allSpaces := true
		for row := 0; row < len(lines); row++ {
			if lines[row][col] != ' ' {
				allSpaces = false
				break
			}
		}
		isBoundaryColumn[col] = allSpaces
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
				problemGrid := make([][]rune, len(lines))
				for row := 0; row < len(lines); row++ {
					problemGrid[row] = make([]rune, col-startCol)
					copy(problemGrid[row], lines[row][startCol:col])
				}
				rawProblems = append(rawProblems, makeMathProblemRaw(problemGrid))
				startCol = -1
			}
		}
	}

	if startCol != -1 {
		problemGrid := make([][]rune, len(lines))
		for row := 0; row < len(lines); row++ {
			problemGrid[row] = make([]rune, maxWidth-startCol)
			copy(problemGrid[row], lines[row][startCol:maxWidth])
		}
		rawProblems = append(rawProblems, makeMathProblemRaw(problemGrid))
	}

	return rawProblems
}

func SolveDay6Part1(input string) int {
	problems := ReadInput(input)
	total := 0
	for _, problem := range problems {
		answer := makeMathProblem(problem.grid).Solve()
		total += answer
	}
	return total
}

func SolveDay6Part2(input string) int {
	return 0
}
