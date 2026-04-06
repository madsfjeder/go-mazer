package draw

import (
	"image/color"

	"maze/config"
	"maze/internal/grid"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
