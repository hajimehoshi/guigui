// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type LeftPanel struct {
	guigui.DefaultWidget

	panel   basicwidget.Panel
	content guigui.WidgetWithSize[*leftPanelContent]
}

func (l *LeftPanel) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&l.panel)
}

func (l *LeftPanel) Build(context *guigui.Context) error {
	l.panel.SetStyle(basicwidget.PanelStyleSide)
	l.panel.SetBorders(basicwidget.PanelBorder{
		End: true,
	})
	l.content.SetFixedSize(context.Bounds(l).Size())
	l.panel.SetContent(&l.content)

	return nil
}

func (l *LeftPanel) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &l.panel:
		return context.Bounds(l)
	}
	return image.Rectangle{}
}

type leftPanelContent struct {
	guigui.DefaultWidget

	text basicwidget.Text
}

func (l *leftPanelContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&l.text)
}

func (l *leftPanelContent) Build(context *guigui.Context) error {
	l.text.SetValue("Left panel: " + dummyText)
	l.text.SetAutoWrap(true)
	l.text.SetSelectable(true)
	return nil
}

func (l *leftPanelContent) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &l.text:
		u := basicwidget.UnitSize(context)
		return context.Bounds(l).Inset(u / 2)
	}
	return image.Rectangle{}
}
