// Copyright 2025 Hajime Hoshi

package guigui

import (
	"image"
	"math"
)

type Constraints struct {
	minSize         image.Point
	maxSizeMinusMax image.Point
}

func (c *Constraints) MinSize() image.Point {
	return c.minSize
}

func (c *Constraints) MaxSize() image.Point {
	return c.maxSizeMinusMax.Add(image.Pt(math.MaxInt, math.MaxInt))
}

func FixedWidthConstraints(w int) Constraints {
	return Constraints{
		minSize:         image.Pt(w, 0),
		maxSizeMinusMax: image.Pt(w, math.MaxInt).Sub(image.Pt(math.MaxInt, math.MaxInt)),
	}
}

func FixedHeightConstraints(h int) Constraints {
	return Constraints{
		minSize:         image.Pt(0, h),
		maxSizeMinusMax: image.Pt(math.MaxInt, h).Sub(image.Pt(math.MaxInt, math.MaxInt)),
	}
}
