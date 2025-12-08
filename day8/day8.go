package day8

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Coordinate struct {
	x int
	y int
	z int
}

type CoordinatePair struct {
	Coordinate1 Coordinate
	Coordinate2 Coordinate
}

func (p CoordinatePair) Equals(other CoordinatePair) bool {
	return (p.Coordinate1 == other.Coordinate1 && p.Coordinate2 == other.Coordinate2) ||
		(p.Coordinate1 == other.Coordinate2 && p.Coordinate2 == other.Coordinate1)
}

func (p CoordinatePair) NotEquals(other CoordinatePair) bool {
	return !p.Equals(other)
}

func (p CoordinatePair) String() string {
	return fmt.Sprintf("{%v %v}", p.Coordinate1, p.Coordinate2)
}

func (c Coordinate) Distance(other Coordinate) float64 {
	dx := float64(c.x - other.x)
	dy := float64(c.y - other.y)
	dz := float64(c.z - other.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c Coordinate) Equals(other Coordinate) bool {
	return c.x == other.x && c.y == other.y && c.z == other.z
}

func (c Coordinate) NotEquals(other Coordinate) bool {
	return !c.Equals(other)
}

func (c Coordinate) String() string {
	return fmt.Sprintf("%d,%d,%d", c.x, c.y, c.z)
}

func CoordinateFromString(s string) (Coordinate, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return Coordinate{}, errors.New("invalid coordinate string, expected 3 parts")
	}
	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Coordinate{}, err
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Coordinate{}, err
	}
	z, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return Coordinate{}, err
	}
	return Coordinate{x: x, y: y, z: z}, nil
}

func ReadInput(input string) ([]Coordinate, error) {
	lines := strings.Split(input, "\n")
	coordinates := []Coordinate{}
	for _, line := range lines {
		coordinate, err := CoordinateFromString(line)
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, coordinate)
	}
	return coordinates, nil
}

type PairWithDistance struct {
	Distance float64
	Pair     CoordinatePair
}

func (p PairWithDistance) Equals(other PairWithDistance) bool {
	return p.Distance == other.Distance && p.Pair.Equals(other.Pair)
}

func (p PairWithDistance) String() string {
	return fmt.Sprintf("{Distance: %.2f, Pair: %v}", p.Distance, p.Pair)
}

type pairMaxHeap []PairWithDistance

func (h pairMaxHeap) Len() int           { return len(h) }
func (h pairMaxHeap) Less(i, j int) bool { return h[i].Distance > h[j].Distance }
func (h pairMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *pairMaxHeap) Push(x interface{}) {
	*h = append(*h, x.(PairWithDistance))
}

func (h *pairMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func PairsByDistance(coordinates []Coordinate, limit int) []PairWithDistance {
	if limit <= 0 {
		return []PairWithDistance{}
	}

	h := &pairMaxHeap{}
	heap.Init(h)

	for i, coordinate := range coordinates {
		for j := i + 1; j < len(coordinates); j++ {
			otherCoordinate := coordinates[j]
			distance := coordinate.Distance(otherCoordinate)
			pair := PairWithDistance{
				Distance: distance,
				Pair:     CoordinatePair{Coordinate1: coordinate, Coordinate2: otherCoordinate},
			}

			if h.Len() < limit {
				heap.Push(h, pair)
			} else if distance < (*h)[0].Distance {
				(*h)[0] = pair
				heap.Fix(h, 0)
			}
		}
	}

	pairs := make([]PairWithDistance, h.Len())
	for i := 0; i < h.Len(); i++ {
		pairs[i] = (*h)[i]
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Distance < pairs[j].Distance
	})

	return pairs
}

func GroupCoordinates(pairs []PairWithDistance, allCoordinates []Coordinate) []map[Coordinate]bool {
	groups := []map[Coordinate]bool{}
	coordinatesInGroups := make(map[Coordinate]bool)

	for _, pair := range pairs {
		if len(groups) == 0 {
			groups = append(groups, map[Coordinate]bool{pair.Pair.Coordinate1: true, pair.Pair.Coordinate2: true})
			coordinatesInGroups[pair.Pair.Coordinate1] = true
			coordinatesInGroups[pair.Pair.Coordinate2] = true
		} else {
			group1Index := -1
			group2Index := -1
			for i, group := range groups {
				hasC1 := group[pair.Pair.Coordinate1]
				hasC2 := group[pair.Pair.Coordinate2]
				if hasC1 {
					group1Index = i
				}
				if hasC2 {
					group2Index = i
				}
			}

			if group1Index == -1 && group2Index == -1 {
				groups = append(groups, map[Coordinate]bool{pair.Pair.Coordinate1: true, pair.Pair.Coordinate2: true})
				coordinatesInGroups[pair.Pair.Coordinate1] = true
				coordinatesInGroups[pair.Pair.Coordinate2] = true
			} else if group1Index != -1 && group2Index == -1 {
				groups[group1Index][pair.Pair.Coordinate2] = true
				coordinatesInGroups[pair.Pair.Coordinate2] = true
			} else if group1Index == -1 && group2Index != -1 {
				groups[group2Index][pair.Pair.Coordinate1] = true
				coordinatesInGroups[pair.Pair.Coordinate1] = true
			} else if group1Index != group2Index {
				for coord := range groups[group2Index] {
					groups[group1Index][coord] = true
				}
				groups = append(groups[:group2Index], groups[group2Index+1:]...)
			}
		}
	}

	for _, coord := range allCoordinates {
		if !coordinatesInGroups[coord] {
			groups = append(groups, map[Coordinate]bool{coord: true})
		}
	}

	sort.Slice(groups, func(i, j int) bool {
		if len(groups[i]) != len(groups[j]) {
			return len(groups[i]) > len(groups[j])
		}
		var coordI, coordJ Coordinate
		for c := range groups[i] {
			coordI = c
			break
		}
		for c := range groups[j] {
			coordJ = c
			break
		}
		if coordI.x != coordJ.x {
			return coordI.x < coordJ.x
		}
		if coordI.y != coordJ.y {
			return coordI.y < coordJ.y
		}
		return coordI.z < coordJ.z
	})

	return groups
}

func getNPairsForPart1(numCoordinates int) int {
	if numCoordinates <= 20 {
		return 10
	}
	return 1000
}

func SolveDay8Part1(input string) interface{} {
	coordinates, err := ReadInput(input)
	if err != nil {
		return -1
	}
	n := getNPairsForPart1(len(coordinates))
	pairs := PairsByDistance(coordinates, n)
	groups := GroupCoordinates(pairs, coordinates)

	if len(groups) < 3 {
		return 0
	}
	return len(groups[0]) * len(groups[1]) * len(groups[2])
}

func FindFirstConnectingPair(initialGroups []map[Coordinate]bool, allPairs []PairWithDistance, pairsProcessed map[CoordinatePair]bool) (CoordinatePair, bool) {
	groups := make([]map[Coordinate]bool, len(initialGroups))
	for i, group := range initialGroups {
		groups[i] = make(map[Coordinate]bool)
		for coord := range group {
			groups[i][coord] = true
		}
	}

	for _, pair := range allPairs {
		if pairsProcessed[pair.Pair] {
			continue
		}

		group1Index := -1
		group2Index := -1
		for i, group := range groups {
			hasC1 := group[pair.Pair.Coordinate1]
			hasC2 := group[pair.Pair.Coordinate2]
			if hasC1 {
				group1Index = i
			}
			if hasC2 {
				group2Index = i
			}
		}

		if group1Index != -1 && group2Index != -1 && group1Index != group2Index {
			for coord := range groups[group2Index] {
				groups[group1Index][coord] = true
			}
			groups = append(groups[:group2Index], groups[group2Index+1:]...)

			if len(groups) == 1 {
				return pair.Pair, true
			}
		} else if group1Index == -1 && group2Index == -1 {
			groups[0][pair.Pair.Coordinate1] = true
			groups[0][pair.Pair.Coordinate2] = true
		} else if group1Index != -1 && group2Index == -1 {
			groups[group1Index][pair.Pair.Coordinate2] = true
		} else if group1Index == -1 && group2Index != -1 {
			groups[group2Index][pair.Pair.Coordinate1] = true
		}
	}

	return CoordinatePair{}, false
}

func SolveDay8Part2(input string) interface{} {
	coordinates, err := ReadInput(input)
	if err != nil {
		return -1
	}

	n := getNPairsForPart1(len(coordinates))
	pairs := PairsByDistance(coordinates, n)
	groups := GroupCoordinates(pairs, coordinates)

	if len(groups) <= 1 {
		return 0
	}

	allPairs := PairsByDistance(coordinates, len(coordinates)*(len(coordinates)-1)/2)

	pairsProcessed := make(map[CoordinatePair]bool)
	for _, pair := range pairs {
		pairsProcessed[pair.Pair] = true
	}

	connectingPair, found := FindFirstConnectingPair(groups, allPairs, pairsProcessed)
	if !found {
		return 0
	}

	return connectingPair.Coordinate1.x * connectingPair.Coordinate2.x
}
