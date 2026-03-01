// Package config - global constants
package config

const (
	EdgeWidth     int32 = 20
	WallWidth     int32 = EdgeWidth / 4
	Width         int32 = 800 + WallWidth
	Height        int32 = 800 - WallWidth
	MenuBarHeight int32 = 30
	MenuBarWidth        = Width - WallWidth
)

const (
	VerticesPerCol = (Height - MenuBarHeight) / (EdgeWidth)
	VerticesPerRow = Width / (EdgeWidth)
)
