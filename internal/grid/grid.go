// Package grid Vertices, edges and helpers functions for drawing
package grid

import (
	"errors"
	"image/color"
	"math/rand"

	"maze/config"
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

		if !t.VisitedBySolver {
			e.VisitedBySolver = true
			nextVertex = t
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
	xPos := (r.Config().X * edgeWidth)
	yPos := (r.Config().Y * edgeWidth) + config.MenuBarHeight

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
	if DEBUG {
		edgeColor = r.Colors().DebugWall
	}

	leftVertex := v.GetConnectedVertex(v.LeftEdge)
	bottomVertex := v.GetConnectedVertex(v.BottomEdge)
	rightVertex := v.GetConnectedVertex(v.RightEdge)
	topVertex := v.GetConnectedVertex(v.TopEdge)

	if cellType == Wall {
		if v.TopEdge == nil || v.TopEdge.IsWall || !v.GetConnectedVertex(v.TopEdge).IsPath {
			r.DrawRectangle(xPos, yPos, edgeWidth+wallWidth, wallWidth, edgeColor)
		}

		if v.RightEdge == nil || v.RightEdge.IsWall || !v.GetConnectedVertex(v.RightEdge).IsPath {
			r.DrawRectangle(xPos+edgeWidth, yPos, wallWidth, edgeWidth+wallWidth, edgeColor)
		}

		if v.BottomEdge == nil || v.BottomEdge.IsWall || !v.GetConnectedVertex(v.BottomEdge).IsPath {
			r.DrawRectangle(xPos, yPos+edgeWidth, edgeWidth+wallWidth, wallWidth, edgeColor)
		}

		if v.LeftEdge == nil || v.LeftEdge.IsWall || !v.GetConnectedVertex(v.LeftEdge).IsPath {
			r.DrawRectangle(xPos, yPos, wallWidth, edgeWidth+wallWidth, edgeColor)
		}
		return
	}

	if cellType == EmptyCell {
		r.DrawRectangle(xPos, yPos, edgeWidth, edgeWidth, cellColor)
		return
	}

	if cellType == Solution {
		hasLeftPath := v.LeftEdge != nil && !v.LeftEdge.IsWall && leftVertex != nil && (leftVertex.IsPartOfSolution || (leftVertex.IsBacktracking && showBacktracking))
		hasBottomPath := v.BottomEdge != nil && !v.BottomEdge.IsWall && bottomVertex != nil && (bottomVertex.IsPartOfSolution || (bottomVertex.IsBacktracking && showBacktracking))
		hasRightPath := v.RightEdge != nil && !v.RightEdge.IsWall && rightVertex != nil && (rightVertex.IsPartOfSolution || (rightVertex.IsBacktracking && showBacktracking))
		hasTopPath := v.TopEdge != nil && !v.TopEdge.IsWall && topVertex != nil && (topVertex.IsPartOfSolution || (topVertex.IsBacktracking && showBacktracking))

		padding := edgeWidth/4 + 2
		pathWidth := edgeWidth / 2

		r.DrawRectangleRounded(xPos+padding, yPos+padding, pathWidth, pathWidth, 50, cellColor)

		if hasLeftPath {
			r.DrawRectangle(xPos, yPos+padding, pathWidth, pathWidth, cellColor)
		}

		if hasBottomPath {
			r.DrawRectangle(xPos+padding, yPos+edgeWidth/2, pathWidth, pathWidth, cellColor)
		}

		if hasRightPath {
			r.DrawRectangle(xPos+edgeWidth/2, yPos+padding, pathWidth, pathWidth, cellColor)
		}

		if hasTopPath {
			r.DrawRectangle(xPos+padding, yPos, pathWidth, pathWidth, cellColor)
		}

		return
	}

	if cellType == Backtracking {
		cellColor = r.Colors().Backtracking

		hasLeftPath := v.LeftEdge != nil && !v.LeftEdge.IsWall && (leftVertex.IsBacktracking || leftVertex.IsPartOfSolution)
		hasBottomPath := v.BottomEdge != nil && !v.BottomEdge.IsWall && (bottomVertex.IsBacktracking || bottomVertex.IsPartOfSolution)
		hasRightPath := v.RightEdge != nil && !v.RightEdge.IsWall && (rightVertex.IsBacktracking || rightVertex.IsPartOfSolution)
		hasTopPath := v.TopEdge != nil && !v.TopEdge.IsWall && (topVertex.IsBacktracking || topVertex.IsPartOfSolution)

		padding := edgeWidth/4 + 2
		pathWidth := edgeWidth / 2

		r.DrawRectangleRounded(xPos+padding, yPos+padding, pathWidth, pathWidth, 50, cellColor)

		if hasLeftPath {
			r.DrawRectangle(xPos, yPos+padding, pathWidth, pathWidth, cellColor)
		}

		if hasBottomPath {
			r.DrawRectangle(xPos+padding, yPos+edgeWidth/2, pathWidth, pathWidth, cellColor)
		}

		if hasRightPath {
			r.DrawRectangle(xPos+edgeWidth/2, yPos+padding, pathWidth, pathWidth, cellColor)
		}

		if hasTopPath {
			r.DrawRectangle(xPos+padding, yPos, pathWidth, pathWidth, cellColor)
		}

		return
	}

	if cellType == CPU {
		r.DrawRectangle(xPos, yPos, edgeWidth, edgeWidth, cellColor)
		return
	}

	if v.IsPath && cellType != Solution {
		r.DrawRectangle(xPos, yPos, edgeWidth, edgeWidth, cellColor)
	}

	if v.IsStart {
		r.DrawCircle(xPos+(edgeWidth/2)+(wallWidth/2), yPos+(edgeWidth/2)+(wallWidth/2), float32(edgeWidth/2)-float32(wallWidth/2), r.Colors().Start)
	}

	if v.IsEnd {
		r.DrawCircle(xPos+(edgeWidth/2)+(wallWidth/2), yPos+(edgeWidth/2)+(wallWidth/2), float32(edgeWidth/2)-float32(wallWidth/2), r.Colors().End)
	}
}

func (v *Vertex) DrawText(r Renderer, s string, fontSize int32) {
	edgeWidth := r.Config().EdgeWidth
	xPos := (r.Config().X * edgeWidth * 2) + edgeWidth
	yPos := (r.Config().Y * edgeWidth * 2) + edgeWidth
	r.DrawText(s, xPos, yPos, fontSize, r.Colors().Text)
}
