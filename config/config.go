// Package config - global constants
package config

const (
	Width         int32 = 800
	Height        int32 = 800
	EdgeWidth     int32 = 20
	WallWidth     int32 = EdgeWidth / 4
	MenuBarHeight int32 = 100
	MenuBarWidth        = Width
)

const (
	VerticesPerCol = (Height - MenuBarHeight) / (EdgeWidth)
	VerticesPerRow = Width / (EdgeWidth)
)
