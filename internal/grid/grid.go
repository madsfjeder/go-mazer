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
	Wall      color.RGBA
	Cell      color.RGBA
	Text      color.RGBA
	DebugWall color.RGBA
}

type Config struct {
	Width  int32
	Height int32
	X      int32
	Y      int32
	Debug  bool
}

type Renderer interface {
	DrawRectangle(xPos, yPos, width, height int32, color color.RGBA)
	DrawText(text string, xPos, yPos, fontSize int32, color color.RGBA)
	Config() Config
	Colors() Colors
}

func (v Vertex) DrawVertex(r Renderer) {
	width := r.Config().Width
	height := r.Config().Height
	xPos := (r.Config().X * height * 2) + width
	yPos := (r.Config().Y * width * 2) + width

	wallColor := r.Colors().Wall

	if r.Config().Debug {
		wallColor = r.Colors().DebugWall
	}

	cellColor := r.Colors().Cell

	if v.Visited {
		r.DrawRectangle(xPos, yPos, width, height, cellColor)
	}

	if v.hasConnectedVertex("top") && !v.TopEdge.IsWall {
		r.DrawRectangle(xPos, yPos-width, width, height, wallColor)
	}

	if v.hasConnectedVertex("right") && !v.RightEdge.IsWall {
		r.DrawRectangle(xPos+width, yPos, width, height, wallColor)
	}

	if v.hasConnectedVertex("bottom") && !v.BottomEdge.IsWall {
		r.DrawRectangle(xPos, yPos+width, width, height, wallColor)
	}

	if v.hasConnectedVertex("left") && !v.LeftEdge.IsWall {
		r.DrawRectangle(xPos-width, yPos, width, height, wallColor)
	}
}

func (v *Vertex) DrawText(r Renderer, s string, fontSize int) {
	width := r.Config().width
	height := r.Config().height
	xPos := (r.Config().x * width * 2) + width
	yPos := (r.Config().y * width * 2) + height
	r.DrawText(s, xPos, yPos, fontSize, r.Colors().Text)
}
