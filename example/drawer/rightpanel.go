// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type RightPanel struct {
	guigui.DefaultWidget

	panel   basicwidget.Panel
	content guigui.WidgetWithSize[*rightPanelContent]
}

func (r *RightPanel) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&r.panel)
}

func (r *RightPanel) Update(context *guigui.Context) error {
	r.panel.SetStyle(basicwidget.PanelStyleSide)
	r.panel.SetBorders(basicwidget.PanelBorder{
		Start: true,
	})
	r.content.SetFixedSize(context.Bounds(r).Size())
	r.panel.SetContent(&r.content)

	return nil
}

func (r *RightPanel) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.panel:
		return context.Bounds(r)
	}
	return image.Rectangle{}
}

type rightPanelContent struct {
	guigui.DefaultWidget

	text basicwidget.Text
}

func (r *rightPanelContent) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&r.text)
}

func (r *rightPanelContent) Update(context *guigui.Context) error {
	r.text.SetValue("Right panel: " + dummyText)
	r.text.SetAutoWrap(true)
	r.text.SetSelectable(true)
	return nil
}

func (r *rightPanelContent) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.text:
		u := basicwidget.UnitSize(context)
		return context.Bounds(r).Inset(u / 2)
	}
	return image.Rectangle{}
}
