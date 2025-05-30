// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"math"
	"math/big"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

var (
	minInt64  big.Int
	maxInt64  big.Int
	maxUint64 big.Int
)

func init() {
	minInt64.SetInt64(math.MinInt64)
	maxInt64.SetInt64(math.MaxInt64)
	maxUint64.SetUint64(math.MaxUint64)
}

type NumberInput struct {
	guigui.DefaultWidget

	textInput  TextInput
	upButton   Button
	downButton Button

	abstractNumberInput abstractNumberInput
	nextValue           *big.Int
}

func (n *NumberInput) IsEditable() bool {
	return n.textInput.IsEditable()
}

func (n *NumberInput) SetEditable(editable bool) {
	n.textInput.SetEditable(editable)
}

func (n *NumberInput) SetOnValueChangedBigInt(f func(value *big.Int, committed bool)) {
	n.abstractNumberInput.SetOnValueChangedBigInt(f)
}

func (n *NumberInput) SetOnValueChangedInt64(f func(value int64, committed bool)) {
	n.abstractNumberInput.SetOnValueChangedInt64(f)
}

func (n *NumberInput) SetOnValueChangedUint64(f func(value uint64, committed bool)) {
	n.abstractNumberInput.SetOnValueChangedUint64(f)
}

func (n *NumberInput) ValueBigInt() *big.Int {
	if n.nextValue != nil {
		return n.nextValue
	}
	return n.abstractNumberInput.ValueBigInt()
}

func (n *NumberInput) ValueInt64() int64 {
	if n.nextValue != nil {
		return n.nextValue.Int64()
	}
	return n.abstractNumberInput.ValueInt64()
}

func (n *NumberInput) ValueUint64() uint64 {
	if n.nextValue != nil {
		return n.nextValue.Uint64()
	}
	return n.abstractNumberInput.ValueUint64()
}

func (n *NumberInput) SetValueBigInt(value *big.Int) {
	if n.nextValue != nil && n.nextValue.Cmp(value) == 0 {
		return
	}
	if n.nextValue == nil {
		n.nextValue = &big.Int{}
	}
	n.nextValue.Set(value)
}

func (n *NumberInput) SetValueInt64(value int64) {
	n.SetValueBigInt((&big.Int{}).SetInt64(value))
}

func (n *NumberInput) SetValueUint64(value uint64) {
	n.SetValueBigInt((&big.Int{}).SetUint64(value))
}

func (n *NumberInput) MinimumValueBigInt() *big.Int {
	return n.abstractNumberInput.MinimumValueBigInt()
}

func (n *NumberInput) SetMinimumValueBigInt(minimum *big.Int) {
	n.abstractNumberInput.SetMinimumValueBigInt(minimum)
}

func (n *NumberInput) SetMinimumValueInt64(minimum int64) {
	n.abstractNumberInput.SetMinimumValueInt64(minimum)
}

func (n *NumberInput) SetMinimumValueUint64(minimum uint64) {
	n.abstractNumberInput.SetMinimumValueUint64(minimum)
}

func (n *NumberInput) MaximumValueBigInt() *big.Int {
	return n.abstractNumberInput.MaximumValueBigInt()
}

func (n *NumberInput) SetMaximumValueBigInt(maximum *big.Int) {
	n.abstractNumberInput.SetMaximumValueBigInt(maximum)
}

func (n *NumberInput) SetMaximumValueInt64(maximum int64) {
	n.abstractNumberInput.SetMaximumValueInt64(maximum)
}

func (n *NumberInput) SetMaximumValueUint64(maximum uint64) {
	n.abstractNumberInput.SetMaximumValueUint64(maximum)
}

func (n *NumberInput) SetStepBigInt(step *big.Int) {
	n.abstractNumberInput.SetStepBigInt(step)
}

func (n *NumberInput) SetStepInt64(step int64) {
	n.abstractNumberInput.SetStepInt64(step)
}

func (n *NumberInput) SetStepUint64(step uint64) {
	n.abstractNumberInput.SetStepUint64(step)
}

func (n *NumberInput) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if n.nextValue != nil && !n.textInput.isFocused(context) && !context.IsFocused(n) {
		n.abstractNumberInput.SetValueBigInt(n.nextValue, true)
		n.nextValue = nil
	}

	n.abstractNumberInput.SetOnValueChangedString(func(text string, force bool) {
		if force {
			n.textInput.ForceSetValue(text)
		} else {
			n.textInput.SetValue(text)
		}
		n.nextValue = nil
	})

	n.textInput.SetValue(n.abstractNumberInput.ValueString())
	n.textInput.SetHorizontalAlign(HorizontalAlignEnd)
	n.textInput.SetTabular(true)
	n.textInput.setPaddingEnd(UnitSize(context) / 2)
	n.textInput.SetOnValueChanged(func(text string, committed bool) {
		n.abstractNumberInput.SetString(text, committed)
		if committed {
			n.nextValue = nil
		}
	})
	appender.AppendChildWidgetWithBounds(&n.textInput, context.Bounds(n))

	imgUp, err := theResourceImages.Get("keyboard_arrow_up", context.ColorMode())
	if err != nil {
		return err
	}
	imgDown, err := theResourceImages.Get("keyboard_arrow_down", context.ColorMode())
	if err != nil {
		return err
	}

	n.upButton.SetIcon(imgUp)
	n.upButton.setSharpenCorners(draw.SharpenCorners{
		LowerStart: true,
		LowerEnd:   true,
	})
	n.upButton.setPairedButton(&n.downButton)
	n.upButton.setOnRepeat(func() {
		n.increment()
	})
	context.SetEnabled(&n.upButton, n.IsEditable() && n.abstractNumberInput.CanIncrement())

	b := context.Bounds(n)
	appender.AppendChildWidgetWithBounds(&n.upButton, image.Rectangle{
		Min: image.Point{
			X: b.Max.X - UnitSize(context)*3/4,
			Y: b.Min.Y,
		},
		Max: image.Point{
			X: b.Max.X,
			Y: b.Min.Y + b.Dy()/2,
		},
	})

	n.downButton.SetIcon(imgDown)
	n.downButton.setSharpenCorners(draw.SharpenCorners{
		UpperStart: true,
		UpperEnd:   true,
	})
	n.downButton.setPairedButton(&n.upButton)
	n.downButton.setOnRepeat(func() {
		n.decrement()
	})
	context.SetEnabled(&n.downButton, n.IsEditable() && n.abstractNumberInput.CanDecrement())

	appender.AppendChildWidgetWithBounds(&n.downButton, image.Rectangle{
		Min: image.Point{
			X: b.Max.X - UnitSize(context)*3/4,
			Y: b.Min.Y + b.Dy()/2,
		},
		Max: b.Max,
	})

	return nil
}

func (n *NumberInput) HandleButtonInput(context *guigui.Context) guigui.HandleInputResult {
	if isKeyRepeating(ebiten.KeyUp) {
		n.increment()
		return guigui.HandleInputByWidget(n)
	}
	if isKeyRepeating(ebiten.KeyDown) {
		n.decrement()
		return guigui.HandleInputByWidget(n)
	}
	return guigui.HandleInputResult{}
}

func (n *NumberInput) DefaultSize(context *guigui.Context) image.Point {
	return n.textInput.DefaultSize(context)
}

func (n *NumberInput) increment() {
	if !n.IsEditable() {
		return
	}
	n.textInput.CommitWithCurrentInputValue()
	n.abstractNumberInput.SetString(n.textInput.Value(), true)
	n.abstractNumberInput.Increment()
	n.nextValue = nil
}

func (n *NumberInput) decrement() {
	if !n.IsEditable() {
		return
	}
	n.textInput.CommitWithCurrentInputValue()
	n.abstractNumberInput.SetString(n.textInput.Value(), true)
	n.abstractNumberInput.Decrement()
	n.nextValue = nil
}
