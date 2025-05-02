// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package main

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	_ "github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.RootWidget

	fill bool
	gap  bool

	configForm basicwidget.Form
	fillText   basicwidget.Text
	fillToggle basicwidget.Toggle
	gapText    basicwidget.Text
	gapToggle  basicwidget.Toggle

	background basicwidget.Background
	buttons    [16]guigui.Widget
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	r.fillText.SetText("Fill Widgets into Grid Cells")
	r.fillToggle.SetValue(r.fill)
	r.fillToggle.SetOnValueChanged(func(value bool) {
		r.fill = value
	})
	r.gapText.SetText("Use Gap")
	r.gapToggle.SetValue(r.gap)
	r.gapToggle.SetOnValueChanged(func(value bool) {
		r.gap = value
	})
	formItems := []*basicwidget.FormItem{
		{
			PrimaryWidget:   &r.fillText,
			SecondaryWidget: &r.fillToggle,
		},
		{
			PrimaryWidget:   &r.gapText,
			SecondaryWidget: &r.gapToggle,
		},
	}
	r.configForm.SetItems(formItems)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(int(u / 2)),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row == 0 {
					return layout.FixedSize(context.Size(&r.configForm).Y)
				}
				return layout.FixedSize(0)
			}),
			layout.FlexibleSize(1),
		},
		RowGap: int(u / 2),
	}
	for i, bounds := range gl.CellBounds() {
		if i == 0 {
			appender.AppendChildWidgetWithBounds(&r.configForm, bounds)
			continue
		}

		for i := range r.buttons {
			if r.buttons[i] == nil {
				r.buttons[i] = &basicwidget.TextButton{}
			}
			t := r.buttons[i].(*basicwidget.TextButton)
			t.SetText(fmt.Sprintf("Button %d", i))
		}

		var firstColumnWidth int
		for j := range 4 {
			firstColumnWidth = max(firstColumnWidth, r.buttons[4*j].DefaultSize(context).X)
		}
		g := layout.GridLayout{
			Bounds: bounds,
			Widths: []layout.Size{
				layout.FixedSize(firstColumnWidth),
				layout.FixedSize(200),
				layout.FlexibleSize(1),
				layout.FlexibleSize(2),
			},
			Heights: []layout.Size{
				layout.LazySize(func(row int) layout.Size {
					var height int
					for i := range 4 {
						height = max(height, r.buttons[4*row+i].DefaultSize(context).Y)
					}
					return layout.FixedSize(height)
				}),
				layout.FixedSize(100),
				layout.FlexibleSize(1),
				layout.FlexibleSize(2),
			},
		}
		if r.gap {
			g.ColumnGap = int(u / 2)
			g.RowGap = int(u / 2)
		}
		for i, bounds := range g.CellBounds() {
			widget := r.buttons[i]
			if r.fill {
				appender.AppendChildWidgetWithBounds(widget, bounds)
			} else {
				pt := bounds.Min
				s := widget.DefaultSize(context)
				pt.X += (bounds.Dx() - s.X) / 2
				pt.Y += (bounds.Dy() - s.Y) / 2
				appender.AppendChildWidgetWithBounds(widget, image.Rectangle{
					Min: pt,
					Max: pt.Add(s),
				})
			}
		}
	}
	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title: "Grid Layout",
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
