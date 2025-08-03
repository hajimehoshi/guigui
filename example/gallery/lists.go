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

	items []basicwidget.ListItem[int]
}

func (l *Lists) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&l.listForm)
	appender.AppendChildWidget(&l.configForm)
}

func (l *Lists) Build(context *guigui.Context) error {
	model := context.Model(l, modelKeyModel).(*Model)

	u := basicwidget.UnitSize(context)

	// Lists
	l.listText.SetValue("Text list")

	l.list.SetStripeVisible(model.Lists().IsStripeVisible())
	if model.Lists().IsHeaderVisible() {
		l.list.SetHeaderHeight(u)
	} else {
		l.list.SetHeaderHeight(0)
	}
	if model.Lists().IsFooterVisible() {
		l.list.SetFooterHeight(u)
	} else {
		l.list.SetFooterHeight(0)
	}
	l.list.SetOnItemsMoved(func(context *guigui.Context, from, count, to int) {
		idx := model.Lists().MoveListItems(from, count, to)
		l.list.SelectItemByIndex(idx)
	})

	l.items = slices.Delete(l.items, 0, len(l.items))
	l.items = model.lists.AppendListItems(l.items)
	l.list.SetItems(l.items)
	context.SetSize(&l.list, image.Pt(guigui.AutoSize, 6*basicwidget.UnitSize(context)), l)
	context.SetEnabled(&l.list, model.Lists().Enabled())

	l.listForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &l.listText,
			SecondaryWidget: &l.list,
		},
	})

	// Configurations
	l.showStripeText.SetValue("Show stripe")
	l.showStripeToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
		model.Lists().SetStripeVisible(value)
	})
	l.showStripeToggle.SetValue(model.Lists().IsStripeVisible())
	l.showHeaderText.SetValue("Show header")
	l.showHeaderToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
		model.Lists().SetHeaderVisible(value)
	})
	l.showHeaderToggle.SetValue(model.Lists().IsHeaderVisible())
	l.showFooterText.SetValue("Show footer")
	l.showFooterToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
		model.Lists().SetFooterVisible(value)
	})
	l.showFooterToggle.SetValue(model.Lists().IsFooterVisible())
	l.movableText.SetValue("Enable to move items")
	l.movableToggle.SetValue(model.Lists().Movable())
	l.movableToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
		model.Lists().SetMovable(value)
	})
	l.enabledText.SetValue("Enabled")
	l.enabledToggle.SetOnValueChanged(func(context *guigui.Context, value bool) {
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

	gl := layout.GridLayout{
		Bounds: context.Bounds(l).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(l.listForm.DefaultSizeInContainer(context, context.Bounds(l).Dx()-u).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(l.configForm.DefaultSizeInContainer(context, context.Bounds(l).Dx()-u).Y),
		},
		RowGap: u / 2,
	}
	context.SetBounds(&l.listForm, gl.CellBounds(0, 0), l)
	context.SetBounds(&l.configForm, gl.CellBounds(0, 2), l)

	return nil
}
