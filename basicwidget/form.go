// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"image"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type FormItem struct {
	PrimaryWidget   guigui.Widget
	SecondaryWidget guigui.Widget
}

type Form struct {
	guigui.DefaultWidget

	items []FormItem

	itemBounds      []image.Rectangle
	primaryBounds   []image.Rectangle
	secondaryBounds []image.Rectangle
}

func formItemPadding(context *guigui.Context) image.Point {
	return image.Pt(UnitSize(context)/2, UnitSize(context)/4)
}

func (f *Form) SetItems(items []FormItem) {
	f.items = slices.Delete(f.items, 0, len(f.items))
	f.items = append(f.items, items...)
}

func (f *Form) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	for _, item := range f.items {
		if item.PrimaryWidget != nil {
			appender.AppendChildWidget(item.PrimaryWidget)
		}
		if item.SecondaryWidget != nil {
			appender.AppendChildWidget(item.SecondaryWidget)
		}
	}
}

func (f *Form) Build(context *guigui.Context) error {
	f.calcItemBounds(context, context.ActualSize(f).X)

	for i, item := range f.items {
		if item.PrimaryWidget != nil {
			context.SetPosition(item.PrimaryWidget, f.primaryBounds[i].Min)
		}
		if item.SecondaryWidget != nil {
			context.SetPosition(item.SecondaryWidget, f.secondaryBounds[i].Min)
		}
	}

	return nil
}

func (f *Form) isItemOmitted(context *guigui.Context, item FormItem) bool {
	return (item.PrimaryWidget == nil || !context.IsVisible(item.PrimaryWidget)) &&
		(item.SecondaryWidget == nil || !context.IsVisible(item.SecondaryWidget))
}

func (f *Form) calcItemBounds(context *guigui.Context, width int) {
	// TODO: Cache the result?

	f.itemBounds = slices.Delete(f.itemBounds, 0, len(f.itemBounds))
	f.primaryBounds = slices.Delete(f.primaryBounds, 0, len(f.primaryBounds))
	f.secondaryBounds = slices.Delete(f.secondaryBounds, 0, len(f.secondaryBounds))

	bounds := context.Bounds(f)
	paddingS := formItemPadding(context)

	var y int
	for i, item := range f.items {
		f.itemBounds = append(f.itemBounds, image.Rectangle{})
		f.primaryBounds = append(f.primaryBounds, image.Rectangle{})
		f.secondaryBounds = append(f.secondaryBounds, image.Rectangle{})

		if f.isItemOmitted(context, item) {
			continue
		}

		var primaryS image.Point
		var secondaryS image.Point
		if item.PrimaryWidget != nil {
			primaryS = context.ActualSize(item.PrimaryWidget)
		}
		if item.SecondaryWidget != nil {
			secondaryS = context.ActualSize(item.SecondaryWidget)
		}
		newLine := item.PrimaryWidget != nil && primaryS.X+secondaryS.X+2*paddingS.X > bounds.Dx()
		var baseH int
		if newLine {
			baseH = max(primaryS.Y+paddingS.Y+secondaryS.Y, minFormItemHeight(context)) + 2*paddingS.Y
		} else {
			baseH = max(primaryS.Y, secondaryS.Y, minFormItemHeight(context)) + 2*paddingS.Y
		}
		p := bounds.Min
		f.itemBounds[i] = image.Rectangle{
			Min: p.Add(image.Pt(0, y)),
			Max: p.Add(image.Pt(width, y+baseH)),
		}

		maxPaddingY := paddingS.Y + int((float64(UnitSize(context))-LineHeight(context))/2)
		if item.PrimaryWidget != nil {
			bounds := f.itemBounds[i]
			bounds.Min.X += paddingS.X
			bounds.Max.X = bounds.Min.X + primaryS.X
			pY := min((baseH-primaryS.Y)/2, maxPaddingY)
			bounds.Min.Y += pY
			bounds.Max.Y += pY
			f.primaryBounds[i] = bounds
		}
		if item.SecondaryWidget != nil {
			bounds := f.itemBounds[i]
			bounds.Min.X = bounds.Max.X - secondaryS.X - paddingS.X
			bounds.Max.X -= paddingS.X
			if newLine {
				bounds.Min.Y += paddingS.Y + primaryS.Y + paddingS.Y
				bounds.Max.Y += paddingS.Y + primaryS.Y + paddingS.Y
			} else {
				pY := min((baseH-secondaryS.Y)/2, maxPaddingY)
				bounds.Min.Y += pY
				bounds.Max.Y += pY
			}
			f.secondaryBounds[i] = bounds
		}

		y += baseH
	}
}

func (f *Form) Draw(context *guigui.Context, dst *ebiten.Image) {
	bgClr := draw.ScaleAlpha(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0), 1/32.0)
	borderClr := draw.ScaleAlpha(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0), 2/32.0)

	bounds := context.Bounds(f)
	draw.DrawRoundedRect(context, dst, bounds, bgClr, RoundedCornerRadius(context))

	// Render borders between items.
	if len(f.items) > 0 {
		paddingS := formItemPadding(context)
		for i := range f.items[:len(f.items)-1] {
			x0 := float32(bounds.Min.X + paddingS.X)
			x1 := float32(bounds.Max.X - paddingS.X)
			y := float32(f.itemBounds[i].Max.Y)
			width := 1 * float32(context.Scale())
			vector.StrokeLine(dst, x0, y, x1, y, width, borderClr, false)
		}
	}

	draw.DrawRoundedRectBorder(context, dst, bounds, borderClr, borderClr, RoundedCornerRadius(context), 1*float32(context.Scale()), draw.RoundedRectBorderTypeRegular)
}

func (f *Form) measureWithoutConstraints(context *guigui.Context) image.Point {
	// Measure without size constraints should return the default size rather than an actual size.
	// Do not use itemBounds, primaryBounds, or secondaryBounds here.

	paddingS := formItemPadding(context)
	gapX := UnitSize(context)

	var s image.Point
	for _, item := range f.items {
		if f.isItemOmitted(context, item) {
			continue
		}
		var primaryS image.Point
		var secondaryS image.Point
		if item.PrimaryWidget != nil {
			primaryS = context.ActualSize(item.PrimaryWidget)
		}
		if item.SecondaryWidget != nil {
			secondaryS = context.ActualSize(item.SecondaryWidget)
		}

		s.X = max(s.X, primaryS.X+secondaryS.X+2*paddingS.X+gapX)
		h := max(primaryS.Y, secondaryS.Y, minFormItemHeight(context)) + 2*paddingS.Y
		s.Y += h
	}
	return s
}

func (f *Form) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	if len(f.itemBounds) == 0 {
		return image.Point{}
	}

	if s := f.measureWithoutConstraints(context); s.X <= constraints.MaxSize().X {
		return image.Pt(max(s.X, constraints.MinSize().X), s.Y)
	}

	f.calcItemBounds(context, constraints.MaxSize().X)
	return f.itemBounds[len(f.itemBounds)-1].Max.Sub(f.itemBounds[0].Min)
}

func minFormItemHeight(context *guigui.Context) int {
	return UnitSize(context)
}
