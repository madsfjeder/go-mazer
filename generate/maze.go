// Package generate - generates the maze
package generate

import (
	"fmt"
	"math/rand"

	"maze/config"
	"maze/internal/grid"
	"maze/internal/stack"
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

func (m *Maze) generate() {
	history := stack.New[*grid.Vertex]()
	allSteps := stack.New[*grid.Vertex]()
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
			allSteps.Push(currentVertex, i)
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
		allSteps.Push(currentVertex, i)
		currentVertex = nextVertex
		isBacktracking = false
	}

	if mazeIncomplete {
		panic("cannot generate maze")
	}

	m.Steps = *allSteps
	m.BacktrackSteps = *backtracking
}

type SolverAlgorithm int

const (
	DFS SolverAlgorithm = iota
	BFS
	GFS
	AStar
)

func (m Maze) Solve(algo SolverAlgorithm) stack.Stack[*grid.Vertex] {
	history := *stack.New[*grid.Vertex]()
	allSteps := *stack.New[*grid.Vertex]()

	currentVertex := m.Matrix[0][0]
	currentVertex.VisitedBySolver = true
	previousVertex := currentVertex
	previousVertex.IsPartOfSolution = true
	isBacktracking := false
	var backtrackingRootVertex *grid.Vertex

	for i := 0; i < 10000; i++ {
		allSteps.Push(currentVertex, i)
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
			previousVertex = currentVertex
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

	return allSteps
}

func getMaze() *Maze {
	m := &Maze{}

	m.setupEmpty()
	m.generate()
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

func Generate() (Maze, error) {
	matrix := getMaze()
	return *matrix, nil
}
