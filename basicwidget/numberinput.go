// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package basicwidget

import (
	"image"
	"math"
	"math/big"
	"strings"

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
	upButton   TextButton
	downButton TextButton

	value   big.Int
	min     big.Int
	minSet  bool
	max     big.Int
	maxSet  bool
	step    big.Int
	stepSet bool

	onValueChangedBigInt func(value *big.Int)
	onValueChangedInt64  func(value int64)
	onValueChangedUint64 func(value uint64)
}

func (n *NumberInput) IsEditable() bool {
	return n.textInput.IsEditable()
}

func (n *NumberInput) SetEditable(editable bool) {
	n.textInput.SetEditable(editable)
}

func (n *NumberInput) SetOnValueChangedBigInt(f func(value *big.Int)) {
	n.onValueChangedBigInt = f
}

func (n *NumberInput) SetOnValueChangedInt64(f func(value int64)) {
	n.onValueChangedInt64 = f
}

func (n *NumberInput) SetOnValueChangedUint64(f func(value uint64)) {
	n.onValueChangedUint64 = f
}

func (n *NumberInput) ValueBigInt() *big.Int {
	var v big.Int
	v.Set(&n.value)
	return &v
}

func (n *NumberInput) ValueInt64() int64 {
	if n.value.IsInt64() {
		return n.value.Int64()
	} else if n.value.Cmp(&maxInt64) > 0 {
		return math.MaxInt64
	} else if n.value.Cmp(&minInt64) < 0 {
		return math.MinInt64
	}
	return 0
}

func (n *NumberInput) ValueUint64() uint64 {
	if n.value.IsUint64() {
		return n.value.Uint64()
	} else if n.value.Cmp(&maxUint64) > 0 {
		return math.MaxUint64
	} else if n.value.Cmp(big.NewInt(0)) < 0 {
		return 0
	}
	return 0
}

func (n *NumberInput) SetValueBigInt(value *big.Int) {
	n.setValue(value, false)
}

func (n *NumberInput) SetValueInt64(value int64) {
	var v big.Int
	v.SetInt64(value)
	n.setValue(&v, false)
}

func (n *NumberInput) SetValueUint64(value uint64) {
	var v big.Int
	v.SetUint64(value)
	n.setValue(&v, false)
}

func (n *NumberInput) setValue(value *big.Int, force bool) {
	n.clamp(value)
	if n.value.Cmp(value) == 0 {
		return
	}
	n.value.Set(value)
	if force {
		n.textInput.ForceSetValue(n.value.String())
	} else {
		n.textInput.SetValue(n.value.String())
	}
	n.fireValueChangeEvents()
}

func (n *NumberInput) MinimumValueBigInt() *big.Int {
	if !n.minSet {
		return nil
	}
	var v big.Int
	v.Set(&n.min)
	return &v
}

func (n *NumberInput) SetMinimumValueBigInt(minimum *big.Int) {
	if minimum == nil {
		n.min = big.Int{}
		n.minSet = false
		return
	}
	n.min.Set(minimum)
	n.minSet = true
	var v big.Int
	v.Set(&n.value)
	n.SetValueBigInt(&v)
}

func (n *NumberInput) SetMinimumValueInt64(minimum int64) {
	n.min.SetInt64(minimum)
	n.minSet = true
	var v big.Int
	v.Set(&n.value)
	n.SetValueBigInt(&v)
}

func (n *NumberInput) SetMinimumValueUint64(minimum uint64) {
	n.min.SetUint64(minimum)
	n.minSet = true
	var v big.Int
	v.Set(&n.value)
	n.SetValueBigInt(&v)
}

func (n *NumberInput) MaximumValueBigInt() *big.Int {
	if !n.maxSet {
		return nil
	}
	var v big.Int
	v.Set(&n.max)
	return &v
}

func (n *NumberInput) SetMaximumValueBigInt(maximum *big.Int) {
	if maximum == nil {
		n.max = big.Int{}
		n.maxSet = false
		return
	}
	n.max.Set(maximum)
	n.maxSet = true
	var v big.Int
	v.Set(&n.value)
	n.SetValueBigInt(&v)
}

func (n *NumberInput) SetMaximumValueInt64(maximum int64) {
	n.max.SetInt64(maximum)
	n.maxSet = true
	var v big.Int
	v.Set(&n.value)
	n.SetValueBigInt(&v)
}

func (n *NumberInput) SetMaximumValueUint64(maximum uint64) {
	n.max.SetUint64(maximum)
	n.maxSet = true
	var v big.Int
	v.Set(&n.value)
	n.SetValueBigInt(&v)
}

func (n *NumberInput) SetStepBigInt(step *big.Int) {
	if step == nil {
		n.step = big.Int{}
		n.stepSet = false
		return
	}
	n.step.Set(step)
	n.stepSet = true
}

func (n *NumberInput) SetStepInt64(step int64) {
	n.step.SetInt64(step)
	n.stepSet = true
}

func (n *NumberInput) SetStepUint64(step uint64) {
	n.step.SetUint64(step)
	n.stepSet = true
}

func (n *NumberInput) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	n.textInput.SetHorizontalAlign(HorizontalAlignEnd)
	n.textInput.SetNumber(true)
	n.textInput.setPaddingRight(UnitSize(context) / 2)
	n.textInput.SetOnValueChanged(func(text string, committed bool) {
		if !committed {
			return
		}
		n.commit(text)
	})
	appender.AppendChildWidgetWithBounds(&n.textInput, context.Bounds(n))
	// HasFocusedChildWidget works after appending the child widget.
	if !context.HasFocusedChildWidget(n) {
		n.textInput.SetValue(n.value.String())
	}

	imgUp, err := theResourceImages.Get("keyboard_arrow_up", context.ColorMode())
	if err != nil {
		return err
	}
	imgDown, err := theResourceImages.Get("keyboard_arrow_down", context.ColorMode())
	if err != nil {
		return err
	}

	n.upButton.SetImage(imgUp)
	n.upButton.setSharpenCorners(draw.SharpenCorners{
		LowerLeft:  true,
		LowerRight: true,
	})
	n.upButton.setPairedButton(&n.downButton)
	n.upButton.setOnRepeat(func() {
		n.increment()
	})
	context.SetEnabled(&n.upButton, n.IsEditable() && (!n.maxSet || n.value.Cmp(&n.max) < 0))

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

	n.downButton.SetImage(imgDown)
	n.downButton.setSharpenCorners(draw.SharpenCorners{
		UpperLeft:  true,
		UpperRight: true,
	})
	n.downButton.setPairedButton(&n.upButton)
	n.downButton.setOnRepeat(func() {
		n.decrement()
	})
	context.SetEnabled(&n.downButton, n.IsEditable() && (!n.minSet || n.value.Cmp(&n.min) > 0))

	appender.AppendChildWidgetWithBounds(&n.downButton, image.Rectangle{
		Min: image.Point{
			X: b.Max.X - UnitSize(context)*3/4,
			Y: b.Min.Y + b.Dy()/2,
		},
		Max: b.Max,
	})

	return nil
}

var numberTextReplacer = strings.NewReplacer(
	"\u2212", "-",
	"\ufe62", "+",
	"\ufe63", "-",
	"\uff0b", "+",
	"\uff0d", "-",
	"\uff10", "0",
	"\uff11", "1",
	"\uff12", "2",
	"\uff13", "3",
	"\uff14", "4",
	"\uff15", "5",
	"\uff16", "6",
	"\uff17", "7",
	"\uff18", "8",
	"\uff19", "9",
)

func (n *NumberInput) commit(text string) {
	text = strings.TrimSpace(text)
	text = numberTextReplacer.Replace(text)

	var v big.Int
	if _, ok := v.SetString(text, 10); !ok {
		return
	}
	n.SetValueBigInt(&v)
	n.fireValueChangeEvents()
}

func (n *NumberInput) fireValueChangeEvents() {
	if n.onValueChangedBigInt != nil {
		n.onValueChangedBigInt(n.ValueBigInt())
	}
	if n.onValueChangedInt64 != nil {
		n.onValueChangedInt64(n.ValueInt64())
	}
	if n.onValueChangedUint64 != nil {
		n.onValueChangedUint64(n.ValueUint64())
	}
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

func (n *NumberInput) increment() {
	if !n.IsEditable() {
		return
	}
	n.commit(n.textInput.Value())
	var step big.Int
	if n.stepSet {
		step.Set(&n.step)
	} else {
		step.SetInt64(1)
	}
	var newValue big.Int
	newValue.Add(&n.value, &step)
	n.setValue(&newValue, true)
}

func (n *NumberInput) decrement() {
	if !n.IsEditable() {
		return
	}
	n.commit(n.textInput.Value())
	var step big.Int
	if n.stepSet {
		step.Set(&n.step)
	} else {
		step.SetInt64(1)
	}
	var newValue big.Int
	newValue.Sub(&n.value, &step)
	n.setValue(&newValue, true)
}

func (n *NumberInput) DefaultSize(context *guigui.Context) image.Point {
	return n.textInput.DefaultSize(context)
}

func (n *NumberInput) clamp(value *big.Int) {
	if n.minSet && value.Cmp(&n.min) < 0 {
		value.Set(&n.min)
		return
	}
	if n.maxSet && value.Cmp(&n.max) > 0 {
		value.Set(&n.max)
		return
	}
}
