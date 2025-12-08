package day8

import "testing"

const kDay8SampleInput = `162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689`

var kClosestPairSampleCoordinates = CoordinatePair{
	Coordinate1: Coordinate{x: 162, y: 817, z: 812},
	Coordinate2: Coordinate{x: 425, y: 690, z: 689},
}
var kSecondClosestPairSampleCoordinates = CoordinatePair{
	Coordinate1: Coordinate{x: 162, y: 817, z: 812},
	Coordinate2: Coordinate{x: 431, y: 825, z: 988},
}

var kGroupSizesFromSampleInput = []int{5, 4, 2, 2, 1, 1, 1, 1, 1, 1, 1}

const kDay8SampleOutputPart1 = 40
const kDay8SampleOutputPart2 = 25272

func TestReadInput(t *testing.T) {
	coordinates, err := ReadInput(kDay8SampleInput)
	if err != nil {
		t.Errorf("ReadInput(%s) unexpected error %v", kDay8SampleInput, err)
	}
	if len(coordinates) != 20 {
		t.Errorf("ReadInput(%s) = %d, expected 20", kDay8SampleInput, len(coordinates))
	}
}

func TestCoordinateFromString(t *testing.T) {

	coordinate, err := CoordinateFromString("162,817,812")
	if err != nil {
		t.Errorf("CoordinateFromString(%s) unexpected error %v", "162,817,812", err)
	}
	if coordinate.x != 162 || coordinate.y != 817 || coordinate.z != 812 {
		t.Errorf("CoordinateFromString(%s) = %v, expected 162,817,812", "162,817,812", coordinate)
	}
}

func TestCoordinatePairEquals(t *testing.T) {
	c1 := Coordinate{x: 1, y: 2, z: 3}
	c2 := Coordinate{x: 4, y: 5, z: 6}
	c3 := Coordinate{x: 7, y: 8, z: 9}

	pair1 := CoordinatePair{Coordinate1: c1, Coordinate2: c2}
	pair2 := CoordinatePair{Coordinate1: c2, Coordinate2: c1} // reversed order
	pair3 := CoordinatePair{Coordinate1: c1, Coordinate2: c3} // different pair

	// Test order-independent equality
	if !pair1.Equals(pair2) {
		t.Errorf("pair1.Equals(pair2) = false, expected true (order-independent)")
	}
	if !pair2.Equals(pair1) {
		t.Errorf("pair2.Equals(pair1) = false, expected true (order-independent)")
	}
	// Test same order equality
	if !pair1.Equals(pair1) {
		t.Errorf("pair1.Equals(pair1) = false, expected true")
	}
	// Test different pairs
	if pair1.Equals(pair3) {
		t.Errorf("pair1.Equals(pair3) = true, expected false")
	}
}

func TestCoordinatePairNotEquals(t *testing.T) {
	c1 := Coordinate{x: 1, y: 2, z: 3}
	c2 := Coordinate{x: 4, y: 5, z: 6}
	c3 := Coordinate{x: 7, y: 8, z: 9}

	pair1 := CoordinatePair{Coordinate1: c1, Coordinate2: c2}
	pair2 := CoordinatePair{Coordinate1: c1, Coordinate2: c3}

	// Test NotEquals
	if pair1.NotEquals(pair1) {
		t.Errorf("pair1.NotEquals(pair1) = true, expected false")
	}
	if !pair1.NotEquals(pair2) {
		t.Errorf("pair1.NotEquals(pair2) = false, expected true")
	}
}

func TestPairWithDistanceEquals(t *testing.T) {
	c1 := Coordinate{x: 1, y: 2, z: 3}
	c2 := Coordinate{x: 4, y: 5, z: 6}
	pair := CoordinatePair{Coordinate1: c1, Coordinate2: c2}

	pwd1 := PairWithDistance{Distance: 5.0, Pair: pair}
	pwd2 := PairWithDistance{Distance: 5.0, Pair: pair}
	pwd3 := PairWithDistance{Distance: 10.0, Pair: pair}

	// Test equality
	if !pwd1.Equals(pwd2) {
		t.Errorf("pwd1.Equals(pwd2) = false, expected true")
	}
	// Test different distance
	if pwd1.Equals(pwd3) {
		t.Errorf("pwd1.Equals(pwd3) = true, expected false")
	}
}

func TestPairsByDistance(t *testing.T) {
	// Arrange
	coordinates, err := ReadInput(kDay8SampleInput)
	if err != nil {
		t.Errorf("ReadInput(%s) unexpected error %v", kDay8SampleInput, err)
	}
	// Act
	pairs := PairsByDistance(coordinates, 2)
	// Assert
	// Access the first pair (shortest distance)
	if len(pairs) == 0 {
		t.Errorf("PairsByDistance(%s) returned empty slice", kDay8SampleInput)
		return
	}
	closestPair := pairs[0].Pair
	// Check if the pair matches
	if !closestPair.Equals(kClosestPairSampleCoordinates) {
		t.Errorf("PairsByDistance(%s) = %v, expected %v", kDay8SampleInput, closestPair, kClosestPairSampleCoordinates)
	}
	// Check the second pair
	secondPair := pairs[1].Pair
	if !secondPair.Equals(kSecondClosestPairSampleCoordinates) {
		t.Errorf("PairsByDistance(%s) = %v, expected %v", kDay8SampleInput, secondPair, kSecondClosestPairSampleCoordinates)
	}
}

func TestPairsByDistanceOrdering(t *testing.T) {
	// Arrange
	coordinates, err := ReadInput(kDay8SampleInput)
	if err != nil {
		t.Errorf("ReadInput(%s) unexpected error %v", kDay8SampleInput, err)
		return
	}

	// Act
	pairs := PairsByDistance(coordinates, 10)

	// Assert
	if len(pairs) == 0 {
		t.Errorf("PairsByDistance returned empty slice")
		return
	}

	// Check that pairs are in ascending order by distance
	for i := 1; i < len(pairs); i++ {
		if pairs[i].Distance < pairs[i-1].Distance {
			t.Errorf("PairsByDistance: pairs not sorted by distance. pairs[%d].Distance (%.2f) < pairs[%d].Distance (%.2f)",
				i, pairs[i].Distance, i-1, pairs[i-1].Distance)
		}
	}

	// Verify specific expected pairs
	expectedPairs := []CoordinatePair{
		{Coordinate1: Coordinate{x: 162, y: 817, z: 812}, Coordinate2: Coordinate{x: 425, y: 690, z: 689}},
		{Coordinate1: Coordinate{x: 162, y: 817, z: 812}, Coordinate2: Coordinate{x: 431, y: 825, z: 988}},
		{Coordinate1: Coordinate{x: 906, y: 360, z: 560}, Coordinate2: Coordinate{x: 805, y: 96, z: 715}},
		{Coordinate1: Coordinate{x: 431, y: 825, z: 988}, Coordinate2: Coordinate{x: 425, y: 690, z: 689}},
	}

	if len(pairs) < len(expectedPairs) {
		t.Errorf("PairsByDistance returned %d pairs, expected at least %d", len(pairs), len(expectedPairs))
		return
	}

	for i, expectedPair := range expectedPairs {
		if !pairs[i].Pair.Equals(expectedPair) {
			t.Errorf("PairsByDistance pairs[%d] = %v, expected %v", i, pairs[i].Pair, expectedPair)
		}
	}
}

func TestGroupCoordinatesFromInput(t *testing.T) {
	// Arrange
	coordinates, err := ReadInput(kDay8SampleInput)
	if err != nil {
		t.Errorf("ReadInput(%s) unexpected error %v", kDay8SampleInput, err)
	}
	pairs := PairsByDistance(coordinates, 10)
	groups := GroupCoordinates(pairs, coordinates)

	// Assert
	if len(groups) != len(kGroupSizesFromSampleInput) {
		t.Errorf("GroupCoordinates(%s) = %d groups, expected %d", kDay8SampleInput, len(groups), len(kGroupSizesFromSampleInput))
	}
	for i, group := range groups {
		if len(group) != kGroupSizesFromSampleInput[i] {
			t.Errorf("GroupCoordinates(%s) group[%d] = %d coordinates, expected %d", kDay8SampleInput, i, len(group), kGroupSizesFromSampleInput[i])
		}
	}
}

func TestSolveDay8Part1(t *testing.T) {
	result := SolveDay8Part1(kDay8SampleInput)
	if result != kDay8SampleOutputPart1 {
		t.Errorf("SolveDay8Part1(%s) = %d, expected %d", kDay8SampleInput, result, kDay8SampleOutputPart1)
	}
}

func TestSolveDay8Part2(t *testing.T) {
	result := SolveDay8Part2(kDay8SampleInput)
	if result != kDay8SampleOutputPart2 {
		t.Errorf("SolveDay8Part2(%s) = %d, expected %d", kDay8SampleInput, result, kDay8SampleOutputPart2)
	}
}

func TestGroupCoordinates(t *testing.T) {
	// Arrange
	c1 := Coordinate{x: 1, y: 2, z: 3}
	c2 := Coordinate{x: 4, y: 5, z: 6}
	c3 := Coordinate{x: 7, y: 8, z: 9}
	c4 := Coordinate{x: 10, y: 11, z: 12}
	coordinates := []Coordinate{c1, c2, c3, c4}

	// Create pairs: (c1,c2), (c2,c3), (c4,c4) - last one should create new group
	pairs := []PairWithDistance{
		{Distance: 1.0, Pair: CoordinatePair{Coordinate1: c1, Coordinate2: c2}},
		{Distance: 2.0, Pair: CoordinatePair{Coordinate1: c2, Coordinate2: c3}},
		{Distance: 3.0, Pair: CoordinatePair{Coordinate1: c4, Coordinate2: c4}},
	}

	// Act
	groups := GroupCoordinates(pairs, coordinates)

	// Assert
	// Should have 2 groups: {c1, c2, c3} and {c4}
	if len(groups) != 2 {
		t.Errorf("GroupCoordinates() = %d groups, expected 2", len(groups))
	}

	// First group should contain c1, c2, c3 (connected through pairs)
	foundGroup1 := false
	for _, group := range groups {
		if group[c1] && group[c2] && group[c3] && !group[c4] {
			foundGroup1 = true
			break
		}
	}
	if !foundGroup1 {
		t.Errorf("GroupCoordinates() did not create group with c1, c2, c3")
	}

	// Second group should contain only c4
	foundGroup2 := false
	for _, group := range groups {
		if group[c4] && !group[c1] && !group[c2] && !group[c3] {
			foundGroup2 = true
			break
		}
	}
	if !foundGroup2 {
		t.Errorf("GroupCoordinates() did not create separate group for c4")
	}
}
