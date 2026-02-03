// package Generates a maze and then solves it
package main

import (
	"flag"

	"maze/generate"
	render "maze/render/draw"
)

func main() {
	flag.Parse()

	maze, err := generate.Generate()
	if err != nil {
		panic(err)
	}

	render.Draw(maze)
}
