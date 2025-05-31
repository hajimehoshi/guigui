// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package layout

import (
	"image"

	"github.com/hajimehoshi/guigui/internal/layoututil"
)

type Size = layoututil.Size

func FixedSize(value int) Size {
	return layoututil.FixedSize(value)
}

func FlexibleSize(value int) Size {
	return layoututil.FlexibleSize(value)
}

func LazySize(f func(rowOrColumn int) Size) Size {
	return layoututil.LazySize(f)
}

type GridLayout struct {
	Bounds    image.Rectangle
	Widths    []Size
	Heights   []Size
	ColumnGap int
	RowGap    int

	widthsInPixels  []int
	heightsInPixels []int
}

var (
	defaultWidths  = []Size{FlexibleSize(1)}
	defaultHeights = []Size{FlexibleSize(1)}
)

func (g *GridLayout) CellBounds(column, row int) image.Rectangle {
	if column < 0 || column >= max(len(g.Widths), 1) {
		return image.Rectangle{}
	}
	if row < 0 {
		return image.Rectangle{}
	}

	var bounds image.Rectangle

	var minX int
	widthCount := max(len(g.Widths), 1)
	if cap(g.widthsInPixels) < widthCount {
		g.widthsInPixels = make([]int, widthCount)
	}
	g.widthsInPixels = g.widthsInPixels[:widthCount]
	widths := g.Widths
	if len(widths) == 0 {
		widths = defaultWidths
	}
	layoututil.WidthsInPixels(g.widthsInPixels, widths, g.Bounds.Dx(), g.ColumnGap)
	for i := range column {
		minX += g.widthsInPixels[i]
		minX += g.ColumnGap

	}
	bounds.Min.X = g.Bounds.Min.X + minX
	bounds.Max.X = g.Bounds.Min.X + minX + g.widthsInPixels[column]

	var minY int
	heightCount := max(len(g.Heights), 1)
	if cap(g.heightsInPixels) < heightCount {
		g.heightsInPixels = make([]int, heightCount)
	}
	g.heightsInPixels = g.heightsInPixels[:heightCount]
	heights := g.Heights
	if len(heights) == 0 {
		heights = defaultHeights
	}
	for loopIndex := range row / heightCount {
		layoututil.HeightsInPixels(g.heightsInPixels, heights, g.Bounds.Dy(), g.RowGap, loopIndex)
		for _, h := range g.heightsInPixels {
			minY += h
			minY += g.RowGap
		}
	}
	layoututil.HeightsInPixels(g.heightsInPixels, heights, g.Bounds.Dy(), g.RowGap, row/heightCount)
	for j := range row % heightCount {
		minY += g.heightsInPixels[j]
		minY += g.RowGap
	}

	bounds.Min.Y = g.Bounds.Min.Y + minY
	bounds.Max.Y = g.Bounds.Min.Y + minY + g.heightsInPixels[row%heightCount]

	return bounds
}
