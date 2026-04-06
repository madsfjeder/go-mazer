package draw

import (
	"fmt"
	"time"

	"maze/config"
	"maze/generate"
	"maze/internal/grid"
	"maze/internal/stack"
)

type AnimationConfig struct {
	colors           grid.Colors
	showBacktracking *bool
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
	matrixToDraw := grid.New(int(config.VerticesPerRow), int(config.VerticesPerRow))

	return GeneratorAnimationData{
		animationData: animationData{
			matrixToRender:  matrixToDraw,
			completedMatrix: maze,
			itemsToRender:   maze.Steps.PopAllWithIdx(),
		},
	}
}

func (g *GeneratorAnimationData) Draw(animationConfig AnimationConfig) {
	for i := range g.completedMatrix.Matrix {
		for j := range g.completedMatrix.Matrix[i] {
			e := g.completedMatrix.Matrix[i][j]

			if e != nil && e.IsPath {
				cellType := grid.Path

				if e.IsSplit {
					cellType = grid.Split
				}

				r := NewRaylibRenderer(i, j, cellType, animationConfig.colors, false)
				e.DrawVertex(r)
			} else {
				r := NewRaylibRenderer(i, j, grid.EmptyCell, animationConfig.colors, false)
				empty := grid.Vertex{
					IsPath:          false,
					VisitedBySolver: false,
				}
				empty.DrawVertex(r)
			}
		}
	}
}

func (g *GeneratorAnimationData) DrawWalls(animationConfig AnimationConfig) {
	for i := range g.completedMatrix.Matrix {
		for j := range g.completedMatrix.Matrix[i] {
			e := g.completedMatrix.Matrix[i][j]
			if e != nil {
				r := NewRaylibRenderer(i, j, grid.Wall, animationConfig.colors, false)
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
	var zero *grid.Vertex
	startItem := CurrentElement{
		element: zero,
		x:       0,
		y:       0,
	}

	for i := range maze.Matrix {
		for j := range maze.Matrix[i] {
			e := maze.Matrix[i][j]
			if e.IsStart {
				startItem.element = e
				startItem.x = i
				startItem.y = j
				break
			}
		}
	}

	return &SolverAnimationData{
		animationData: animationData{
			itemsToRender:   make([]stack.StackItem[*grid.Vertex], 0),
			matrixToRender:  grid.New(int(config.VerticesPerRow), int(config.VerticesPerCol)),
			completedMatrix: maze,
		},
		path:             solution,
		fullStackOfItems: solverItems,
		currentElement:   startItem,
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
