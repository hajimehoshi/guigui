// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type RightPanel struct {
	guigui.DefaultWidget

	panel   basicwidget.Panel
	content rightPanelContent
}

func (r *RightPanel) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.panel.SetStyle(basicwidget.PanelStyleSide)
	r.panel.SetBorder(basicwidget.PanelBorder{
		Start: true,
	})
	context.SetSize(&r.content, context.ActualSize(r))
	r.panel.SetContent(&r.content)

	appender.AppendChildWidgetWithBounds(&r.panel, context.Bounds(r))
	return nil
}

type rightPanelContent struct {
	guigui.DefaultWidget

	text basicwidget.Text
}

func (r *rightPanelContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.text.SetValue("Right panel: " + dummyText)
	r.text.SetAutoWrap(true)
	r.text.SetSelectable(true)
	u := basicwidget.UnitSize(context)
	appender.AppendChildWidgetWithBounds(&r.text, context.Bounds(r).Inset(u/2))
	return nil
}
