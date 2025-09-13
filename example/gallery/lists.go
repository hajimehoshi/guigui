// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"
	"slices"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Lists struct {
	guigui.DefaultWidget

	listForm basicwidget.Form
	listText basicwidget.Text
	list     guigui.WidgetWithSize[*basicwidget.List[int]]

	configForm       basicwidget.Form
	showStripeText   basicwidget.Text
	showStripeToggle basicwidget.Toggle
	showHeaderText   basicwidget.Text
	showHeaderToggle basicwidget.Toggle
	showFooterText   basicwidget.Text
	showFooterToggle basicwidget.Toggle
	movableText      basicwidget.Text
	movableToggle    basicwidget.Toggle
	enabledText      basicwidget.Text
	enabledToggle    basicwidget.Toggle

	items []basicwidget.ListItem[int]

	layout layout.GridLayout
}

func (l *Lists) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&l.listForm)
	appender.AppendChildWidget(&l.configForm)
}

func (l *Lists) Update(context *guigui.Context) error {
	model := context.Model(l, modelKeyModel).(*Model)

	u := basicwidget.UnitSize(context)

	// Lists
	l.listText.SetValue("Text list")

	list := l.list.Widget()
	list.SetStripeVisible(model.Lists().IsStripeVisible())
	if model.Lists().IsHeaderVisible() {
		list.SetHeaderHeight(u)
	} else {
		list.SetHeaderHeight(0)
	}
	if model.Lists().IsFooterVisible() {
		list.SetFooterHeight(u)
	} else {
		list.SetFooterHeight(0)
	}
	list.SetOnItemsMoved(func(from, count, to int) {
		idx := model.Lists().MoveListItems(from, count, to)
		list.SelectItemByIndex(idx)
	})

	l.items = slices.Delete(l.items, 0, len(l.items))
	l.items = model.lists.AppendListItems(l.items)
	list.SetItems(l.items)
	context.SetEnabled(&l.list, model.Lists().Enabled())
	l.list.SetFixedHeight(6 * u)

	l.listForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &l.listText,
			SecondaryWidget: &l.list,
		},
	})

	// Configurations
	l.showStripeText.SetValue("Show stripe")
	l.showStripeToggle.SetOnValueChanged(func(value bool) {
		model.Lists().SetStripeVisible(value)
	})
	l.showStripeToggle.SetValue(model.Lists().IsStripeVisible())
	l.showHeaderText.SetValue("Show header")
	l.showHeaderToggle.SetOnValueChanged(func(value bool) {
		model.Lists().SetHeaderVisible(value)
	})
	l.showHeaderToggle.SetValue(model.Lists().IsHeaderVisible())
	l.showFooterText.SetValue("Show footer")
	l.showFooterToggle.SetOnValueChanged(func(value bool) {
		model.Lists().SetFooterVisible(value)
	})
	l.showFooterToggle.SetValue(model.Lists().IsFooterVisible())
	l.movableText.SetValue("Enable to move items")
	l.movableToggle.SetValue(model.Lists().Movable())
	l.movableToggle.SetOnValueChanged(func(value bool) {
		model.Lists().SetMovable(value)
	})
	l.enabledText.SetValue("Enabled")
	l.enabledToggle.SetOnValueChanged(func(value bool) {
		model.Lists().SetEnabled(value)
	})
	l.enabledToggle.SetValue(model.Lists().Enabled())

	l.configForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &l.showStripeText,
			SecondaryWidget: &l.showStripeToggle,
		},
		{
			PrimaryWidget:   &l.showHeaderText,
			SecondaryWidget: &l.showHeaderToggle,
		},
		{
			PrimaryWidget:   &l.showFooterText,
			SecondaryWidget: &l.showFooterToggle,
		},
		{
			PrimaryWidget:   &l.movableText,
			SecondaryWidget: &l.movableToggle,
		},
		{
			PrimaryWidget:   &l.enabledText,
			SecondaryWidget: &l.enabledToggle,
		},
	})

	l.layout = layout.GridLayout{
		Bounds: context.Bounds(l).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(l.listForm.Measure(context, guigui.FixedWidthConstraints(context.Bounds(l).Dx()-u)).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(l.configForm.Measure(context, guigui.FixedWidthConstraints(context.Bounds(l).Dx()-u)).Y),
		},
		RowGap: u / 2,
	}

	return nil
}

func (l *Lists) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &l.listForm:
		return l.layout.CellBounds(0, 0)
	case &l.configForm:
		return l.layout.CellBounds(0, 2)
	}
	return image.Rectangle{}
}
