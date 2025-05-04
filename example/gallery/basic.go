// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Basic struct {
	guigui.DefaultWidget

	form            basicwidget.Form
	textButtonText  basicwidget.Text
	textButton      basicwidget.TextButton
	toggleText      basicwidget.Text
	toggle          basicwidget.Toggle
	textInputText   basicwidget.Text
	textInput       basicwidget.TextInput
	numberInputText basicwidget.Text
	numberInput     basicwidget.NumberInput
	sliderText      basicwidget.Text
	slider          basicwidget.Slider
	textListText    basicwidget.Text
	textList        basicwidget.TextList[int]
}

func (b *Basic) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	b.textButtonText.SetValue("Text Button")
	b.textButton.SetText("Click Me!")
	b.toggleText.SetValue("Toggle")
	b.textInputText.SetValue("Text Input")
	b.textInput.SetHorizontalAlign(basicwidget.HorizontalAlignEnd)
	b.numberInputText.SetValue("Number Input")
	b.sliderText.SetValue("Slider")
	b.slider.SetMinimumValueInt64(0)
	b.slider.SetMaximumValueInt64(100)
	b.textListText.SetValue("Text List")
	b.textList.SetItemsByStrings([]string{"Item 1", "Item 2", "Item 3"})

	b.form.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &b.textButtonText,
			SecondaryWidget: &b.textButton,
		},
		{
			PrimaryWidget:   &b.toggleText,
			SecondaryWidget: &b.toggle,
		},
		{
			PrimaryWidget:   &b.textInputText,
			SecondaryWidget: &b.textInput,
		},
		{
			PrimaryWidget:   &b.numberInputText,
			SecondaryWidget: &b.numberInput,
		},
		{
			PrimaryWidget:   &b.sliderText,
			SecondaryWidget: &b.slider,
		},
		{
			PrimaryWidget:   &b.textListText,
			SecondaryWidget: &b.textList,
		},
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(b).Inset(u / 2),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row >= 1 {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(b.form.DefaultSize(context).Y)
			}),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&b.form, gl.CellBounds(0, 0))

	return nil
}
