// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/guigui"
)

const (
	popupMenuEventItemSelected = "itemSelected"
)

type PopupMenuItem[T comparable] struct {
	Text         string
	TextColor    color.Color
	Header       bool
	Content      guigui.Widget
	Unselectable bool
	Border       bool
	ID           T
}

type PopupMenu[T comparable] struct {
	guigui.DefaultWidget

	list  List[T]
	popup Popup
}

func (p *PopupMenu[T]) SetOnItemSelected(f func(index int)) {
	guigui.RegisterEventHandler(p, popupMenuEventItemSelected, f)
}

func (p *PopupMenu[T]) SetCheckmarkIndex(index int) {
	p.list.SetCheckmarkIndex(index)
}

func (p *PopupMenu[T]) IsWidgetOrBackgroundHitAtCursor(context *guigui.Context, widget guigui.Widget) bool {
	return p.popup.IsWidgetOrBackgroundHitAtCursor(context, widget)
}


func (p *PopupMenu[T]) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&p.popup)
}

func (p *PopupMenu[T]) Build(context *guigui.Context) error {
	p.list.SetStyle(ListStyleMenu)
	p.list.list.SetOnItemSelected(func(index int) {
		p.popup.Close()
		guigui.InvokeEventHandler(p, popupMenuEventItemSelected, index)
	})

	p.popup.SetContent(&p.list)
	p.popup.SetCloseByClickingOutside(true)
	bounds := p.contentBounds(context)
	context.SetSize(&p.list, bounds.Size(), p)
	context.SetBounds(&p.popup, bounds, p)

	return nil
}

func (p *PopupMenu[T]) contentBounds(context *guigui.Context) image.Rectangle {
	pos := context.Position(p)
	// List size can dynamically change based on the items. Use the default size.
	s := p.list.DefaultSize(context)
	s.Y = min(s.Y, 24*UnitSize(context))
	r := image.Rectangle{
		Min: pos,
		Max: pos.Add(s),
	}
	as := context.AppSize()
	if r.Max.X > as.X {
		r.Min.X = as.X - s.X
		r.Max.X = as.X
	}
	if r.Min.X < 0 {
		r.Min.X = 0
		r.Max.X = s.X
	}
	if r.Max.Y > as.Y {
		r.Min.Y = as.Y - s.Y
		r.Max.Y = as.Y
	}
	if r.Min.Y < 0 {
		r.Min.Y = 0
		r.Max.Y = s.Y
	}
	return r
}

func (p *PopupMenu[T]) Open(context *guigui.Context) {
	p.popup.Open(context)
}

func (p *PopupMenu[T]) Close() {
	p.popup.Close()
}

func (p *PopupMenu[T]) IsOpen() bool {
	return p.popup.IsOpen()
}

func (p *PopupMenu[T]) SetItems(items []PopupMenuItem[T]) {
	var listItems []ListItem[T]
	for _, item := range items {
		listItems = append(listItems, ListItem[T]{
			Text:         item.Text,
			TextColor:    item.TextColor,
			Header:       item.Header,
			Content:      item.Content,
			Unselectable: item.Unselectable,
			Border:       item.Border,
			ID:           item.ID,
		})
	}
	p.list.SetItems(listItems)
}

func (p *PopupMenu[T]) SetItemsByStrings(items []string) {
	p.list.SetItemsByStrings(items)
}

func (p *PopupMenu[T]) SelectedItem() (PopupMenuItem[T], bool) {
	listItem, ok := p.list.SelectedItem()
	if !ok {
		return PopupMenuItem[T]{}, false
	}
	return PopupMenuItem[T]{
		Text:         listItem.Text,
		TextColor:    listItem.TextColor,
		Header:       listItem.Header,
		Content:      listItem.Content,
		Unselectable: listItem.Unselectable,
		Border:       listItem.Border,
		ID:           listItem.ID,
	}, true
}

func (p *PopupMenu[T]) ItemByIndex(index int) (PopupMenuItem[T], bool) {
	listItem, ok := p.list.ItemByIndex(index)
	if !ok {
		return PopupMenuItem[T]{}, false
	}
	return PopupMenuItem[T]{
		Text:         listItem.Text,
		TextColor:    listItem.TextColor,
		Header:       listItem.Header,
		Content:      listItem.Content,
		Unselectable: listItem.Unselectable,
		Border:       listItem.Border,
		ID:           listItem.ID,
	}, true
}

func (p *PopupMenu[T]) SelectedItemIndex() int {
	return p.list.SelectedItemIndex()
}

func (p *PopupMenu[T]) SelectItemByIndex(index int) {
	p.list.SelectItemByIndex(index)
}

func (p *PopupMenu[T]) SelectItemByID(id T) {
	p.list.SelectItemByID(id)
}

func (p *PopupMenu[T]) ItemTextColor(context *guigui.Context, index int) color.Color {
	return p.list.ItemTextColor(context, index)
}
