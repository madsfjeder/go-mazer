// Package draw - draws the maze and solved path
package draw

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

type State int

const (
	StateGeneration State = iota
	StateSolving
	StateDone
)

var (
	backgroundColor color.RGBA = rl.NewColor(98, 123, 183, 100)
	borderColor     color.RGBA = rl.NewColor(47, 62, 99, 100)
)

func drawGui(elements []GuiElement, stats Stats) {
	buttonsContainerHeight := int32(25)
	padding := config.Padding
	xPos := padding
	yPos := padding

	rl.DrawRectangle(
		xPos,
		yPos,
		config.MenuBarWidth,
		config.MenuBarHeight-padding,
		backgroundColor,
	)

	rl.DrawRectangleLinesEx(
		rl.NewRectangle(
			0,
			0,
			float32(config.Width),
			float32(buttonsContainerHeight+(padding*2)),
		),
		3,
		borderColor,
	)

	rl.DrawRectangleLinesEx(
		rl.NewRectangle(
			0,
			float32(buttonsContainerHeight),
			float32(config.Width),
			35,
		),
		3,
		borderColor,
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

	statsOffset := buttonsContainerHeight + padding + 5

	rl.DrawText(totalStepsText, 0, statsOffset, 14, rl.Black)
	rl.DrawText(solutionStepsText, 150, statsOffset, 14, rl.Black)
	rl.DrawText(runTimeText, 300, statsOffset, 14, rl.Black)
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
	animationData GeneratorAnimationData,
	animationConfig AnimationConfig,
) {
	animationData.DrawWalls(animationConfig)
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

func Draw() {
	solverAlgorithm := generate.DFS
	selectedLevel := generate.RandomMaze
	solvingPlaying := true
	showBacktracking := true
	var solveDrawInterval int64 = 15

	maze, err := generate.Generate(selectedLevel)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	solution := maze.Solve(solverAlgorithm)
	runTimeMicroseconds := time.Since(now).Microseconds()

	debugPtr := flag.Bool("debug", false, "turns debugging on")
	flag.Parse()

	if *debugPtr {
		DEBUG = true
	}

	colors := grid.Colors{
		Start:        rl.DarkPurple,
		End:          rl.Green,
		Wall:         rl.NewColor(30, 40, 64, 255),
		Backtracking: rl.Gray,
		Split:        rl.Magenta,
		EmptyCell:    rl.Brown,
		Cell:         rl.LightGray,
		Solution:     rl.NewColor(240, 125, 151, 255),
		Text:         rl.Red,
		CPU:          rl.Blue,
		DebugWall:    rl.Beige,
	}

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

	animationConfig := AnimationConfig{
		colors,
		&showBacktracking,
	}

	state := StateGeneration
	stats := generateStats(solution, int(runTimeMicroseconds))

	reset := func(shouldRegenerateMaze bool) {
		if shouldRegenerateMaze {
			maze, err = generate.Generate(selectedLevel)
		}

		generatorAnimationData = getGeneratorAnimationData(maze)

		now = time.Now()
		solution = maze.Solve(solverAlgorithm)
		runTimeMicroseconds = time.Since(now).Microseconds()
		stats = generateStats(solution, int(runTimeMicroseconds))

		solverAnimationData = getSolverAnimationData(
			solution,
			maze,
		)

		animationTiming.Reset()

		if err != nil {
			panic(0)
		}

		state = StateGeneration
	}

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
				animationTiming.running = false
			} else {
				playBtnText = "Pause"
				animationTiming.running = true
			}
		},
	}

	resetBtnText := "Reset"
	resetBtn := &Button{
		baseElement: baseElement{
			width:  100,
			height: 20,
		},
		text: &resetBtnText,
		onClick: func() {
			reset(true)
		},
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
		onChange: func() {
			reset(false)
		},
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
		onChange: func() {
			reset(true)
		},
	}

	backtrackingToggle := &Toggle{
		baseElement: baseElement{
			width:  100,
			height: 20,
		},
		text:   "Show backtracking",
		active: &showBacktracking,
	}

	guiElements = append(guiElements, playBtn, resetBtn, levelSelectDropdown, algoDropdown, backtrackingToggle)

	rl.InitWindow(config.Width, config.Height, "Mazen")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.DrawRectangle(0, 0, config.Width, config.Height, borderColor)
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
			generatorAnimationData,
			animationConfig,
		)

		drawGui(guiElements, stats)

		animationTiming.prevTime = time.Now()

		rl.EndDrawing()
	}
}
