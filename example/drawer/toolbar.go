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

func (t *Toolbar) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&t.panel)
}

func (t *Toolbar) Build(context *guigui.Context) error {
	t.panel.SetStyle(basicwidget.PanelStyleSide)
	t.panel.SetBorder(basicwidget.PanelBorder{
		Bottom: true,
	})
	context.SetSize(&t.content, context.ActualSize(t), t)
	t.panel.SetContent(&t.content)
	context.SetBounds(&t.panel, context.Bounds(t), t)

	return nil
}

func (t *Toolbar) DefaultSize(context *guigui.Context) image.Point {
	return t.content.DefaultSize(context)
}

type toolbarContent struct {
	guigui.DefaultWidget

	leftPanelButton  basicwidget.Button
	rightPanelButton basicwidget.Button
}

func (t *toolbarContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&t.leftPanelButton)
	appender.AppendChildWidget(&t.rightPanelButton)
}

func (t *toolbarContent) Build(context *guigui.Context) error {
	model := context.Model(t, modelKeyModel).(*Model)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 4),
		Widths: []layout.Size{
			layout.FixedSize(u * 3 / 2),
			layout.FlexibleSize(1),
			layout.FixedSize(u * 3 / 2),
		},
	}
	if model.IsLeftPanelOpen() {
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
	if model.IsRightPanelOpen() {
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
		model.SetLeftPanelOpen(!model.IsLeftPanelOpen())
	})
	t.rightPanelButton.SetOnDown(func() {
		model.SetRightPanelOpen(!model.IsRightPanelOpen())
	})
	context.SetBounds(&t.leftPanelButton, gl.CellBounds(0, 0), t)
	context.SetBounds(&t.rightPanelButton, gl.CellBounds(2, 0), t)

	return nil
}

func (t *toolbarContent) DefaultSize(context *guigui.Context) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(t.DefaultWidget.DefaultSize(context).X, 2*u)
}
