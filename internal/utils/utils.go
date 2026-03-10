// Package utils - simple utilities for other packages
package utils

import "math"

func GetCartesianDistance(x1, y1, x2, y2 int32) float64 {
	return math.Sqrt((math.Pow(float64(x2)-float64(x1), 2)) + (math.Pow(float64(y2)-float64(y1), 2)))
}
