// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	_ "github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.DefaultWidget

	fill bool
	gap  bool

	configForm basicwidget.Form
	fillText   basicwidget.Text
	fillToggle basicwidget.Toggle
	gapText    basicwidget.Text
	gapToggle  basicwidget.Toggle

	background basicwidget.Background
	buttons    [16]guigui.Widget

	mainLayout    layout.GridLayout
	contentLayout layout.GridLayout
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.background)
	appender.AppendChildWidget(&r.configForm)
	for i := range r.buttons {
		if r.buttons[i] != nil {
			appender.AppendChildWidget(r.buttons[i])
		}
	}
}

func (r *Root) Build(context *guigui.Context) error {
	r.fillText.SetValue("Fill Widgets into Grid Cells")
	r.fillToggle.SetValue(r.fill)
	r.fillToggle.SetOnValueChanged(func(value bool) {
		r.fill = value
	})
	r.gapText.SetValue("Use Gap")
	r.gapToggle.SetValue(r.gap)
	r.gapToggle.SetOnValueChanged(func(value bool) {
		r.gap = value
	})
	r.configForm.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &r.fillText,
			SecondaryWidget: &r.fillToggle,
		},
		{
			PrimaryWidget:   &r.gapText,
			SecondaryWidget: &r.gapToggle,
		},
	})

	u := basicwidget.UnitSize(context)
	r.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(r).Inset(int(u / 2)),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row == 0 {
					return layout.FixedSize(r.configForm.Measure(context, guigui.FixedWidthConstraints(context.Bounds(r).Dx()-u)).Y)
				}
				return layout.FixedSize(0)
			}),
			layout.FlexibleSize(1),
		},
		RowGap: int(u / 2),
	}

	for i := range r.buttons {
		if r.buttons[i] == nil {
			r.buttons[i] = &basicwidget.Button{}
		}
		t := r.buttons[i].(*basicwidget.Button)
		t.SetText(fmt.Sprintf("Button %d", i))
	}

	var firstColumnWidth int
	for j := range 4 {
		firstColumnWidth = max(firstColumnWidth, r.buttons[4*j].Measure(context, guigui.Constraints{}).X)
	}
	r.contentLayout = layout.GridLayout{
		Bounds: r.mainLayout.CellBounds(0, 1),
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
					height = max(height, r.buttons[4*row+i].Measure(context, guigui.Constraints{}).Y)
				}
				return layout.FixedSize(height)
			}),
			layout.FixedSize(100),
			layout.FlexibleSize(1),
			layout.FlexibleSize(2),
		},
	}
	if r.gap {
		r.contentLayout.ColumnGap = int(u / 2)
		r.contentLayout.RowGap = int(u / 2)
	}

	return nil
}

func (r *Root) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.background:
		return context.Bounds(r)
	case &r.configForm:
		return r.mainLayout.CellBounds(0, 0)
	}
	// TODO: This is not efficient. Define a better layout utility to get bounds from a widget directly.
	for i := range r.buttons {
		if widget != r.buttons[i] {
			continue
		}
		b := r.contentLayout.CellBounds(i%4, i/4)
		if r.fill {
			return b
		}
		pt := b.Min
		s := widget.Measure(context, guigui.Constraints{})
		pt.X += (b.Dx() - s.X) / 2
		pt.Y += (b.Dy() - s.Y) / 2
		return image.Rectangle{
			Min: pt,
			Max: pt.Add(s),
		}
	}
	return image.Rectangle{}
}

func main() {
	op := &guigui.RunOptions{
		Title: "Grid Layout",
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
