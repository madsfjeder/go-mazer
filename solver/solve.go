// Package solver - solves the provided maze
package solver

import (
	"fmt"

	"maze/internal/grid"
	"maze/internal/stack"
)

func Solve(maze [][]*grid.Vertex) stack.Stack[*grid.Vertex] {
	history := *stack.New[*grid.Vertex]()

	currentVertex := maze[0][0]
	currentVertex.VisitedBySolver = true
	isBacktracking := false

	return history
	for i := 0; i < 100; i++ {
		if currentVertex.IsEnd {
			history.Push(currentVertex)
			break
		}

		nextVertex := currentVertex.VisitNextVertex()

		if nextVertex == nil {
			v, err := history.Pop()
			if err != nil {
				fmt.Println("No more history!")
				break
			}
			currentVertex = v
			isBacktracking = true
			currentVertex.VisitedBySolver = true
			continue
		}

		if !isBacktracking {
			history.Push(currentVertex)
		}

		currentVertex.VisitedBySolver = true
		currentVertex = nextVertex
		isBacktracking = false
	}

	return history
}
