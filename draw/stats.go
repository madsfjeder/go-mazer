package draw

import (
	"maze/internal/grid"
	"maze/internal/stack"
)

type Stats struct {
	runTimeMicroseconds int
	totalSteps          int
	solutionSteps       int
}

func generateStats(generatedSolution stack.Stack[*grid.Vertex], runTimeMicroseconds int) Stats {
	totalSteps := generatedSolution.Length()
	solutionSteps := 0

	for _, val := range generatedSolution.Items() {
		if val.IsPartOfSolution {
			solutionSteps++
		}
	}

	return Stats{
		runTimeMicroseconds,
		totalSteps,
		solutionSteps,
	}
}
