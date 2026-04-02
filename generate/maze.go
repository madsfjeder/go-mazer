// Package generate - generates the maze
package generate

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"maze/config"
	"maze/internal/grid"
	"maze/internal/pqueue"
	"maze/internal/queue"
	"maze/internal/stack"
	"maze/internal/utils"
)

type Maze struct {
	Matrix         [][]*grid.Vertex
	Steps          stack.Stack[*grid.Vertex]
	BacktrackSteps stack.Stack[*grid.Vertex]
}

func (m *Maze) setupEmpty() {
	matrix := make([][]*grid.Vertex, config.VerticesPerRow)

	for i := range matrix {
		matrix[i] = make([]*grid.Vertex, config.VerticesPerCol)
	}

	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = &grid.Vertex{}
		}
	}

	for i := range matrix {
		for j := range matrix[i] {
			currentVertex := matrix[i][j]

			var topVertex *grid.Vertex
			var rightVertex *grid.Vertex
			var botVertex *grid.Vertex
			var leftVertex *grid.Vertex

			if j > 0 {
				topVertex = matrix[i][j-1]
				topVertexBottomConnection := topVertex.GetConnectedVertex(topVertex.BottomEdge)

				if topVertexBottomConnection == nil && currentVertex.TopEdge == nil {
					if currentVertex.TopEdge == nil {
						edge := grid.Edge{
							IsWall:  true,
							Vertex1: currentVertex,
							Vertex2: topVertex,
						}
						currentVertex.TopEdge = &edge
						topVertex.BottomEdge = &edge
					}
				}
			}

			if i < int(config.VerticesPerRow-1) {
				rightVertex = matrix[i+1][j]
				rightVertexLeftConnection := rightVertex.GetConnectedVertex(rightVertex.LeftEdge)

				if rightVertexLeftConnection == nil && currentVertex.RightEdge == nil {
					edge := grid.Edge{
						IsWall:  true,
						Vertex1: currentVertex,
						Vertex2: rightVertex,
					}

					currentVertex.RightEdge = &edge
					rightVertex.LeftEdge = &edge
				}
			}

			if j < int(config.VerticesPerCol-1) {
				botVertex = matrix[i][j+1]
				botVertexTopConnection := botVertex.GetConnectedVertex(botVertex.TopEdge)

				if botVertexTopConnection == nil && currentVertex.BottomEdge == nil {
					edge := grid.Edge{
						IsWall:  true,
						Vertex1: currentVertex,
						Vertex2: botVertex,
					}
					currentVertex.BottomEdge = &edge
					botVertex.TopEdge = &edge
				}
			}

			if i > 0 {
				leftVertex = matrix[i-1][j]
				leftVertexRightConnection := leftVertex.GetConnectedVertex(leftVertex.RightEdge)

				if leftVertexRightConnection == nil && currentVertex.LeftEdge == nil {
					edge := grid.Edge{
						IsWall:  true,
						Vertex1: currentVertex,
						Vertex2: leftVertex,
					}

					currentVertex.LeftEdge = &edge
					leftVertex.RightEdge = &edge
				}
			}
		}
	}

	m.Matrix = matrix
}

func (m *Maze) generateEmptyTest() {
	steps := stack.New[*grid.Vertex]()
	currentVertex := m.Matrix[10][5]
	currentVertex.IsStart = true

	idx := 0
	for i := range m.Matrix {
		for j := range m.Matrix[i] {
			currentVertex := m.Matrix[i][j]
			steps.Push(currentVertex, idx)
			currentVertex.IsPath = true

			if currentVertex.LeftEdge != nil {
				currentVertex.LeftEdge.IsWall = false
			}

			if currentVertex.BottomEdge != nil {
				currentVertex.BottomEdge.IsWall = false
			}

			if currentVertex.RightEdge != nil {
				currentVertex.RightEdge.IsWall = false
			}

			if currentVertex.TopEdge != nil {
				currentVertex.TopEdge.IsWall = false
			}
			idx += 1
		}
	}

	lastVertex := m.Matrix[config.VerticesPerRow-10][config.VerticesPerCol-5]
	lastVertex.IsEnd = true
	m.Steps = *steps
}

func (m *Maze) generateRandomDFS() {
	history := stack.New[*grid.Vertex]()
	steps := stack.New[*grid.Vertex]()
	backtracking := stack.New[*grid.Vertex]()

	currentVertex := m.Matrix[0][0]
	currentVertex.IsStart = true
	currentVertex.IsPath = true

	var currentSplitVertex *grid.Vertex

	cartesianDistanceFromStart := 0
	endGoalPlaced := false

	mazeIncomplete := false
	isBacktracking := false
	for i := 0; i < 10000; i++ {
		shouldSplit := false

		if currentSplitVertex == nil {
			shouldSplit = rollArbitraryDice(10)

			if shouldSplit {
				currentSplitVertex = currentVertex
				currentSplitVertex.IsSplit = true
			}
		}

		if currentSplitVertex != nil {
			shouldRevertToSplit := rollArbitraryDice(30)
			if shouldRevertToSplit {
				var zero *grid.Vertex
				currentVertex = currentSplitVertex
				currentSplitVertex = zero
			}
		}

		nextVertex, err := currentVertex.GetNextVertex()

		if nextVertex == nil || err != nil {
			currentVertex, err = history.Pop()
			if err != nil || currentVertex == nil {
				break
			}
			steps.Push(currentVertex, i)
			backtracking.Push(currentVertex, i)
			isBacktracking = true
			continue
		}

		cartesianDistanceFromStart = getDistanceFromStart(currentVertex, m.Matrix)

		if !endGoalPlaced && int32(cartesianDistanceFromStart) > (config.VerticesPerCol+config.VerticesPerRow-5) {
			currentVertex.IsEnd = true
			endGoalPlaced = true
		}

		currentVertex.IsPath = true
		if !isBacktracking {
			history.Push(currentVertex, i)
		}
		steps.Push(currentVertex, i)
		currentVertex = nextVertex
		isBacktracking = false
	}

	if mazeIncomplete {
		panic("cannot generate maze")
	}

	m.Steps = *steps
	m.BacktrackSteps = *backtracking
}

func (m *Maze) generate(selectedLevel SelectedLevel) {
	switch selectedLevel {
	case EmptyTest:
		{
			m.generateEmptyTest()
		}

	case RandomMaze:
		{
			m.generateRandomDFS()
		}
	}
}

func (m *Maze) resetSolution() {
	for i := range m.Matrix {
		for j := range m.Matrix[i] {
			vertex := m.Matrix[i][j]
			vertex.Closed = false
			vertex.VisitedBySolver = false
			vertex.IsBacktracking = false
			vertex.IsPartOfSolution = false
		}
	}
}

func (m *Maze) solveDFS() stack.Stack[*grid.Vertex] {
	steps := *stack.New[*grid.Vertex]()

	var currentVertex *grid.Vertex

	for i := range m.Matrix {
		for _, v := range m.Matrix[i] {
			if v.IsStart {
				currentVertex = v
				break
			}
		}
	}

	currentVertex.VisitedBySolver = true
	previousVertex := currentVertex
	previousVertex.IsPartOfSolution = true
	isBacktracking := false
	var backtrackingRootVertex *grid.Vertex

	history := *stack.New[*grid.Vertex]()

	for i := 0; i < 10000; i++ {
		steps.Push(currentVertex, i)
		if currentVertex.IsEnd {
			history.Push(currentVertex, i)
			break
		}

		nextVertex := currentVertex.VisitNextVertex()

		if nextVertex == nil {
			v, err := history.Pop()
			if err != nil {
				fmt.Println("No more history!")
				break
			}
			previousVertex := currentVertex
			previousVertex.IsPartOfSolution = false
			previousVertex.IsBacktracking = true
			currentVertex = v
			backtrackingRootVertex = v
			isBacktracking = true
			continue
		}

		if isBacktracking && backtrackingRootVertex != nil {
			var zero *grid.Vertex
			history.Push(backtrackingRootVertex, i)
			backtrackingRootVertex = zero
		}

		if !isBacktracking {
			history.Push(currentVertex, i)
		}

		currentVertex = nextVertex
		isBacktracking = false
		currentVertex.VisitedBySolver = true
		currentVertex.IsPartOfSolution = true
	}

	return steps
}

func (m *Maze) solveBFS() stack.Stack[*grid.Vertex] {
	var startVertex *grid.Vertex

	for i := range m.Matrix {
		for _, v := range m.Matrix[i] {
			if v.IsStart {
				startVertex = v
				break
			}
		}
	}

	startVertex.IsPartOfSolution = true

	history := *stack.New[*grid.Vertex]()
	q := *queue.New[*grid.Vertex]()
	q.Push(startVertex, 0)

	parentMap := make(map[*grid.Vertex]*grid.Vertex)

	var currentVertex *grid.Vertex

	count := 0
	for !q.IsEmpty() {
		count++
		nextVertex, idx := q.Pop()
		if nextVertex != nil && nextVertex.VisitedBySolver {
			continue
		}

		currentVertex = nextVertex
		currentVertex.VisitedBySolver = true
		history.Push(currentVertex, idx)

		if currentVertex.IsEnd {
			break
		}

		currentVertex.IsBacktracking = true
		neighbours := currentVertex.GetNeighbours(false)

		for _, v := range neighbours {
			if v.VisitedBySolver {
				continue
			}

			parentMap[v] = currentVertex
			q.Push(v, idx+1)
		}
	}

	previousVertex := parentMap[currentVertex]
	currentVertex.IsPartOfSolution = true

	stepBackCount := 0
	for previousVertex != nil {
		previousVertex.IsBacktracking = false
		previousVertex.IsPartOfSolution = true
		previousVertex = parentMap[previousVertex]
		stepBackCount++
	}

	return history
}

// Does not guarantee a solution - by design
func (m *Maze) solveGFS() stack.Stack[*grid.Vertex] {
	goalPositionX := 0
	goalPositionY := 0

	for x := range m.Matrix {
		for y, v := range m.Matrix[x] {
			if v.IsEnd {
				goalPositionX = x
				goalPositionY = y
				break
			}
		}
	}

	steps := *stack.New[*grid.Vertex]()

	var startVertex *grid.Vertex

	for i := range m.Matrix {
		for _, v := range m.Matrix[i] {
			if v.IsStart {
				startVertex = v
				break
			}
		}
	}

	startVertex.IsPartOfSolution = true
	currentVertex := startVertex

	evaluationFn := func(vertex *grid.Vertex) float64 {
		currentX := 0
		currentY := 0

		for x := range m.Matrix {
			for y, v := range m.Matrix[x] {
				if v == vertex {
					currentX = x
					currentY = y
					break
				}
			}
		}

		return utils.GetCartesianDistance(int32(currentX), int32(currentY), int32(goalPositionX), int32(goalPositionY))
	}

	count := 0
	for currentVertex != nil {
		steps.Push(currentVertex, count)
		if currentVertex.IsEnd {
			break
		}

		neighbours := currentVertex.GetNeighbours(false)

		if len(neighbours) == 0 {
			fmt.Println("no more ways to go!")
			break
		}

		sort.Slice(neighbours, func(i, j int) bool {
			return evaluationFn(neighbours[i]) < evaluationFn(neighbours[j])
		})

		currentVertex = neighbours[0]
		currentVertex.IsPartOfSolution = true
		currentVertex.VisitedBySolver = true
		count++
	}

	return steps
}

func (m *Maze) solveAStar() stack.Stack[*grid.Vertex] {
	steps := *stack.New[*grid.Vertex]()
	open := pqueue.New[*grid.Vertex]()
	gMap := make(map[*grid.Vertex]float32)
	parentMap := make(map[*grid.Vertex]*grid.Vertex)
	goalPositionX := 0
	goalPositionY := 0

	for x := range m.Matrix {
		for y, v := range m.Matrix[x] {
			if v.IsEnd {
				goalPositionX = x
				goalPositionY = y
				break
			}
		}
	}

	evaluationFn := func(vertex *grid.Vertex) float64 {
		currentX := 0
		currentY := 0

		for x := range m.Matrix {
			for y, v := range m.Matrix[x] {
				if v == vertex {
					currentX = x
					currentY = y
					break
				}
			}
		}

		return math.Abs(float64(goalPositionX-currentX)) + math.Abs(float64(goalPositionY-currentY))
	}

	var startVertex *grid.Vertex

	for i := range m.Matrix {
		for _, v := range m.Matrix[i] {
			if v.IsStart {
				startVertex = v
				break
			}
		}
	}

	startVertex.IsPartOfSolution = true
	startVertex.VisitedBySolver = true
	startVertex.G = 0
	startVertex.H = float32(evaluationFn(startVertex))
	startVertex.F = startVertex.G + startVertex.H
	currentVertex := startVertex
	open.Insert(startVertex, startVertex.F)

	count := 0
	for open.Length() > 0 {
		currentVertex = open.Pop()

		if currentVertex.Closed {
			continue
		}

		steps.Push(currentVertex, count)

		neighbours := currentVertex.GetNeighbours(false)

		for _, neighbour := range neighbours {
			if neighbour.Closed {
				continue
			}

			tentativeG := currentVertex.G + 1
			currentG, ok := gMap[neighbour]

			if ok && currentG <= tentativeG {
				continue
			}

			parentMap[neighbour] = currentVertex
			if neighbour.IsEnd {
				tentativeG = 0
			}

			gMap[neighbour] = tentativeG
			neighbour.G = tentativeG
			neighbour.H = float32(evaluationFn(neighbour))
			neighbour.F = neighbour.G + neighbour.H

			open.Insert(neighbour, neighbour.F)
			parentMap[neighbour] = currentVertex
		}

		currentVertex.Closed = true
		currentVertex.VisitedBySolver = true
		currentVertex.IsBacktracking = true
		count++

		if currentVertex.IsEnd {
			currentVertex.IsBacktracking = false
			break
		}
	}

	fmt.Println("Total count:", count)

	previousVertex := parentMap[currentVertex]
	currentVertex.IsPartOfSolution = true

	stepBackCount := 0
	for previousVertex != nil {
		previousVertex.IsBacktracking = false
		previousVertex.IsPartOfSolution = true
		previousVertex = parentMap[previousVertex]
		stepBackCount++
	}

	return steps
}

type SolverAlgorithm = int32

const (
	DFS SolverAlgorithm = iota
	BFS
	GFS
	AStar
)

type SelectedLevel = int32

const (
	RandomMaze SelectedLevel = iota
	EmptyTest
)

func (m *Maze) Solve(algo SolverAlgorithm) stack.Stack[*grid.Vertex] {
	m.resetSolution()
	switch algo {
	default:
	case DFS:
		{
			return m.solveDFS()
		}
	case BFS:
		{
			return m.solveBFS()
		}
	case GFS:
		{
			return m.solveGFS()
		}
	case AStar:
		{
			return m.solveAStar()
		}
	}

	emptySteps := stack.New[*grid.Vertex]()
	return *emptySteps
}

func getMaze(selectedLevel SelectedLevel) *Maze {
	m := &Maze{}

	m.setupEmpty()
	m.generate(selectedLevel)
	return m
}

// For making random decisions
// Ex. if you want a 10% chance of returning true - rollArbitraryDice(n = 10)
func rollArbitraryDice(n int) bool {
	num := rand.Intn(n)

	return num == n-1
}

func getDistanceFromStart(currentVertex *grid.Vertex, matrix [][]*grid.Vertex) int {
	for i := range matrix {
		for j := range matrix[i] {
			elem := matrix[i][j]
			if elem == currentVertex {
				return i + j
			}
		}
	}

	return 0
}

func Generate(selectedLevel SelectedLevel) (Maze, error) {
	matrix := getMaze(selectedLevel)
	return *matrix, nil
}
