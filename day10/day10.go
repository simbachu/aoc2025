package day10

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type Solvable struct {
	Workspace []int
	Goal      []int
	Buttons   []Button
}

func MakeSolvable(input string) (Solvable, error) {
	solvable := Solvable{}
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		buttons, err := parseButtons(line)
		if err != nil {
			return solvable, fmt.Errorf("parsing buttons: %w", err)
		}
		solvable.Buttons = buttons

		goal, err := parseJoltage(line)
		if err != nil {
			return solvable, fmt.Errorf("parsing goal: %w", err)
		}
		solvable.Goal = goal
		solvable.Workspace = make([]int, len(goal))
	}
	return solvable, nil
}

func (s *Solvable) Solve() ([]int, bool) {
	numIndicators := len(s.Workspace)
	numButtons := len(s.Buttons)

	if numButtons == 0 || numIndicators == 0 {
		return nil, false
	}

	data := make([]float64, numIndicators*numButtons)
	for j, button := range s.Buttons {
		for _, pos := range button {
			if pos < numIndicators {
				data[pos*numButtons+j] = 1
			}
		}
	}
	A := mat.NewDense(numIndicators, numButtons, data)

	goalData := make([]float64, numIndicators)
	for i, g := range s.Goal {
		goalData[i] = float64(g)
	}
	b := mat.NewVecDense(numIndicators, goalData)

	var svd mat.SVD
	if !svd.Factorize(A, mat.SVDFull) {
		return nil, false
	}

	var vMat mat.Dense
	svd.VTo(&vMat)

	var uMat mat.Dense
	svd.UTo(&uMat)

	values := svd.Values(nil)
	tol := 1e-10
	rank := 0
	for _, v := range values {
		if v > tol {
			rank++
		}
	}

	particular := s.computePseudoinverseSolution(&uMat, &vMat, values, b, tol, numButtons, numIndicators)

	intSol := make([]int, numButtons)
	allNonNeg := true
	for i := 0; i < numButtons; i++ {
		intSol[i] = int(math.Round(particular[i]))
		if intSol[i] < 0 {
			allNonNeg = false
		}
	}

	if rank >= numButtons {
		if allNonNeg && s.verifySolution(intSol) {
			return intSol, true
		}
		if numButtons <= 5 {
			return s.searchAroundSolution(intSol, 3)
		}
		for delta := 1; delta <= 2; delta++ {
			for i := 0; i < numButtons; i++ {
				testSol := make([]int, numButtons)
				copy(testSol, intSol)
				testSol[i] += delta
				if testSol[i] >= 0 && s.verifySolution(testSol) {
					return testSol, true
				}
				testSol[i] = intSol[i] - delta
				if testSol[i] >= 0 && s.verifySolution(testSol) {
					return testSol, true
				}
			}
		}
		if numButtons <= 5 {
			return s.bruteForceSearch(30)
		}
		return nil, false
	}

	nullDim := numButtons - rank
	nullVecs := make([][]float64, nullDim)
	for k := 0; k < nullDim; k++ {
		nullVecs[k] = make([]float64, numButtons)
		for i := 0; i < numButtons; i++ {
			nullVecs[k][i] = vMat.At(i, rank+k)
		}
	}

	return s.searchNullSpace(particular, nullVecs, nullDim, numButtons)
}

func (s *Solvable) computePseudoinverseSolution(uMat, vMat *mat.Dense, values []float64, b *mat.VecDense, tol float64, numButtons, numIndicators int) []float64 {
	utb := make([]float64, min(numIndicators, numButtons))
	for i := 0; i < len(utb); i++ {
		sum := 0.0
		for j := 0; j < numIndicators; j++ {
			sum += uMat.At(j, i) * b.AtVec(j)
		}
		utb[i] = sum
	}

	sigmaInvUtb := make([]float64, len(utb))
	for i := 0; i < len(utb); i++ {
		if values[i] > tol {
			sigmaInvUtb[i] = utb[i] / values[i]
		}
	}

	particular := make([]float64, numButtons)
	for i := 0; i < numButtons; i++ {
		sum := 0.0
		for j := 0; j < len(sigmaInvUtb); j++ {
			sum += vMat.At(i, j) * sigmaInvUtb[j]
		}
		particular[i] = sum
	}
	return particular
}

func (s *Solvable) searchNullSpace(particular []float64, nullVecs [][]float64, nullDim, numButtons int) ([]int, bool) {
	bestSum := math.Inf(1)
	var bestSol []int

	switch nullDim {
	case 1:
		bestSol, _ = s.search1DNullSpace(particular, nullVecs[0], numButtons, bestSum)
	case 2:
		bestSol, _ = s.search2DNullSpace(particular, nullVecs, numButtons, bestSum)
	case 3:
		bestSol, _ = s.search3DNullSpaceParallel(particular, nullVecs, numButtons, bestSum)
	default:
		bestSol, _ = s.searchNDNullSpace(particular, nullVecs, nullDim, numButtons, bestSum)
	}

	if bestSol != nil {
		return bestSol, true
	}
	return nil, false
}

func (s *Solvable) search1DNullSpace(particular []float64, nullVec []float64, numButtons int, bestSum float64) ([]int, float64) {
	var bestSol []int

	minMult, maxMult := -500, 500
	for j := 0; j < numButtons; j++ {
		if math.Abs(nullVec[j]) > 0.01 {
			bound := -particular[j] / nullVec[j]
			if nullVec[j] > 0 {
				if int(math.Ceil(bound)) > minMult {
					minMult = int(math.Ceil(bound))
				}
			} else {
				if int(math.Floor(bound)) < maxMult {
					maxMult = int(math.Floor(bound))
				}
			}
		}
	}

	if minMult > -100 {
		minMult = -100
	}
	if maxMult < 100 {
		maxMult = 100
	}

	for mult := minMult - 10; mult <= maxMult+10; mult++ {
		intSol, sum, ok := s.tryNullSpaceSolution(particular, nullVec, float64(mult), numButtons)
		if ok && float64(sum) < bestSum {
			bestSum = float64(sum)
			bestSol = intSol
		}
	}

	return bestSol, bestSum
}

func (s *Solvable) search2DNullSpace(particular []float64, nullVecs [][]float64, numButtons int, bestSum float64) ([]int, float64) {
	var bestSol []int
	searchRange := 100

	for m1 := -searchRange; m1 <= searchRange; m1++ {
		for m2 := -searchRange; m2 <= searchRange; m2++ {
			sol := make([]float64, numButtons)
			for i := 0; i < numButtons; i++ {
				sol[i] = particular[i] + float64(m1)*nullVecs[0][i] + float64(m2)*nullVecs[1][i]
			}

			intSol, sum, ok := s.roundAndVerify(sol, numButtons)
			if ok && float64(sum) < bestSum {
				bestSum = float64(sum)
				bestSol = intSol
			}
		}
	}

	return bestSol, bestSum
}

func (s *Solvable) search3DNullSpaceParallel(particular []float64, nullVecs [][]float64, numButtons int, bestSum float64) ([]int, float64) {
	var bestSol []int
	searchRange := 100

	type result struct {
		sol []int
		sum int
	}

	numWorkers := 8
	chunkSize := (2*searchRange + 1) / numWorkers
	if chunkSize < 1 {
		chunkSize = 1
	}

	results := make(chan result, numWorkers)

	for w := 0; w < numWorkers; w++ {
		startM1 := -searchRange + w*chunkSize
		endM1 := startM1 + chunkSize
		if w == numWorkers-1 {
			endM1 = searchRange + 1
		}

		go func(start, end int) {
			localBestSum := math.MaxInt
			var localBestSol []int

			for m1 := start; m1 < end; m1++ {
				for m2 := -searchRange; m2 <= searchRange; m2++ {
					for m3 := -searchRange; m3 <= searchRange; m3++ {
						sol := make([]float64, numButtons)
						for i := 0; i < numButtons; i++ {
							sol[i] = particular[i] +
								float64(m1)*nullVecs[0][i] +
								float64(m2)*nullVecs[1][i] +
								float64(m3)*nullVecs[2][i]
						}

						intSol, sum, ok := s.roundAndVerify(sol, numButtons)
						if ok && sum < localBestSum {
							localBestSum = sum
							localBestSol = intSol
						}
					}
				}
			}

			if localBestSol != nil {
				results <- result{localBestSol, localBestSum}
			} else {
				results <- result{nil, math.MaxInt}
			}
		}(startM1, endM1)
	}

	for w := 0; w < numWorkers; w++ {
		r := <-results
		if r.sol != nil && float64(r.sum) < bestSum {
			bestSum = float64(r.sum)
			bestSol = r.sol
		}
	}

	return bestSol, bestSum
}

func (s *Solvable) searchNDNullSpace(particular []float64, nullVecs [][]float64, nullDim, numButtons int, bestSum float64) ([]int, float64) {
	var bestSol []int
	searchRange := 30

	var search func(depth int, mult []int)
	search = func(depth int, mult []int) {
		if depth == nullDim {
			sol := make([]float64, numButtons)
			for i := 0; i < numButtons; i++ {
				sol[i] = particular[i]
				for k := 0; k < nullDim; k++ {
					sol[i] += float64(mult[k]) * nullVecs[k][i]
				}
			}

			intSol, sum, ok := s.roundAndVerify(sol, numButtons)
			if ok && float64(sum) < bestSum {
				bestSum = float64(sum)
				bestSol = intSol
			}
			return
		}

		for m := -searchRange; m <= searchRange; m++ {
			mult[depth] = m
			search(depth+1, mult)
		}
	}
	search(0, make([]int, nullDim))

	return bestSol, bestSum
}

func (s *Solvable) tryNullSpaceSolution(particular []float64, nullVec []float64, mult float64, numButtons int) ([]int, int, bool) {
	sol := make([]float64, numButtons)
	for i := 0; i < numButtons; i++ {
		sol[i] = particular[i] + mult*nullVec[i]
	}
	return s.roundAndVerify(sol, numButtons)
}

func (s *Solvable) roundAndVerify(sol []float64, numButtons int) ([]int, int, bool) {
	intSol := make([]int, numButtons)
	sum := 0
	for i := 0; i < numButtons; i++ {
		intSol[i] = int(math.Round(sol[i]))
		if intSol[i] < 0 {
			return nil, 0, false
		}
		sum += intSol[i]
	}
	if s.verifySolution(intSol) {
		return intSol, sum, true
	}
	return nil, 0, false
}

func (s *Solvable) verifySolution(sol []int) bool {
	result := make([]int, len(s.Goal))
	for j, presses := range sol {
		if presses < 0 {
			return false
		}
		for _, pos := range s.Buttons[j] {
			if pos < len(result) {
				result[pos] += presses
			}
		}
	}
	for i, g := range s.Goal {
		if result[i] != g {
			return false
		}
	}
	return true
}

func (s *Solvable) bruteForceSearch(maxPresses int) ([]int, bool) {
	numButtons := len(s.Buttons)
	bestSum := math.MaxInt
	var bestSol []int

	var search func(depth int, current []int, currentSum int, partial []int)
	search = func(depth int, current []int, currentSum int, partial []int) {
		if currentSum >= bestSum {
			return
		}
		if depth == numButtons {
			for i, g := range s.Goal {
				if partial[i] != g {
					return
				}
			}
			if currentSum < bestSum {
				bestSum = currentSum
				bestSol = make([]int, numButtons)
				copy(bestSol, current)
			}
			return
		}

		maxUseful := maxPresses
		for _, pos := range s.Buttons[depth] {
			if pos < len(s.Goal) {
				remaining := s.Goal[pos] - partial[pos]
				if remaining < maxUseful {
					maxUseful = remaining
				}
			}
		}
		if maxUseful < 0 {
			maxUseful = 0
		}

		for p := 0; p <= maxUseful; p++ {
			current[depth] = p
			newPartial := make([]int, len(partial))
			copy(newPartial, partial)
			for _, pos := range s.Buttons[depth] {
				if pos < len(newPartial) {
					newPartial[pos] += p
				}
			}
			exceeded := false
			for i, g := range s.Goal {
				if newPartial[i] > g {
					exceeded = true
					break
				}
			}
			if !exceeded {
				search(depth+1, current, currentSum+p, newPartial)
			}
		}
	}

	search(0, make([]int, numButtons), 0, make([]int, len(s.Goal)))

	if bestSol != nil {
		return bestSol, true
	}
	return nil, false
}

func (s *Solvable) searchAroundSolution(start []int, radius int) ([]int, bool) {
	numButtons := len(start)
	bestSum := math.MaxInt
	var bestSol []int

	var search func(depth int, current []int, currentSum int)
	search = func(depth int, current []int, currentSum int) {
		if currentSum >= bestSum {
			return
		}
		if depth == numButtons {
			if s.verifySolution(current) && currentSum < bestSum {
				bestSum = currentSum
				bestSol = make([]int, numButtons)
				copy(bestSol, current)
			}
			return
		}

		minVal := start[depth] - radius
		if minVal < 0 {
			minVal = 0
		}
		maxVal := start[depth] + radius

		for v := minVal; v <= maxVal; v++ {
			current[depth] = v
			search(depth+1, current, currentSum+v)
		}
	}

	search(0, make([]int, numButtons), 0)

	if bestSol != nil {
		return bestSol, true
	}
	return nil, false
}

type Indicator bool

type Button []int

type Joltage []int

type Machine struct {
	Indicators []Indicator
	Buttons    []Button
	Joltages   Joltage
	Goal       []Indicator
}

func (m *Machine) Toggle(b Button) {
	for _, indicator := range b {
		m.Indicators[indicator] = !m.Indicators[indicator]
	}
}

func (m *Machine) MakeButtonMatrix() [][]bool {
	matrix := make([][]bool, len(m.Buttons))
	for i, button := range m.Buttons {
		matrix[i] = make([]bool, len(m.Indicators))
		for _, indicator := range button {
			matrix[i][indicator] = true
		}
	}
	return matrix
}

func (m *Machine) GetGoalVector() []bool {
	vector := make([]bool, len(m.Indicators))
	for i, indicator := range m.Indicators {
		vector[i] = bool(indicator)
	}
	return vector
}

func (m *Machine) Solve() ([]bool, bool) {
	matrix := m.MakeButtonMatrix()
	goalVector := m.GetGoalVector()
	augmentedMatrix := buildAugmentedMatrix(matrix, goalVector)
	numButtons := len(m.Buttons)
	numIndicators := len(m.Indicators)

	numPivots := forwardElimination(augmentedMatrix, numButtons)
	if !checkConsistency(augmentedMatrix, numPivots, numIndicators, numButtons) {
		return nil, false
	}

	return findMinimalSolution(augmentedMatrix, numPivots, numButtons)
}

func buildAugmentedMatrix(matrix [][]bool, goalVector []bool) [][]bool {
	numButtons := len(matrix)
	if numButtons == 0 {
		return nil
	}
	numIndicators := len(matrix[0])

	augmented := make([][]bool, numIndicators)
	for i := range augmented {
		augmented[i] = make([]bool, numButtons+1)
		for j := 0; j < numButtons; j++ {
			augmented[i][j] = matrix[j][i]
		}
		augmented[i][numButtons] = goalVector[i]
	}
	return augmented
}

func forwardElimination(augmentedMatrix [][]bool, numberOfButtons int) int {
	numRows := len(augmentedMatrix)
	pivotRow := 0

	for col := 0; col < numberOfButtons && pivotRow < numRows; col++ {
		pivot := findPivot(augmentedMatrix, col, pivotRow)
		if pivot == -1 {
			continue
		}
		swapRows(augmentedMatrix, pivotRow, pivot)
		eliminateColumn(augmentedMatrix, col, pivotRow)
		pivotRow++
	}
	return pivotRow
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
		if row != pivotRow && augmentedMatrix[row][col] {
			for c := 0; c < len(augmentedMatrix[row]); c++ {
				augmentedMatrix[row][c] = augmentedMatrix[row][c] != augmentedMatrix[pivotRow][c]
			}
		}
	}
}

func checkConsistency(augmentedMatrix [][]bool, numberOfPivots, _, numberOfButtons int) bool {
	for row := numberOfPivots; row < len(augmentedMatrix); row++ {
		allZero := true
		for col := 0; col < numberOfButtons; col++ {
			if augmentedMatrix[row][col] {
				allZero = false
				break
			}
		}
		if allZero && augmentedMatrix[row][numberOfButtons] {
			return false
		}
	}
	return true
}

func findMinimalSolution(augmentedMatrix [][]bool, numberOfPivots, numberOfButtons int) ([]bool, bool) {
	pivotCols := make(map[int]int)
	for row := 0; row < numberOfPivots; row++ {
		for col := 0; col < numberOfButtons; col++ {
			if augmentedMatrix[row][col] {
				pivotCols[col] = row
				break
			}
		}
	}

	freeVars := []int{}
	for col := 0; col < numberOfButtons; col++ {
		if _, isPivot := pivotCols[col]; !isPivot {
			freeVars = append(freeVars, col)
		}
	}

	numFree := len(freeVars)
	bestCount := numberOfButtons + 1
	var bestSol []bool

	for mask := 0; mask < (1 << numFree); mask++ {
		sol := make([]bool, numberOfButtons)

		for i, col := range freeVars {
			sol[col] = (mask & (1 << i)) != 0
		}

		for col := numberOfButtons - 1; col >= 0; col-- {
			if row, isPivot := pivotCols[col]; isPivot {
				sol[col] = augmentedMatrix[row][numberOfButtons]
				for c := col + 1; c < numberOfButtons; c++ {
					if augmentedMatrix[row][c] {
						sol[col] = sol[col] != sol[c]
					}
				}
			}
		}

		count := 0
		for _, v := range sol {
			if v {
				count++
			}
		}
		if count < bestCount {
			bestCount = count
			bestSol = make([]bool, numberOfButtons)
			copy(bestSol, sol)
		}
	}

	return bestSol, bestSol != nil
}

func MakeButton(input string) (Button, error) {
	parts := strings.Split(input, ",")
	button := make(Button, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid button value: %s", part)
		}
		button[i] = val
	}
	return button, nil
}

func MakeJoltage(input string) (Joltage, error) {
	parts := strings.Split(input, ",")
	joltage := make(Joltage, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid joltage value: %s", part)
		}
		joltage[i] = val
	}
	return joltage, nil
}

func ReadInput(input string) []Machine {
	machines := []Machine{}
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		machine, err := parseMachine(line)
		if err != nil {
			continue
		}
		machines = append(machines, machine)
	}
	return machines
}

func parseMachine(line string) (Machine, error) {
	machine := Machine{}

	indicators, err := parseGoalState(line)
	if err != nil {
		return machine, err
	}
	machine.Indicators = indicators
	machine.Goal = indicators // Goal is same as Indicators for boolean problems

	buttons, err := parseButtons(line)
	if err != nil {
		return machine, err
	}
	machine.Buttons = buttons

	joltage, err := parseJoltage(line)
	if err != nil {
		return machine, err
	}
	machine.Joltages = joltage

	return machine, nil
}

func parseGoalState(line string) ([]Indicator, error) {
	re := regexp.MustCompile(`\[([^\]]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no goal state found in brackets")
	}

	indicators := make([]Indicator, len(matches[1]))
	for i, char := range matches[1] {
		indicators[i] = char == '#'
	}
	return indicators, nil
}

func parseButtons(line string) ([]Button, error) {
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no buttons found in parentheses")
	}

	buttons := make([]Button, len(matches))
	for i, match := range matches {
		button, err := MakeButton(match[1])
		if err != nil {
			return nil, fmt.Errorf("parsing button %d: %w", i, err)
		}
		buttons[i] = button
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
		for _, pressed := range vector {
			if pressed {
				result++
			}
		}
	}
	return result
}

func SolveDay10Part2(input string) int {
	result := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		solvable, err := MakeSolvable(line)
		if err != nil {
			continue
		}
		vector, success := solvable.Solve()
		if !success {
			continue
		}
		for _, presses := range vector {
			result += presses
		}
	}
	return result
}
