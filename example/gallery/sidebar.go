// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	sidebar        basicwidget.Sidebar
	sidebarContent sidebarContent
}

func (s *Sidebar) SetModel(model *Model) {
	s.sidebarContent.SetModel(model)
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetSize(&s.sidebarContent, context.Size(s))
	s.sidebar.SetContent(&s.sidebarContent)

	appender.AppendChildWidgetWithBounds(&s.sidebar, context.Bounds(s))

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.TextList[string]

	model *Model
}

func (s *sidebarContent) SetModel(model *Model) {
	s.model = model
}

func (s *sidebarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.TextListItem[string]{
		{
			Text: "Settings",
			Tag:  "settings",
		},
		{
			Text: "Basic",
			Tag:  "basic",
		},
		{
			Text: "Buttons",
			Tag:  "buttons",
		},
		{
			Text: "Texts",
			Tag:  "texts",
		},
		{
			Text: "Text Inputs",
			Tag:  "textinputs",
		},
		{
			Text: "Number Inputs",
			Tag:  "numberinputs",
		},
		{
			Text: "Lists",
			Tag:  "lists",
		},
		{
			Text: "Popups",
			Tag:  "popups",
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByTag(s.model.Mode())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			s.model.SetMode("")
			return
		}
		s.model.SetMode(item.Tag)
	})

	appender.AppendChildWidgetWithBounds(&s.list, context.Bounds(s))

	return nil
}
