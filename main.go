// package Generates a maze and then solves it
package main

import (
	"maze/generate"
	render "maze/render/draw"
)

func main() {
	maze, err := generate.Generate()
	if err != nil {
		panic(err)
	}

	render.Draw(maze)
}
