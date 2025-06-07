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

func (s *Sidebar) SetModel(model *Model) {
	s.panelContent.SetModel(model)
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetSize(&s.panelContent, context.ActualSize(s))
	s.panel.SetContent(&s.panelContent)

	appender.AppendChildWidgetWithBounds(&s.panel, context.Bounds(s))

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.List[string]

	model *Model
}

func (s *sidebarContent) SetModel(model *Model) {
	s.model = model
}

func (s *sidebarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.ListItem[string]{
		{
			Text: "Settings",
			ID:   "settings",
		},
		{
			Text: "Basic",
			ID:   "basic",
		},
		{
			Text: "Buttons",
			ID:   "buttons",
		},
		{
			Text: "Texts",
			ID:   "texts",
		},
		{
			Text: "Text Inputs",
			ID:   "textinputs",
		},
		{
			Text: "Number Inputs",
			ID:   "numberinputs",
		},
		{
			Text: "Lists",
			ID:   "lists",
		},
		{
			Text: "Tables",
			ID:   "tables",
		},
		{
			Text: "Popups",
			ID:   "popups",
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByID(s.model.Mode())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			s.model.SetMode("")
			return
		}
		s.model.SetMode(item.ID)
	})

	appender.AppendChildWidgetWithBounds(&s.list, context.Bounds(s))

	return nil
}
