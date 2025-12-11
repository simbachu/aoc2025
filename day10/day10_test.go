package day10

import (
	"fmt"
	"testing"
)

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

const sampleInput = `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`
const sampleOutPutPart1 = 7
const sampleOutPutPart2 = 33

func TestReadInput(t *testing.T) {
	t.Run("Sample Input", func(t *testing.T) {
		machines := ReadInput(sampleInput)
		if len(machines) != 3 {
			t.Fatalf("Expected 3 machines, got %d", len(machines))
		}

		expected := []struct {
			indicators int
			buttons    int
			joltages   int
		}{
			{4, 6, 4},
			{5, 5, 5},
			{6, 4, 6},
		}

		for i, machine := range machines {
			exp := expected[i]
			if len(machine.Indicators) != exp.indicators {
				t.Errorf("Machine %d: Expected %d indicators, got %d", i+1, exp.indicators, len(machine.Indicators))
			}
			if len(machine.Goal) != exp.indicators {
				t.Errorf("Machine %d: Expected %d goal states, got %d", i+1, exp.indicators, len(machine.Goal))
			}
			if len(machine.Buttons) != exp.buttons {
				t.Errorf("Machine %d: Expected %d buttons, got %d", i+1, exp.buttons, len(machine.Buttons))
			}
			if len(machine.Joltages) != exp.joltages {
				t.Errorf("Machine %d: Expected %d joltages, got %d", i+1, exp.joltages, len(machine.Joltages))
			}
		}
	})
}

func TestMakeButtonMatrix(t *testing.T) {
	machines := ReadInput(sampleInput)
	for i, machine := range machines {
		matrix := machine.MakeButtonMatrix()
		if len(matrix) != len(machine.Buttons) {
			t.Errorf("Machine %d: Expected %d rows, got %d", i+1, len(machine.Buttons), len(matrix))
		}
		if len(matrix[0]) != len(machine.Indicators) {
			t.Errorf("Machine %d: Expected %d columns, got %d", i+1, len(machine.Indicators), len(matrix[0]))
		}
		for j, row := range matrix {
			for k, cell := range row {
				expected := contains(machine.Buttons[j], k)
				if cell != expected {
					t.Errorf("Machine %d: Row %d, Col %d: Expected %t, got %t", i+1, j, k, expected, cell)
				}
			}
		}
	}
}

func TestGetGoalVector(t *testing.T) {
	machines := ReadInput(sampleInput)

	tests := []struct {
		name     string
		machine  Machine
		expected []bool
	}{
		{
			name:     "Machine 1",
			machine:  machines[0],
			expected: []bool{false, true, true, false},
		},
		{
			name:     "Machine 2",
			machine:  machines[1],
			expected: []bool{false, false, false, true, false},
		},
		{
			name:     "Machine 3",
			machine:  machines[2],
			expected: []bool{false, true, true, true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.machine.GetGoalVector()
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, exp := range tt.expected {
				if result[i] != exp {
					t.Errorf("Index %d: Expected %t, got %t", i, exp, result[i])
				}
			}
		})
	}
}

func TestBuildAugmentedMatrix(t *testing.T) {
	matrix := [][]bool{
		{true, true, false},
		{false, true, true},
	}
	goalVector := []bool{true, false, true}

	augmented := buildAugmentedMatrix(matrix, goalVector)

	expected := [][]bool{
		{true, false, true},
		{true, true, false},
		{false, true, true},
	}

	if len(augmented) != len(expected) {
		t.Fatalf("Expected %d rows, got %d", len(expected), len(augmented))
	}

	for i, row := range augmented {
		if len(row) != len(expected[i]) {
			t.Fatalf("Row %d: Expected %d columns, got %d", i, len(expected[i]), len(row))
		}
		for j, val := range row {
			if val != expected[i][j] {
				t.Errorf("Row %d, Col %d: Expected %t, got %t", i, j, expected[i][j], val)
			}
		}
	}
}

func TestForwardElimination(t *testing.T) {
	augmentedMatrix := [][]bool{
		{true, true, true},
		{true, false, false},
	}

	numPivots := forwardElimination(augmentedMatrix, 2)

	if numPivots != 2 {
		t.Errorf("Expected 2 pivots, got %d", numPivots)
	}

	if !augmentedMatrix[0][0] {
		t.Error("Expected pivot at [0][0]")
	}
	if augmentedMatrix[1][0] {
		t.Error("Expected column 0 to be eliminated in row 1")
	}
}

func TestCheckConsistency(t *testing.T) {
	tests := []struct {
		name          string
		augmented     [][]bool
		numPivots     int
		numIndicators int
		numButtons    int
		expected      bool
		description   string
	}{
		{
			name: "Consistent system",
			augmented: [][]bool{
				{true, false, true},
				{false, true, false},
			},
			numPivots:     2,
			numIndicators: 2,
			numButtons:    2,
			expected:      true,
			description:   "2 pivots, 2 equations, consistent",
		},
		{
			name: "Inconsistent system - 0=1",
			augmented: [][]bool{
				{true, false, true},
				{false, false, true},
			},
			numPivots:     1,
			numIndicators: 2,
			numButtons:    2,
			expected:      false,
			description:   "1 pivot, row 1 has 0=1, inconsistent",
		},
		{
			name: "Consistent overdetermined",
			augmented: [][]bool{
				{true, false, true},
				{false, true, false},
				{false, false, false},
			},
			numPivots:     2,
			numIndicators: 3,
			numButtons:    2,
			expected:      true,
			description:   "2 pivots, 3 equations, extra row is 0=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkConsistency(tt.augmented, tt.numPivots, tt.numIndicators, tt.numButtons)
			if result != tt.expected {
				t.Errorf("%s: Expected %t, got %t", tt.description, tt.expected, result)
			}
		})
	}
}

func TestFindFreeVariables(t *testing.T) {
	machines := ReadInput(sampleInput)

	for i, machine := range machines {
		t.Run(fmt.Sprintf("Machine %d", i+1), func(t *testing.T) {
			matrix := machine.MakeButtonMatrix()
			goalVector := machine.GetGoalVector()
			augmentedMatrix := buildAugmentedMatrix(matrix, goalVector)
			numButtons := len(machine.Buttons)
			numIndicators := len(machine.Indicators)

			numPivots := forwardElimination(augmentedMatrix, numButtons)
			t.Logf("Machine %d: %d pivots out of %d buttons, %d indicators",
				i+1, numPivots, numButtons, numIndicators)

			pivotColumns := make(map[int]bool)
			for row := 0; row < numPivots; row++ {
				for col := 0; col < numButtons; col++ {
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

			freeVars := numButtons - len(pivotColumns)
			t.Logf("Machine %d: %d pivot columns, %d free variables",
				i+1, len(pivotColumns), freeVars)
		})
	}
}

func TestVerifySolution(t *testing.T) {
	machines := ReadInput(sampleInput)

	for i, machine := range machines {
		t.Run(fmt.Sprintf("Machine %d", i+1), func(t *testing.T) {
			vector, success := machine.Solve()
			if !success {
				t.Fatalf("Machine %d: Solve failed", i+1)
			}

			testMachine := machine
			testMachine.Indicators = make([]Indicator, len(machine.Indicators))
			for j, pressed := range vector {
				if pressed {
					testMachine.Toggle(machine.Buttons[j])
				}
			}

			for j, goal := range machine.Goal {
				if testMachine.Indicators[j] != goal {
					t.Errorf("Machine %d, Indicator %d: Expected %t, got %t. Solution: %v",
						i+1, j, bool(goal), bool(testMachine.Indicators[j]), vector)
				}
			}
		})
	}
}

func TestSolveIndividualMachines(t *testing.T) {
	machines := ReadInput(sampleInput)

	expectedPresses := []int{2, 3, 2}

	for i, machine := range machines {
		t.Run(fmt.Sprintf("Machine %d", i+1), func(t *testing.T) {
			vector, success := machine.Solve()
			if !success {
				t.Fatalf("Machine %d: Solve failed", i+1)
			}

			pressCount := 0
			for _, pressed := range vector {
				if pressed {
					pressCount++
				}
			}

			if pressCount != expectedPresses[i] {
				t.Errorf("Machine %d: Expected %d presses, got %d. Vector: %v",
					i+1, expectedPresses[i], pressCount, vector)
			}
		})
	}
}

func TestSolveDay10Part1(t *testing.T) {
	result := SolveDay10Part1(sampleInput)
	if result != sampleOutPutPart1 {
		t.Errorf("Expected %d, got %d", sampleOutPutPart1, result)
	}
}

func TestSolveDay10Part2(t *testing.T) {
	result := SolveDay10Part2(sampleInput)
	if result != sampleOutPutPart2 {
		t.Errorf("Expected %d, got %d", sampleOutPutPart2, result)
	}
}

func TestSolvableSolve(t *testing.T) {
	tests := []struct {
		name          string
		solvable      Solvable
		shouldSucceed bool
		description   string
	}{
		{
			name: "Simple solvable system",
			solvable: Solvable{
				Workspace: []int{0, 0},
				Goal:      []int{2, 3},
				Buttons: []Button{
					{0},
					{1},
				},
			},
			shouldSucceed: true,
			description:   "Each button affects one position",
		},
		{
			name: "System with overlapping buttons",
			solvable: Solvable{
				Workspace: []int{0, 0},
				Goal:      []int{3, 4},
				Buttons: []Button{
					{0, 1},
					{1},
				},
			},
			shouldSucceed: true,
			description:   "Overlapping buttons",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vector, success := tt.solvable.Solve()
			if success != tt.shouldSucceed {
				t.Errorf("%s: Expected success=%t, got %t", tt.description, tt.shouldSucceed, success)
				return
			}
			if success {
				testWorkspace := make([]int, len(tt.solvable.Workspace))
				copy(testWorkspace, tt.solvable.Workspace)
				for i, presses := range vector {
					for j := 0; j < presses; j++ {
						for _, indicator := range tt.solvable.Buttons[i] {
							testWorkspace[indicator]++
						}
					}
				}
				for i, goal := range tt.solvable.Goal {
					if testWorkspace[i] != goal {
						t.Errorf("%s: Position %d: Expected %d, got %d. Solution: %v",
							tt.description, i, goal, testWorkspace[i], vector)
					}
				}
			}
		})
	}
}
