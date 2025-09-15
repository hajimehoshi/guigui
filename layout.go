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
	WidgetBounds(context *Context, bounds image.Rectangle, widget Widget) image.Rectangle
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
}

// linearLayoutItemCacheIdentity represents the identity of a cache.
// If and only if two linearLayoutItemCacheIdentity values are equal, the cache can be reused.
type linearLayoutItemCacheIdentity struct {
	widgetIntrinsicSize int

	// widgetState is needed to make widgetIndices valid.
	widgetState *widgetState

	size Size
}

func (l LinearLayout) WidgetBounds(context *Context, bounds image.Rectangle, widget Widget) image.Rectangle {
	if b, ok := theCachedLinearLayouts.widgetBounds(context, &l, bounds, widget); ok {
		return b
	}
	for i, item := range l.Items {
		if item.Layout == nil {
			continue
		}
		b := theCachedLinearLayouts.itemBounds(context, &l, bounds, i)
		if r := item.Layout.WidgetBounds(context, b, widget); !r.Empty() {
			return r
		}
	}
	return image.Rectangle{}
}

func (l *LinearLayout) alongSize(bounds image.Rectangle) int {
	switch l.Direction {
	case LayoutDirectionHorizontal:
		return bounds.Dx() - l.Padding.Start - l.Padding.End
	case LayoutDirectionVertical:
		return bounds.Dy() - l.Padding.Top - l.Padding.Bottom
	}
	return 0
}

func (l *LinearLayout) acrossSize(bounds image.Rectangle) int {
	switch l.Direction {
	case LayoutDirectionHorizontal:
		return bounds.Dy() - l.Padding.Top - l.Padding.Bottom
	case LayoutDirectionVertical:
		return bounds.Dx() - l.Padding.Start - l.Padding.End
	}
	return 0
}

type positionAndSize struct {
	position int
	size     int
}

func (l *LinearLayout) appendWidgetAlongPositionAndSizes(widgetAlongPositions []positionAndSize, context *Context, alongSize, acrossSize int) []positionAndSize {
	sizesInPixels := l.appendSizesInPixels(nil, context, alongSize, acrossSize)

	var progress int
	for i := range l.Items {
		widgetAlongPositions = append(widgetAlongPositions, positionAndSize{
			position: progress,
			size:     sizesInPixels[i],
		})
		progress += sizesInPixels[i] + l.Gap
	}

	return widgetAlongPositions
}

func (l *LinearLayout) appendSizesInPixels(sizesInPixels []int, context *Context, alongSize, acrossSize int) []int {
	rest := alongSize
	rest -= (len(l.Items) - 1) * l.Gap
	if rest < 0 {
		rest = 0
	}
	var denom int

	origLen := len(sizesInPixels)
	for i, item := range l.Items {
		switch item.Size.typ {
		case sizeTypeIntrinsic:
			switch l.Direction {
			case LayoutDirectionHorizontal:
				sizesInPixels = append(sizesInPixels, item.Widget.Measure(context, FixedHeightConstraints(acrossSize)).X)
			case LayoutDirectionVertical:
				sizesInPixels = append(sizesInPixels, item.Widget.Measure(context, FixedWidthConstraints(acrossSize)).Y)
			}
		case sizeTypeFixed:
			sizesInPixels = append(sizesInPixels, item.Size.value)
		case sizeTypeFlexible:
			sizesInPixels = append(sizesInPixels, 0)
			denom += item.Size.value
		}
		rest -= sizesInPixels[origLen+i]
	}

	if denom > 0 {
		origRest := rest
		for i, item := range l.Items {
			if item.Size.typ != sizeTypeFlexible {
				continue
			}
			w := int(float64(origRest) * float64(item.Size.value) / float64(denom))
			sizesInPixels[origLen+i] = w
			rest -= w
		}
		// TODO: Use a better algorithm to distribute the rest.
		for rest > 0 {
			for i := len(sizesInPixels) - origLen - 1; i >= 0; i-- {
				if l.Items[i].Size.typ != sizeTypeFlexible {
					continue
				}
				sizesInPixels[origLen+i]++
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

	return sizesInPixels
}

type LinearLayoutItem struct {
	Widget Widget
	Size   Size
	Layout Layout
}

func (l *LinearLayoutItem) cacheIdentity(direction LayoutDirection, acrossSize int) linearLayoutItemCacheIdentity {
	identity := linearLayoutItemCacheIdentity{
		size: l.Size,
	}
	if l.Widget != nil {
		identity.widgetState = l.Widget.widgetState()
	}
	if l.Size.typ == sizeTypeIntrinsic {
		switch direction {
		case LayoutDirectionHorizontal:
			identity.widgetIntrinsicSize = l.Widget.Measure(nil, FixedHeightConstraints(acrossSize)).X
		case LayoutDirectionVertical:
			identity.widgetIntrinsicSize = l.Widget.Measure(nil, FixedWidthConstraints(acrossSize)).Y
		}
	}
	return identity
}

type cachedLinearLayoutValues struct {
	itemAlongPositionAndSizes []positionAndSize
	widgetIndices             map[Widget]int

	direction  LayoutDirection
	alongSize  int
	acrossSize int
	items      []linearLayoutItemCacheIdentity
	gap        int

	atime int64
}

func (c *cachedLinearLayoutValues) matches(linearLayout *LinearLayout, alongSize, acrossSize int) bool {
	if c.alongSize != alongSize {
		return false
	}
	if c.acrossSize != acrossSize {
		return false
	}
	if c.direction != linearLayout.Direction {
		return false
	}
	if len(c.items) != len(linearLayout.Items) {
		return false
	}
	for i, item := range linearLayout.Items {
		if c.items[i] != item.cacheIdentity(linearLayout.Direction, acrossSize) {
			return false
		}
	}
	if c.gap != linearLayout.Gap {
		return false
	}
	return true
}

type cachedLinearLayouts struct {
	values []*cachedLinearLayoutValues

	m sync.Mutex
}

var theCachedLinearLayouts cachedLinearLayouts

func (c *cachedLinearLayouts) itemBounds(context *Context, linearLayout *LinearLayout, bounds image.Rectangle, index int) image.Rectangle {
	c.m.Lock()
	defer c.m.Unlock()

	v := c.get(context, linearLayout, bounds)
	ps := v.itemAlongPositionAndSizes[index]
	return positionAndSizeToBounds(linearLayout, bounds, ps)
}

func (c *cachedLinearLayouts) widgetBounds(context *Context, linearLayout *LinearLayout, bounds image.Rectangle, widget Widget) (image.Rectangle, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	v := c.get(context, linearLayout, bounds)
	idx, ok := v.widgetIndices[widget]
	if !ok {
		return image.Rectangle{}, false
	}
	ps := v.itemAlongPositionAndSizes[idx]
	return positionAndSizeToBounds(linearLayout, bounds, ps), true
}

func positionAndSizeToBounds(linearLayout *LinearLayout, bounds image.Rectangle, ps positionAndSize) image.Rectangle {
	pt := bounds.Min.Add(image.Pt(linearLayout.Padding.Start, linearLayout.Padding.Top))
	acrossSize := linearLayout.acrossSize(bounds)
	switch linearLayout.Direction {
	case LayoutDirectionHorizontal:
		pt.X += ps.position
		return image.Rectangle{
			Min: pt,
			Max: pt.Add(image.Pt(ps.size, acrossSize)),
		}
	case LayoutDirectionVertical:
		pt.Y += ps.position
		return image.Rectangle{
			Min: pt,
			Max: pt.Add(image.Pt(acrossSize, ps.size)),
		}
	}
	return image.Rectangle{}
}

func (c *cachedLinearLayouts) get(context *Context, linearLayout *LinearLayout, bounds image.Rectangle) *cachedLinearLayoutValues {
	alongSize := linearLayout.alongSize(bounds)
	acrossSize := linearLayout.acrossSize(bounds)

	for _, v := range c.values {
		if !v.matches(linearLayout, alongSize, acrossSize) {
			continue
		}
		v.atime = ebiten.Tick()
		return v
	}

	// GC old results.
	now := ebiten.Tick()
	for i := len(c.values) - 1; i >= 0; i-- {
		if now-c.values[i].atime > int64(ebiten.TPS()) {
			c.values = slices.Delete(c.values, i, i+1)
		}
	}

	v := &cachedLinearLayoutValues{
		alongSize:  alongSize,
		acrossSize: acrossSize,
		direction:  linearLayout.Direction,
		gap:        linearLayout.Gap,
		atime:      now,
	}

	if len(linearLayout.Items) > 0 {
		v.items = make([]linearLayoutItemCacheIdentity, len(linearLayout.Items))
		for i, item := range linearLayout.Items {
			v.items[i] = item.cacheIdentity(linearLayout.Direction, alongSize)
			if item.Widget != nil {
				if v.widgetIndices == nil {
					v.widgetIndices = map[Widget]int{}
				}
				v.widgetIndices[item.Widget] = i
			}
		}

		v.itemAlongPositionAndSizes = make([]positionAndSize, 0, len(linearLayout.Items))
		v.itemAlongPositionAndSizes = linearLayout.appendWidgetAlongPositionAndSizes(v.itemAlongPositionAndSizes, context, alongSize, acrossSize)
	}
	c.values = append(c.values, v)

	return v
}
