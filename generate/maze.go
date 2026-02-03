// Package generate - generates the maze
package generate

import (
	"errors"

	"maze/config"
	"maze/internal/grid"
	"maze/internal/stack"
)

type Maze struct {
	Matrix         [][]*grid.Vertex
	Steps          stack.Stack[*grid.Vertex]
	BacktrackSteps stack.Stack[*grid.Vertex]
}

func Generate() (Maze, error) {
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

	history := stack.New[*grid.Vertex]()
	allSteps := stack.New[*grid.Vertex]()
	backtracking := stack.New[*grid.Vertex]()
	currentVertex := matrix[0][0]
	mazeIncomplete := false

	for i := 0; i < 10000; i++ {
		currentVertex.Visited = true
		nextVertex, err := currentVertex.GetNextVertex()
		if nextVertex == nil || err != nil {
			currentVertex, err = history.Pop()
			if err != nil || currentVertex == nil {
				break
			}
			allSteps.Push(currentVertex)
			backtracking.Push(currentVertex)
			continue
		}

		history.Push(currentVertex)
		allSteps.Push(currentVertex)
		currentVertex = nextVertex
	}

	if mazeIncomplete {
		return Maze{}, errors.New("cannot generate maze")
	}

	return Maze{
		Matrix:         matrix,
		Steps:          *allSteps,
		BacktrackSteps: *backtracking,
	}, nil
}
