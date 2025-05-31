// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package layoututil

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

func WidthsInPixels(widthsInPixels []int, widths []Size, totalWidth int, columnGap int) {
	if len(widths) == 0 {
		widths = defaultWidths
	}

	// Calculate widths in pixels.
	restW := totalWidth
	restW -= (len(widths) - 1) * columnGap
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
			panic("layout: only FixedSize and FlexibleSize are supported for getWidthsInPixels")
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

func HeightsInPixels(heightsInPixels []int, heights []Size, totalHeight int, rowGap int, loopIndex int) {
	if len(heights) == 0 {
		heights = defaultHeights
	}

	// Calculate hights in pixels.
	// This is needed for each loop since the index starts with widgetBaseIdx for sizeTypeMaxContent.
	restH := totalHeight
	if restH < 0 {
		restH = 0
	}
	restH -= (len(heights) - 1) * rowGap
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
