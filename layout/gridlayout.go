// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package layout

import (
	"image"
)

type Size struct {
	typ   sizeType
	value int
	lazy  func(rowOrColumn int) Size
}

type sizeType int

const (
	sizeTypeFixed sizeType = iota
	sizeTypeFlexible
	sizeTypeLazy
)

func FixedSize(value int) Size {
	return Size{
		typ:   sizeTypeFixed,
		value: value,
	}
}

func FlexibleSize(value int) Size {
	return Size{
		typ:   sizeTypeFlexible,
		value: value,
	}
}

func LazySize(f func(rowOrColumn int) Size) Size {
	return Size{
		typ:   sizeTypeLazy,
		value: 0,
		lazy:  f,
	}
}

var (
	defaultWidths  = []Size{FlexibleSize(1)}
	defaultHeights = []Size{FlexibleSize(1)}
)

type GridLayout struct {
	Bounds    image.Rectangle
	Widths    []Size
	Heights   []Size
	ColumnGap int
	RowGap    int

	widthsInPixels  []int
	heightsInPixels []int
}

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
	g.getWidthsInPixels(g.widthsInPixels)
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
	for loopIndex := range row / heightCount {
		g.getHeightsInPixels(g.heightsInPixels, loopIndex)
		for _, h := range g.heightsInPixels {
			minY += h
			minY += g.RowGap
		}
	}
	g.getHeightsInPixels(g.heightsInPixels, row/heightCount)
	for j := range row % heightCount {
		minY += g.heightsInPixels[j]
		minY += g.RowGap
	}

	bounds.Min.Y = g.Bounds.Min.Y + minY
	bounds.Max.Y = g.Bounds.Min.Y + minY + g.heightsInPixels[row%heightCount]

	return bounds
}

func (g *GridLayout) getWidthsInPixels(widthsInPixels []int) {
	widths := g.Widths
	if len(widths) == 0 {
		widths = defaultWidths
	}

	// Calculate widths in pixels.
	restW := g.Bounds.Dx()
	restW -= (len(widths) - 1) * g.ColumnGap
	if restW < 0 {
		restW = 0
	}
	var denomW int

	for i, width := range widths {
		switch width.typ {
		case sizeTypeFixed:
			widthsInPixels[i] = width.value
		case sizeTypeFlexible:
			widthsInPixels[i] = 0
			denomW += width.value
		default:
			panic("layout: only FixedSize and FlexibleSize are supported for widths")
		}
		restW -= widthsInPixels[i]
	}

	if denomW > 0 {
		origRestW := restW
		for i, width := range widths {
			if width.typ != sizeTypeFlexible {
				continue
			}
			w := int(float64(origRestW) * float64(width.value) / float64(denomW))
			widthsInPixels[i] = w
			restW -= w
		}
		// TODO: Use a better algorithm to distribute the rest.
		for restW > 0 {
			for i := len(widthsInPixels) - 1; i >= 0; i-- {
				if widths[i].typ != sizeTypeFlexible {
					continue
				}
				widthsInPixels[i]++
				restW--
				if restW <= 0 {
					break
				}
			}
			if restW <= 0 {
				break
			}
		}
	}
}

func (g *GridLayout) getHeightsInPixels(heightsInPixels []int, loopIndex int) {
	heights := g.Heights
	if len(heights) == 0 {
		heights = defaultHeights
	}

	// Calculate hights in pixels.
	// This is needed for each loop since the index starts with widgetBaseIdx for sizeTypeMaxContent.
	restH := g.Bounds.Dy()
	if restH < 0 {
		restH = 0
	}
	restH -= (len(heights) - 1) * g.RowGap
	var denomH int

	for j, height := range heights {
		switch height.typ {
		case sizeTypeFixed:
			heightsInPixels[j] = height.value
		case sizeTypeFlexible:
			heightsInPixels[j] = 0
			denomH += height.value
		case sizeTypeLazy:
			if height.lazy != nil {
				size := height.lazy(loopIndex*len(heights) + j)
				switch size.typ {
				case sizeTypeFixed:
					heightsInPixels[j] = size.value
				case sizeTypeFlexible:
					heightsInPixels[j] = 0
					denomH += size.value
				default:
					panic("layout: only FixedSize and FlexibleSize are supported for LazySize")
				}
			} else {
				heightsInPixels[j] = 0
			}
		}
		restH -= heightsInPixels[j]
	}

	if denomH > 0 {
		origRestH := restH
		for j, height := range heights {
			if height.typ != sizeTypeFlexible {
				continue
			}
			h := int(float64(origRestH) * float64(height.value) / float64(denomH))
			heightsInPixels[j] = h
			restH -= h
		}
		for restH > 0 {
			for j := len(heightsInPixels) - 1; j >= 0; j-- {
				if heights[j].typ != sizeTypeFlexible {
					continue
				}
				heightsInPixels[j]++
				restH--
				if restH <= 0 {
					break
				}
			}
			if restH <= 0 {
				break
			}
		}
	}
}
