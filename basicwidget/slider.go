// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"math"
	"math/big"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"

	"github.com/hajimehoshi/guigui"
)

type Slider struct {
	guigui.DefaultWidget

	abstractNumberInput abstractNumberInput

	dragging           bool
	draggingStartValue big.Int
	draggingStartX     int

	prevThumbHovered bool

	onValueChangedBigInt func(value *big.Int)
}

func (s *Slider) SetOnValueChangedBigInt(f func(value *big.Int)) {
	s.onValueChangedBigInt = f
}

func (s *Slider) SetOnValueChangedInt64(f func(value int64)) {
	if f == nil {
		s.abstractNumberInput.SetOnValueChangedInt64(nil)
		return
	}
	s.abstractNumberInput.SetOnValueChangedInt64(func(value int64, committed bool) {
		f(value)
	})
}

func (s *Slider) SetOnValueChangedUint64(f func(value uint64)) {
	if f == nil {
		s.abstractNumberInput.SetOnValueChangedUint64(nil)
		return
	}
	s.abstractNumberInput.SetOnValueChangedUint64(func(value uint64, committed bool) {
		f(value)
	})
}

func (s *Slider) ValueBigInt() *big.Int {
	return s.abstractNumberInput.ValueBigInt()
}

func (s *Slider) ValueInt64() int64 {
	return s.abstractNumberInput.ValueInt64()
}

func (s *Slider) ValueUint64() uint64 {
	return s.abstractNumberInput.ValueUint64()
}

func (s *Slider) SetValueBigInt(value *big.Int) {
	s.abstractNumberInput.SetValueBigInt(value, true)
}

func (s *Slider) SetValueInt64(value int64) {
	s.abstractNumberInput.SetValueInt64(value, true)
}

func (s *Slider) SetValueUint64(value uint64) {
	s.abstractNumberInput.SetValueUint64(value, true)
}

func (s *Slider) MinimumValueBigInt() *big.Int {
	return s.abstractNumberInput.MinimumValueBigInt()
}

func (s *Slider) SetMinimumValueBigInt(minimum *big.Int) {
	s.abstractNumberInput.SetMinimumValueBigInt(minimum)
}

func (s *Slider) SetMinimumValueInt64(minimum int64) {
	s.abstractNumberInput.SetMinimumValueInt64(minimum)
}

func (s *Slider) SetMinimumValueUint64(minimum uint64) {
	s.abstractNumberInput.SetMinimumValueUint64(minimum)
}

func (s *Slider) MaximumValueBigInt() *big.Int {
	return s.abstractNumberInput.MaximumValueBigInt()
}

func (s *Slider) SetMaximumValueBigInt(maximum *big.Int) {
	s.abstractNumberInput.SetMaximumValueBigInt(maximum)
}

func (s *Slider) SetMaximumValueInt64(maximum int64) {
	s.abstractNumberInput.SetMaximumValueInt64(maximum)
}

func (s *Slider) SetMaximumValueUint64(maximum uint64) {
	s.abstractNumberInput.SetMaximumValueUint64(maximum)
}

func (s *Slider) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if hovered := s.isThumbHovered(context); s.prevThumbHovered != hovered {
		s.prevThumbHovered = hovered
		guigui.RequestRedraw(s)
	}
	return nil
}

func (s *Slider) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	s.abstractNumberInput.SetOnValueChangedBigInt(func(value *big.Int, committed bool) {
		if s.onValueChangedBigInt != nil {
			s.onValueChangedBigInt(value)
		}
		guigui.RequestRedraw(s)
	})

	max := s.abstractNumberInput.MaximumValueBigInt()
	min := s.abstractNumberInput.MinimumValueBigInt()
	if max == nil || min == nil {
		return guigui.HandleInputResult{}
	}

	if context.IsEnabled(s) && context.IsWidgetHitAtCursor(s) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !s.dragging {
		context.SetFocused(s, true)
		if !s.isThumbHovered(context) {
			s.setValueFromCursor(context)
		}
		s.dragging = true
		x, _ := ebiten.CursorPosition()
		s.draggingStartX = x
		s.draggingStartValue.Set(s.abstractNumberInput.ValueBigInt())
		guigui.RequestRedraw(s)
		return guigui.HandleInputByWidget(s)
	}

	if !context.IsEnabled(s) || !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if s.dragging {
			guigui.RequestRedraw(s)
		}
		s.dragging = false
		s.draggingStartX = 0
		s.draggingStartValue = big.Int{}
		return guigui.HandleInputResult{}
	}

	if context.IsEnabled(s) && s.dragging && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.setValueFromCursorDelta(context)
		return guigui.HandleInputByWidget(s)
	}

	return guigui.HandleInputResult{}
}

func (s *Slider) setValueFromCursorDelta(context *guigui.Context) {
	s.setValue(context, &s.draggingStartValue, s.draggingStartX)
}

func (s *Slider) setValueFromCursor(context *guigui.Context) {
	min := s.abstractNumberInput.MinimumValueBigInt()
	if min == nil {
		return
	}

	b := context.Bounds(s)
	minX := b.Min.X + (b.Dx()-s.barWidth(context))/2
	s.setValue(context, min, minX)
}

func (s *Slider) setValue(context *guigui.Context, originValue *big.Int, originX int) {
	max := s.abstractNumberInput.MaximumValueBigInt()
	min := s.abstractNumberInput.MinimumValueBigInt()
	if max == nil || min == nil {
		return
	}

	c := image.Pt(ebiten.CursorPosition())
	var v big.Int
	v.Sub(max, min)
	v.Mul(&v, (&big.Int{}).SetInt64(int64(c.X-originX)))
	v.Div(&v, (&big.Int{}).SetInt64(int64(s.barWidth(context))))
	v.Add(&v, originValue)
	s.abstractNumberInput.SetValueBigInt(&v, true)
}

func (s *Slider) barWidth(context *guigui.Context) int {
	w := context.Bounds(s).Dx()
	return w - 2*sliderThumbRadius(context)
}

func sliderThumbRadius(context *guigui.Context) int {
	return int(UnitSize(context) * 7 / 16)
}

func (s *Slider) thumbBounds(context *guigui.Context) image.Rectangle {
	rate := s.abstractNumberInput.Rate()
	if math.IsNaN(rate) {
		return image.Rectangle{}
	}
	bounds := context.Bounds(s)
	x := bounds.Min.X + int(rate*float64(s.barWidth(context)))
	y := bounds.Min.Y + (bounds.Dy()-2*sliderThumbRadius(context))/2
	w := 2 * sliderThumbRadius(context)
	h := 2 * sliderThumbRadius(context)
	return image.Rect(x, y, x+w, y+h)
}

func (s *Slider) CursorShape(context *guigui.Context) (ebiten.CursorShapeType, bool) {
	if s.canPress(context) || s.dragging {
		return ebiten.CursorShapePointer, true
	}
	return 0, true
}

func (s *Slider) Draw(context *guigui.Context, dst *ebiten.Image) {
	rate := s.abstractNumberInput.Rate()

	b := context.Bounds(s)
	x0 := b.Min.X + sliderThumbRadius(context)
	x1 := x0
	if !math.IsNaN(rate) {
		x1 += int(float64(s.barWidth(context)) * float64(rate))
	}
	x2 := b.Max.X - sliderThumbRadius(context)
	strokeWidth := int(5 * context.Scale())
	r := strokeWidth / 2
	y0 := (b.Min.Y+b.Max.Y)/2 - r
	y1 := (b.Min.Y+b.Max.Y)/2 + r

	bgColorOn := draw.Color(context.ColorMode(), draw.ColorTypeAccent, 0.5)
	bgColorOff := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.8)
	if !context.IsEnabled(s) {
		bgColorOn = bgColorOff
	}

	if x0 < x1 {
		b := image.Rect(x0, y0, x1, y1)
		draw.DrawRoundedRect(context, dst, b, bgColorOn, r)

		if !context.IsEnabled(s) {
			borderClr1, borderClr2 := draw.BorderColors(context.ColorMode(), draw.RoundedRectBorderTypeInset, false)
			draw.DrawRoundedRectBorder(context, dst, b, borderClr1, borderClr2, r, float32(1*context.Scale()), draw.RoundedRectBorderTypeInset)
		}
	}

	if x1 < x2 {
		b := image.Rect(x1, y0, x2, y1)
		draw.DrawRoundedRect(context, dst, b, bgColorOff, r)

		borderClr1, borderClr2 := draw.BorderColors(context.ColorMode(), draw.RoundedRectBorderTypeInset, false)
		draw.DrawRoundedRectBorder(context, dst, b, borderClr1, borderClr2, r, float32(1*context.Scale()), draw.RoundedRectBorderTypeInset)
	}

	if thumbBounds := s.thumbBounds(context); !thumbBounds.Empty() {
		cm := context.ColorMode()
		thumbColor := draw.ThumbColor(context.ColorMode(), context.IsEnabled(s))
		if s.isActive(context) {
			thumbColor = draw.Color2(cm, draw.ColorTypeBase, 0.95, 0.55)
		} else if s.canPress(context) {
			thumbColor = draw.Color2(cm, draw.ColorTypeBase, 0.975, 0.575)
		}
		thumbClr1, thumbClr2 := draw.BorderColors(context.ColorMode(), draw.RoundedRectBorderTypeOutset, false)
		r := thumbBounds.Dy() / 2
		draw.DrawRoundedRect(context, dst, thumbBounds, thumbColor, r)
		draw.DrawRoundedRectBorder(context, dst, thumbBounds, thumbClr1, thumbClr2, r, float32(1*context.Scale()), draw.RoundedRectBorderTypeOutset)
	}
}

func (s *Slider) canPress(context *guigui.Context) bool {
	return context.IsEnabled(s) && s.isThumbHovered(context) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !s.dragging
}

func (s *Slider) isThumbHovered(context *guigui.Context) bool {
	return context.IsWidgetHitAtCursor(s) && image.Pt(ebiten.CursorPosition()).In(s.thumbBounds(context))
}

func (s *Slider) isActive(context *guigui.Context) bool {
	return context.IsEnabled(s) && s.isThumbHovered(context) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && s.dragging
}

func (s *Slider) DefaultSize(context *guigui.Context) image.Point {
	return image.Pt(6*UnitSize(context), UnitSize(context))
}
