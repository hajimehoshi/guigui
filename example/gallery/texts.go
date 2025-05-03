// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Texts struct {
	guigui.DefaultWidget

	form                            basicwidget.Form
	horizontalAlignText             basicwidget.Text
	horizontalAlignSegmentedControl basicwidget.SegmentedControl[basicwidget.HorizontalAlign]
	verticalAlignText               basicwidget.Text
	verticalAlignSegmentedControl   basicwidget.SegmentedControl[basicwidget.VerticalAlign]
	autoWrapText                    basicwidget.Text
	autoWrapToggle                  basicwidget.Toggle
	boldText                        basicwidget.Text
	boldToggle                      basicwidget.Toggle
	selectableText                  basicwidget.Text
	selectableToggle                basicwidget.Toggle
	editableText                    basicwidget.Text
	editableToggle                  basicwidget.Toggle
	sampleText                      basicwidget.Text

	model *Model
}

func (t *Texts) SetModel(model *Model) {
	t.model = model
}

func (t *Texts) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
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
			t.model.Texts().SetHorizontalAlign(basicwidget.HorizontalAlignStart)
			return
		}
		t.model.Texts().SetHorizontalAlign(item.Tag)
	})
	t.horizontalAlignSegmentedControl.SelectItemByTag(t.model.Texts().HorizontalAlign())

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
			t.model.Texts().SetVerticalAlign(basicwidget.VerticalAlignTop)
			return
		}
		t.model.Texts().SetVerticalAlign(item.Tag)
	})
	t.verticalAlignSegmentedControl.SelectItemByTag(t.model.Texts().VerticalAlign())

	t.autoWrapText.SetValue("Auto Wrap")
	t.autoWrapToggle.SetOnValueChanged(func(value bool) {
		t.model.Texts().SetAutoWrap(value)
	})
	t.autoWrapToggle.SetValue(t.model.Texts().AutoWrap())

	t.boldText.SetValue("Bold")
	t.boldToggle.SetOnValueChanged(func(value bool) {
		t.model.Texts().SetBold(value)
	})
	t.boldToggle.SetValue(t.model.Texts().Bold())

	t.selectableText.SetValue("Selectable")
	t.selectableToggle.SetOnValueChanged(func(checked bool) {
		t.model.Texts().SetSelectable(checked)
	})
	t.selectableToggle.SetValue(t.model.Texts().Selectable())

	t.editableText.SetValue("Editable")
	t.editableToggle.SetOnValueChanged(func(value bool) {
		t.model.Texts().SetEditable(value)
	})
	t.editableToggle.SetValue(t.model.Texts().Editable())

	t.form.SetItems([]*basicwidget.FormItem{
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
			PrimaryWidget:   &t.boldText,
			SecondaryWidget: &t.boldToggle,
		},
		{
			PrimaryWidget:   &t.selectableText,
			SecondaryWidget: &t.selectableToggle,
		},
		{
			PrimaryWidget:   &t.editableText,
			SecondaryWidget: &t.editableToggle,
		},
	})

	t.sampleText.SetMultiline(true)
	t.sampleText.SetHorizontalAlign(t.model.Texts().HorizontalAlign())
	t.sampleText.SetVerticalAlign(t.model.Texts().VerticalAlign())
	t.sampleText.SetAutoWrap(t.model.Texts().AutoWrap())
	t.sampleText.SetBold(t.model.Texts().Bold())
	t.sampleText.SetSelectable(t.model.Texts().Selectable())
	t.sampleText.SetEditable(t.model.Texts().Editable())
	t.sampleText.SetOnValueChanged(func(text string, committed bool) {
		if committed {
			t.model.Texts().SetText(text)
		}
	})
	if !context.HasFocusedChildWidget(&t.sampleText) {
		t.sampleText.SetValue(t.model.Texts().Text())
	}

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(t.form.DefaultSize(context).Y),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&t.sampleText, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&t.form, gl.CellBounds(0, 1))

	return nil
}
