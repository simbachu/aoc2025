package day9

import (
	"fmt"
	"strings"
	"testing"
)

const kDay9SampleInput = `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

const kDay9SampleOutputPart1 = 50 // 11,1 and 2,5
const kDay9SampleOutputPart2 = 24 // 9,5 and 2,3

func TestSolveDay9Part1(t *testing.T) {
	result := SolveDay9Part1(kDay9SampleInput)
	if result != kDay9SampleOutputPart1 {
		t.Errorf("SolveDay9Part1(%s) = %d, expected %d", kDay9SampleInput, result, kDay9SampleOutputPart1)
	}
}

func TestSolveDay9Part2(t *testing.T) {
	result := SolveDay9Part2(kDay9SampleInput)
	if result != kDay9SampleOutputPart2 {
		t.Errorf("SolveDay9Part2(%s) = %d, expected %d", kDay9SampleInput, result, kDay9SampleOutputPart2)
	}
}

func TestMakeRectangle(t *testing.T) {
	tests := []struct {
		name     string
		c1       Coordinate
		c2       Coordinate
		expected Rectangle
		wantErr  bool
	}{
		{
			name:     "(11,1) and (2,5)",
			c1:       Coordinate{x: 11, y: 1},
			c2:       Coordinate{x: 2, y: 5},
			expected: Rectangle{C1: Coordinate{x: 2, y: 1}, C2: Coordinate{x: 11, y: 5}},
			wantErr:  false,
		},
		{
			name:     "(2,5) and (11,1)",
			c1:       Coordinate{x: 2, y: 5},
			c2:       Coordinate{x: 11, y: 1},
			expected: Rectangle{C1: Coordinate{x: 2, y: 1}, C2: Coordinate{x: 11, y: 5}},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeRectangle(tt.c1, tt.c2)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeRectangle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.C1.x != tt.expected.C1.x || got.C1.y != tt.expected.C1.y ||
					got.C2.x != tt.expected.C2.x || got.C2.y != tt.expected.C2.y {
					t.Errorf("MakeRectangle() = %v, expected %v", got, tt.expected)
				}
				area := got.Area()
				t.Logf("Rectangle %v has area %d", got, area)
			}
		})
	}
}

func TestMakePerimeter(t *testing.T) {
	coordinates := ReadInput(kDay9SampleInput, false)
	perimeter := MakePerimeter(coordinates)
	t.Logf("Perimeter: %v", perimeter)
	drawing := perimeter.String()
	fmt.Printf("\n=== Perimeter Drawing ===\n%s\n========================\n\n", drawing)
	t.Logf("Perimeter drawing:\n%s", drawing)

	// Verify that all perimeter nodes are marked with #
	for _, c := range perimeter {
		lines := strings.Split(drawing, "\n")
		if c.y < len(lines) && c.x < len(lines[c.y]) {
			char := string(lines[c.y][c.x])
			if char != "#" {
				t.Errorf("Perimeter node %v should be marked with #, but found %s", c, char)
			}
		}
	}
}

func TestFindLargestRectangleInPerimeter(t *testing.T) {
	coordinates, perimeter := ReadInputAndMakePerimeter(kDay9SampleInput)
	largestRectangle, err := FindLargestRectangleInPerimeter(perimeter, coordinates)
	if err != nil {
		t.Fatalf("FindLargestRectangleInPerimeter(%v, %v) returned error: %v", perimeter, coordinates, err)
	}
	t.Logf("Largest rectangle in perimeter: %v", largestRectangle)
}

func TestAreaCalculation(t *testing.T) {
	c1 := Coordinate{x: 11, y: 1}
	c2 := Coordinate{x: 2, y: 5}

	rectangle, err := MakeRectangle(c1, c2)
	if err != nil {
		t.Fatalf("MakeRectangle(%v, %v) returned error: %v", c1, c2, err)
	}

	area := rectangle.Area()
	t.Logf("Rectangle from (%d,%d) and (%d,%d): %v", c1.x, c1.y, c2.x, c2.y, rectangle)
	t.Logf("Area = %d", area)

	if area != 50 {
		t.Errorf("Area = %d, expected 50", area)
	}
}
