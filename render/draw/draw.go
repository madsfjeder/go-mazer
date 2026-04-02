// Package render - draws the maze and solved path
package render

import (
	"flag"
	"fmt"
	"image/color"
	"strconv"
	"time"

	"maze/config"
	"maze/generate"
	"maze/internal/grid"
	"maze/internal/stack"

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

func (r RaylibRenderer) DrawTile(x, y, width, height int32, color color.RGBA) {
	r.DrawRectangle(x, y, width, height, color)

	// Top
	r.DrawRectangle(x, y, width, 1, rl.NewColor(230, 230, 230, 200))

	// Left
	r.DrawRectangle(x, y, 1, height, rl.NewColor(230, 230, 230, 200))

	// Bottom
	r.DrawRectangle(x, y+height-1, width, 1, rl.NewColor(0, 0, 0, 125))

	// Right
	r.DrawRectangle(x+width-1, y, 1, height, rl.NewColor(0, 0, 0, 125))
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

func NewRaylibRenderer(x, y int, cellType grid.CellType, colors grid.Colors, showBacktracking bool) RaylibRenderer {
	return RaylibRenderer{
		config: grid.Config{
			X:                int32(x),
			Y:                int32(y),
			EdgeWidth:        config.EdgeWidth,
			CellType:         cellType,
			ShowBacktracking: showBacktracking,
		},
		colors: colors,
	}
}

type ElementBounds struct {
	width  int32
	height int32
}

type GuiElement interface {
	Render(xPos, yPos float32)
	Bounds() ElementBounds
}

type baseElement struct {
	width  int32
	height int32
}

func (b *baseElement) Bounds() ElementBounds {
	return ElementBounds{
		width:  b.width,
		height: b.height,
	}
}

type Button struct {
	baseElement
	text    *string
	onClick func()
}

func (b *Button) Render(xPos, yPos float32) {
	clicked := gui.Button(
		rl.NewRectangle(
			xPos,
			yPos,
			float32(b.width),
			float32(b.height)),
		*b.text,
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
			xPos+50,
			yPos,
			// Account for the labels also taking up space
			float32(s.width-100),
			float32(s.height),
		),
		s.textLeft,
		s.textRight,
		s.value,
		0,
		300,
	)

	s.value = interval
}

type Toggle struct {
	baseElement
	text   string
	active *bool
}

func (t *Toggle) Render(xPos, yPos float32) {
	active := gui.Toggle(
		rl.NewRectangle(
			xPos,
			yPos,
			float32(t.width),
			float32(t.height),
		),
		t.text,
		*t.active,
	)

	*t.active = active
}

type Dropdown struct {
	baseElement
	text           string
	active         *int32
	previousActive int32
	editMode       *bool
	onChange       func()
}

func (s *Dropdown) Render(xPos, yPos float32) {
	editMode := gui.DropdownBox(
		rl.NewRectangle(
			xPos,
			yPos,
			float32(s.width),
			float32(s.height),
		),
		s.text,
		s.active,
		*s.editMode,
	)

	if editMode {
		*s.editMode = !*s.editMode
	}

	if *s.active != s.previousActive {
		s.previousActive = *s.active
		s.onChange()
	}
}

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

type AnimationTiming struct {
	timeAcc  int64
	prevTime time.Time
	idxCount int
	interval int64
	running  bool
}

func (a *AnimationTiming) Increment() {
	if !a.running {
		return
	}

	if a.interval > 0 && a.idxCount < 100_000 {
		delta := time.Since(a.prevTime)
		a.timeAcc += delta.Milliseconds()

		for a.timeAcc >= a.interval {
			a.timeAcc -= a.interval
			a.idxCount++
		}
	} else {
		a.idxCount = 100_000
	}
}

func (a *AnimationTiming) Reset() {
	a.idxCount = 0
	a.prevTime = time.Now()
	a.timeAcc = 0
}

func formatTime(runTimeMicroseconds int) string {
	duration := time.Duration(runTimeMicroseconds) * time.Microsecond
	milliseconds := float64(duration) / float64(time.Millisecond)
	return fmt.Sprintf("%02.3f", milliseconds)
}

func drawGui(elements []GuiElement, stats Stats) {
	padding := config.Padding
	xPos := padding
	yPos := padding

	rl.DrawRectangle(
		xPos,
		yPos,
		config.MenuBarWidth,
		config.MenuBarHeight-padding,
		rl.White,
	)

	xOffset := xPos + padding
	yOffset := yPos + padding

	for _, element := range elements {
		element.Render(
			float32(xOffset),
			float32(yOffset),
		)
		xOffset += element.Bounds().width + padding
	}

	runTimeText := "Total run time: " + formatTime(stats.runTimeMicroseconds) + "ms"
	totalStepsText := "Total steps: " + strconv.Itoa(stats.totalSteps)
	solutionStepsText := "Solution steps: " + strconv.Itoa(stats.solutionSteps)

	rl.DrawText(totalStepsText, 0, 25, 14, rl.Black)
	rl.DrawText(solutionStepsText, 150, 25, 14, rl.Black)
	rl.DrawText(runTimeText, 300, 25, 14, rl.Black)
}

func drawGeneration(
	animationData GeneratorAnimationData,
	animationConfig AnimationConfig,
	onComplete func(),
) {
	animationData.Draw(animationConfig)
	onComplete()
}

func drawWalls(
	matrixToDraw [][]*grid.Vertex,
	colors grid.Colors,
) {
	for i := range matrixToDraw {
		for j := range matrixToDraw[i] {
			e := matrixToDraw[i][j]
			if e != nil {
				r := NewRaylibRenderer(i, j, grid.Wall, colors, false)
				e.DrawVertex(r)
			}
		}
	}
}

type animationData struct {
	itemsToRender   []stack.StackItem[*grid.Vertex]
	completedMatrix generate.Maze
	matrixToRender  [][]*grid.Vertex
}

type GeneratorAnimationData struct {
	animationData
}

func getGeneratorAnimationData(
	maze generate.Maze,
) GeneratorAnimationData {
	steps := maze.Steps.Copy()

	matrixToDraw := grid.New(int(config.VerticesPerRow), int(config.VerticesPerRow))

	return GeneratorAnimationData{
		animationData: animationData{
			matrixToRender:  matrixToDraw,
			completedMatrix: maze,
			itemsToRender:   steps.PopAllWithIdx(),
		},
	}
}

func (g *GeneratorAnimationData) Draw(animationConfig AnimationConfig) {
	for i := range g.completedMatrix.Matrix {
		for j := range g.completedMatrix.Matrix[i] {
			e := g.completedMatrix.Matrix[i][j]

			if e != nil {
				cellType := grid.Path

				if e.IsSplit {
					cellType = grid.Split
				}

				r := NewRaylibRenderer(i, j, cellType, animationConfig.colors, false)
				e.DrawVertex(r)
			}
		}
	}
}

type CurrentElement struct {
	element *grid.Vertex
	x       int
	y       int
}

type SolverAnimationData struct {
	animationData
	path             stack.Stack[*grid.Vertex]
	fullStackOfItems []stack.StackItem[*grid.Vertex]
	currentElement   CurrentElement
}

func getSolverAnimationData(
	solution stack.Stack[*grid.Vertex],
	maze generate.Maze,
) *SolverAnimationData {
	solverItems := solution.PopAllWithIdx()
	var element *grid.Vertex
	return &SolverAnimationData{
		animationData: animationData{
			itemsToRender:   solverItems,
			matrixToRender:  grid.New(int(config.VerticesPerRow), int(config.VerticesPerCol)),
			completedMatrix: maze,
		},
		path:             solution,
		fullStackOfItems: solverItems,
		currentElement: CurrentElement{
			element: element,
			x:       0,
			y:       0,
		},
	}
}

func (a *SolverAnimationData) Filter(idx int) {
	filteredItems := make([]stack.StackItem[*grid.Vertex], 0, len(a.itemsToRender))
	for _, v := range a.fullStackOfItems {
		shouldAppend := v.Item != nil && v.Index <= idx
		if shouldAppend {
			filteredItems = append(filteredItems, v)
		}
	}

	a.itemsToRender = filteredItems
}

func (a *SolverAnimationData) Refresh() {
	for i := range a.completedMatrix.Matrix {
		for j := range a.completedMatrix.Matrix[i] {
			v := a.completedMatrix.Matrix[i][j]
			for _, toDraw := range a.itemsToRender {
				if v == toDraw.Item {
					a.matrixToRender[i][j] = toDraw.Item
				}
			}

			if a.currentElement.element != nil && v == a.currentElement.element {
				a.currentElement.x = i
				a.currentElement.y = j
			}
		}
	}

	if len(a.itemsToRender) > 0 {
		a.currentElement.element = a.itemsToRender[len(a.itemsToRender)-1].Item
	}
}

func (a *SolverAnimationData) Draw(animationConfig AnimationConfig) {
	for i := range a.matrixToRender {
		for j := range a.matrixToRender[i] {
			e := a.matrixToRender[i][j]

			if e != nil {
				cellType := grid.Solution
				if e.IsBacktracking {
					cellType = grid.Backtracking
				}

				r := NewRaylibRenderer(i, j, cellType, animationConfig.colors, *animationConfig.showBacktracking)
				e.DrawVertex(r)
			}
		}
	}

	if a.currentElement.element != nil {
		r := NewRaylibRenderer(a.currentElement.x, a.currentElement.y, grid.CPU, animationConfig.colors, *animationConfig.showBacktracking)
		a.currentElement.element.DrawVertex(r)
	}
}

type AnimationConfig struct {
	colors           grid.Colors
	showBacktracking *bool
}

func drawSolver(
	animationTiming *AnimationTiming,
	animationData *SolverAnimationData,
	animationConfig AnimationConfig,
) {
	animationTiming.Increment()
	animationData.Filter(animationTiming.idxCount)
	animationData.Refresh()
	animationData.Draw(animationConfig)
}

type State int

const (
	StateGeneration State = iota
	StateSolving
	StateDone
)

func Draw() {
	solverAlgorithm := generate.DFS
	selectedLevel := generate.RandomMaze
	solvingPlaying := true

	maze, err := generate.Generate(selectedLevel)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	solution := maze.Solve(solverAlgorithm)
	runTimeMicroseconds := time.Since(now).Microseconds()

	generatedMaze := maze
	generatedSolution := solution

	debugPtr := flag.Bool("debug", false, "turns debugging on")
	flag.Parse()

	if *debugPtr {
		DEBUG = true
	}

	colors := grid.Colors{
		Start:        rl.DarkPurple,
		End:          rl.Green,
		Wall:         rl.NewColor(30, 30, 30, 255),
		Backtracking: rl.Gray,
		Split:        rl.Magenta,
		EmptyCell:    rl.Brown,
		Cell:         rl.LightGray,
		Solution:     rl.NewColor(240, 125, 151, 255),
		Text:         rl.Red,
		CPU:          rl.Blue,
		DebugWall:    rl.Beige,
	}

	var solveDrawInterval int64 = 15

	animationTiming := AnimationTiming{
		idxCount: 0,
		timeAcc:  0,
		prevTime: time.Now(),
		interval: solveDrawInterval,
		running:  solvingPlaying,
	}

	generatorAnimationData := getGeneratorAnimationData(maze)

	solverAnimationData := getSolverAnimationData(
		solution,
		maze,
	)

	showBacktracking := true

	animationConfig := AnimationConfig{
		colors,
		&showBacktracking,
	}

	state := StateGeneration
	stats := generateStats(generatedSolution, int(runTimeMicroseconds))

	guiElements := make([]GuiElement, 0)

	playBtnText := "Pause"
	playBtn := &Button{
		baseElement: baseElement{
			width:  100,
			height: 20,
		},
		text: &playBtnText,
		onClick: func() {
			if solvingPlaying {
				playBtnText = "Play"
				solvingPlaying = false
			} else {
				playBtnText = "Pause"
				solvingPlaying = true
			}
		},
	}

	reset := func() {
		generatedMaze, err := generate.Generate(selectedLevel)

		now = time.Now()
		generatedSolution = generatedMaze.Solve(solverAlgorithm)
		generatedSolution.Reverse()
		runTimeMicroseconds = time.Since(now).Microseconds()
		stats = generateStats(generatedSolution, int(runTimeMicroseconds))

		generatorAnimationData = getGeneratorAnimationData(generatedMaze)

		solverAnimationData = getSolverAnimationData(
			generatedSolution,
			generatedMaze,
		)

		animationTiming.Reset()

		if err != nil {
			panic(0)
		}

		state = StateGeneration
	}

	changeAlgorithm := func() {
		now = time.Now()
		generatedSolution = generatedMaze.Solve(solverAlgorithm)
		runTimeMicroseconds = int64(time.Since(now).Microseconds())
		generatedSolution.Reverse()

		generatorAnimationData = getGeneratorAnimationData(generatedMaze)

		solverAnimationData = getSolverAnimationData(
			generatedSolution,
			generatedMaze,
		)

		animationTiming.Reset()

		stats = generateStats(generatedSolution, int(runTimeMicroseconds))

		if err != nil {
			panic(0)
		}
	}

	resetBtnText := "Reset"
	resetBtn := &Button{
		baseElement: baseElement{
			width:  100,
			height: 20,
		},
		text:    &resetBtnText,
		onClick: reset,
	}

	algoDropdownOpen := false
	algoDropdown := &Dropdown{
		baseElement: baseElement{
			width:  100,
			height: 20,
		},
		active:         &solverAlgorithm,
		previousActive: solverAlgorithm,
		text:           "DFS;BFS;GFS;AStar",
		editMode:       &algoDropdownOpen,
		onChange:       changeAlgorithm,
	}

	levelSelectDropdownOpen := false
	levelSelectDropdown := &Dropdown{
		baseElement: baseElement{
			width:  150,
			height: 20,
		},
		active:         &selectedLevel,
		previousActive: selectedLevel,
		text:           "Random maze;Empty test",
		editMode:       &levelSelectDropdownOpen,
		onChange:       reset,
	}

	slider := &Slider{
		baseElement: baseElement{
			width:  150,
			height: 20,
		},
		value:     float32(solveDrawInterval),
		textLeft:  "Slower",
		textRight: "Faster",
	}

	toggle := &Toggle{
		baseElement: baseElement{
			width:  100,
			height: 20,
		},
		text:   "Show backtracking",
		active: &showBacktracking,
	}

	guiElements = append(guiElements, playBtn, resetBtn, levelSelectDropdown, algoDropdown, slider, toggle)

	rl.InitWindow(config.Width, config.Height, "Mazen")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.DrawRectangle(0, 0, config.Width, config.Height, rl.Black)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		drawGeneration(
			generatorAnimationData,
			animationConfig,
			func() {
				if state == StateGeneration {
					state = StateSolving
				}
			},
		)

		if state == StateSolving {
			drawSolver(
				&animationTiming,
				solverAnimationData,
				animationConfig,
			)
		}

		drawWalls(
			generatorAnimationData.completedMatrix.Matrix,
			colors,
		)

		drawGui(guiElements, stats)

		animationTiming.prevTime = time.Now()

		rl.EndDrawing()
	}
}
