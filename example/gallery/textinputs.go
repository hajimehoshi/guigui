// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type TextInputs struct {
	guigui.DefaultWidget

	textInputForm       basicwidget.Form
	singleLineText      basicwidget.Text
	singleLineTextInput basicwidget.TextInput
	multilineText       basicwidget.Text
	multilineTextInput  basicwidget.TextInput

	configForm                      basicwidget.Form
	horizontalAlignText             basicwidget.Text
	horizontalAlignSegmentedControl basicwidget.SegmentedControl[basicwidget.HorizontalAlign]
	verticalAlignText               basicwidget.Text
	verticalAlignSegmentedControl   basicwidget.SegmentedControl[basicwidget.VerticalAlign]
	autoWrapText                    basicwidget.Text
	autoWrapToggle                  basicwidget.Toggle
	editableText                    basicwidget.Text
	editableToggle                  basicwidget.Toggle
	enabledText                     basicwidget.Text
	enabledToggle                   basicwidget.Toggle

	model *Model
}

func (t *TextInputs) SetModel(model *Model) {
	t.model = model
}

func (t *TextInputs) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	imgAlignStart, err := theImageCache.Get("format_align_left", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignCenter, err := theImageCache.Get("format_align_center", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignEnd, err := theImageCache.Get("format_align_right", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignTop, err := theImageCache.Get("vertical_align_top", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignMiddle, err := theImageCache.Get("vertical_align_center", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignBottom, err := theImageCache.Get("vertical_align_bottom", context.ColorMode())
	if err != nil {
		return err
	}

	u := basicwidget.UnitSize(context)

	// Text Inputs
	width := 12 * u

	t.singleLineText.SetValue("Single Line")
	t.singleLineTextInput.SetOnValueChanged(func(text string, committed bool) {
		if committed {
			t.model.TextInputs().SetSingleLineText(text)
		}
	})
	t.singleLineTextInput.SetValue(t.model.TextInputs().SingleLineText())
	t.singleLineTextInput.SetHorizontalAlign(t.model.TextInputs().HorizontalAlign())
	t.singleLineTextInput.SetVerticalAlign(t.model.TextInputs().VerticalAlign())
	t.singleLineTextInput.SetEditable(t.model.TextInputs().Editable())
	context.SetEnabled(&t.singleLineTextInput, t.model.TextInputs().Enabled())
	context.SetSize(&t.singleLineTextInput, image.Pt(width, guigui.DefaultSize))

	t.multilineText.SetValue("Multiline")
	t.multilineTextInput.SetOnValueChanged(func(text string, committed bool) {
		if committed {
			t.model.TextInputs().SetMultilineText(text)
		}
	})
	t.multilineTextInput.SetValue(t.model.TextInputs().MultilineText())
	t.multilineTextInput.SetMultiline(true)
	t.multilineTextInput.SetHorizontalAlign(t.model.TextInputs().HorizontalAlign())
	t.multilineTextInput.SetVerticalAlign(t.model.TextInputs().VerticalAlign())
	t.multilineTextInput.SetAutoWrap(t.model.TextInputs().AutoWrap())
	t.multilineTextInput.SetEditable(t.model.TextInputs().Editable())
	context.SetEnabled(&t.multilineTextInput, t.model.TextInputs().Enabled())
	context.SetSize(&t.multilineTextInput, image.Pt(width, 4*u))

	t.textInputForm.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &t.singleLineText,
			SecondaryWidget: &t.singleLineTextInput,
		},
		{
			PrimaryWidget:   &t.multilineText,
			SecondaryWidget: &t.multilineTextInput,
		},
	})

	// Configurations
	t.horizontalAlignText.SetValue("Horizontal Align")
	t.horizontalAlignSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[basicwidget.HorizontalAlign]{
		{
			Image: imgAlignStart,
			Tag:   basicwidget.HorizontalAlignStart,
		},
		{
			Image: imgAlignCenter,
			Tag:   basicwidget.HorizontalAlignCenter,
		},
		{
			Image: imgAlignEnd,
			Tag:   basicwidget.HorizontalAlignEnd,
		},
	})
	t.horizontalAlignSegmentedControl.SetOnItemSelected(func(index int) {
		item, ok := t.horizontalAlignSegmentedControl.ItemByIndex(index)
		if !ok {
			t.model.TextInputs().SetHorizontalAlign(basicwidget.HorizontalAlignStart)
			return
		}
		t.model.TextInputs().SetHorizontalAlign(item.Tag)
	})
	t.horizontalAlignSegmentedControl.SelectItemByTag(t.model.TextInputs().HorizontalAlign())

	t.verticalAlignText.SetValue("Vertical Align")
	t.verticalAlignSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[basicwidget.VerticalAlign]{
		{
			Image: imgAlignTop,
			Tag:   basicwidget.VerticalAlignTop,
		},
		{
			Image: imgAlignMiddle,
			Tag:   basicwidget.VerticalAlignMiddle,
		},
		{
			Image: imgAlignBottom,
			Tag:   basicwidget.VerticalAlignBottom,
		},
	})
	t.verticalAlignSegmentedControl.SetOnItemSelected(func(index int) {
		item, ok := t.verticalAlignSegmentedControl.ItemByIndex(index)
		if !ok {
			t.model.TextInputs().SetVerticalAlign(basicwidget.VerticalAlignTop)
			return
		}
		t.model.TextInputs().SetVerticalAlign(item.Tag)
	})
	t.verticalAlignSegmentedControl.SelectItemByTag(t.model.TextInputs().VerticalAlign())

	t.autoWrapText.SetValue("Auto Wrap")
	t.autoWrapToggle.SetOnValueChanged(func(value bool) {
		t.model.TextInputs().SetAutoWrap(value)
	})
	t.autoWrapToggle.SetValue(t.model.TextInputs().AutoWrap())

	t.editableText.SetValue("Editable")
	t.editableToggle.SetOnValueChanged(func(value bool) {
		t.model.TextInputs().SetEditable(value)
	})
	t.editableToggle.SetValue(t.model.TextInputs().Editable())

	t.enabledText.SetValue("Enabled")
	t.enabledToggle.SetOnValueChanged(func(value bool) {
		t.model.TextInputs().SetEnabled(value)
	})
	t.enabledToggle.SetValue(t.model.TextInputs().Enabled())

	t.configForm.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &t.horizontalAlignText,
			SecondaryWidget: &t.horizontalAlignSegmentedControl,
		},
		{
			PrimaryWidget:   &t.verticalAlignText,
			SecondaryWidget: &t.verticalAlignSegmentedControl,
		},
		{
			PrimaryWidget:   &t.autoWrapText,
			SecondaryWidget: &t.autoWrapToggle,
		},
		{
			PrimaryWidget:   &t.editableText,
			SecondaryWidget: &t.editableToggle,
		},
		{
			PrimaryWidget:   &t.enabledText,
			SecondaryWidget: &t.enabledToggle,
		},
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(t.textInputForm.DefaultSize(context).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(t.configForm.DefaultSize(context).Y),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&t.textInputForm, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&t.configForm, gl.CellBounds(0, 2))
	return nil
}
