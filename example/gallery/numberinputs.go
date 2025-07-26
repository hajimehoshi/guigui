// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"
	"math"
	"math/big"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type NumberInputs struct {
	guigui.DefaultWidget

	numberInputForm       basicwidget.Form
	numberInput1Text      basicwidget.Text
	numberInput1          basicwidget.NumberInput
	numberInput2Text      basicwidget.Text
	numberInput2          basicwidget.NumberInput
	numberInput3Text      basicwidget.Text
	numberInput3          basicwidget.NumberInput
	sliderText            basicwidget.Text
	slider                basicwidget.Slider
	slierWithoutRangeText basicwidget.Text
	sliderWithoutRange    basicwidget.Slider

	configForm     basicwidget.Form
	editableText   basicwidget.Text
	editableToggle basicwidget.Toggle
	enabledText    basicwidget.Text
	enabledToggle  basicwidget.Toggle
}

func (n *NumberInputs) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&n.numberInputForm)
	appender.AppendChildWidget(&n.configForm)
}

func (n *NumberInputs) Build(context *guigui.Context) error {
	model := context.Model(n, modelKeyModel).(*Model)

	u := basicwidget.UnitSize(context)

	// Number Inputs
	width := 12 * u

	n.numberInput1Text.SetValue("Number input")
	n.numberInput1.SetOnValueChangedBigInt(func(value *big.Int, committed bool) {
		if !committed {
			return
		}
		model.NumberInputs().SetNumberInputValue1(value)
	})
	n.numberInput1.SetValueBigInt(model.NumberInputs().NumberInputValue1())
	n.numberInput1.SetEditable(model.NumberInputs().Editable())
	context.SetEnabled(&n.numberInput1, model.NumberInputs().Enabled())
	context.SetSize(&n.numberInput1, image.Pt(width, guigui.AutoSize), n)

	n.numberInput2Text.SetValue("Number input (uint64)")
	n.numberInput2.SetOnValueChangedUint64(func(value uint64, committed bool) {
		if !committed {
			return
		}
		model.NumberInputs().SetNumberInputValue2(value)
	})
	n.numberInput2.SetMinimumValueUint64(0)
	n.numberInput2.SetMaximumValueUint64(math.MaxUint64)
	n.numberInput2.SetValueUint64(model.NumberInputs().NumberInputValue2())
	n.numberInput2.SetEditable(model.NumberInputs().Editable())
	context.SetEnabled(&n.numberInput2, model.NumberInputs().Enabled())
	context.SetSize(&n.numberInput2, image.Pt(width, guigui.AutoSize), n)

	n.numberInput3Text.SetValue("Number input (Range: [-100, 100], Step: 5)")
	n.numberInput3.SetOnValueChangedInt64(func(value int64, committed bool) {
		if !committed {
			return
		}
		model.NumberInputs().SetNumberInputValue3(int(value))
	})
	n.numberInput3.SetMinimumValueInt64(-100)
	n.numberInput3.SetMaximumValueInt64(100)
	n.numberInput3.SetStepInt64(5)
	n.numberInput3.SetValueInt64(int64(model.NumberInputs().NumberInputValue3()))
	n.numberInput3.SetEditable(model.NumberInputs().Editable())
	context.SetEnabled(&n.numberInput3, model.NumberInputs().Enabled())
	context.SetSize(&n.numberInput3, image.Pt(width, guigui.AutoSize), n)

	n.sliderText.SetValue("Slider (Range: [-100, 100])")
	n.slider.SetOnValueChangedInt64(func(value int64) {
		model.NumberInputs().SetNumberInputValue3(int(value))
	})
	n.slider.SetMinimumValueInt64(-100)
	n.slider.SetMaximumValueInt64(100)
	n.slider.SetValueInt64(int64(model.NumberInputs().NumberInputValue3()))
	context.SetEnabled(&n.slider, model.NumberInputs().Enabled())
	context.SetSize(&n.slider, image.Pt(width, guigui.AutoSize), n)

	n.slierWithoutRangeText.SetValue("Slider w/o range")
	n.sliderWithoutRange.SetOnValueChangedInt64(func(value int64) {
	})
	context.SetEnabled(&n.sliderWithoutRange, model.NumberInputs().Enabled())
	context.SetSize(&n.sliderWithoutRange, image.Pt(width, guigui.AutoSize), n)

	n.numberInputForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &n.numberInput1Text,
			SecondaryWidget: &n.numberInput1,
		},
		{
			PrimaryWidget:   &n.numberInput2Text,
			SecondaryWidget: &n.numberInput2,
		},
		{
			PrimaryWidget:   &n.numberInput3Text,
			SecondaryWidget: &n.numberInput3,
		},
		{
			PrimaryWidget:   &n.sliderText,
			SecondaryWidget: &n.slider,
		},
		{
			PrimaryWidget:   &n.slierWithoutRangeText,
			SecondaryWidget: &n.sliderWithoutRange,
		},
	})

	// Configurations
	n.editableText.SetValue("Editable (for number inputs)")
	n.editableToggle.SetOnValueChanged(func(value bool) {
		model.NumberInputs().SetEditable(value)
	})
	n.editableToggle.SetValue(model.NumberInputs().Editable())

	n.enabledText.SetValue("Enabled")
	n.enabledToggle.SetOnValueChanged(func(value bool) {
		model.NumberInputs().SetEnabled(value)
	})
	n.enabledToggle.SetValue(model.NumberInputs().Enabled())

	n.configForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &n.editableText,
			SecondaryWidget: &n.editableToggle,
		},
		{
			PrimaryWidget:   &n.enabledText,
			SecondaryWidget: &n.enabledToggle,
		},
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(n).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(n.numberInputForm.DefaultSizeInContainer(context, context.Bounds(n).Dx()-u).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(n.configForm.DefaultSizeInContainer(context, context.Bounds(n).Dx()-u).Y),
		},
		RowGap: u / 2,
	}
	context.SetBounds(&n.numberInputForm, gl.CellBounds(0, 0), n)
	context.SetBounds(&n.configForm, gl.CellBounds(0, 2), n)

	return nil
}
