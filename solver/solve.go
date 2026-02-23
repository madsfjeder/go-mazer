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
	var backtrackingRootVertex *grid.Vertex

	for i := 0; i < 10000; i++ {
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
			backtrackingRootVertex = v
			isBacktracking = true
			continue
		}

		if isBacktracking && backtrackingRootVertex != nil {
			var zero *grid.Vertex
			history.Push(backtrackingRootVertex)
			backtrackingRootVertex = zero
		}

		if !isBacktracking {
			history.Push(currentVertex)
		}

		currentVertex = nextVertex
		isBacktracking = false
		currentVertex.VisitedBySolver = true
	}

	fmt.Println("solution len", history.Length())

	return history
}
