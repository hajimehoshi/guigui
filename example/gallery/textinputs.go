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

	textInputForm               basicwidget.Form
	singleLineText              basicwidget.Text
	singleLineTextInput         basicwidget.TextInput
	singleLineWithIconText      basicwidget.Text
	singleLineWithIconTextInput basicwidget.TextInput
	multilineText               basicwidget.Text
	multilineTextInput          basicwidget.TextInput
	inlineText                  basicwidget.Text
	inlineTextInput             inlineTextInputContainer

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
}

func (t *TextInputs) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&t.textInputForm)
	appender.AppendChildWidget(&t.configForm)
}

func (t *TextInputs) Build(context *guigui.Context) error {
	model := context.Model(t, modelKeyModel).(*Model)

	imgAlignStart, err := theImageCache.GetMonochrome("format_align_left", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignCenter, err := theImageCache.GetMonochrome("format_align_center", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignEnd, err := theImageCache.GetMonochrome("format_align_right", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignTop, err := theImageCache.GetMonochrome("vertical_align_top", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignMiddle, err := theImageCache.GetMonochrome("vertical_align_center", context.ColorMode())
	if err != nil {
		return err
	}
	imgAlignBottom, err := theImageCache.GetMonochrome("vertical_align_bottom", context.ColorMode())
	if err != nil {
		return err
	}
	imgSearch, err := theImageCache.GetMonochrome("search", context.ColorMode())
	if err != nil {
		return err
	}

	u := basicwidget.UnitSize(context)

	// Text Inputs
	width := 12 * u

	t.singleLineText.SetValue("Single line")
	t.singleLineTextInput.SetOnValueChanged(func(context *guigui.Context, text string, committed bool) {
		if committed {
			model.TextInputs().SetSingleLineText(text)
		}
	})
	t.singleLineTextInput.SetValue(model.TextInputs().SingleLineText())
	t.singleLineTextInput.SetHorizontalAlign(model.TextInputs().HorizontalAlign())
	t.singleLineTextInput.SetVerticalAlign(model.TextInputs().VerticalAlign())
	t.singleLineTextInput.SetEditable(model.TextInputs().Editable())
	context.SetEnabled(&t.singleLineTextInput, model.TextInputs().Enabled())
	context.SetSize(&t.singleLineTextInput, image.Pt(width, guigui.AutoSize), t)

	t.singleLineWithIconText.SetValue("Single line with icon")
	t.singleLineWithIconTextInput.SetHorizontalAlign(model.TextInputs().HorizontalAlign())
	t.singleLineWithIconTextInput.SetVerticalAlign(model.TextInputs().VerticalAlign())
	t.singleLineWithIconTextInput.SetEditable(model.TextInputs().Editable())
	t.singleLineWithIconTextInput.SetIcon(imgSearch)
	context.SetEnabled(&t.singleLineWithIconTextInput, model.TextInputs().Enabled())
	context.SetSize(&t.singleLineWithIconTextInput, image.Pt(width, guigui.AutoSize), t)

	t.multilineText.SetValue("Multiline")
	t.multilineTextInput.SetOnValueChanged(func(context *guigui.Context, text string, committed bool) {
		if committed {
			model.TextInputs().SetMultilineText(text)
		}
	})
	t.multilineTextInput.SetValue(model.TextInputs().MultilineText())
	t.multilineTextInput.SetMultiline(true)
	t.multilineTextInput.SetHorizontalAlign(model.TextInputs().HorizontalAlign())
	t.multilineTextInput.SetVerticalAlign(model.TextInputs().VerticalAlign())
	t.multilineTextInput.SetAutoWrap(model.TextInputs().AutoWrap())
	t.multilineTextInput.SetEditable(model.TextInputs().Editable())
	context.SetEnabled(&t.multilineTextInput, model.TextInputs().Enabled())
	context.SetSize(&t.multilineTextInput, image.Pt(width, 4*u), t)

	t.inlineText.SetValue("Inline")
	t.inlineTextInput.SetHorizontalAlign(model.TextInputs().HorizontalAlign())
	t.inlineTextInput.textInput.SetVerticalAlign(model.TextInputs().VerticalAlign())
	t.inlineTextInput.textInput.SetEditable(model.TextInputs().Editable())
	context.SetEnabled(&t.inlineTextInput, model.TextInputs().Enabled())
	context.SetSize(&t.inlineTextInput, image.Pt(width, guigui.AutoSize), t)

	t.textInputForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &t.singleLineText,
			SecondaryWidget: &t.singleLineTextInput,
		},
		{
			PrimaryWidget:   &t.singleLineWithIconText,
			SecondaryWidget: &t.singleLineWithIconTextInput,
		},
		{
			PrimaryWidget:   &t.multilineText,
			SecondaryWidget: &t.multilineTextInput,
		},
		{
			PrimaryWidget:   &t.inlineText,
			SecondaryWidget: &t.inlineTextInput,
		},
	})

	// Configurations
	t.horizontalAlignText.SetValue("Horizontal align")
	t.horizontalAlignSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[basicwidget.HorizontalAlign]{
		{
			Icon: imgAlignStart,
			ID:   basicwidget.HorizontalAlignStart,
		},
		{
			Icon: imgAlignCenter,
			ID:   basicwidget.HorizontalAlignCenter,
		},
		{
			Icon: imgAlignEnd,
			ID:   basicwidget.HorizontalAlignEnd,
		},
	})
	t.horizontalAlignSegmentedControl.SetOnItemSelected(func(context *guigui.Context, index int) {
		item, ok := t.horizontalAlignSegmentedControl.ItemByIndex(index)
		if !ok {
			model.TextInputs().SetHorizontalAlign(basicwidget.HorizontalAlignStart)
			return
		}
		model.TextInputs().SetHorizontalAlign(item.ID)
	})
	t.horizontalAlignSegmentedControl.SelectItemByID(model.TextInputs().HorizontalAlign())

	t.verticalAlignText.SetValue("Vertical align")
	t.verticalAlignSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[basicwidget.VerticalAlign]{
		{
			Icon: imgAlignTop,
			ID:   basicwidget.VerticalAlignTop,
		},
		{
			Icon: imgAlignMiddle,
			ID:   basicwidget.VerticalAlignMiddle,
		},
		{
			Icon: imgAlignBottom,
			ID:   basicwidget.VerticalAlignBottom,
		},
	})
	t.verticalAlignSegmentedControl.SetOnItemSelected(func(context *guigui.Context, index int) {
		item, ok := t.verticalAlignSegmentedControl.ItemByIndex(index)
		if !ok {
			model.TextInputs().SetVerticalAlign(basicwidget.VerticalAlignTop)
			return
		}
		model.TextInputs().SetVerticalAlign(item.ID)
	})
	t.verticalAlignSegmentedControl.SelectItemByID(model.TextInputs().VerticalAlign())

	t.autoWrapText.SetValue("Auto wrap")
	t.autoWrapToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
		model.TextInputs().SetAutoWrap(value)
	})
	t.autoWrapToggle.SetValue(model.TextInputs().AutoWrap())

	t.editableText.SetValue("Editable")
	t.editableToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
		model.TextInputs().SetEditable(value)
	})
	t.editableToggle.SetValue(model.TextInputs().Editable())

	t.enabledText.SetValue("Enabled")
	t.enabledToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
		model.TextInputs().SetEnabled(value)
	})
	t.enabledToggle.SetValue(model.TextInputs().Enabled())

	t.configForm.SetItems([]basicwidget.FormItem{
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
			layout.FixedSize(t.textInputForm.DefaultSizeInContainer(context, context.Bounds(t).Dx()-u).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(t.configForm.DefaultSizeInContainer(context, context.Bounds(t).Dx()-u).Y),
		},
		RowGap: u / 2,
	}
	context.SetBounds(&t.textInputForm, gl.CellBounds(0, 0), t)
	context.SetBounds(&t.configForm, gl.CellBounds(0, 2), t)
	return nil
}

type inlineTextInputContainer struct {
	guigui.DefaultWidget

	textInput       basicwidget.TextInput
	horizontalAlign basicwidget.HorizontalAlign
}

func (c *inlineTextInputContainer) SetHorizontalAlign(align basicwidget.HorizontalAlign) {
	c.horizontalAlign = align
	c.textInput.SetHorizontalAlign(align)
}

func (c *inlineTextInputContainer) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&c.textInput)
}

func (c *inlineTextInputContainer) Build(context *guigui.Context) error {
	c.textInput.SetStyle(basicwidget.TextInputStyleInline)
	if c.textInput.DefaultSize(context).X > context.ActualSize(c).X {
		context.SetSize(&c.textInput, image.Pt(context.ActualSize(c).X, guigui.AutoSize), c)
	} else {
		context.SetSize(&c.textInput, image.Pt(guigui.AutoSize, guigui.AutoSize), c)
	}

	pos := context.Position(c)
	switch c.horizontalAlign {
	case basicwidget.HorizontalAlignStart:
	case basicwidget.HorizontalAlignCenter:
		pos.X += (context.ActualSize(c).X - context.ActualSize(&c.textInput).X) / 2
	case basicwidget.HorizontalAlignEnd:
		pos.X += context.ActualSize(c).X - context.ActualSize(&c.textInput).X
	}
	context.SetPosition(&c.textInput, pos)
	return nil
}

func (c *inlineTextInputContainer) DefaultSize(context *guigui.Context) image.Point {
	return c.textInput.DefaultSize(context)
}
