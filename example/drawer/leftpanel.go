// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type LeftPanel struct {
	guigui.DefaultWidget

	panel   basicwidget.Panel
	content leftPanelContent
}

func (l *LeftPanel) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&l.panel)
	l.panel.SetContent(&l.content)
}

func (l *LeftPanel) Build(context *guigui.Context) error {
	l.panel.SetStyle(basicwidget.PanelStyleSide)
	l.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetSize(&l.content, context.ActualSize(l), l)

	context.SetBounds(&l.panel, context.Bounds(l), l)
	return nil
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
	u := basicwidget.UnitSize(context)
	context.SetBounds(&l.text, context.Bounds(l).Inset(u/2), l)
	return nil
}
