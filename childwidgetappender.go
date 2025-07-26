// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package guigui

type ChildWidgetAppender struct {
	app    *app
	widget Widget
}

func (c *ChildWidgetAppender) AppendChildWidget(widget Widget) {
	widgetState := widget.widgetState()
	widgetState.parent = c.widget
	widgetState.builtAt = c.app.buildCount
	cWidgetState := c.widget.widgetState()
	cWidgetState.children = append(cWidgetState.children, widget)
}
