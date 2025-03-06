// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Björn De Meyer, Hajime Hoshi

package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Tables struct {
	guigui.DefaultWidget
	table basicwidget.Table
}

func (t *Tables) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	var columns [5]basicwidget.Column

	for j := 0; j < len(columns); j++ {
		var cells [100]basicwidget.Cell
		var items [100]basicwidget.Text
		for i := 0; i < 100; i++ {
			items[i].SetText(fmt.Sprintf("Cell %c%d", 'A'+j, i))
			cells[i].Content = &items[i]
			cells[i].Selectable = true
		}
		columns[j].SetCells(cells[:])
		columns[j].SetCaption(fmt.Sprintf("Column %c", 'A'+j))
		columns[j].ShowCellBorders(true)
		columns[j].Selectable = true
	}

	t.table.SetColumns(columns[:])

	u := float64(basicwidget.UnitSize(context))
	w, _ := t.Size(context)
	t.table.SetWidth(w - int(1*u))

	p := guigui.Position(t).Add(image.Pt(int(0.5*u), int(0.5*u)))
	guigui.SetPosition(&t.table, p)
	appender.AppendChildWidget(&t.table)
}

func (l *Tables) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(l).Size(context)
	w -= sidebarWidth(context)
	return w, h
}
