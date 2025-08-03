// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type List[T comparable] struct {
	guigui.DefaultWidget

	list            baseList[T]
	baseListItems   []baseListItem[T]
	listItems       []ListItem[T]
	listItemWidgets []listItemWidget[T]

	listItemHeightPlus1 int
}

type ListItem[T comparable] struct {
	Text         string
	TextColor    color.Color
	Header       bool
	Content      guigui.Widget
	Unselectable bool
	Border       bool
	Movable      bool
	Value        T
}

func (l *ListItem[T]) selectable() bool {
	return !l.Header && !l.Unselectable && !l.Border
}

func (l *List[T]) SetStripeVisible(visible bool) {
	l.list.SetStripeVisible(visible)
}

func (l *List[T]) SetItemHeight(height int) {
	if l.listItemHeightPlus1 == height+1 {
		return
	}
	l.listItemHeightPlus1 = height + 1
	guigui.RequestRedraw(l)
}

func (l *List[T]) SetOnItemSelected(f func(index int)) {
	l.list.SetOnItemSelected(f)
}

func (l *List[T]) SetOnItemsMoved(f func(from, count, to int)) {
	l.list.SetOnItemsMoved(f)
}

func (l *List[T]) SetCheckmarkIndex(index int) {
	l.list.SetCheckmarkIndex(index)
}

func (l *List[T]) SetHeaderHeight(height int) {
	l.list.SetHeaderHeight(height)
}

func (l *List[T]) SetFooterHeight(height int) {
	l.list.SetFooterHeight(height)
}

func (l *List[T]) updateListItems() {
	l.listItemWidgets = adjustSliceSize(l.listItemWidgets, len(l.listItems))
	l.baseListItems = adjustSliceSize(l.baseListItems, len(l.listItems))

	for i, item := range l.listItems {
		l.listItemWidgets[i].setListItem(item)
		l.listItemWidgets[i].setHeight(l.listItemHeightPlus1 - 1)
		l.listItemWidgets[i].setStyle(l.list.style)
		l.baseListItems[i] = l.listItemWidgets[i].listItem()
	}
	l.list.SetItems(l.baseListItems)
}

func (l *List[T]) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&l.list)
}

func (l *List[T]) Build(context *guigui.Context) error {
	context.SetSize(&l.list, context.ActualSize(l), l)

	l.updateListItems()

	context.SetPosition(&l.list, context.Position(l))

	itemSize := image.Pt(guigui.AutoSize, guigui.AutoSize)
	if l.list.style != ListStyleMenu {
		itemSize.X = context.ActualSize(l).X - 2*listItemPadding(context)
	}
	if l.listItemHeightPlus1 > 0 {
		itemSize.Y = l.listItemHeightPlus1 - 1
	}
	for i := range l.listItemWidgets {
		item := &l.listItemWidgets[i]
		item.text.SetBold(item.item.Header || l.list.style == ListStyleSidebar && l.SelectedItemIndex() == i)
		item.text.SetColor(l.ItemTextColor(context, i))
		context.SetSize(item, itemSize, l)
	}

	return nil
}

func (l *List[T]) ItemTextColor(context *guigui.Context, index int) color.Color {
	item := &l.listItemWidgets[index]
	switch {
	case l.list.style == ListStyleNormal && l.list.SelectedItemIndex() == index && item.selectable() && context.IsEnabled(item):
		return DefaultActiveListItemTextColor(context)
	case l.list.style == ListStyleSidebar && l.list.SelectedItemIndex() == index && item.selectable() && context.IsEnabled(item):
		return DefaultActiveListItemTextColor(context)
	case l.list.style == ListStyleMenu && l.list.isHoveringVisible() && l.list.hoveredItemIndex(context) == index && item.selectable() && context.IsEnabled(item):
		return DefaultActiveListItemTextColor(context)
	case item.item.TextColor != nil:
		return item.item.TextColor
	default:
		return draw.TextColor(context.ColorMode(), context.IsEnabled(item))
	}
}

func (l *List[T]) SelectedItemIndex() int {
	return l.list.SelectedItemIndex()
}

func (l *List[T]) SelectedItem() (ListItem[T], bool) {
	if l.list.SelectedItemIndex() < 0 || l.list.SelectedItemIndex() >= len(l.listItemWidgets) {
		return ListItem[T]{}, false
	}
	return l.listItemWidgets[l.list.SelectedItemIndex()].item, true
}

func (l *List[T]) ItemByIndex(index int) (ListItem[T], bool) {
	if index < 0 || index >= len(l.listItemWidgets) {
		return ListItem[T]{}, false
	}
	return l.listItemWidgets[index].item, true
}

func (l *List[T]) SetItemsByStrings(strs []string) {
	items := make([]ListItem[T], len(strs))
	for i, str := range strs {
		items[i].Text = str
	}
	l.SetItems(items)
}

func (l *List[T]) SetItems(items []ListItem[T]) {
	l.listItems = adjustSliceSize(l.listItems, len(items))
	copy(l.listItems, items)

	// Updating list items at Build might be too late, when the text list is not visible like a dropdown menu.
	// Update it here.
	l.updateListItems()
}

func (l *List[T]) ItemsCount() int {
	return len(l.listItemWidgets)
}

func (l *List[T]) ID(index int) any {
	return l.listItemWidgets[index].item.Value
}

func (l *List[T]) SelectItemByIndex(index int) {
	l.list.SelectItemByIndex(index)
}

func (l *List[T]) SelectItemByValue(value T) {
	l.list.SelectItemByValue(value)
}

func (l *List[T]) JumpToItemIndex(index int) {
	l.list.JumpToItemIndex(index)
}

func (l *List[T]) SetStyle(style ListStyle) {
	l.list.SetStyle(style)
}

func (l *List[T]) SetItemString(str string, index int) {
	l.listItemWidgets[index].item.Text = str
}

func (l *List[T]) DefaultSize(context *guigui.Context) image.Point {
	return l.list.DefaultSize(context)
}

type listItemWidget[T comparable] struct {
	guigui.DefaultWidget

	item ListItem[T]

	text        Text
	heightPlus1 int
	style       ListStyle
}

func (l *listItemWidget[T]) setListItem(listItem ListItem[T]) {
	l.item = listItem
	l.text.SetValue(listItem.Text)
}

func (l *listItemWidget[T]) setHeight(height int) {
	if l.heightPlus1 == height+1 {
		return
	}
	l.heightPlus1 = height + 1
}

func (l *listItemWidget[T]) setStyle(style ListStyle) {
	l.style = style
}

func (l *listItemWidget[T]) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	if l.item.Content != nil {
		appender.AppendChildWidget(l.item.Content)
	}
	appender.AppendChildWidget(&l.text)
}

func (l *listItemWidget[T]) Build(context *guigui.Context) error {
	if l.item.Content != nil {
		s := image.Pt(guigui.AutoSize, guigui.AutoSize)
		if l.style != ListStyleMenu {
			s.X = context.ActualSize(l).X
		}
		if l.heightPlus1 > 0 {
			s.Y = l.heightPlus1 - 1
		}
		context.SetSize(l.item.Content, s, l)
		context.SetPosition(l.item.Content, context.Bounds(l).Min)
	}

	l.text.SetValue(l.item.Text)
	l.text.SetVerticalAlign(VerticalAlignMiddle)
	s := image.Pt(guigui.AutoSize, guigui.AutoSize)
	if l.style != ListStyleMenu {
		s.X = context.ActualSize(l).X
	}
	if l.heightPlus1 > 0 {
		s.Y = l.heightPlus1 - 1
	}
	context.SetSize(&l.text, s, l)
	context.SetPosition(&l.text, context.Bounds(l).Min)

	return nil
}

func (l *listItemWidget[T]) Draw(context *guigui.Context, dst *ebiten.Image) {
	if l.item.Border {
		p := context.Position(l)
		s := context.ActualSize(l)
		x0 := float32(p.X)
		x1 := float32(p.X + s.X)
		y := float32(p.Y) + float32(s.Y)/2
		width := float32(1 * context.Scale())
		vector.StrokeLine(dst, x0, y, x1, y, width, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.8), false)
		return
	}
	/*if l.item.Header {
		bounds := context.Bounds(l)
		draw.DrawRoundedRect(context, dst, bounds, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.8), RoundedCornerRadius(context))
	}*/
}

func (l *listItemWidget[T]) DefaultSize(context *guigui.Context) image.Point {
	var w, h int
	if l.item.Content != nil {
		s := context.ActualSize(l.item.Content)
		w, h = s.X, s.Y
	}

	// Assume that every item can use a bold font.
	w = max(w, l.text.boldTextSize(context, 0).X)
	h = max(h, int(LineHeight(context)))
	if l.item.Border {
		h = UnitSize(context) / 2
	} else if l.item.Header {
		h = UnitSize(context) * 3 / 2
	}
	return image.Pt(w, h)
}

func (l *listItemWidget[T]) selectable() bool {
	return l.item.selectable() && !l.item.Border
}

func (l *listItemWidget[T]) listItem() baseListItem[T] {
	return baseListItem[T]{
		Content:    l,
		Selectable: l.selectable(),
		Movable:    l.item.Movable,
		Value:      l.item.Value,
	}
}
