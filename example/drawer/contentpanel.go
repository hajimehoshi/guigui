// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type ContentPanel struct {
	guigui.DefaultWidget

	panel   basicwidget.Panel
	content contentPanelContent
}

func (c *ContentPanel) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&c.panel)
}

func (c *ContentPanel) Build(context *guigui.Context) error {
	context.SetSize(&c.content, context.ActualSize(c), c)
	c.panel.SetContent(&c.content)

	context.SetBounds(&c.panel, context.Bounds(c), c)
	return nil
}

type contentPanelContent struct {
	guigui.DefaultWidget

	text basicwidget.Text
}

func (c *contentPanelContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&c.text)
}

func (c *contentPanelContent) Build(context *guigui.Context) error {
	c.text.SetValue("Content panel: " + dummyText)
	c.text.SetAutoWrap(true)
	c.text.SetSelectable(true)
	u := basicwidget.UnitSize(context)
	context.SetBounds(&c.text, context.Bounds(c).Inset(u/2), c)
	return nil
}
