// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package main

import (
	"image"
	"sync"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	sidebar         basicwidget.Sidebar
	list            basicwidget.List
	listItemWidgets []basicwidget.ListItem

	initOnce sync.Once
}

func sidebarWidth(context *guigui.Context) int {
	return 8 * basicwidget.UnitSize(context)
}

func (s *Sidebar) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	_, h := s.Size(context)
	s.sidebar.SetSize(context, sidebarWidth(context), h)
	s.sidebar.SetContent(context, func(context *guigui.Context, childAppender *basicwidget.ContainerChildWidgetAppender, offsetX, offsetY float64) {
		s.list.SetWidth(sidebarWidth(context))
		s.list.SetHeight(h)
		guigui.SetPosition(&s.list, guigui.Position(s).Add(image.Pt(int(offsetX), int(offsetY))))
		childAppender.AppendChildWidget(&s.list)
	})
	guigui.SetPosition(&s.sidebar, guigui.Position(s))
	appender.AppendChildWidget(&s.sidebar)

	s.list.SetStyle(basicwidget.ListStyleSidebar)
	if len(s.listItemWidgets) == 0 {
		{
			var t basicwidget.Text
			t.SetText("Settings")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetHeight(basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "settings",
			})
		}
		{
			var t basicwidget.Text
			t.SetText("Basic")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetHeight(basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "basic",
			})
		}
		{
			var t basicwidget.Text
			t.SetText("Buttons")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetHeight(basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "buttons",
			})
		}
		{
			var t basicwidget.Text
			t.SetText("Lists")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetHeight(basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "lists",
			})
		}
		{
			var t basicwidget.Text
			t.SetText("Popups")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetHeight(basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "popups",
			})
		}
		{
			var t basicwidget.Text
			t.SetText("Tables")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetHeight(basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "tables",
			})
		}
	}
	s.list.SetItems(s.listItemWidgets)

	s.initOnce.Do(func() {
		s.list.SetSelectedItemIndex(0)
	})
}

func (s *Sidebar) Update(context *guigui.Context) error {
	for i, w := range s.listItemWidgets {
		t := w.Content.(*basicwidget.Text)
		if s.list.SelectedItemIndex() == i {
			t.SetColor(basicwidget.DefaultActiveListItemTextColor(context))
		} else {
			t.SetColor(basicwidget.DefaultTextColor(context))
		}
	}
	return nil
}

func (s *Sidebar) Size(context *guigui.Context) (int, int) {
	_, h := guigui.Parent(s).Size(context)
	return sidebarWidth(context), h
}

func (s *Sidebar) SelectedItemTag() string {
	item, ok := s.list.SelectedItem()
	if !ok {
		return ""
	}
	return item.Tag.(string)
}

func (s *Sidebar) SetSelectedItemIndex(index int) {
	s.list.SetSelectedItemIndex(index)
}
