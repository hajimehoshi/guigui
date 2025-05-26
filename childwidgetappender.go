// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package guigui

import "image"

type ChildWidgetAppender struct {
	app    *app
	widget Widget
}

func (c *ChildWidgetAppender) AppendChildWidget(widget Widget) {
	c.appendChildWidget(widget)
}

func (c *ChildWidgetAppender) AppendChildWidgetWithPosition(widget Widget, position image.Point) {
	c.app.context.SetPosition(widget, position)
	c.appendChildWidget(widget)
}

func (c *ChildWidgetAppender) AppendChildWidgetWithBounds(widget Widget, bounds image.Rectangle) {
	c.app.context.SetPosition(widget, bounds.Min)
	c.app.context.SetSize(widget, bounds.Size())
	c.appendChildWidget(widget)
}

func (c *ChildWidgetAppender) appendChildWidget(widget Widget) {
	widgetState := widget.widgetState()
	widgetState.parent = c.widget
	widgetState.builtAt = c.app.buildCount
	cWidgetState := c.widget.widgetState()
	cWidgetState.children = append(cWidgetState.children, widget)
}
