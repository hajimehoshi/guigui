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
	buttonText      basicwidget.Text
	button          basicwidget.Button
	toggleText      basicwidget.Text
	toggle          basicwidget.Toggle
	textInputText   basicwidget.Text
	textInput       basicwidget.TextInput
	numberInputText basicwidget.Text
	numberInput     basicwidget.NumberInput
	sliderText      basicwidget.Text
	slider          basicwidget.Slider
	listText        basicwidget.Text
	list            basicwidget.List[int]
}

func (b *Basic) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&b.form)
}

func (b *Basic) Build(context *guigui.Context) error {
	b.buttonText.SetValue("Button")
	b.button.SetText("Click me!")
	b.toggleText.SetValue("Toggle")
	b.textInputText.SetValue("Text input")
	b.textInput.SetHorizontalAlign(basicwidget.HorizontalAlignEnd)
	b.numberInputText.SetValue("Number input")
	b.sliderText.SetValue("Slider")
	b.slider.SetMinimumValueInt64(0)
	b.slider.SetMaximumValueInt64(100)
	b.listText.SetValue("Text list")
	b.list.SetItemsByStrings([]string{"Item 1", "Item 2", "Item 3"})

	b.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &b.buttonText,
			SecondaryWidget: &b.button,
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
			PrimaryWidget:   &b.listText,
			SecondaryWidget: &b.list,
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
				return layout.FixedSize(b.form.DefaultSizeInContainer(context, context.Bounds(b).Dx()-u).Y)
			}),
		},
		RowGap: u / 2,
	}
	context.SetBounds(&b.form, gl.CellBounds(0, 0), b)

	return nil
}
