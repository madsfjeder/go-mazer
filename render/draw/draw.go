// Package render - draws the maze and solved path
package render

import (
	"flag"
	"image/color"
	"strconv"
	"time"

	"maze/config"
	"maze/generate"
	"maze/internal/grid"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var DEBUG = false

type RaylibRenderer struct {
	config grid.Config
	colors grid.Colors
}

func (r RaylibRenderer) DrawRectangle(x, y, width, height int32, color color.RGBA) {
	rl.DrawRectangle(x, y, width, height, color)
}

func (r RaylibRenderer) DrawText(text string, x, y, fontSize int32, color color.RGBA) {
	rl.DrawText(text, x, y, fontSize, color)
}

func (r RaylibRenderer) Colors() grid.Colors {
	return r.colors
}

func (r RaylibRenderer) Config() grid.Config {
	return r.config
}

func NewRaylibRenderer(x, y int, colors grid.Colors) RaylibRenderer {
	return RaylibRenderer{
		config: grid.Config{
			X:      int32(x),
			Y:      int32(y),
			Width:  config.Width,
			Height: config.Height,
		},
		colors: colors,
	}
}

func Draw(maze generate.Maze) {
	debugPtr := flag.Bool("debug", false, "turns debugging on")
	intervalPtr := flag.Int("interval", 5, "set the rendering interval")

	if *debugPtr {
		DEBUG = true
	}

	colors := grid.Colors{
		Wall:      rl.Black,
		Cell:      rl.White,
		Text:      rl.Red,
		DebugWall: rl.Beige,
	}

	rl.InitWindow(config.Width+config.EdgeWidth, config.Height+config.EdgeWidth, "Mazen")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	reversedSteps := maze.Steps.Copy()
	reversedSteps.Reverse()

	var timeAcc int64 = 0
	var interval int64 = int64(*intervalPtr)
	prevTime := time.Now()
	matrixToDraw := make([][]*grid.Vertex, config.VerticesPerRow)

	for i := range matrixToDraw {
		matrixToDraw[i] = make([]*grid.Vertex, config.VerticesPerCol)
	}

	var elementToDraw *grid.Vertex

	rl.DrawRectangle(0, 0, config.Width+config.EdgeWidth, config.Height+config.EdgeWidth, rl.Black)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		delta := time.Now().Sub(prevTime)
		timeAcc += delta.Milliseconds()
		if timeAcc >= interval {
			timeAcc -= interval
			e, _ := reversedSteps.Pop()
			elementToDraw = e
		}

		var x, y int
		for i := range maze.Matrix {
			for j := range maze.Matrix {
				v := maze.Matrix[i][j]
				if v == elementToDraw {
					x = i
					y = j
					matrixToDraw[i][j] = elementToDraw
					break
				}
			}
		}

		for i := range matrixToDraw {
			for j := range matrixToDraw[i] {
				e := matrixToDraw[i][j]
				if e != nil {
					r := NewRaylibRenderer(i, j, colors)
					e.DrawVertex(r)
					if DEBUG {
						backTracked := maze.BacktrackSteps.FindOrder(e)
						if backTracked != -1 {
							e.DrawVertex(r)
							e.DrawText(i, j, strconv.Itoa(backTracked))
						}
					}
				}
			}
		}

		if elementToDraw != nil {
			elementToDraw.DrawVertex(x, y, rl.Blue)
		}

		prevTime = time.Now()
		rl.EndDrawing()
	}
}
