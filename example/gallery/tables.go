// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Tables struct {
	guigui.DefaultWidget

	table basicwidget.Table[int]

	configForm       basicwidget.Form
	showFooterText   basicwidget.Text
	showFooterToggle basicwidget.Toggle
	movableText      basicwidget.Text
	movableToggle    basicwidget.Toggle
	enabledText      basicwidget.Text
	enabledToggle    basicwidget.Toggle

	tableItems       []basicwidget.TableItem[int]
	tableItemWidgets []guigui.Widget
}

func (t *Tables) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	model := context.Model(t, modelKeyModel).(*Model)

	u := basicwidget.UnitSize(context)
	t.table.SetColumns([]basicwidget.TableColumn{
		{
			HeaderText:                "ID",
			HeaderTextHorizontalAlign: basicwidget.HorizontalAlignRight,
			Width:                     layout.FlexibleSize(1),
			MinWidth:                  2 * u,
		},
		{
			HeaderText: "Name",
			Width:      layout.FlexibleSize(2),
			MinWidth:   4 * u,
		},
		{
			HeaderText:                "Amount",
			HeaderTextHorizontalAlign: basicwidget.HorizontalAlignRight,
			Width:                     layout.FlexibleSize(1),
			MinWidth:                  2 * u,
		},
		{
			HeaderText:                "Cost",
			HeaderTextHorizontalAlign: basicwidget.HorizontalAlignRight,
			Width:                     layout.FlexibleSize(1),
			MinWidth:                  2 * u,
		},
	})
	t.tableItems = slices.Delete(t.tableItems, 0, len(t.tableItems))

	// Prepare widgets for table items.
	// Use slices.Grow and slices.Delete not to break the existing widgets.
	n := 4
	if newNum := n * model.Tables().TableItemCount(); len(t.tableItemWidgets) < newNum {
		t.tableItemWidgets = slices.Grow(t.tableItemWidgets, newNum-len(t.tableItemWidgets))[:newNum]
	} else {
		t.tableItemWidgets = slices.Delete(t.tableItemWidgets, newNum, len(t.tableItemWidgets))
	}
	for i := range model.Tables().TableItemCount() {
		for j := range n {
			if t.tableItemWidgets[n*i+j] == nil {
				t.tableItemWidgets[n*i+j] = &basicwidget.Text{}
			}
		}
	}
	for i, item := range model.Tables().TableItems() {
		text := t.tableItemWidgets[n*i].(*basicwidget.Text)
		text.SetValue(strconv.Itoa(item.ID))
		text.SetHorizontalAlign(basicwidget.HorizontalAlignRight)
		text.SetTabular(true)

		text = t.tableItemWidgets[n*i+1].(*basicwidget.Text)
		text.SetValue(item.Name)

		text = t.tableItemWidgets[n*i+2].(*basicwidget.Text)
		text.SetValue(strconv.Itoa(item.Amount))
		text.SetHorizontalAlign(basicwidget.HorizontalAlignRight)
		text.SetTabular(true)

		text = t.tableItemWidgets[n*i+3].(*basicwidget.Text)
		text.SetValue(fmt.Sprintf("%d.%02d", item.Cost/100, item.Cost%100))
		text.SetHorizontalAlign(basicwidget.HorizontalAlignRight)
		text.SetTabular(true)

		t.tableItems = append(t.tableItems, basicwidget.TableItem[int]{
			Contents: t.tableItemWidgets[n*i : n*(i+1)],
			Movable:  model.Tables().Movable(),
			ID:       item.ID,
		})
	}
	t.table.SetItems(t.tableItems)
	// Set the text colors after setting the items, or ItemTextColor will not work correctly.
	for i := range model.Tables().TableItemCount() {
		for j := range n {
			text := t.tableItemWidgets[n*i+j].(*basicwidget.Text)
			text.SetColor(t.table.ItemTextColor(context, i))
		}
	}
	if model.Tables().IsFooterVisible() {
		t.table.SetFooterHeight(u)
	} else {
		t.table.SetFooterHeight(0)
	}
	context.SetEnabled(&t.table, model.Tables().Enabled())
	t.table.SetOnItemsMoved(func(from, count, to int) {
		idx := model.Tables().MoveTableItems(from, count, to)
		t.table.SelectItemByIndex(idx)
	})

	// Configurations
	t.showFooterText.SetValue("Show footer")
	t.showFooterToggle.SetOnValueChanged(func(value bool) {
		model.Tables().SetFooterVisible(value)
	})
	t.movableText.SetValue("Enable to move items")
	t.movableToggle.SetValue(model.Tables().Movable())
	t.movableToggle.SetOnValueChanged(func(value bool) {
		model.Tables().SetMovable(value)
	})
	t.enabledText.SetValue("Enabled")
	t.enabledToggle.SetOnValueChanged(func(value bool) {
		model.Tables().SetEnabled(value)
	})
	t.enabledToggle.SetValue(model.Tables().Enabled())

	t.configForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &t.showFooterText,
			SecondaryWidget: &t.showFooterToggle,
		},
		{
			PrimaryWidget:   &t.movableText,
			SecondaryWidget: &t.movableToggle,
		},
		{
			PrimaryWidget:   &t.enabledText,
			SecondaryWidget: &t.enabledToggle,
		},
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(12 * u),
			layout.FlexibleSize(1),
			layout.FixedSize(t.configForm.DefaultSizeInContainer(context, context.Bounds(t).Dx()-u).Y),
		},
	}

	appender.AppendChildWidgetWithBounds(&t.table, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&t.configForm, gl.CellBounds(0, 2))
	return nil
}
