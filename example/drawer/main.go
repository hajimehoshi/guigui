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
	"github.com/hajimehoshi/guigui/layout"
)

type modelKey int

const (
	modelKeyModel modelKey = iota
)

const dummyText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

type Root struct {
	guigui.DefaultWidget

	background   basicwidget.Background
	toolbar      Toolbar
	leftPanel    LeftPanel
	contentPanel ContentPanel
	rightPanel   RightPanel

	model Model

	mainLayout    layout.GridLayout
	contentLayout layout.GridLayout
}

func (r *Root) Model(key any) any {
	switch key {
	case modelKeyModel:
		return &r.model
	default:
		return nil
	}
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.background)
	appender.AppendChildWidget(&r.toolbar)
	appender.AppendChildWidget(&r.leftPanel)
	appender.AppendChildWidget(&r.contentPanel)
	appender.AppendChildWidget(&r.rightPanel)
}

func (r *Root) Build(context *guigui.Context) error {
	r.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(r),
		Heights: []layout.Size{
			layout.FixedSize(r.toolbar.Measure(context, guigui.Constraints{}).Y),
			layout.FlexibleSize(1),
		},
	}
	r.contentLayout = layout.GridLayout{
		Bounds: r.mainLayout.CellBounds(0, 1),
		Widths: []layout.Size{
			layout.FixedSize(r.model.LeftPanelWidth(context)),
			layout.FlexibleSize(1),
			layout.FixedSize(r.model.RightPanelWidth(context)),
		},
	}
	return nil
}

func (r *Root) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.background:
		return context.Bounds(r)
	case &r.toolbar:
		return r.mainLayout.CellBounds(0, 0)
	case &r.leftPanel:
		b := r.contentLayout.CellBounds(0, 0)
		b.Min.X = b.Max.X - r.model.DefaultPanelWidth(context)
		return b
	case &r.contentPanel:
		return r.contentLayout.CellBounds(1, 0)
	case &r.rightPanel:
		b := r.contentLayout.CellBounds(2, 0)
		b.Max.X = b.Min.X + r.model.DefaultPanelWidth(context)
		return b
	}
	return image.Rectangle{}
}

func (r *Root) Tick(context *guigui.Context) error {
	r.model.Tick()
	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title:      "Drawers",
		WindowSize: image.Pt(800, 600),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
