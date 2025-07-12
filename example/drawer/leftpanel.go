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

func (l *LeftPanel) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	l.panel.SetStyle(basicwidget.PanelStyleSide)
	l.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetSize(&l.content, context.ActualSize(l), l)
	l.panel.SetContent(&l.content)

	appender.AppendChildWidgetWithBounds(&l.panel, context.Bounds(l))
	return nil
}

type leftPanelContent struct {
	guigui.DefaultWidget

	text basicwidget.Text
}

func (l *leftPanelContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	l.text.SetValue("Left panel: " + dummyText)
	l.text.SetAutoWrap(true)
	l.text.SetSelectable(true)
	u := basicwidget.UnitSize(context)
	appender.AppendChildWidgetWithBounds(&l.text, context.Bounds(l).Inset(u/2))
	return nil
}
