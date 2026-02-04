// Package grid Vertices, edges and helpers functions for drawing
package grid

import (
	"errors"
	"image/color"
	"math/rand"
)

type Edge struct {
	IsWall  bool
	Vertex1 *Vertex
	Vertex2 *Vertex
}

type Vertex struct {
	Visited    bool
	IsStart    bool
	IsEnd      bool
	TopEdge    *Edge
	RightEdge  *Edge
	BottomEdge *Edge
	LeftEdge   *Edge
}

func (v *Vertex) GetNextVertex() (*Vertex, error) {
	var upVertex *Vertex
	var rightVertex *Vertex
	var bottomVertex *Vertex
	var leftVertex *Vertex

	options := make([]string, 0)

	if v.TopEdge != nil {
		upVertex = v.GetConnectedVertex(v.TopEdge)

		if upVertex != nil && !upVertex.Visited {
			options = append(options, "up")
		}
	}

	if v.RightEdge != nil {
		rightVertex = v.GetConnectedVertex(v.RightEdge)
		if rightVertex != nil && !rightVertex.Visited {
			options = append(options, "right")
		}
	}

	if v.BottomEdge != nil {
		bottomVertex = v.GetConnectedVertex(v.BottomEdge)

		if bottomVertex != nil && !bottomVertex.Visited {
			options = append(options, "down")
		}
	}

	if v.LeftEdge != nil {
		leftVertex = v.GetConnectedVertex(v.LeftEdge)

		if leftVertex != nil && !leftVertex.Visited {
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
	if e == nil {
		var zero *Vertex
		return zero
	}
	otherVertex := e.Vertex1
	if otherVertex == v {
		otherVertex = e.Vertex2
	}

	return otherVertex
}

func (v *Vertex) hasConnectedVertex(dir string) bool {
	switch dir {
	case "top":
		topVertex := v.GetConnectedVertex(v.TopEdge)
		return topVertex != nil

	case "right":
		rightVertex := v.GetConnectedVertex(v.RightEdge)
		return rightVertex != nil

	case "bottom":
		bottomVertex := v.GetConnectedVertex(v.BottomEdge)
		return bottomVertex != nil

	case "left":
		leftVertex := v.GetConnectedVertex(v.LeftEdge)
		return leftVertex != nil
	}

	return false
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
		Visited:    v.Visited,
		TopEdge:    &topEdge,
		RightEdge:  &rightEdge,
		BottomEdge: &bottomEdge,
		LeftEdge:   &leftEdge,
	}
}

type Colors struct {
	Start     color.RGBA
	End       color.RGBA
	Wall      color.RGBA
	Cell      color.RGBA
	CPU       color.RGBA
	Text      color.RGBA
	DebugWall color.RGBA
}

type CellType int

const (
	Wall CellType = iota
	Path
	CPU
)

type Config struct {
	EdgeWidth int32
	X         int32
	Y         int32
	Debug     bool
	CellType  CellType
}

type Renderer interface {
	DrawRectangle(xPos, yPos, width, height int32, color color.RGBA)
	DrawText(text string, xPos, yPos, fontSize int32, color color.RGBA)
	DrawCircle(x, y int32, radius float32, color color.RGBA)
	Config() Config
	Colors() Colors
}

func (v Vertex) DrawVertex(r Renderer) {
	cellType := r.Config().CellType
	edgeWidth := r.Config().EdgeWidth
	xPos := (r.Config().X * edgeWidth * 2) + edgeWidth
	yPos := (r.Config().Y * edgeWidth * 2) + edgeWidth

	cellColor := r.Colors().Wall

	switch cellType {
	case Wall:
		break

	case Path:
		cellColor = r.Colors().Cell

	case CPU:
		cellColor = r.Colors().CPU
	}

	if v.Visited {
		r.DrawRectangle(xPos, yPos, edgeWidth, edgeWidth, cellColor)
	}

	if v.IsStart {
		r.DrawCircle(xPos+(edgeWidth/2), yPos+(edgeWidth/2), float32(edgeWidth/2), r.Colors().Start)
	}

	if v.IsEnd {
		r.DrawCircle(xPos+(edgeWidth/2), yPos+(edgeWidth/2), float32(edgeWidth/2), r.Colors().End)
	}

	if v.hasConnectedVertex("top") && !v.TopEdge.IsWall {
		r.DrawRectangle(xPos, yPos-edgeWidth, edgeWidth, edgeWidth, cellColor)
	}

	if v.hasConnectedVertex("right") && !v.RightEdge.IsWall {
		r.DrawRectangle(xPos+edgeWidth, yPos, edgeWidth, edgeWidth, cellColor)
	}

	if v.hasConnectedVertex("bottom") && !v.BottomEdge.IsWall {
		r.DrawRectangle(xPos, yPos+edgeWidth, edgeWidth, edgeWidth, cellColor)
	}

	if v.hasConnectedVertex("left") && !v.LeftEdge.IsWall {
		r.DrawRectangle(xPos-edgeWidth, yPos, edgeWidth, edgeWidth, cellColor)
	}
}

func (v *Vertex) DrawText(r Renderer, s string, fontSize int32) {
	edgeWidth := r.Config().EdgeWidth
	xPos := (r.Config().X * edgeWidth * 2) + edgeWidth
	yPos := (r.Config().Y * edgeWidth * 2) + edgeWidth
	r.DrawText(s, xPos, yPos, fontSize, r.Colors().Text)
}
