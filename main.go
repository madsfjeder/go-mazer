// package Generates a maze and then solves it
package main

import (
	"maze/generate"
	render "maze/render/draw"
	"maze/solver"
)

func main() {
	maze, err := generate.Generate()
	if err != nil {
		panic(err)
	}

	solution := solver.Solve(maze.Matrix)

	render.Draw(maze, solution)
}
