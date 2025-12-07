package day7

import (
	"strings"
	"testing"
)

const kDay7SampleInput = `.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............`
const kDay7SampleOutputPart1 = 21
const kDay7SampleOutputPart2 = 40

func TestDiagramRow_AddExit(t *testing.T) {
	tests := []struct {
		name          string
		exitsToAdd    []int
		expectedAdded []bool
		expectedExits map[int]bool
	}{
		{
			name:          "Add single exit",
			exitsToAdd:    []int{5},
			expectedAdded: []bool{true},
			expectedExits: map[int]bool{5: true},
		},
		{
			name:          "Add multiple unique exits",
			exitsToAdd:    []int{3, 7, 10},
			expectedAdded: []bool{true, true, true},
			expectedExits: map[int]bool{3: true, 7: true, 10: true},
		},
		{
			name:          "Add duplicate exit",
			exitsToAdd:    []int{5, 5},
			expectedAdded: []bool{true, false},
			expectedExits: map[int]bool{5: true},
		},
		{
			name:          "Add multiple with duplicates",
			exitsToAdd:    []int{3, 7, 3, 10, 7},
			expectedAdded: []bool{true, true, false, true, false},
			expectedExits: map[int]bool{3: true, 7: true, 10: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			row := &DiagramRow{}
			for i, exitPos := range test.exitsToAdd {
				result := row.AddExit(exitPos)
				if result != test.expectedAdded[i] {
					t.Errorf("AddExit(%d) = %v, expected %v", exitPos, result, test.expectedAdded[i])
				}
			}
			if len(row.exits) != len(test.expectedExits) {
				t.Errorf("Expected %d exits, got %d", len(test.expectedExits), len(row.exits))
			}
			for pos, expected := range test.expectedExits {
				if actual, exists := row.exits[pos]; !exists || actual != expected {
					t.Errorf("Exit at position %d: expected %v, got %v (exists: %v)", pos, expected, actual, exists)
				}
			}
		})
	}
}

func TestNewDiagramRow_FirstRow(t *testing.T) {
	tests := []struct {
		name        string
		gridRow     string
		expectedPos int
	}{
		{
			name:        "S at position 7",
			gridRow:     ".......S.......",
			expectedPos: 7,
		},
		{
			name:        "S at position 0",
			gridRow:     "S..............",
			expectedPos: 0,
		},
		{
			name:        "No S in row",
			gridRow:     "...............",
			expectedPos: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gridRunes := []rune(test.gridRow)
			row := NewDiagramRow(gridRunes, 0, nil, func(int) {})
			if test.expectedPos == -1 {
				if len(row.exits) != 0 {
					t.Errorf("Expected no exits, got %d", len(row.exits))
				}
			} else {
				if len(row.exits) != 1 {
					t.Errorf("Expected 1 exit, got %d", len(row.exits))
				}
				if !row.exits[test.expectedPos] {
					t.Errorf("Expected exit at position %d, but it doesn't exist", test.expectedPos)
				}
			}
		})
	}
}

func TestNewDiagramRow_CarryThroughEmpty(t *testing.T) {
	// First row has S at position 7
	firstRow := NewDiagramRow([]rune(".......S......."), 0, nil, func(int) {})

	// Second row is all empty, should carry the exit
	secondRow := NewDiagramRow([]rune("..............."), 1, firstRow, func(int) {})

	if len(secondRow.exits) != 1 {
		t.Errorf("Expected 1 exit, got %d", len(secondRow.exits))
	}
	if !secondRow.exits[7] {
		t.Errorf("Expected exit at position 7, but it doesn't exist")
	}
}

func TestNewDiagramRow_Splitter(t *testing.T) {
	// First row has S at position 7
	firstRow := NewDiagramRow([]rune(".......S......."), 0, nil, func(int) {})

	// Second row is all empty, should carry the exit
	secondRow := NewDiagramRow([]rune("..............."), 1, firstRow, func(int) {})

	// Third row has splitter at position 7
	splitCount := 0
	splitPositions := make(map[int]bool)
	thirdRow := NewDiagramRow([]rune(".......^......."), 2, secondRow, func(splitterPos int) {
		splitCount++
		splitPositions[splitterPos] = true
	})

	// Should have exits at positions 6 and 8 (left and right of splitter)
	if len(thirdRow.exits) != 2 {
		t.Errorf("Expected 2 exits, got %d", len(thirdRow.exits))
	}
	if !thirdRow.exits[6] {
		t.Errorf("Expected exit at position 6, but it doesn't exist")
	}
	if !thirdRow.exits[8] {
		t.Errorf("Expected exit at position 8, but it doesn't exist")
	}

	// Should have counted 1 split (the splitter position at 7)
	if splitCount != 1 {
		t.Errorf("Expected 1 split to be counted, got %d", splitCount)
	}
	if len(splitPositions) != 1 {
		t.Errorf("Expected 1 unique split position, got %d", len(splitPositions))
	}
	if !splitPositions[7] {
		t.Errorf("Expected split position 7, got %v", splitPositions)
	}
}

func TestNewDiagramRow_SplitterWithBlockedSide(t *testing.T) {
	// First row has S at position 1
	firstRow := NewDiagramRow([]rune(".S."), 0, nil, func(int) {})

	// Second row is all empty, should carry the exit
	secondRow := NewDiagramRow([]rune("..."), 1, firstRow, func(int) {})

	// Third row has splitter at position 1
	// Left side (position 0) is valid and empty, right side (position 2) is valid and empty
	splitCount := 0
	splitPositions := make(map[int]bool)
	thirdRow := NewDiagramRow([]rune(".^."), 2, secondRow, func(splitterPos int) {
		splitCount++
		splitPositions[splitterPos] = true
	})

	// Should have exits at positions 0 and 2 (left and right of splitter)
	if len(thirdRow.exits) != 2 {
		t.Errorf("Expected 2 exits, got %d", len(thirdRow.exits))
	}
	if !thirdRow.exits[0] {
		t.Errorf("Expected exit at position 0, but it doesn't exist")
	}
	if !thirdRow.exits[2] {
		t.Errorf("Expected exit at position 2, but it doesn't exist")
	}

	// Should have counted 1 split (the splitter position at 1)
	if splitCount != 1 {
		t.Errorf("Expected 1 split to be counted, got %d", splitCount)
	}
	if !splitPositions[1] {
		t.Errorf("Expected split position 1, got %v", splitPositions)
	}
}

func TestNewDiagramRow_MultipleExitsToSplitter(t *testing.T) {
	// Create a previous row with multiple exits, some of which will hit the same splitter
	prevRow := &DiagramRow{}
	prevRow.AddExit(6)
	prevRow.AddExit(7)
	prevRow.AddExit(8)

	// Current row has splitter at position 7
	// Exit 6 -> position 6 ('.') -> no splitter
	// Exit 7 -> position 7 ('^') -> splitter -> callback
	// Exit 8 -> position 8 ('.') -> no splitter
	splitCount := 0
	splitPositions := make(map[int]int) // track how many times each position is called
	currentRow := NewDiagramRow([]rune(".......^......."), 1, prevRow, func(splitterPos int) {
		splitCount++
		splitPositions[splitterPos]++
	})

	// Exit 6 continues to position 6, exit 7 splits to 6 and 8, exit 8 continues to position 8
	// So we should have exits at 6 and 8
	if len(currentRow.exits) != 2 {
		t.Errorf("Expected 2 exits, got %d", len(currentRow.exits))
	}

	// Only exit 7 hits the splitter, so callback should be called once
	if splitCount != 1 {
		t.Errorf("Expected 1 split callback, got %d", splitCount)
	}
	if splitPositions[7] != 1 {
		t.Errorf("Expected split position 7 to be called 1 time, got %d", splitPositions[7])
	}
}

func TestSolveDay7Part1(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: kDay7SampleInput, expected: kDay7SampleOutputPart1},
	}
	for _, test := range tests {
		result := SolveDay7Part1(test.input)
		if result != test.expected {
			t.Errorf("SolveDay7Part1(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestDiagram_CountPathsFromS(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Sample input",
			input:    kDay7SampleInput,
			expected: kDay7SampleOutputPart2,
		},
		{
			name: "Simple S to splitter to end",
			input: `S
.
^`,
			expected: 1, // One path from S to the end (through the splitter)
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lines := strings.Split(test.input, "\n")
			linesAsRunes := make([][]rune, len(lines))
			for i, line := range lines {
				linesAsRunes[i] = []rune(line)
			}
			diagram := &Diagram{rows: linesAsRunes}
			result := diagram.CountPathsFromS()
			if result != test.expected {
				t.Errorf("CountPathsFromS() = %d, expected %d", result, test.expected)
			}
		})
	}
}

func TestSolveDay7Part2(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: kDay7SampleInput, expected: kDay7SampleOutputPart2},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := SolveDay7Part2(test.input)
			if result != test.expected {
				t.Errorf("SolveDay7Part2(%s) = %d, expected %d", test.input, result, test.expected)
			}
		})
	}
}

func TestCountPathsFromS_Optimized(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Sample input",
			input:    kDay7SampleInput,
			expected: kDay7SampleOutputPart2,
		},
		{
			name: "Simple S to splitter to end",
			input: `S
.
^`,
			expected: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lines := strings.Split(test.input, "\n")
			linesAsRunes := make([][]rune, len(lines))
			for i, line := range lines {
				linesAsRunes[i] = []rune(line)
			}
			diagram := &Diagram{rows: linesAsRunes}
			result := diagram.CountPathsFromS()
			if result != test.expected {
				t.Errorf("CountPathsFromS_Optimized() = %d, expected %d", result, test.expected)
			}
		})
	}
}
