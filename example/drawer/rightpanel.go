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

func (r *RightPanel) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.panel)
	r.panel.SetContent(&r.content)
}

func (r *RightPanel) Build(context *guigui.Context) error {
	r.panel.SetStyle(basicwidget.PanelStyleSide)
	r.panel.SetBorder(basicwidget.PanelBorder{
		Start: true,
	})
	context.SetSize(&r.content, context.ActualSize(r), r)

	context.SetBounds(&r.panel, context.Bounds(r), r)
	return nil
}

type rightPanelContent struct {
	guigui.DefaultWidget

	text basicwidget.Text
}

func (r *rightPanelContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.text)
}

func (r *rightPanelContent) Build(context *guigui.Context) error {
	r.text.SetValue("Right panel: " + dummyText)
	r.text.SetAutoWrap(true)
	r.text.SetSelectable(true)
	u := basicwidget.UnitSize(context)
	context.SetBounds(&r.text, context.Bounds(r).Inset(u/2), r)
	return nil
}
