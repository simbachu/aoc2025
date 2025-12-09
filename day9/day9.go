package day9

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Coordinate struct {
	x int
	y int
}

type Perimeter []Coordinate

type PerimeterIndex struct {
	perimeter          Perimeter
	bounds             Rectangle
	insideCache        map[Coordinate]bool
	outsideCache       map[Coordinate]bool
	inputCoordinateSet map[Coordinate]bool
	cacheMu            sync.RWMutex
}

func NewPerimeterIndex(perimeter Perimeter) *PerimeterIndex {
	if len(perimeter) < 4 {
		return &PerimeterIndex{
			perimeter:          perimeter,
			insideCache:        make(map[Coordinate]bool),
			outsideCache:       make(map[Coordinate]bool),
			inputCoordinateSet: make(map[Coordinate]bool),
			cacheMu:            sync.RWMutex{},
		}
	}

	minX, minY := perimeter[0].x, perimeter[0].y
	maxX, maxY := perimeter[0].x, perimeter[0].y
	for _, c := range perimeter {
		if c.x < minX {
			minX = c.x
		}
		if c.x > maxX {
			maxX = c.x
		}
		if c.y < minY {
			minY = c.y
		}
		if c.y > maxY {
			maxY = c.y
		}
	}

	bounds := Rectangle{
		C1: Coordinate{x: minX, y: minY},
		C2: Coordinate{x: maxX, y: maxY},
	}

	return &PerimeterIndex{
		perimeter:          perimeter,
		bounds:             bounds,
		insideCache:        make(map[Coordinate]bool),
		outsideCache:       make(map[Coordinate]bool),
		inputCoordinateSet: make(map[Coordinate]bool),
		cacheMu:            sync.RWMutex{},
	}
}

func (pi *PerimeterIndex) PrecomputeCoordinates(coordinates []Coordinate) {
	for _, c := range coordinates {
		pi.inputCoordinateSet[c] = true
	}

	for _, c := range coordinates {
		if _, inInside := pi.insideCache[c]; inInside {
			continue
		}
		if _, inOutside := pi.outsideCache[c]; inOutside {
			continue
		}

		if !pi.bounds.Contains(c) {
			pi.outsideCache[c] = true
			continue
		}

		result := pi.perimeter.Contains(c)
		if result {
			pi.insideCache[c] = true
		} else {
			pi.outsideCache[c] = true
		}
	}
}

func (pi *PerimeterIndex) Contains(c Coordinate) bool {
	pi.cacheMu.RLock()
	if result, ok := pi.insideCache[c]; ok {
		pi.cacheMu.RUnlock()
		return result
	}
	if _, ok := pi.outsideCache[c]; ok {
		pi.cacheMu.RUnlock()
		return false
	}
	pi.cacheMu.RUnlock()

	if pi.inputCoordinateSet[c] {
		if !pi.bounds.Contains(c) {
			pi.cacheMu.Lock()
			pi.outsideCache[c] = true
			pi.cacheMu.Unlock()
			return false
		}
		result := pi.perimeter.Contains(c)
		pi.cacheMu.Lock()
		if result {
			pi.insideCache[c] = true
		} else {
			pi.outsideCache[c] = true
		}
		pi.cacheMu.Unlock()
		return result
	}

	if !pi.bounds.Contains(c) {
		pi.cacheMu.Lock()
		pi.outsideCache[c] = true
		pi.cacheMu.Unlock()
		return false
	}

	result := pi.perimeter.Contains(c)
	pi.cacheMu.Lock()
	if result {
		pi.insideCache[c] = true
	} else {
		pi.outsideCache[c] = true
	}
	pi.cacheMu.Unlock()
	return result
}

func MakePerimeter(coordinates []Coordinate) Perimeter {
	if len(coordinates) == 0 {
		return Perimeter{}
	}
	perimeter := make(Perimeter, 0, len(coordinates)+1)
	perimeter = append(perimeter, coordinates...)
	perimeter = append(perimeter, coordinates[0])
	return perimeter
}

func (p Perimeter) Contains(c Coordinate) bool {
	if len(p) < 4 {
		return false
	}

	for i := 0; i < len(p)-1; i++ {
		p1 := p[i]
		p2 := p[i+1]

		dx := p2.x - p1.x
		dy := p2.y - p1.y
		dx1 := c.x - p1.x
		dy1 := c.y - p1.y

		cross := dx1*dy - dy1*dx
		if cross == 0 {
			if dx != 0 {
				if (dx > 0 && c.x >= p1.x && c.x <= p2.x) || (dx < 0 && c.x >= p2.x && c.x <= p1.x) {
					return true
				}
			} else {
				if (dy > 0 && c.y >= p1.y && c.y <= p2.y) || (dy < 0 && c.y >= p2.y && c.y <= p1.y) {
					return true
				}
			}
		}
	}

	inside := false
	for i := 0; i < len(p)-1; i++ {
		p1 := p[i]
		p2 := p[i+1]

		if (p1.y > c.y) != (p2.y > c.y) {
			if p2.y != p1.y {
				xIntersect := float64(p1.x) + float64(c.y-p1.y)*float64(p2.x-p1.x)/float64(p2.y-p1.y)
				if float64(c.x) < xIntersect {
					inside = !inside
				}
			}
		}
	}

	return inside
}

func (p Perimeter) IsOnEdge(c Coordinate) bool {
	if len(p) < 4 {
		return false
	}

	for i := 0; i < len(p)-1; i++ {
		p1 := p[i]
		p2 := p[i+1]

		dx := p2.x - p1.x
		dy := p2.y - p1.y
		dx1 := c.x - p1.x
		dy1 := c.y - p1.y

		cross := dx1*dy - dy1*dx
		if cross == 0 {
			if dx != 0 {
				if (dx > 0 && c.x >= p1.x && c.x <= p2.x) || (dx < 0 && c.x >= p2.x && c.x <= p1.x) {
					return true
				}
			} else {
				if (dy > 0 && c.y >= p1.y && c.y <= p2.y) || (dy < 0 && c.y >= p2.y && c.y <= p1.y) {
					return true
				}
			}
		}
	}
	return false
}

// ..............
// .......#XXX#..
// .......XXXXX..
// ..#XXXX#XXXX..
// ..XXXXXXXXXX..
// ..#XXXXXX#XX..
// .........XXX..
// .........#X#..
// ..............
func (p Perimeter) String() string {
	maxX := 0
	maxY := 0
	for _, c := range p {
		if c.x > maxX {
			maxX = c.x
		}
		if c.y > maxY {
			maxY = c.y
		}
	}

	nodeSet := make(map[Coordinate]bool)
	for _, c := range p {
		nodeSet[c] = true
	}

	grid := make([][]string, maxY+1)
	for i := range grid {
		grid[i] = make([]string, maxX+1)
		for j := range grid[i] {
			coord := Coordinate{x: j, y: i}
			if nodeSet[coord] {
				grid[i][j] = "#"
			} else if p.Contains(coord) {
				grid[i][j] = "X"
			} else {
				grid[i][j] = "."
			}
		}
	}
	gridStrings := make([]string, len(grid))
	for i, row := range grid {
		gridStrings[i] = strings.Join(row, "")
	}
	return strings.Join(gridStrings, "\n")
}

func CoordinateFromString(s string) (Coordinate, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return Coordinate{}, errors.New("invalid coordinate string, expected 2 parts")
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return Coordinate{}, err
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return Coordinate{}, err
	}
	return Coordinate{x: x, y: y}, nil
}

func (c Coordinate) Min(other Coordinate) Coordinate {
	return Coordinate{x: min(c.x, other.x), y: min(c.y, other.y)}
}

func (c Coordinate) Max(other Coordinate) Coordinate {
	return Coordinate{x: max(c.x, other.x), y: max(c.y, other.y)}
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.x, c.y)
}

type Rectangle struct {
	C1 Coordinate
	C2 Coordinate
}

func (r Rectangle) Contains(c Coordinate) bool {
	return c.x >= r.C1.x && c.x <= r.C2.x && c.y >= r.C1.y && c.y <= r.C2.y
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle{C1: %v, C2: %v}", r.C1, r.C2)
}

func MakeRectangle(c1 Coordinate, c2 Coordinate) (Rectangle, error) {
	minCoord := c1.Min(c2)
	maxCoord := c1.Max(c2)

	if minCoord.x == maxCoord.x || minCoord.y == maxCoord.y {
		return Rectangle{}, errors.New("coordinates must form a valid rectangle with non-zero width and height")
	}

	return Rectangle{C1: minCoord, C2: maxCoord}, nil
}

func (r Rectangle) Area() int {
	return (r.C2.x - r.C1.x + 1) * (r.C2.y - r.C1.y + 1)
}

func (r Rectangle) Corners() []Coordinate {
	return []Coordinate{
		r.C1,
		{x: r.C2.x, y: r.C1.y},
		{x: r.C1.x, y: r.C2.y},
		r.C2,
	}
}

func (r Rectangle) IsFullyContainedInPerimeter(perimeter Perimeter) bool {
	for x := r.C1.x; x <= r.C2.x; x++ {
		point := Coordinate{x: x, y: r.C2.y}
		if !perimeter.Contains(point) {
			return false
		}
	}
	for x := r.C1.x; x <= r.C2.x; x++ {
		point := Coordinate{x: x, y: r.C1.y}
		if !perimeter.Contains(point) {
			return false
		}
	}
	for y := r.C1.y + 1; y < r.C2.y; y++ {
		point := Coordinate{x: r.C1.x, y: y}
		if !perimeter.Contains(point) {
			return false
		}
	}
	for y := r.C1.y + 1; y < r.C2.y; y++ {
		point := Coordinate{x: r.C2.x, y: y}
		if !perimeter.Contains(point) {
			return false
		}
	}
	return true
}

func (r Rectangle) IsFullyContainedInPerimeterIndex(index *PerimeterIndex) bool {
	if r.C1.x < index.bounds.C1.x || r.C2.x > index.bounds.C2.x ||
		r.C1.y < index.bounds.C1.y || r.C2.y > index.bounds.C2.y {
		return false
	}

	for x := r.C1.x; x <= r.C2.x; x++ {
		point := Coordinate{x: x, y: r.C2.y}
		if !index.Contains(point) {
			return false
		}
	}
	for x := r.C1.x; x <= r.C2.x; x++ {
		point := Coordinate{x: x, y: r.C1.y}
		if !index.Contains(point) {
			return false
		}
	}
	for y := r.C1.y + 1; y < r.C2.y; y++ {
		point := Coordinate{x: r.C1.x, y: y}
		if !index.Contains(point) {
			return false
		}
	}
	for y := r.C1.y + 1; y < r.C2.y; y++ {
		point := Coordinate{x: r.C2.x, y: y}
		if !index.Contains(point) {
			return false
		}
	}
	return true
}

func ReadInput(input string, sortCoordinates bool) []Coordinate {
	lines := strings.Split(input, "\n")
	coordinates := make([]Coordinate, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		c, err := CoordinateFromString(line)
		if err != nil {
			continue
		}
		coordinates = append(coordinates, c)
	}
	if sortCoordinates {
		sort.SliceStable(coordinates, func(i, j int) bool {
			if coordinates[i].x == coordinates[j].x {
				return coordinates[i].y < coordinates[j].y
			}
			return coordinates[i].x < coordinates[j].x
		})
	}
	return coordinates
}

func ReadInputAndMakePerimeter(input string) (coordinates []Coordinate, perimeter Perimeter) {
	coordinates = ReadInput(input, false)
	perimeter = MakePerimeter(coordinates)
	sort.SliceStable(coordinates, func(i, j int) bool {
		if coordinates[i].x == coordinates[j].x {
			return coordinates[i].y < coordinates[j].y
		}
		return coordinates[i].x < coordinates[j].x
	})
	return coordinates, perimeter
}

func FindLargestRectangle(coordinates []Coordinate) Rectangle {
	if len(coordinates) < 2 {
		return Rectangle{}
	}

	var largestRectangle Rectangle
	largestArea := 0

	for i := 0; i < len(coordinates); i++ {
		for j := len(coordinates) - 1; j > i; j-- {
			rectangle, err := MakeRectangle(coordinates[i], coordinates[j])
			if err != nil {
				continue
			}
			area := rectangle.Area()
			if area > largestArea {
				largestRectangle = rectangle
				largestArea = area
			}
		}
	}
	return largestRectangle
}

func FindLargestRectangleInPerimeter(perimeter Perimeter, coordinates []Coordinate) (Rectangle, error) {
	if len(perimeter) < 4 {
		return Rectangle{}, errors.New("perimeter must have at least 4 coordinates")
	}

	index := NewPerimeterIndex(perimeter)
	index.PrecomputeCoordinates(coordinates)

	numWorkers := runtime.NumCPU()
	if len(coordinates) < 100 {
		numWorkers = 1
	}

	if numWorkers == 1 {
		return findLargestRectangleSequential(index, coordinates)
	}

	return findLargestRectangleParallel(index, coordinates, numWorkers)
}

func findLargestRectangleSequential(index *PerimeterIndex, coordinates []Coordinate) (Rectangle, error) {
	var largestRectangle Rectangle
	largestArea := 0
	found := false

	for i := 0; i < len(coordinates); i++ {
		for j := len(coordinates) - 1; j > i; j-- {
			rectangle, err := MakeRectangle(coordinates[i], coordinates[j])
			if err != nil {
				continue
			}

			area := rectangle.Area()
			if area <= largestArea {
				continue
			}

			if rectangle.IsFullyContainedInPerimeterIndex(index) {
				largestRectangle = rectangle
				largestArea = area
				found = true
			}
		}
	}
	if !found {
		return Rectangle{}, errors.New("no rectangle found in perimeter")
	}
	return largestRectangle, nil
}

func findLargestRectangleParallel(index *PerimeterIndex, coordinates []Coordinate, numWorkers int) (Rectangle, error) {
	type result struct {
		rectangle Rectangle
		area      int
	}
	results := make(chan result, numWorkers*10)

	type workItem struct {
		i, j int
	}
	work := make(chan workItem, numWorkers*10)

	go func() {
		defer close(work)
		for i := 0; i < len(coordinates); i++ {
			for j := len(coordinates) - 1; j > i; j-- {
				work <- workItem{i: i, j: j}
			}
		}
	}()

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range work {
				rectangle, err := MakeRectangle(coordinates[item.i], coordinates[item.j])
				if err != nil {
					continue
				}

				area := rectangle.Area()
				if rectangle.IsFullyContainedInPerimeterIndex(index) {
					results <- result{rectangle: rectangle, area: area}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var largestRectangle Rectangle
	largestArea := 0
	found := false

	for res := range results {
		if res.area > largestArea {
			largestRectangle = res.rectangle
			largestArea = res.area
			found = true
		}
	}

	if !found {
		return Rectangle{}, errors.New("no rectangle found in perimeter")
	}
	return largestRectangle, nil
}

func SolveDay9Part1(input string) interface{} {
	coordinates := ReadInput(input, true)
	largestRectangle := FindLargestRectangle(coordinates)
	return largestRectangle.Area()
}

func SolveDay9Part2(input string) interface{} {
	coordinates, perimeter := ReadInputAndMakePerimeter(input)
	largestRectangle, err := FindLargestRectangleInPerimeter(perimeter, coordinates)
	if err != nil {
		return 0
	}
	return largestRectangle.Area()
}
