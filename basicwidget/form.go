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

func (f *Form) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	f.calcItemBounds(context)

	for i, item := range f.items {
		if item.PrimaryWidget != nil {
			appender.AppendChildWidgetWithPosition(item.PrimaryWidget, f.primaryBounds[i].Min)
		}
		if item.SecondaryWidget != nil {
			appender.AppendChildWidgetWithPosition(item.SecondaryWidget, f.secondaryBounds[i].Min)
		}
	}

	return nil
}

func (f *Form) isItemOmitted(context *guigui.Context, item FormItem) bool {
	return (item.PrimaryWidget == nil || !context.IsVisible(item.PrimaryWidget)) &&
		(item.SecondaryWidget == nil || !context.IsVisible(item.SecondaryWidget))
}

func (f *Form) calcItemBounds(context *guigui.Context) {
	f.primaryBounds = slices.Delete(f.primaryBounds, 0, len(f.primaryBounds))
	f.secondaryBounds = slices.Delete(f.secondaryBounds, 0, len(f.secondaryBounds))

	paddingS := formItemPadding(context)

	var y int
	for i, item := range f.items {
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
		baseH := max(primaryS.Y, secondaryS.Y, minFormItemHeight(context))
		baseBounds := context.Bounds(f)
		baseBounds.Min.X += paddingS.X
		baseBounds.Max.X -= paddingS.X
		baseBounds.Min.Y += y
		baseBounds.Max.Y = baseBounds.Min.Y + baseH

		maxPaddingY := paddingS.Y + int((float64(UnitSize(context))-LineHeight(context))/2)
		if item.PrimaryWidget != nil {
			bounds := baseBounds
			bounds.Max.X = bounds.Min.X + primaryS.X
			pY := min((baseH+2*paddingS.Y-primaryS.Y)/2, maxPaddingY)
			bounds.Min.Y += pY
			bounds.Max.Y += pY
			f.primaryBounds[i] = bounds
		}
		if item.SecondaryWidget != nil {
			bounds := baseBounds
			bounds.Min.X = bounds.Max.X - secondaryS.X
			pY := min((baseH+2*paddingS.Y-secondaryS.Y)/2, maxPaddingY)
			bounds.Min.Y += pY
			bounds.Max.Y += pY
			f.secondaryBounds[i] = bounds
		}

		y += baseH + 2*paddingS.Y
	}
}

func (f *Form) Draw(context *guigui.Context, dst *ebiten.Image) {
	bgClr := draw.ScaleAlpha(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0), 1/32.0)
	borderClr := draw.ScaleAlpha(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0), 2/32.0)

	bounds := context.Bounds(f)
	bounds.Max.Y = bounds.Min.Y + f.DefaultSize(context).Y
	draw.DrawRoundedRect(context, dst, bounds, bgClr, RoundedCornerRadius(context))

	if len(f.items) > 0 {
		paddingS := formItemPadding(context)
		y := bounds.Min.Y
		for _, item := range f.items[:len(f.items)-1] {
			var primaryH int
			var secondaryH int
			if item.PrimaryWidget != nil {
				primaryH = context.ActualSize(item.PrimaryWidget).Y
			}
			if item.SecondaryWidget != nil {
				secondaryH = context.ActualSize(item.SecondaryWidget).Y
			}
			h := max(primaryH, secondaryH, minFormItemHeight(context))
			y += h + paddingS.Y

			x0 := float32(bounds.Min.X + paddingS.X)
			x1 := float32(bounds.Max.X - paddingS.X)
			yy := float32(y) + float32(paddingS.Y)
			width := 1 * float32(context.Scale())
			vector.StrokeLine(dst, x0, yy, x1, yy, width, borderClr, false)

			y += paddingS.Y
		}
	}

	draw.DrawRoundedRectBorder(context, dst, bounds, borderClr, borderClr, RoundedCornerRadius(context), 1*float32(context.Scale()), draw.RoundedRectBorderTypeRegular)
}

func (f *Form) DefaultSize(context *guigui.Context) image.Point {
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
		h := max(primaryS.Y, secondaryS.Y, minFormItemHeight(context))
		s.Y += h + 2*paddingS.Y
	}
	return s
}

func minFormItemHeight(context *guigui.Context) int {
	return UnitSize(context)
}
