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
	list     basicwidget.List[int]

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

	model *Model
	items []basicwidget.ListItem[int]
}

func (l *Lists) SetModel(model *Model) {
	l.model = model
}

func (l *Lists) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)

	// Lists
	l.listText.SetValue("Text list")

	l.list.SetStripeVisible(l.model.Lists().IsStripeVisible())
	if l.model.Lists().IsHeaderVisible() {
		l.list.SetHeaderHeight(u)
	} else {
		l.list.SetHeaderHeight(0)
	}
	if l.model.Lists().IsFooterVisible() {
		l.list.SetFooterHeight(u)
	} else {
		l.list.SetFooterHeight(0)
	}
	l.list.SetOnItemsMoved(func(from, count, to int) {
		idx := l.model.Lists().MoveListItems(from, count, to)
		l.list.SelectItemByIndex(idx)
	})

	l.items = slices.Delete(l.items, 0, len(l.items))
	l.items = l.model.lists.AppendListItems(l.items)
	l.list.SetItems(l.items)
	context.SetSize(&l.list, image.Pt(guigui.AutoSize, 6*basicwidget.UnitSize(context)), l)
	context.SetEnabled(&l.list, l.model.Lists().Enabled())

	l.listForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &l.listText,
			SecondaryWidget: &l.list,
		},
	})

	// Configurations
	l.showStripeText.SetValue("Show stripe")
	l.showStripeToggle.SetOnValueChanged(func(value bool) {
		l.model.Lists().SetStripeVisible(value)
	})
	l.showStripeToggle.SetValue(l.model.Lists().IsStripeVisible())
	l.showHeaderText.SetValue("Show header")
	l.showHeaderToggle.SetOnValueChanged(func(value bool) {
		l.model.Lists().SetHeaderVisible(value)
	})
	l.showHeaderToggle.SetValue(l.model.Lists().IsHeaderVisible())
	l.showFooterText.SetValue("Show footer")
	l.showFooterToggle.SetOnValueChanged(func(value bool) {
		l.model.Lists().SetFooterVisible(value)
	})
	l.showFooterToggle.SetValue(l.model.Lists().IsFooterVisible())
	l.movableText.SetValue("Enable to move items")
	l.movableToggle.SetValue(l.model.Lists().Movable())
	l.movableToggle.SetOnValueChanged(func(value bool) {
		l.model.Lists().SetMovable(value)
	})
	l.enabledText.SetValue("Enabled")
	l.enabledToggle.SetOnValueChanged(func(value bool) {
		l.model.Lists().SetEnabled(value)
	})
	l.enabledToggle.SetValue(l.model.Lists().Enabled())

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

	gl := layout.GridLayout{
		Bounds: context.Bounds(l).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(l.listForm.DefaultSizeInContainer(context, context.Bounds(l).Dx()-u).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(l.configForm.DefaultSizeInContainer(context, context.Bounds(l).Dx()-u).Y),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&l.listForm, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&l.configForm, gl.CellBounds(0, 2))

	return nil
}
