// package Generates a maze and then solves it
package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const EDGE_WIDTH int32 = 20

var DEBUG = false

type Vertex struct {
	visited    bool
	topEdge    *Edge
	rightEdge  *Edge
	bottomEdge *Edge
	leftEdge   *Edge
}

func (v *Vertex) getNextVertex() (*Vertex, error) {
	var upVertex *Vertex
	var rightVertex *Vertex
	var bottomVertex *Vertex
	var leftVertex *Vertex

	options := make([]string, 0)

	if v.topEdge != nil {
		upVertex = v.getConnectedVertex(v.topEdge)

		if upVertex != nil && !upVertex.visited {
			options = append(options, "up")
		}
	}

	if v.rightEdge != nil {
		rightVertex = v.getConnectedVertex(v.rightEdge)
		if rightVertex != nil && !rightVertex.visited {
			options = append(options, "right")
		}
	}

	if v.bottomEdge != nil {
		bottomVertex = v.getConnectedVertex(v.bottomEdge)

		if bottomVertex != nil && !bottomVertex.visited {
			options = append(options, "down")
		}
	}

	if v.leftEdge != nil {
		leftVertex = v.getConnectedVertex(v.leftEdge)

		if leftVertex != nil && !leftVertex.visited {
			options = append(options, "left")
		}
	}

	var next *Vertex

	var dir string
	idx := -1

	if len(options) == 0 {
		return nil, errors.New("doomed!")
	}

	idx = rand.Intn(len(options))
	dir = options[idx]

	switch dir {
	case "up":
		next = upVertex
		v.topEdge.isWall = false
		break

	case "right":
		next = rightVertex
		v.rightEdge.isWall = false
		break

	case "down":
		next = bottomVertex
		v.bottomEdge.isWall = false
		break

	case "left":
		next = leftVertex
		v.leftEdge.isWall = false
		break
	}

	return next, nil
}

func (v *Vertex) getConnectedVertex(e *Edge) *Vertex {
	if e == nil {
		var zero *Vertex
		return zero
	}
	otherVertex := e.v1
	if otherVertex == v {
		otherVertex = e.v2
	}

	return otherVertex
}

func (v *Vertex) hasConnectedTopVertex() bool {
	topVertex := v.getConnectedVertex(v.topEdge)
	if topVertex != nil {
		return true
	}
	return false
}

func (v *Vertex) hasConnectedRightVertex() bool {
	rightVertex := v.getConnectedVertex(v.rightEdge)
	if rightVertex != nil {
		return true
	}
	return false
}

func (v *Vertex) hasConnectedBottomVertex() bool {
	bottomVertex := v.getConnectedVertex(v.bottomEdge)
	if bottomVertex != nil {
		return true
	}
	return false
}

func (v *Vertex) hasConnectedLeftVertex() bool {
	leftVertex := v.getConnectedVertex(v.leftEdge)
	if leftVertex != nil {
		return true
	}
	return false
}

func (v Vertex) drawVertex(x int, y int, color rl.Color) {
	xPos := (int32(x) * EDGE_WIDTH * 2) + EDGE_WIDTH
	yPos := (int32(y) * EDGE_WIDTH * 2) + EDGE_WIDTH

	wallColor := rl.White

	if DEBUG {
		wallColor = rl.Beige
	}

	cellColor := rl.Black
	if v.visited {
		cellColor = color
	}

	if v.visited {
		rl.DrawRectangle(xPos, yPos, EDGE_WIDTH, EDGE_WIDTH, cellColor)
	}

	if v.hasConnectedTopVertex() && !v.topEdge.isWall {
		rl.DrawRectangle(xPos, yPos-EDGE_WIDTH, EDGE_WIDTH, EDGE_WIDTH, wallColor)
	}

	if v.hasConnectedRightVertex() && !v.rightEdge.isWall {
		rl.DrawRectangle(xPos+EDGE_WIDTH, yPos, EDGE_WIDTH, EDGE_WIDTH, wallColor)
	}

	if v.hasConnectedBottomVertex() && !v.bottomEdge.isWall {
		rl.DrawRectangle(xPos, yPos+EDGE_WIDTH, EDGE_WIDTH, EDGE_WIDTH, wallColor)
	}

	if v.hasConnectedLeftVertex() && !v.leftEdge.isWall {
		rl.DrawRectangle(xPos-EDGE_WIDTH, yPos, EDGE_WIDTH, EDGE_WIDTH, wallColor)
	}
}

func (v *Vertex) drawText(x int, y int, s string) {
	xPos := (int32(x) * EDGE_WIDTH * 2) + EDGE_WIDTH
	yPos := (int32(y) * EDGE_WIDTH * 2) + EDGE_WIDTH
	rl.DrawText(s, xPos, yPos, 13, rl.Green)
}

type Edge struct {
	isWall bool
	v1     *Vertex
	v2     *Vertex
}

type Stack[T comparable] struct {
	items []T
}

func (s *Stack[T]) push(e T) {
	s.items = append(s.items, e)
}

func (s *Stack[T]) pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, errors.New("no more history")
	}
	lastIndex := len(s.items) - 1
	element := s.items[lastIndex]

	var zero T

	s.items[lastIndex] = zero
	s.items = s.items[:lastIndex]

	return element, nil
}

func (s Stack[T]) print() {
	for _, v := range s.items {
		fmt.Println(v)
	}
}

func (s Stack[T]) copy() Stack[T] {
	itemsSlice := make([]T, len(s.items))
	copy(itemsSlice, s.items)

	return Stack[T]{
		items: itemsSlice,
	}
}

func (s *Stack[T]) reverse() {
	slices.Reverse(s.items)
}

func (s *Stack[T]) findOrder(v T) int {
	for i, e := range s.items {
		if v == e {
			return i
		}
	}

	return -1
}

func (s *Stack[T]) length() int {
	return len(s.items)
}

func newStack[T comparable]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

func main() {
	var width int32 = 800
	var height int32 = 800
	intervalPtr := flag.Int("interval", 5, "set the rendering interval")
	debugPtr := flag.Bool("debug", false, "turns debugging on")
	flag.Parse()

	if *debugPtr {
		DEBUG = true
	}

	verticesPerCol := height / (EDGE_WIDTH * 2)
	verticesPerRow := width / (EDGE_WIDTH * 2)

	matrix := make([][]*Vertex, verticesPerRow)

	for i := range matrix {
		matrix[i] = make([]*Vertex, verticesPerCol)
	}

	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = &Vertex{}
		}
	}

	for i := range matrix {
		for j := range matrix[i] {
			currentVertex := matrix[i][j]

			var topVertex *Vertex
			var rightVertex *Vertex
			var botVertex *Vertex
			var leftVertex *Vertex

			if j > 0 {
				topVertex = matrix[i][j-1]
				topVertexBottomConnection := topVertex.getConnectedVertex(topVertex.bottomEdge)

				if topVertexBottomConnection != nil && topVertexBottomConnection == currentVertex {
				} else if currentVertex.topEdge == nil {
					if currentVertex.topEdge == nil {
						edge := Edge{
							isWall: true,
							v1:     currentVertex,
							v2:     topVertex,
						}
						currentVertex.topEdge = &edge
						topVertex.bottomEdge = &edge
					}
				}
			}

			if i < int(verticesPerRow-1) {
				rightVertex = matrix[i+1][j]
				rightVertexLeftConnection := rightVertex.getConnectedVertex(rightVertex.leftEdge)

				if rightVertexLeftConnection != nil && rightVertexLeftConnection == currentVertex {
				} else if currentVertex.rightEdge == nil {
					edge := Edge{
						isWall: true,
						v1:     currentVertex,
						v2:     rightVertex,
					}

					currentVertex.rightEdge = &edge
					rightVertex.leftEdge = &edge
				}
			}

			if j < int(verticesPerCol-1) {
				botVertex = matrix[i][j+1]
				botVertexTopConnection := botVertex.getConnectedVertex(botVertex.topEdge)

				if botVertexTopConnection != nil && botVertexTopConnection == currentVertex {
				} else if currentVertex.bottomEdge == nil {
					edge := Edge{
						isWall: true,
						v1:     currentVertex,
						v2:     botVertex,
					}
					currentVertex.bottomEdge = &edge
					botVertex.topEdge = &edge
				}
			}

			if i > 0 {
				leftVertex = matrix[i-1][j]
				leftVertexRightConnection := leftVertex.getConnectedVertex(leftVertex.rightEdge)

				if leftVertexRightConnection != nil && leftVertexRightConnection == currentVertex {
				} else if currentVertex.leftEdge == nil {
					edge := Edge{
						isWall: true,
						v1:     currentVertex,
						v2:     leftVertex,
					}

					currentVertex.leftEdge = &edge
					leftVertex.rightEdge = &edge
				}
			}
		}
	}

	history := newStack[*Vertex]()
	allSteps := newStack[*Vertex]()
	backtracking := newStack[*Vertex]()
	currentVertex := matrix[0][0]
	currentVertex.visited = true

	for i := 0; i < 10000; i++ {
		currentVertex.visited = true
		nextVertex, err := currentVertex.getNextVertex()
		if nextVertex == nil {
			currentVertex, err = history.pop()
			if err != nil || currentVertex == nil {
				break
			}
			allSteps.push(currentVertex)
			backtracking.push(currentVertex)
			continue
		}

		history.push(currentVertex)
		allSteps.push(currentVertex)
		currentVertex = nextVertex
	}

	rl.InitWindow(width+EDGE_WIDTH, height+EDGE_WIDTH, "Mazen")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	reversedSteps := allSteps.copy()
	reversedSteps.reverse()

	animationComplete := false
	var timeAcc int64 = 0
	var interval int64 = int64(*intervalPtr)
	prevTime := time.Now()
	matrixToDraw := make([][]*Vertex, verticesPerRow)

	for i := range matrixToDraw {
		matrixToDraw[i] = make([]*Vertex, verticesPerCol)
	}

	var elementToDraw *Vertex

	rl.DrawRectangle(0, 0, width+EDGE_WIDTH, height+EDGE_WIDTH, rl.Black)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		delta := time.Now().Sub(prevTime)

		if !animationComplete {
			timeAcc += delta.Milliseconds()
			if timeAcc >= interval {
				timeAcc -= interval
				e, _ := reversedSteps.pop()
				elementToDraw = e
			}
		}

		var x, y int
		for i := range matrix {
			for j := range matrix {
				v := matrix[i][j]
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
					e.drawVertex(i, j, rl.White)
					if DEBUG {
						backTracked := backtracking.findOrder(e)
						if backTracked != -1 {
							e.drawVertex(i, j, rl.Red)
							e.drawText(i, j, strconv.Itoa(backTracked))
						}
					}
				}
			}
		}

		if elementToDraw != nil {
			elementToDraw.drawVertex(x, y, rl.Blue)
		}

		prevTime = time.Now()
		rl.EndDrawing()
	}
}
