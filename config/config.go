// Package config - global constants
package config

const (
	Width     int32 = 800
	Height    int32 = 800
	EdgeWidth int32 = 10
)

const (
	VerticesPerCol = Height / (EdgeWidth * 2)
	VerticesPerRow = Width / (EdgeWidth * 2)
)
