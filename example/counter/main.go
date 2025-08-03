// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package main

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.DefaultWidget

	background  basicwidget.Background
	resetButton basicwidget.Button
	decButton   basicwidget.Button
	incButton   basicwidget.Button
	counterText basicwidget.Text

	counter int
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.background)
	appender.AppendChildWidget(&r.counterText)
	appender.AppendChildWidget(&r.resetButton)
	appender.AppendChildWidget(&r.decButton)
	appender.AppendChildWidget(&r.incButton)
}

func (r *Root) Build(context *guigui.Context) error {
	context.SetBounds(&r.background, context.Bounds(r), r)

	r.counterText.SetSelectable(true)
	r.counterText.SetBold(true)
	r.counterText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	r.counterText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	r.counterText.SetScale(4)
	r.counterText.SetValue(fmt.Sprintf("%d", r.counter))

	r.resetButton.SetText("Reset")
	r.resetButton.SetOnUp(func(context *guigui.Context) {
		r.counter = 0
	})
	context.SetEnabled(&r.resetButton, r.counter != 0)

	r.decButton.SetText("Decrement")
	r.decButton.SetOnUp(func(context *guigui.Context) {
		r.counter--
	})

	r.incButton.SetText("Increment")
	r.incButton.SetOnUp(func(context *guigui.Context) {
		r.counter++
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(u),
		},
		RowGap: u,
	}
	context.SetBounds(&r.counterText, gl.CellBounds(0, 0), r)
	{
		gl := layout.GridLayout{
			Bounds: gl.CellBounds(0, 1),
			Widths: []layout.Size{
				layout.FixedSize(6 * u),
				layout.FlexibleSize(1),
				layout.FixedSize(6 * u),
				layout.FixedSize(6 * u),
			},
			ColumnGap: u / 2,
		}
		context.SetBounds(&r.resetButton, gl.CellBounds(0, 0), r)
		context.SetBounds(&r.decButton, gl.CellBounds(2, 0), r)
		context.SetBounds(&r.incButton, gl.CellBounds(3, 0), r)
	}

	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title:         "Counter",
		WindowMinSize: image.Pt(600, 300),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
