// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Buttons struct {
	guigui.DefaultWidget

	buttonsForm           basicwidget.Form
	buttonText            basicwidget.Text
	button                basicwidget.Button
	textIconButton1Text   basicwidget.Text
	textIconButton1       basicwidget.Button
	textIconButton2Text   basicwidget.Text
	textIconButton2       basicwidget.Button
	imageButtonText       basicwidget.Text
	imageButton           basicwidget.Button
	segmentedControlHText basicwidget.Text
	segmentedControlH     basicwidget.SegmentedControl[int]
	segmentedControlVText basicwidget.Text
	segmentedControlV     basicwidget.SegmentedControl[int]
	toggleText            basicwidget.Text
	toggle                basicwidget.Toggle

	configForm    basicwidget.Form
	enabledText   basicwidget.Text
	enabledToggle basicwidget.Toggle

	model *Model
}

func (b *Buttons) SetModel(model *Model) {
	b.model = model
}

func (b *Buttons) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)

	b.buttonText.SetValue("Button")
	b.button.SetText("Button")
	context.SetEnabled(&b.button, b.model.Buttons().Enabled())

	b.textIconButton1Text.SetValue("Button w/ text and icon (1)")
	b.textIconButton1.SetText("Button")
	img, err := theImageCache.GetMonochrome("check", context.ColorMode())
	if err != nil {
		return err
	}
	b.textIconButton1.SetIcon(img)
	context.SetEnabled(&b.textIconButton1, b.model.Buttons().Enabled())
	context.SetSize(&b.textIconButton1, image.Pt(6*u, guigui.AutoSize), b)

	b.textIconButton2Text.SetValue("Button w/ text and icon (2)")
	b.textIconButton2.SetText("Button")
	b.textIconButton2.SetIcon(img)
	b.textIconButton2.SetIconAlign(basicwidget.IconAlignEnd)
	context.SetEnabled(&b.textIconButton2, b.model.Buttons().Enabled())
	context.SetSize(&b.textIconButton2, image.Pt(6*u, guigui.AutoSize), b)

	b.imageButtonText.SetValue("Image button")
	img, err = theImageCache.Get("gopher")
	if err != nil {
		return err
	}
	b.imageButton.SetIcon(img)
	context.SetEnabled(&b.imageButton, b.model.Buttons().Enabled())
	context.SetSize(&b.imageButton, image.Pt(2*u, 2*u), b)

	b.segmentedControlHText.SetValue("Segmented control (Horizontal)")
	b.segmentedControlH.SetItems([]basicwidget.SegmentedControlItem[int]{
		{
			Text: "One",
		},
		{
			Text: "Two",
		},
		{
			Text: "Three",
		},
	})
	b.segmentedControlH.SetDirection(basicwidget.SegmentedControlDirectionHorizontal)
	context.SetEnabled(&b.segmentedControlH, b.model.Buttons().Enabled())

	b.segmentedControlVText.SetValue("Segmented control (Vertical)")
	b.segmentedControlV.SetItems([]basicwidget.SegmentedControlItem[int]{
		{
			Text: "One",
		},
		{
			Text: "Two",
		},
		{
			Text: "Three",
		},
	})
	b.segmentedControlV.SetDirection(basicwidget.SegmentedControlDirectionVertical)
	context.SetEnabled(&b.segmentedControlV, b.model.Buttons().Enabled())

	b.toggleText.SetValue("Toggle")
	context.SetEnabled(&b.toggle, b.model.Buttons().Enabled())

	b.buttonsForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &b.buttonText,
			SecondaryWidget: &b.button,
		},
		{
			PrimaryWidget:   &b.textIconButton1Text,
			SecondaryWidget: &b.textIconButton1,
		},
		{
			PrimaryWidget:   &b.textIconButton2Text,
			SecondaryWidget: &b.textIconButton2,
		},
		{
			PrimaryWidget:   &b.imageButtonText,
			SecondaryWidget: &b.imageButton,
		},
		{
			PrimaryWidget:   &b.segmentedControlHText,
			SecondaryWidget: &b.segmentedControlH,
		},
		{
			PrimaryWidget:   &b.segmentedControlVText,
			SecondaryWidget: &b.segmentedControlV,
		},
		{
			PrimaryWidget:   &b.toggleText,
			SecondaryWidget: &b.toggle,
		},
	})

	b.enabledText.SetValue("Enabled")
	b.enabledToggle.SetOnValueChanged(func(enabled bool) {
		b.model.Buttons().SetEnabled(enabled)
	})
	b.enabledToggle.SetValue(b.model.Buttons().Enabled())

	b.configForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &b.enabledText,
			SecondaryWidget: &b.enabledToggle,
		},
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(b).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(b.buttonsForm.DefaultSizeInContainer(context, context.Bounds(b).Dx()-u).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(b.configForm.DefaultSizeInContainer(context, context.Bounds(b).Dx()-u).Y),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&b.buttonsForm, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&b.configForm, gl.CellBounds(0, 2))

	return nil
}
