package day10

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Indicator bool

type Button []int

func MakeButton(input string) (Button, error) {
	button := Button{}
	var err error
	if strings.Contains(input, ",") {
		numbers := strings.Split(input, ",")
		for _, number := range numbers {
			numberInt, err := strconv.Atoi(number)
			if err != nil {
				err = fmt.Errorf("invalid button input: %w", err)
				return button, err
			}
			button = append(button, numberInt)
		}
	} else {
		numberInt, err := strconv.Atoi(input)
		if err != nil {
			err = fmt.Errorf("invalid button input: %w", err)
			return button, err
		}
		button = append(button, numberInt)
	}
	return button, err
}

func (m *Machine) MakeButtonMatrix() [][]bool {
	width := len(m.Indicators)
	height := len(m.Buttons)
	matrix := make([][]bool, height)
	for i, button := range m.Buttons {
		matrix[i] = make([]bool, width)
		for _, indicator := range button {
			matrix[i][indicator] = true
		}
	}
	return matrix
}

func (m *Machine) GetGoalVector() []bool {
	vector := make([]bool, len(m.Goal))
	for i, goal := range m.Goal {
		vector[i] = bool(goal)
	}
	return vector
}

func buildAugmentedMatrix(matrix [][]bool, goalVector []bool) [][]bool {
	if len(matrix) == 0 {
		return nil
	}
	numIndicators := len(matrix[0])
	numButtons := len(matrix)

	augmentedMatrix := make([][]bool, numIndicators)
	for i := 0; i < numIndicators; i++ {
		augmentedMatrix[i] = make([]bool, numButtons+1)
		for j := 0; j < numButtons; j++ {
			augmentedMatrix[i][j] = matrix[j][i]
		}
		augmentedMatrix[i][numButtons] = goalVector[i]
	}
	return augmentedMatrix
}

func findPivot(augmentedMatrix [][]bool, col, startRow int) int {
	for row := startRow; row < len(augmentedMatrix); row++ {
		if augmentedMatrix[row][col] {
			return row
		}
	}
	return -1
}

func swapRows(augmentedMatrix [][]bool, row1, row2 int) {
	augmentedMatrix[row1], augmentedMatrix[row2] = augmentedMatrix[row2], augmentedMatrix[row1]
}

func eliminateColumn(augmentedMatrix [][]bool, col int, pivotRow int) {
	for row := 0; row < len(augmentedMatrix); row++ {
		if row == pivotRow {
			continue
		}
		if augmentedMatrix[row][col] {
			for i := 0; i < len(augmentedMatrix[row]); i++ {
				augmentedMatrix[row][i] = augmentedMatrix[row][i] != augmentedMatrix[pivotRow][i]
			}
		}
	}
}

func forwardElimination(augmentedMatrix [][]bool, numberOfButtons int) int {
	row := 0

	for col := 0; col < numberOfButtons && row < len(augmentedMatrix); col++ {
		pivotRow := findPivot(augmentedMatrix, col, row)
		if pivotRow == -1 {
			continue
		}
		if pivotRow != row {
			swapRows(augmentedMatrix, row, pivotRow)
		}
		eliminateColumn(augmentedMatrix, col, row)
		row++
	}
	return row
}

func checkConsistency(augmentedMatrix [][]bool, numberOfPivots int, numberOfIndicators int, numberOfButtons int) bool {
	for i := numberOfPivots; i < numberOfIndicators; i++ {
		allZeros := true
		for j := 0; j < numberOfButtons; j++ {
			if augmentedMatrix[i][j] {
				allZeros = false
				break
			}
		}
		if allZeros && augmentedMatrix[i][numberOfButtons] {
			return false
		}
	}
	return true
}

func findPivotColumn(row []bool, numberOfButtons int) int {
	for i := 0; i < numberOfButtons; i++ {
		if row[i] {
			return i
		}
	}
	return -1
}

func backSubstitution(augmentedMatrix [][]bool, numberOfPivots int, numberOfButtons int) (vector []bool, success bool) {
	vector = make([]bool, numberOfButtons)

	for i := numberOfPivots - 1; i >= 0; i-- {
		pivotColumn := findPivotColumn(augmentedMatrix[i], numberOfButtons)
		if pivotColumn == -1 {
			continue
		}
		value := augmentedMatrix[i][numberOfButtons]
		for j := pivotColumn + 1; j < numberOfButtons; j++ {
			if augmentedMatrix[i][j] && vector[j] {
				value = !value
			}
		}
		vector[pivotColumn] = value
	}
	return vector, true
}

func findMinimalSolution(augmentedMatrix [][]bool, numberOfPivots int, numberOfButtons int) (vector []bool, success bool) {
	particularSolution, ok := backSubstitution(augmentedMatrix, numberOfPivots, numberOfButtons)
	if !ok {
		return nil, false
	}

	pivotColumns := make(map[int]bool)
	for row := 0; row < numberOfPivots; row++ {
		for col := 0; col < numberOfButtons; col++ {
			if augmentedMatrix[row][col] {
				isPivot := true
				for c := 0; c < col; c++ {
					if augmentedMatrix[row][c] {
						isPivot = false
						break
					}
				}
				if isPivot {
					pivotColumns[col] = true
					break
				}
			}
		}
	}

	freeVars := []int{}
	for col := 0; col < numberOfButtons; col++ {
		if !pivotColumns[col] {
			freeVars = append(freeVars, col)
		}
	}

	if len(freeVars) == 0 {
		return particularSolution, true
	}

	nullSpaceVectors := make([][]bool, len(freeVars))
	for i, freeVar := range freeVars {
		nullSpaceVector := make([]bool, numberOfButtons)
		nullSpaceVector[freeVar] = true
		for row := numberOfPivots - 1; row >= 0; row-- {
			pivotColumn := findPivotColumn(augmentedMatrix[row], numberOfButtons)
			if pivotColumn == -1 {
				continue
			}
			value := false
			for j := pivotColumn + 1; j < numberOfButtons; j++ {
				if augmentedMatrix[row][j] && nullSpaceVector[j] {
					value = !value
				}
			}
			nullSpaceVector[pivotColumn] = value
		}
		nullSpaceVectors[i] = nullSpaceVector
	}

	bestSolution := make([]bool, numberOfButtons)
	copy(bestSolution, particularSolution)
	bestWeight := countTrue(particularSolution)

	numCombinations := 1 << len(freeVars)
	for combination := 0; combination < numCombinations; combination++ {
		currentSolution := make([]bool, numberOfButtons)
		copy(currentSolution, particularSolution)

		for i := range freeVars {
			if combination&(1<<i) != 0 {
				nullVec := nullSpaceVectors[i]
				for j := 0; j < numberOfButtons; j++ {
					if nullVec[j] {
						currentSolution[j] = !currentSolution[j]
					}
				}
			}
		}

		weight := countTrue(currentSolution)
		if weight < bestWeight {
			bestWeight = weight
			copy(bestSolution, currentSolution)
		}
	}

	return bestSolution, true
}

func countTrue(vector []bool) int {
	count := 0
	for _, v := range vector {
		if v {
			count++
		}
	}
	return count
}

func (m *Machine) Solve() (vector []bool, success bool) {
	matrix := m.MakeButtonMatrix()
	targetVector := m.GetGoalVector()
	augmentedMatrix := buildAugmentedMatrix(matrix, targetVector)
	numberOfIndicators := len(m.Indicators)
	numberOfButtons := len(m.Buttons)
	numberOfPivots := forwardElimination(augmentedMatrix, numberOfButtons)
	if !checkConsistency(augmentedMatrix, numberOfPivots, numberOfIndicators, numberOfButtons) {
		return nil, false
	}
	vector, success = findMinimalSolution(augmentedMatrix, numberOfPivots, numberOfButtons)
	if !success {
		return nil, false
	}
	return vector, success
}

type Joltage []int

func MakeJoltage(input string) (Joltage, error) {
	joltage := Joltage{}
	var err error
	if strings.Contains(input, ",") {
		numbers := strings.Split(input, ",")
		for _, number := range numbers {
			numberInt, err := strconv.Atoi(number)
			if err != nil {
				err = fmt.Errorf("invalid joltage input: %w", err)
				return joltage, err
			}
			joltage = append(joltage, numberInt)
		}
	} else {
		numberInt, err := strconv.Atoi(input)
		if err != nil {
			err = fmt.Errorf("invalid joltage input: %w", err)
			return joltage, err
		}
		joltage = append(joltage, numberInt)
	}
	return joltage, err
}

type Machine struct {
	Indicators []Indicator
	Buttons    []Button
	Joltage    Joltage
	Goal       []Indicator
}

func (m *Machine) Toggle(b Button) {
	for _, button := range b {
		m.Indicators[button] = !m.Indicators[button]
	}
}

func ReadInput(input string) []Machine {
	var machines []Machine
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		machine, err := parseMachine(line)
		if err != nil {
			fmt.Printf("warning: skipping invalid line %q: %v\n", line, err)
			continue
		}
		machines = append(machines, machine)
	}

	return machines
}

func parseMachine(line string) (Machine, error) {
	var m Machine

	goal, err := parseGoalState(line)
	if err != nil {
		return m, fmt.Errorf("parsing goal: %w", err)
	}
	m.Goal = goal
	m.Indicators = make([]Indicator, len(goal))

	buttons, err := parseButtons(line)
	if err != nil {
		return m, fmt.Errorf("parsing buttons: %w", err)
	}
	m.Buttons = buttons

	joltage, err := parseJoltage(line)
	if err != nil {
		return m, fmt.Errorf("parsing joltage: %w", err)
	}
	m.Joltage = joltage

	return m, nil
}

func parseGoalState(line string) ([]Indicator, error) {
	re := regexp.MustCompile(`\[([.#]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no goal pattern found in brackets")
	}

	pattern := matches[1]
	goal := make([]Indicator, 0, len(pattern))
	for _, char := range pattern {
		switch char {
		case '.':
			goal = append(goal, false)
		case '#':
			goal = append(goal, true)
		default:
			return nil, fmt.Errorf("invalid character %q in goal pattern", char)
		}
	}

	return goal, nil
}

func parseButtons(line string) ([]Button, error) {
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no buttons found")
	}

	var buttons []Button
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		button, err := MakeButton(match[1])
		if err != nil {
			return nil, fmt.Errorf("invalid button %q: %w", match[1], err)
		}
		buttons = append(buttons, button)
	}

	return buttons, nil
}

func parseJoltage(line string) (Joltage, error) {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no joltage found in braces")
	}

	return MakeJoltage(matches[1])
}

func SolveDay10Part1(input string) int {
	machines := ReadInput(input)
	result := 0
	for _, machine := range machines {
		vector, success := machine.Solve()
		if !success {
			continue
		}
		for _, button := range vector {
			if button {
				result++
			}
		}
	}
	return result
}

func SolveDay10Part2(input string) int {
	return 0
}
