// Package grid Vertices, edges and helpers functions for drawing
package grid

import (
	"errors"
	"image/color"
	"math/rand"
	"strconv"

	"maze/config"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Edge struct {
	// Consider IsPath instead to streamline
	IsWall          bool
	VisitedBySolver bool
	Vertex1         *Vertex
	Vertex2         *Vertex
}

type Vertex struct {
	IsPath           bool
	IsStart          bool
	IsEnd            bool
	IsPartOfSolution bool
	IsBacktracking   bool
	VisitedBySolver  bool
	IsSplit          bool
	TopEdge          *Edge
	RightEdge        *Edge
	BottomEdge       *Edge
	LeftEdge         *Edge
	// Actual travel cost from start to this vertex
	G float32
	// Heuristic value to goal
	H float32
	// G + H
	F      float32
	Closed bool
}

func New(width, height int) [][]*Vertex {
	g := make([][]*Vertex, config.VerticesPerRow)
	for i := range g {
		g[i] = make([]*Vertex, config.VerticesPerCol)
	}

	return g
}

func (v *Vertex) CanSplit() bool {
	var upVertex *Vertex
	var rightVertex *Vertex
	var bottomVertex *Vertex
	var leftVertex *Vertex

	options := make([]string, 0)

	if v.TopEdge != nil {
		upVertex = v.GetConnectedVertex(v.TopEdge)

		if upVertex != nil && !upVertex.IsPath && !upVertex.VisitedBySolver {
			options = append(options, "up")
		}
	}

	if v.RightEdge != nil {
		rightVertex = v.GetConnectedVertex(v.RightEdge)
		if rightVertex != nil && !rightVertex.IsPath && !rightVertex.VisitedBySolver {
			options = append(options, "right")
		}
	}

	if v.BottomEdge != nil {
		bottomVertex = v.GetConnectedVertex(v.BottomEdge)

		if bottomVertex != nil && !bottomVertex.IsPath && !bottomVertex.VisitedBySolver {
			options = append(options, "down")
		}
	}

	if v.LeftEdge != nil {
		leftVertex = v.GetConnectedVertex(v.LeftEdge)

		if leftVertex != nil && !leftVertex.IsPath && !leftVertex.VisitedBySolver {
			options = append(options, "left")
		}
	}

	return len(options) > 1
}

// GetNextVertex - For generation
func (v *Vertex) GetNextVertex() (*Vertex, error) {
	var upVertex *Vertex
	var rightVertex *Vertex
	var bottomVertex *Vertex
	var leftVertex *Vertex

	options := make([]string, 0)

	if v.TopEdge != nil {
		upVertex = v.GetConnectedVertex(v.TopEdge)

		if upVertex != nil && !upVertex.IsPath {
			options = append(options, "up")
		}
	}

	if v.RightEdge != nil {
		rightVertex = v.GetConnectedVertex(v.RightEdge)
		if rightVertex != nil && !rightVertex.IsPath {
			options = append(options, "right")
		}
	}

	if v.BottomEdge != nil {
		bottomVertex = v.GetConnectedVertex(v.BottomEdge)

		if bottomVertex != nil && !bottomVertex.IsPath {
			options = append(options, "down")
		}
	}

	if v.LeftEdge != nil {
		leftVertex = v.GetConnectedVertex(v.LeftEdge)

		if leftVertex != nil && !leftVertex.IsPath {
			options = append(options, "left")
		}
	}

	var next *Vertex

	var dir string

	if len(options) == 0 {
		return nil, errors.New("doomed")
	}

	idx := rand.Intn(len(options))
	dir = options[idx]

	switch dir {
	case "up":
		next = upVertex
		v.TopEdge.IsWall = false

	case "right":
		next = rightVertex
		v.RightEdge.IsWall = false

	case "down":
		next = bottomVertex
		v.BottomEdge.IsWall = false

	case "left":
		next = leftVertex
		v.LeftEdge.IsWall = false
	}

	return next, nil
}

func (v *Vertex) GetConnectedVertex(e *Edge) *Vertex {
	var zero *Vertex
	if e == nil {
		return zero
	}

	if e.Vertex1 != v && e.Vertex2 != v {
		return zero
	}

	otherVertex := e.Vertex1
	if otherVertex == v {
		otherVertex = e.Vertex2
	}

	return otherVertex
}

func (v *Vertex) hasConnectedVertex(dir string, shouldBePath bool) bool {
	var connectedVertex *Vertex

	switch dir {
	case "top":
		connectedVertex = v.GetConnectedVertex(v.TopEdge)

	case "right":
		connectedVertex = v.GetConnectedVertex(v.RightEdge)

	case "bottom":
		connectedVertex = v.GetConnectedVertex(v.BottomEdge)

	case "left":
		connectedVertex = v.GetConnectedVertex(v.LeftEdge)
	}

	if shouldBePath {
		return connectedVertex != nil && connectedVertex.IsPath
	}

	return connectedVertex != nil
}

func (v *Vertex) GetNeighbours(allowVisited bool) []*Vertex {
	neighbours := make([]*Vertex, 0)
	hasLeft := v.hasConnectedVertex("left", true)
	hasBottom := v.hasConnectedVertex("bottom", true)
	hasRight := v.hasConnectedVertex("right", true)
	hasTop := v.hasConnectedVertex("top", true)

	if hasTop && !v.TopEdge.IsWall {
		top := v.GetConnectedVertex(v.TopEdge)
		if allowVisited || !top.VisitedBySolver {
			neighbours = append(neighbours, top)
		}
	}

	if hasRight && !v.RightEdge.IsWall {
		right := v.GetConnectedVertex(v.RightEdge)
		if allowVisited || !right.VisitedBySolver {
			neighbours = append(neighbours, right)
		}
	}

	if hasBottom && !v.BottomEdge.IsWall {
		bottom := v.GetConnectedVertex(v.BottomEdge)
		if allowVisited || !bottom.VisitedBySolver {
			neighbours = append(neighbours, bottom)
		}
	}

	if hasLeft && !v.LeftEdge.IsWall {
		left := v.GetConnectedVertex(v.LeftEdge)
		if allowVisited || !left.VisitedBySolver {
			neighbours = append(neighbours, left)
		}
	}

	return neighbours
}

// VisitNextVertex - always goes left first is possible, so it's like depth first search
func (v *Vertex) VisitNextVertex() *Vertex {
	options := []string{"left", "bottom", "right", "top"}

	var nextVertex *Vertex

	for _, dir := range options {
		if !v.hasConnectedVertex(dir, true) {
			continue
		}

		var e *Edge

		switch dir {
		case "left":
			e = v.LeftEdge
		case "bottom":
			e = v.BottomEdge
		case "right":
			e = v.RightEdge
		case "top":
			e = v.TopEdge
		}

		if e == nil || e.IsWall {
			continue
		}

		t := v.GetConnectedVertex(e)

		if t != nil && !t.VisitedBySolver {
			e.VisitedBySolver = true
			nextVertex = t
			break
		}
	}

	return nextVertex
}

func (v *Vertex) Copy() Vertex {
	var topEdge Edge
	var rightEdge Edge
	var bottomEdge Edge
	var leftEdge Edge

	if v.TopEdge != nil {
		topEdge = *v.TopEdge
	}

	if v.RightEdge != nil {
		rightEdge = *v.RightEdge
	}

	if v.BottomEdge != nil {
		bottomEdge = *v.BottomEdge
	}

	if v.LeftEdge != nil {
		leftEdge = *v.LeftEdge
	}

	return Vertex{
		IsPath:     v.IsPath,
		TopEdge:    &topEdge,
		RightEdge:  &rightEdge,
		BottomEdge: &bottomEdge,
		LeftEdge:   &leftEdge,
	}
}

type Colors struct {
	Start        color.RGBA
	End          color.RGBA
	Split        color.RGBA
	Backtracking color.RGBA
	Solution     color.RGBA
	Wall         color.RGBA
	EmptyCell    color.RGBA
	Cell         color.RGBA
	CPU          color.RGBA
	Text         color.RGBA
	DebugWall    color.RGBA
}

type CellType int

const (
	Wall CellType = iota
	Split
	EmptyCell
	Path
	CPU
	Solution
	Backtracking
)

type Config struct {
	EdgeWidth        int32
	X                int32
	Y                int32
	Debug            bool
	CellType         CellType
	ShowBacktracking bool
}

type Renderer interface {
	DrawRectangle(xPos, yPos, width, height int32, color color.RGBA)
	DrawRectangleRounded(xPos, yPos, width, height int32, roundness float32, color color.RGBA)
	DrawTile(xPos, yPos, width, height int32, color color.RGBA)
	DrawText(text string, xPos, yPos, fontSize int32, color color.RGBA)
	DrawCircle(x, y int32, radius float32, color color.RGBA)
	Config() Config
	Colors() Colors
}

func (v *Vertex) DrawVertex(r Renderer) {
	cellType := r.Config().CellType
	edgeWidth := config.EdgeWidth
	wallWidth := config.WallWidth
	showBacktracking := r.Config().ShowBacktracking
	xPos := (r.Config().X * edgeWidth) + config.Padding
	yPos := (r.Config().Y * edgeWidth) + config.MenuBarHeight + config.Padding

	cellColor := r.Colors().Wall
	DEBUG := false

	switch cellType {
	case Wall:
		break

	case Path:
		cellColor = r.Colors().Cell

	case CPU:
		cellColor = r.Colors().CPU

	case Solution:
		cellColor = r.Colors().Solution

	case EmptyCell:
		cellColor = r.Colors().EmptyCell

	case Split:
		// Change to highlight splits
		cellColor = r.Colors().Cell
	}

	edgeColor := r.Colors().Wall

	if cellType == Wall {
		shadowColor := rl.NewColor(0, 0, 0, 90)
		if v.TopEdge != nil && (v.TopEdge.IsWall && !v.GetConnectedVertex(v.TopEdge).IsPath) {
			r.DrawRectangle(xPos, yPos-(wallWidth/2), edgeWidth, wallWidth-config.Padding, edgeColor)
			r.DrawRectangle(xPos+1, yPos+1, edgeWidth, 2, shadowColor)
		}

		if v.RightEdge != nil && (v.RightEdge.IsWall || !v.GetConnectedVertex(v.RightEdge).IsPath) {
			r.DrawRectangle(xPos+edgeWidth-(wallWidth/2), yPos, wallWidth-config.Padding, edgeWidth, edgeColor)

			r.DrawRectangle(xPos+edgeWidth+1, yPos+1, 2, edgeWidth, shadowColor)
		}

		if v.BottomEdge != nil && (v.BottomEdge.IsWall || !v.GetConnectedVertex(v.BottomEdge).IsPath) {
			r.DrawRectangle(xPos, yPos+edgeWidth-(wallWidth/2), edgeWidth, wallWidth-config.Padding, edgeColor)
			r.DrawRectangle(xPos+1, yPos+edgeWidth+1, edgeWidth, 2, shadowColor)
		}

		if v.LeftEdge != nil && (v.LeftEdge.IsWall || !v.GetConnectedVertex(v.LeftEdge).IsPath) {
			r.DrawRectangle(xPos-(wallWidth/2), yPos, wallWidth-config.Padding, edgeWidth, edgeColor)
			r.DrawRectangle(xPos+1, yPos+1, 2, edgeWidth, shadowColor)
		}
		return
	}

	if cellType == EmptyCell {
		r.DrawRectangle(xPos, yPos, edgeWidth, edgeWidth, cellColor)
		return
	}

	if cellType == Solution {
		r.DrawTile(xPos, yPos, edgeWidth, edgeWidth, cellColor)
		if DEBUG {
			r.DrawText(strconv.FormatFloat(float64(v.H), 'f', 2, 64), xPos+3, yPos+3, 8, rl.Black)
		}
	}

	if cellType == Backtracking && showBacktracking {
		cellColor = r.Colors().Backtracking
		r.DrawTile(xPos, yPos, edgeWidth, edgeWidth, cellColor)
		return
	}

	if cellType == CPU {
		r.DrawRectangle(xPos, yPos, edgeWidth, edgeWidth, cellColor)
		return
	}
	if v.IsStart {
		r.DrawTile(xPos, yPos, edgeWidth, edgeWidth, r.Colors().Start)
		return
	}

	if v.IsEnd {
		r.DrawTile(xPos, yPos, edgeWidth, edgeWidth, r.Colors().End)
		return
	}

	if v.IsPath && cellType != Solution {
		r.DrawTile(xPos, yPos, edgeWidth, edgeWidth, r.Colors().Cell)
	}
}

func (v *Vertex) DrawText(r Renderer, s string, fontSize int32) {
	edgeWidth := r.Config().EdgeWidth
	xPos := (r.Config().X * edgeWidth * 2) + edgeWidth
	yPos := (r.Config().Y * edgeWidth * 2) + edgeWidth
	r.DrawText(s, xPos, yPos, fontSize, r.Colors().Text)
}
