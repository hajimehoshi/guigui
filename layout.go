// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package guigui

import (
	"image"
	"slices"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Layout interface {
	WidgetBounds(bounds image.Rectangle, widget Widget) image.Rectangle
}

type Size struct {
	typ   sizeType
	value int
}

type sizeType int

const (
	sizeTypeIntrinsic sizeType = iota
	sizeTypeFixed
	sizeTypeFlexible
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

type LayoutDirection int

const (
	LayoutDirectionHorizontal LayoutDirection = iota
	LayoutDirectionVertical
)

type Padding struct {
	Start  int
	Top    int
	End    int
	Bottom int
}

type LinearLayout struct {
	Direction LayoutDirection
	Items     []LinearLayoutItem
	Gap       int
	Padding   Padding

	tmpSizes []int
}

type linearLayoutItemCacheInfo struct {
	widgetIntrinsicSize int
	size                Size
}

func (l LinearLayout) WidgetBounds(bounds image.Rectangle, widget Widget) image.Rectangle {
	// Do a breadth-first search.
	for i, item := range l.Items {
		if item.Widget == widget {
			return theCachedLinearLayouts.get(&l, bounds, i)
		}
	}
	for i, item := range l.Items {
		b := theCachedLinearLayouts.get(&l, bounds, i)
		if r := item.Layout.WidgetBounds(b, widget); !r.Empty() {
			return r
		}
	}
	return image.Rectangle{}
}

func (l *LinearLayout) fixedSideSize(bounds image.Rectangle) int {
	switch l.Direction {
	case LayoutDirectionHorizontal:
		return bounds.Dy() - l.Padding.Top - l.Padding.Bottom
	case LayoutDirectionVertical:
		return bounds.Dx() - l.Padding.Start - l.Padding.End
	}
	return 0
}

func (l *LinearLayout) appendWidgetBounds(widgetBounds []image.Rectangle, bounds image.Rectangle) []image.Rectangle {
	if n := len(l.Items) - len(l.tmpSizes); n > 0 {
		l.tmpSizes = slices.Grow(l.tmpSizes, n)[:len(l.Items)]
	} else {
		l.tmpSizes = slices.Delete(l.tmpSizes, len(l.Items), len(l.tmpSizes))
	}
	l.sizesInPixels(bounds, l.tmpSizes)

	fixedSideSize := l.fixedSideSize(bounds)

	var progress int
	for i := range l.Items {
		var itemBounds image.Rectangle
		pt := bounds.Min.Add(image.Pt(l.Padding.Start, l.Padding.Top))
		s := l.tmpSizes[i]
		switch l.Direction {
		case LayoutDirectionHorizontal:
			pt.X += progress + i*l.Gap
			itemBounds = image.Rectangle{
				Min: pt,
				Max: pt.Add(image.Pt(s, fixedSideSize)),
			}
		case LayoutDirectionVertical:
			pt.Y += progress + i*l.Gap
			itemBounds = image.Rectangle{
				Min: pt,
				Max: pt.Add(image.Pt(fixedSideSize, s)),
			}
		}
		progress += s

		widgetBounds = append(widgetBounds, itemBounds)
	}

	return widgetBounds
}

func (l *LinearLayout) sizesInPixels(bounds image.Rectangle, sizesInPixels []int) {
	var rest int
	switch l.Direction {
	case LayoutDirectionHorizontal:
		rest = bounds.Dx() - l.Padding.Start - l.Padding.End
	case LayoutDirectionVertical:
		rest = bounds.Dy() - l.Padding.Top - l.Padding.Bottom
	}
	rest -= (len(l.Items) - 1) * l.Gap
	if rest < 0 {
		rest = 0
	}
	var denom int

	for i, item := range l.Items {
		switch item.Size.typ {
		case sizeTypeIntrinsic:
			switch l.Direction {
			case LayoutDirectionHorizontal:
				h := bounds.Dy() - l.Padding.Top - l.Padding.Bottom
				sizesInPixels[i] = item.Widget.Measure(nil, FixedHeightConstraints(h)).X
			case LayoutDirectionVertical:
				w := bounds.Dx() - l.Padding.Start - l.Padding.End
				sizesInPixels[i] = item.Widget.Measure(nil, FixedWidthConstraints(w)).Y
			}
		case sizeTypeFixed:
			sizesInPixels[i] = item.Size.value
		case sizeTypeFlexible:
			sizesInPixels[i] = 0
			denom += item.Size.value
		}
		rest -= sizesInPixels[i]
	}

	if denom > 0 {
		origRest := rest
		for i, item := range l.Items {
			if item.Size.typ != sizeTypeFlexible {
				continue
			}
			w := int(float64(origRest) * float64(item.Size.value) / float64(denom))
			sizesInPixels[i] = w
			rest -= w
		}
		// TODO: Use a better algorithm to distribute the rest.
		for rest > 0 {
			for i := len(sizesInPixels) - 1; i >= 0; i-- {
				if l.Items[i].Size.typ != sizeTypeFlexible {
					continue
				}
				sizesInPixels[i]++
				rest--
				if rest <= 0 {
					break
				}
			}
			if rest <= 0 {
				break
			}
		}
	}
}

type LinearLayoutItem struct {
	Widget Widget
	Size   Size
	Layout Layout
}

func (l *LinearLayoutItem) cacheInfo(direction LayoutDirection, fixedSideSize int) linearLayoutItemCacheInfo {
	info := linearLayoutItemCacheInfo{
		size: l.Size,
	}
	if l.Size.typ == sizeTypeIntrinsic {
		switch direction {
		case LayoutDirectionHorizontal:
			info.widgetIntrinsicSize = l.Widget.Measure(nil, FixedHeightConstraints(fixedSideSize)).X
		case LayoutDirectionVertical:
			info.widgetIntrinsicSize = l.Widget.Measure(nil, FixedWidthConstraints(fixedSideSize)).Y
		}
	}
	return info
}

type cachedLinearLayoutValues struct {
	itemBounds []image.Rectangle

	bounds    image.Rectangle
	direction LayoutDirection
	items     []linearLayoutItemCacheInfo
	gap       int
	padding   Padding

	atime int64
}

func (c *cachedLinearLayoutValues) matches(linearLayout *LinearLayout, bounds image.Rectangle) bool {
	if c.bounds != bounds {
		return false
	}
	if c.direction != linearLayout.Direction {
		return false
	}
	if len(c.items) != len(linearLayout.Items) {
		return false
	}
	fixedSideSize := linearLayout.fixedSideSize(bounds)
	for i, item := range linearLayout.Items {
		if c.items[i] != item.cacheInfo(linearLayout.Direction, fixedSideSize) {
			return false
		}
	}
	if c.gap != linearLayout.Gap {
		return false
	}
	if c.padding != linearLayout.Padding {
		return false
	}
	return true
}

type cachedLinearLayouts struct {
	values []*cachedLinearLayoutValues

	m sync.Mutex
}

var theCachedLinearLayouts cachedLinearLayouts

func (c *cachedLinearLayouts) get(linearLayout *LinearLayout, bounds image.Rectangle, index int) image.Rectangle {
	c.m.Lock()
	defer c.m.Unlock()

	for _, v := range c.values {
		if v.matches(linearLayout, bounds) {
			return v.itemBounds[index]
		}
	}

	// GC old results.
	now := ebiten.Tick()
	for i := len(c.values) - 1; i >= 0; i-- {
		if now-c.values[i].atime > int64(ebiten.TPS()) {
			c.values = slices.Delete(c.values, i, i+1)
		}
	}

	v := &cachedLinearLayoutValues{
		bounds:    bounds,
		direction: linearLayout.Direction,
		gap:       linearLayout.Gap,
		padding:   linearLayout.Padding,
		atime:     now,
	}

	if len(linearLayout.Items) > 0 {
		fixedSideSize := linearLayout.fixedSideSize(bounds)
		v.items = make([]linearLayoutItemCacheInfo, len(linearLayout.Items))
		for i, item := range linearLayout.Items {
			v.items[i] = item.cacheInfo(linearLayout.Direction, fixedSideSize)
		}

		v.itemBounds = make([]image.Rectangle, 0, len(linearLayout.Items))
		v.itemBounds = linearLayout.appendWidgetBounds(v.itemBounds, bounds)
	}
	c.values = append(c.values, v)

	return v.itemBounds[index]
}
