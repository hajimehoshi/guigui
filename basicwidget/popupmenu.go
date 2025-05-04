// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/guigui"
)

type PopupMenuItem[T comparable] struct {
	Text     string
	Color    color.Color
	Header   bool
	Disabled bool
	Border   bool
	Tag      T
}

type PopupMenu[T comparable] struct {
	guigui.DefaultWidget

	textList TextList[T]
	popup    Popup

	onClosed func(index int)
}

func (p *PopupMenu[T]) SetOnClosed(f func(index int)) {
	p.onClosed = f
}

func (p *PopupMenu[T]) SetCheckmarkIndex(index int) {
	p.textList.SetCheckmarkIndex(index)
}

func (p *PopupMenu[T]) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	p.textList.SetStyle(ListStyleMenu)
	p.textList.list.SetOnItemSelected(func(index int) {
		p.popup.Close()
		if p.onClosed != nil {
			p.onClosed(index)
		}
	})

	p.popup.SetContent(&p.textList)
	p.popup.SetCloseByClickingOutside(true)
	p.popup.SetOnClosed(func(reason PopupClosedReason) {
		if p.onClosed != nil {
			p.onClosed(p.textList.SelectedItemIndex())
		}
	})
	bounds := p.contentBounds(context)
	context.SetSize(&p.textList, bounds.Size())
	appender.AppendChildWidgetWithBounds(&p.popup, bounds)

	return nil
}

func (p *PopupMenu[T]) contentBounds(context *guigui.Context) image.Rectangle {
	pos := context.Position(p)
	// textList's size is updated at Build so do not call guigui.Size to determine the content size.
	// Call DefaultSize instead.
	s := p.textList.DefaultSize(context)
	if s.Y > 24*UnitSize(context) {
		s.Y = 24 * UnitSize(context)
	}
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
	var textListItems []TextListItem[T]
	for _, item := range items {
		textListItems = append(textListItems, TextListItem[T]{
			Text:     item.Text,
			Color:    item.Color,
			Header:   item.Header,
			Disabled: item.Disabled,
			Border:   item.Border,
			Tag:      item.Tag,
		})
	}
	p.textList.SetItems(textListItems)
}

func (p *PopupMenu[T]) SetItemsByStrings(items []string) {
	p.textList.SetItemsByStrings(items)
}

func (p *PopupMenu[T]) SelectedItem() (PopupMenuItem[T], bool) {
	textListItem, ok := p.textList.SelectedItem()
	if !ok {
		return PopupMenuItem[T]{}, false
	}
	return PopupMenuItem[T]{
		Text:     textListItem.Text,
		Color:    textListItem.Color,
		Header:   textListItem.Header,
		Disabled: textListItem.Disabled,
		Border:   textListItem.Border,
		Tag:      textListItem.Tag,
	}, true
}

func (p *PopupMenu[T]) ItemByIndex(index int) (PopupMenuItem[T], bool) {
	textListItem, ok := p.textList.ItemByIndex(index)
	if !ok {
		return PopupMenuItem[T]{}, false
	}
	return PopupMenuItem[T]{
		Text:     textListItem.Text,
		Color:    textListItem.Color,
		Header:   textListItem.Header,
		Disabled: textListItem.Disabled,
		Border:   textListItem.Border,
		Tag:      textListItem.Tag,
	}, true
}

func (p *PopupMenu[T]) SelectedItemIndex() int {
	return p.textList.SelectedItemIndex()
}

func (p *PopupMenu[T]) SelectItemByIndex(index int) {
	p.textList.SelectItemByIndex(index)
}

func (p *PopupMenu[T]) SelectItemByTag(tag T) {
	p.textList.SelectItemByTag(tag)
}
