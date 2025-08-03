// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	panel        basicwidget.Panel
	panelContent sidebarContent
}

func (s *Sidebar) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.panel)
}

func (s *Sidebar) Build(context *guigui.Context) error {
	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetSize(&s.panelContent, context.ActualSize(s), s)
	s.panel.SetContent(&s.panelContent)

	context.SetBounds(&s.panel, context.Bounds(s), s)

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.List[string]
}

func (s *sidebarContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.list)
}

func (s *sidebarContent) Build(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*Model)

	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.ListItem[string]{
		{
			Text:  "Settings",
			Value: "settings",
		},
		{
			Text:  "Basic",
			Value: "basic",
		},
		{
			Text:  "Buttons",
			Value: "buttons",
		},
		{
			Text:  "Texts",
			Value: "texts",
		},
		{
			Text:  "Text Inputs",
			Value: "textinputs",
		},
		{
			Text:  "Number Inputs",
			Value: "numberinputs",
		},
		{
			Text:  "Lists",
			Value: "lists",
		},
		{
			Text:  "Tables",
			Value: "tables",
		},
		{
			Text:  "Popups",
			Value: "popups",
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByValue(model.Mode())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			model.SetMode("")
			return
		}
		model.SetMode(item.Value)
	})

	context.SetBounds(&s.list, context.Bounds(s), s)

	return nil
}
