// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Toolbar struct {
	guigui.DefaultWidget

	panel   basicwidget.Panel
	content toolbarContent
}

func (t *Toolbar) SetModel(model *Model) {
	t.content.model = model
}

func (t *Toolbar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	t.panel.SetStyle(basicwidget.PanelStyleSide)
	t.panel.SetBorder(basicwidget.PanelBorder{
		Bottom: true,
	})
	context.SetSize(&t.content, context.ActualSize(t))
	t.panel.SetContent(&t.content)
	appender.AppendChildWidgetWithBounds(&t.panel, context.Bounds(t))

	return nil
}

func (t *Toolbar) DefaultSize(context *guigui.Context) image.Point {
	return t.content.DefaultSize(context)
}

type toolbarContent struct {
	guigui.DefaultWidget

	leftPanelButton  basicwidget.Button
	rightPanelButton basicwidget.Button

	model *Model
}

func (t *toolbarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 4),
		Widths: []layout.Size{
			layout.FixedSize(u * 3 / 2),
			layout.FlexibleSize(1),
			layout.FixedSize(u * 3 / 2),
		},
	}
	if t.model.IsLeftPanelOpen() {
		img, err := theImageCache.GetMonochrome("left_panel_close", context.ColorMode())
		if err != nil {
			return err
		}
		t.leftPanelButton.SetIcon(img)
	} else {
		img, err := theImageCache.GetMonochrome("left_panel_open", context.ColorMode())
		if err != nil {
			return err
		}
		t.leftPanelButton.SetIcon(img)
	}
	if t.model.IsRightPanelOpen() {
		img, err := theImageCache.GetMonochrome("right_panel_close", context.ColorMode())
		if err != nil {
			return err
		}
		t.rightPanelButton.SetIcon(img)
	} else {
		img, err := theImageCache.GetMonochrome("right_panel_open", context.ColorMode())
		if err != nil {
			return err
		}
		t.rightPanelButton.SetIcon(img)
	}
	t.leftPanelButton.SetOnDown(func() {
		t.model.SetLeftPanelOpen(!t.model.IsLeftPanelOpen())
	})
	t.rightPanelButton.SetOnDown(func() {
		t.model.SetRightPanelOpen(!t.model.IsRightPanelOpen())
	})
	appender.AppendChildWidgetWithBounds(&t.leftPanelButton, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&t.rightPanelButton, gl.CellBounds(2, 0))

	return nil
}

func (t *toolbarContent) DefaultSize(context *guigui.Context) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(t.DefaultWidget.DefaultSize(context).X, 2*u)
}
