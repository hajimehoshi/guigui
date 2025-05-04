// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package basicwidget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type TextList[T comparable] struct {
	guigui.DefaultWidget

	list                List[T]
	listItems           []ListItem[T]
	textListItems       []TextListItem[T]
	textListItemWidgets []textListItemWidget[T]

	listItemHeightPlus1 int
}

/*type TextListCallback struct {
	OnItemSelected    func(index int)
	OnItemEditStarted func(index int, str string) (from int)
	OnItemEditEnded   func(index int, str string)
	OnItemDropped     func(from, to int)
	OnContextMenu     func(index int, x, y int)
}*/

type TextListItem[T comparable] struct {
	Text      string
	DummyText string
	Color     color.Color
	Header    bool
	Disabled  bool
	Border    bool
	Draggable bool
	Tag       T
}

func (t *TextListItem[T]) selectable() bool {
	return !t.Header && !t.Disabled && !t.Border
}

/*func NewTextList(settings *model.Settings, callback *TextListCallback) *TextList {
	t := &TextList{
		settings: settings,
		callback: callback,
	}
	t.list = NewList(settings, &ListCallback{
		OnItemSelected: func(index int) {
			if callback != nil && callback.OnItemSelected != nil {
				callback.OnItemSelected(index)
			}
		},
		OnItemEditStarted: func(index int) {
			if index < 0 || index >= len(t.textListItems) {
				return
			}
			if !t.textListItems[index].selectable() {
				return
			}
			item, ok := t.textListItems[index].listItem.WidgetWithHeight.(*textListTextItem)
			if !ok {
				return
			}
			if callback != nil && callback.OnItemEditStarted != nil {
				item.edit(callback.OnItemEditStarted(index, item.textListItem.Text))
			}
		},
		OnItemDropped: func(from int, to int) {
			if callback != nil && callback.OnItemDropped != nil {
				callback.OnItemDropped(from, to)
			}
		},
		OnContextMenu: func(index int, x, y int) {
			if callback != nil && callback.OnContextMenu != nil {
				callback.OnContextMenu(index, x, y)
			}
		},
	})
	t.AddChild(t.list, &view.IdentityLayouter{})
	return t
}*/

func (t *TextList[T]) SetItemHeight(height int) {
	if t.listItemHeightPlus1 == height+1 {
		return
	}
	t.listItemHeightPlus1 = height + 1
	guigui.RequestRedraw(t)
}

func (t *TextList[T]) SetOnItemSelected(callback func(index int)) {
	t.list.SetOnItemSelected(callback)
}

func (t *TextList[T]) SetCheckmarkIndex(index int) {
	t.list.SetCheckmarkIndex(index)
}

func (t *TextList[T]) updateListItems() {
	t.textListItemWidgets = adjustSliceSize(t.textListItemWidgets, len(t.textListItems))
	t.listItems = adjustSliceSize(t.listItems, len(t.textListItems))

	for i, item := range t.textListItems {
		t.textListItemWidgets[i].setTextListItem(item)
		t.listItems[i] = t.textListItemWidgets[i].listItem()
	}
	t.list.SetItems(t.listItems)
}

func (t *TextList[T]) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetSize(&t.list, context.Size(t))

	t.updateListItems()

	// To use HasFocusedChildWidget correctly, create the tree first.
	appender.AppendChildWidgetWithPosition(&t.list, context.Position(t))

	focused := context.HasFocusedChildWidget(t)
	for i := range t.textListItemWidgets {
		item := &t.textListItemWidgets[i]
		item.text.SetBold(item.textListItem.Header || t.list.style == ListStyleSidebar && t.SelectedItemIndex() == i)
		switch {
		case t.list.style == ListStyleNormal && focused && t.list.SelectedItemIndex() == i && item.selectable():
			item.text.SetColor(DefaultActiveListItemTextColor(context))
		case t.list.style == ListStyleSidebar && t.list.SelectedItemIndex() == i && item.selectable():
			item.text.SetColor(DefaultActiveListItemTextColor(context))
		case t.list.style == ListStyleMenu && t.list.isHoveringVisible() && t.list.HoveredItemIndex(context) == i && item.selectable():
			item.text.SetColor(DefaultActiveListItemTextColor(context))
		case !item.selectable() && !item.textListItem.Header:
			item.text.SetColor(DefaultDisabledListItemTextColor(context))
		default:
			item.text.SetColor(item.textListItem.Color)
		}

		if t.listItemHeightPlus1 > 0 {
			context.SetSize(item, image.Pt(guigui.DefaultSize, t.listItemHeightPlus1-1))
		} else {
			context.SetSize(item, image.Pt(guigui.DefaultSize, guigui.DefaultSize))
		}
	}

	return nil
}

func (t *TextList[T]) SelectedItemIndex() int {
	return t.list.SelectedItemIndex()
}

func (t *TextList[T]) SelectedItem() (TextListItem[T], bool) {
	if t.list.SelectedItemIndex() < 0 || t.list.SelectedItemIndex() >= len(t.textListItemWidgets) {
		return TextListItem[T]{}, false
	}
	return t.textListItemWidgets[t.list.SelectedItemIndex()].textListItem, true
}

func (t *TextList[T]) ItemByIndex(index int) (TextListItem[T], bool) {
	if index < 0 || index >= len(t.textListItemWidgets) {
		return TextListItem[T]{}, false
	}
	return t.textListItemWidgets[index].textListItem, true
}

func (t *TextList[T]) SetItemsByStrings(strs []string) {
	items := make([]TextListItem[T], len(strs))
	for i, str := range strs {
		items[i].Text = str
	}
	t.SetItems(items)
}

func (t *TextList[T]) SetItems(items []TextListItem[T]) {
	t.textListItems = adjustSliceSize(t.textListItems, len(items))
	copy(t.textListItems, items)

	// Updating list items at Build might be too late, when the text list is not visible like a dropdown menu.
	// Update it here.
	t.updateListItems()
}

func (t *TextList[T]) ItemsCount() int {
	return len(t.textListItemWidgets)
}

func (t *TextList[T]) Tag(index int) any {
	return t.textListItemWidgets[index].textListItem.Tag
}

func (t *TextList[T]) SelectItemByIndex(index int) {
	t.list.SelectItemByIndex(index)
}

func (t *TextList[T]) SelectItemByTag(tag T) {
	t.list.SelectItemByTag(tag)
}

func (t *TextList[T]) JumpToItemIndex(index int) {
	t.list.JumpToItemIndex(index)
}

func (t *TextList[T]) SetStyle(style ListStyle) {
	t.list.SetStyle(style)
}

func (t *TextList[T]) SetItemString(str string, index int) {
	t.textListItemWidgets[index].textListItem.Text = str
}

func (t *TextList[T]) DefaultSize(context *guigui.Context) image.Point {
	return t.list.DefaultSize(context)
}

type textListItemWidget[T comparable] struct {
	guigui.DefaultWidget

	textListItem TextListItem[T]

	text Text
}

func (t *textListItemWidget[T]) setTextListItem(textListItem TextListItem[T]) {
	t.textListItem = textListItem
	t.text.SetValue(t.textString())
}

func (t *textListItemWidget[T]) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	p := context.Position(t)
	if t.textListItem.Header {
		p.X += UnitSize(context) / 2
		context.SetSize(&t.text, context.Size(t).Add(image.Pt(-UnitSize(context), 0)))
	} else {
		context.SetSize(&t.text, context.Size(t))
	}
	t.text.SetValue(t.textString())
	t.text.SetVerticalAlign(VerticalAlignMiddle)
	appender.AppendChildWidgetWithPosition(&t.text, p)

	return nil
}

func (t *textListItemWidget[T]) textString() string {
	if t.textListItem.DummyText != "" {
		return t.textListItem.DummyText
	}
	return t.textListItem.Text
}

func (t *textListItemWidget[T]) Draw(context *guigui.Context, dst *ebiten.Image) {
	if t.textListItem.Border {
		p := context.Position(t)
		s := context.Size(t)
		x0 := float32(p.X)
		x1 := float32(p.X + s.X)
		y := float32(p.Y) + float32(s.Y)/2
		width := float32(1 * context.Scale())
		vector.StrokeLine(dst, x0, y, x1, y, width, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.8), false)
		return
	}
	if t.textListItem.Header {
		bounds := context.Bounds(t)
		draw.DrawRoundedRect(context, dst, bounds, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.6), RoundedCornerRadius(context))
	}
}

func (t *textListItemWidget[T]) DefaultSize(context *guigui.Context) image.Point {
	// Assume that every item can use a bold font.
	w := t.text.boldTextSize(context).X
	if t.textListItem.Border {
		return image.Pt(w, UnitSize(context)/2)
	}
	return image.Pt(w, int(LineHeight(context)))
}

/*func (t *textListItemWidget[T]) index() int {
	for i, tt := range t.textList.textListItemWidgets {
		if tt == t {
			return i
		}
	}
	return -1
}*/

func (t *textListItemWidget[T]) selectable() bool {
	return t.textListItem.selectable() && !t.textListItem.Border
}

func (t *textListItemWidget[T]) listItem() ListItem[T] {
	return ListItem[T]{
		Content:    t,
		Selectable: t.selectable(),
		Draggable:  t.textListItem.Draggable,
		Tag:        t.textListItem.Tag,
	}
}

/*func (t *textListTextItem) edit(from int) {
	t.label.Hide()
	t0 := t.textListItem.Text[:from]
	var l0 *Label
	if t0 != "" {
		l0 = NewLabel(t.settings)
		l0.SetText(t0)
		t.AddChild(l0, &view.IdentityLayouter{})
	}
	var tf *TextInput
	tf = NewTextInput(t.settings, &TextInputCallback{
		OnTextUpdated: func(value string) {
			t.textListItem.Text = t0 + value
		},
		OnTextConfirmed: func(value string) {
			t.textListItem.Text = t0 + value

			t.label.SetText(t.labelText())
			if l0 != nil {
				l0.RemoveSelf()
			}
			tf.RemoveSelf()
			t.label.Show()

			if t.textList.callback != nil && t.textList.callback.OnItemEditEnded != nil {
				t.textList.callback.OnItemEditEnded(t.index(), t0+value)
			}
		},
	})
	tf.SetText(t.textListItem.Text[from:])
	tf.SetHorizontalAlign(HorizontalAlignStart)
	t.AddChild(tf, view.LayoutFunc(func(args view.WidgetArgs) image.Rectangle {
		bounds := args.Bounds
		if l0 != nil {
			bounds.Min.X += l0.Width(args.Scale)
		}
		return bounds
	}))
	tf.SelectAll()
	tf.Focus()
}*/
