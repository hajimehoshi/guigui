// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type ContentPanel struct {
	guigui.DefaultWidget

	panel   basicwidget.Panel
	content guigui.WidgetWithSize[*contentPanelContent]
}

func (c *ContentPanel) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&c.panel)
}

func (c *ContentPanel) Update(context *guigui.Context) error {
	c.content.SetFixedSize(context.Bounds(c).Size())
	c.panel.SetContent(&c.content)
	return nil
}

func (c *ContentPanel) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &c.panel:
		return context.Bounds(c)
	}
	return image.Rectangle{}
}

type contentPanelContent struct {
	guigui.DefaultWidget

	text basicwidget.Text
}

func (c *contentPanelContent) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&c.text)
}

func (c *contentPanelContent) Update(context *guigui.Context) error {
	c.text.SetValue("Content panel: " + dummyText)
	c.text.SetAutoWrap(true)
	c.text.SetSelectable(true)
	return nil
}

func (c *contentPanelContent) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &c.text:
		u := basicwidget.UnitSize(context)
		return context.Bounds(c).Inset(u / 2)
	}
	return image.Rectangle{}
}
