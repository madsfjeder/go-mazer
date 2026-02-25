// Package render - draws the maze and solved path
package render

import (
	"flag"
	"image/color"
	"time"

	"maze/config"
	"maze/generate"
	"maze/internal/grid"
	"maze/internal/stack"
	"maze/solver"

	gui "github.com/gen2brain/raylib-go/raygui"
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

func (r RaylibRenderer) DrawRectangleRounded(x, y, width, height int32, roundness float32, color color.RGBA) {
	rl.DrawRectangleRounded(rl.NewRectangle(float32(x), float32(y), float32(width), float32(height)), roundness, 4, color)
}

func (r RaylibRenderer) DrawText(text string, x, y, fontSize int32, color color.RGBA) {
	rl.DrawText(text, x, y, fontSize, color)
}

func (r RaylibRenderer) DrawCircle(x, y int32, radius float32, color color.RGBA) {
	rl.DrawCircle(x, y, radius, color)
}

func (r RaylibRenderer) Colors() grid.Colors {
	return r.colors
}

func (r RaylibRenderer) Config() grid.Config {
	return r.config
}

func NewRaylibRenderer(x, y int, cellType grid.CellType, colors grid.Colors) RaylibRenderer {
	return RaylibRenderer{
		config: grid.Config{
			X:         int32(x),
			Y:         int32(y),
			EdgeWidth: config.EdgeWidth,
			CellType:  cellType,
		},
		colors: colors,
	}
}

type GuiElement interface {
	Render(xPos, yPos float32)
}

type baseElement struct {
	width  int32
	height int32
}

type Button struct {
	baseElement
	text    string
	onClick func()
}

func (b *Button) Render(xPos, yPos float32) {
	clicked := gui.Button(
		rl.NewRectangle(
			xPos,
			yPos,
			float32(b.baseElement.width),
			float32(b.baseElement.height)),
		b.text,
	)

	if clicked {
		b.onClick()
	}
}

type Slider struct {
	baseElement
	textLeft  string
	textRight string
	value     float32
}

func (s *Slider) Render(xPos, yPos float32) {
	interval := gui.Slider(
		rl.NewRectangle(
			xPos+30,
			yPos,
			float32(s.baseElement.width),
			float32(s.baseElement.height),
		),
		s.textLeft,
		s.textRight,
		s.value,
		0,
		2000,
	)

	s.value = interval
}

func drawGui(elements []GuiElement) {
	xPos := config.EdgeWidth / 2
	yPos := config.EdgeWidth / 2
	yOffset := 50

	rl.DrawRectangle(
		xPos,
		yPos,
		config.MenuBarWidth-(config.EdgeWidth/8),
		config.MenuBarHeight-(config.EdgeWidth/2),
		rl.White,
	)

	for i, element := range elements {
		y := int32(i)*int32(yOffset) + yPos + 10
		element.Render(
			float32(xPos+20),
			float32(y),
		)

	}
}

func setup(
	maze generate.Maze,
) ([][]*grid.Vertex, stack.Stack[*grid.Vertex], *grid.Vertex) {
	reversedSteps := maze.Steps.Copy()
	reversedSteps.Reverse()

	matrixToDraw := make([][]*grid.Vertex, config.VerticesPerRow)

	for i := range matrixToDraw {
		matrixToDraw[i] = make([]*grid.Vertex, config.VerticesPerCol)
	}

	var elementToDraw *grid.Vertex

	return matrixToDraw, reversedSteps, elementToDraw
}

func drawGeneration(
	completedMaze generate.Maze,
	matrixToDraw [][]*grid.Vertex,
	steps *stack.Stack[*grid.Vertex],
	currentSolverVertex *grid.Vertex,
	colors grid.Colors,
	interval int64,
	timeAcc *int64,
	prevTime time.Time,
	onComplete func(),
) {
	drawEveryLoop := make([]*grid.Vertex, 0)

	if interval > 0 {
		delta := time.Since(prevTime)
		*timeAcc += delta.Milliseconds()

		for *timeAcc >= interval {
			*timeAcc -= interval
			e, _ := steps.Pop()
			drawEveryLoop = append(drawEveryLoop, e)
		}
	} else {
		drawEveryLoop = steps.PopAll()
	}

	if steps.Length() == 0 {
		onComplete()
		return
	}

	if len(drawEveryLoop) > 0 {
		currentSolverVertex = drawEveryLoop[len(drawEveryLoop)-1]
	}

	var x, y int
	for i := range completedMaze.Matrix {
		for j := range completedMaze.Matrix[i] {
			v := completedMaze.Matrix[i][j]
			for _, toDraw := range drawEveryLoop {
				if v == toDraw {
					x = i
					y = j
					matrixToDraw[i][j] = toDraw
				}
			}
		}
	}

	for i := range matrixToDraw {
		for j := range matrixToDraw[i] {
			e := matrixToDraw[i][j]

			if e != nil {
				if e == currentSolverVertex {
					continue
				}

				r := NewRaylibRenderer(i, j, grid.Path, colors)
				e.DrawVertex(r)
			} else {
				// Draw empty cell
				r := NewRaylibRenderer(i, j, grid.EmptyCell, colors)
				empty := grid.Vertex{
					IsPath:          false,
					VisitedBySolver: false,
				}
				empty.DrawVertex(r)
			}
		}
	}

	for i := range matrixToDraw {
		for j := range matrixToDraw[i] {
			e := matrixToDraw[i][j]
			if e != nil {
				if e == currentSolverVertex {
					continue
				}
				r := NewRaylibRenderer(i, j, grid.Wall, colors)
				e.DrawVertex(r)
			}
		}
	}

	if currentSolverVertex != nil {
		r := NewRaylibRenderer(x, y, grid.CPU, colors)
		if false {
			currentSolverVertex.DrawVertex(r)
		}
	}
}

func drawSolver(
	completedMaze generate.Maze,
	matrixToDraw [][]*grid.Vertex,
	solution *stack.Stack[*grid.Vertex],
	colors grid.Colors,
	interval int64,
	prevTime time.Time,
	timeAcc *int64,
) {
	drawEveryLoop := make([]*grid.Vertex, 0)

	if interval > 0 {
		delta := time.Since(prevTime)
		*timeAcc += delta.Milliseconds()

		for *timeAcc >= interval {
			*timeAcc -= interval
			e, _ := solution.Pop()
			drawEveryLoop = append(drawEveryLoop, e)
		}
	} else {
		drawEveryLoop = solution.PopAll()
	}

	for i := range completedMaze.Matrix {
		for j := range completedMaze.Matrix[i] {
			v := completedMaze.Matrix[i][j]
			for _, toDraw := range drawEveryLoop {
				if v == toDraw {
					matrixToDraw[i][j] = toDraw
				}
			}
		}
	}

	for i := range matrixToDraw {
		for j := range matrixToDraw[i] {
			e := matrixToDraw[i][j]

			if e != nil {
				r := NewRaylibRenderer(i, j, grid.Solution, colors)
				e.DrawVertex(r)
			}
		}
	}
}

type State int

const (
	StateGeneration State = iota
	StateSolving
	StateDone
)

func Draw(maze generate.Maze, solution stack.Stack[*grid.Vertex]) {
	generatedMaze := maze
	generatedSolution := solution
	generatedSolution.Reverse()
	debugPtr := flag.Bool("debug", false, "turns debugging on")
	intervalPtr := flag.Int("interval", 1, "set the rendering interval")
	flag.Parse()

	if *debugPtr {
		DEBUG = true
	}

	colors := grid.Colors{
		Start:     rl.DarkPurple,
		End:       rl.Green,
		Wall:      rl.Black,
		EmptyCell: rl.Brown,
		Cell:      rl.White,
		Solution:  rl.Red,
		Text:      rl.Red,
		CPU:       rl.Blue,
		DebugWall: rl.Beige,
	}

	interval := int64(*intervalPtr)

	matrixToDraw, steps, currentSolverVertex := setup(maze)
	solutionToDraw := make([][]*grid.Vertex, config.VerticesPerRow)
	for i := range solutionToDraw {
		solutionToDraw[i] = make([]*grid.Vertex, config.VerticesPerCol)
	}

	state := StateGeneration
	var generationTimeAcc int64 = 0

	guiElements := make([]GuiElement, 0)

	btn := &Button{
		baseElement: baseElement{
			width:  100,
			height: 40,
		},
		text: "Reset",
		onClick: func() {
			newMaze, err := generate.Generate()
			generatedMaze = newMaze
			matrixToDraw, steps, currentSolverVertex = setup(newMaze)
			generatedSolution = solver.Solve(generatedMaze.Matrix)
			generatedSolution.Reverse()
			solutionToDraw = make([][]*grid.Vertex, config.VerticesPerRow)
			for i := range solutionToDraw {
				solutionToDraw[i] = make([]*grid.Vertex, config.VerticesPerCol)
			}

			if err != nil {
				panic(0)
			}
			generationTimeAcc = 0
			state = StateGeneration
		},
	}

	slider := &Slider{
		baseElement: baseElement{
			width:  100,
			height: 20,
		},
		value:     float32(interval),
		textLeft:  "Slower",
		textRight: "Faster",
	}

	guiElements = append(guiElements, btn, slider)

	prevTime := time.Now()
	rl.InitWindow(config.Width+config.EdgeWidth, config.Height+config.EdgeWidth, "Mazen")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.DrawRectangle(0, 0, config.Width+config.EdgeWidth, config.Height+config.EdgeWidth, rl.Black)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		drawGui(guiElements)

		switch state {
		case StateGeneration:
			{
				drawGeneration(
					generatedMaze,
					matrixToDraw,
					&steps,
					currentSolverVertex,
					colors,
					interval,
					&generationTimeAcc,
					prevTime,
					func() {
						state = StateSolving
						generationTimeAcc = 0
					},
				)
			}

		case StateSolving:
			{
				drawSolver(
					generatedMaze,
					solutionToDraw,
					&generatedSolution,
					colors,
					interval,
					prevTime,
					&generationTimeAcc,
				)
			}

		case StateDone:
			{
			}
		}

		prevTime = time.Now()
		rl.EndDrawing()
	}
}
