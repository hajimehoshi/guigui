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
	content guigui.WidgetWithSize[*toolbarContent]
}

func (t *Toolbar) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&t.panel)
}

func (t *Toolbar) Update(context *guigui.Context) error {
	t.panel.SetStyle(basicwidget.PanelStyleSide)
	t.panel.SetBorders(basicwidget.PanelBorder{
		Bottom: true,
	})
	t.content.SetFixedSize(context.Bounds(t).Size())
	t.panel.SetContent(&t.content)

	return nil
}

func (t *Toolbar) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &t.panel:
		return context.Bounds(t)
	}
	return image.Rectangle{}
}

func (t *Toolbar) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(t.DefaultWidget.Measure(context, constraints).X, 2*u)
}

type toolbarContent struct {
	guigui.DefaultWidget

	leftPanelButton  basicwidget.Button
	rightPanelButton basicwidget.Button

	layout layout.GridLayout
}

func (t *toolbarContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&t.leftPanelButton)
	appender.AppendChildWidget(&t.rightPanelButton)
}

func (t *toolbarContent) Update(context *guigui.Context) error {
	model := context.Model(t, modelKeyModel).(*Model)

	u := basicwidget.UnitSize(context)
	t.layout = layout.GridLayout{
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

	return nil
}

func (t *toolbarContent) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &t.leftPanelButton:
		return t.layout.CellBounds(0, 0)
	case &t.rightPanelButton:
		return t.layout.CellBounds(2, 0)
	}
	return image.Rectangle{}
}
