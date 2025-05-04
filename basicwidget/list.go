// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"image/color"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type ListStyle int

const (
	ListStyleNormal ListStyle = iota
	ListStyleSidebar
	ListStyleMenu
)

type ListItem[T comparable] struct {
	Content    guigui.Widget
	Selectable bool
	Draggable  bool
	Tag        T
}

func (l ListItem[T]) tag() T {
	return l.Tag
}

func DefaultActiveListItemTextColor(context *guigui.Context) color.Color {
	return draw.Color2(context.ColorMode(), draw.ColorTypeBase, 1, 1)
}

func DefaultDisabledListItemTextColor(context *guigui.Context) color.Color {
	return draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.5)
}

type List[T comparable] struct {
	guigui.DefaultWidget

	checkmark       Image
	listFrame       listFrame[T]
	scrollOverlay   ScrollOverlay
	dragDropOverlay DragDropOverlay

	abstractList               abstractList[T, ListItem[T]]
	itemBorderVisible          bool
	style                      ListStyle
	checkmarkIndexPlus1        int
	lastHoverredItemIndexPlus1 int
	lastSelectingItemTime      time.Time // TODO: Use ebiten.Tick.

	indexToJumpPlus1        int
	dropSrcIndexPlus1       int
	dropDstIndexPlus1       int
	pressStartX             int
	pressStartY             int
	startPressingIndexPlus1 int
	startPressingLeft       bool

	cachedDefaultWidth  int
	cachedDefaultHeight int
}

func listItemPadding(context *guigui.Context) int {
	return UnitSize(context) / 4
}

func (l *List[T]) SetOnItemSelected(f func(index int)) {
	l.abstractList.SetOnItemSelected(f)
}

func (l *List[T]) SetCheckmarkIndex(index int) {
	if index < 0 {
		index = -1
	}
	if l.checkmarkIndexPlus1 == index+1 {
		return
	}
	l.checkmarkIndexPlus1 = index + 1
	guigui.RequestRedraw(l)
}

func (l *List[T]) contentSize(context *guigui.Context) image.Point {
	return image.Pt(context.Size(l).X, l.defaultHeight(context))
}

func (l *List[T]) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if l.style != ListStyleSidebar && l.style != ListStyleMenu {
		l.listFrame.list = l
		appender.AppendChildWidgetWithPosition(&l.listFrame, context.Position(l))
	}

	l.scrollOverlay.SetContentSize(context, l.contentSize(context))

	if idx := l.indexToJumpPlus1 - 1; idx >= 0 {
		y := l.itemYFromIndex(context, idx) - RoundedCornerRadius(context)
		l.scrollOverlay.SetOffset(context, l.contentSize(context), 0, float64(-y))
		l.indexToJumpPlus1 = 0
	}

	appender.AppendChildWidgetWithBounds(&l.scrollOverlay, context.Bounds(l))

	hoveredItemIndex := l.HoveredItemIndex(context)
	p := context.Position(l)
	_, offsetY := l.scrollOverlay.Offset()
	p.X += RoundedCornerRadius(context) + listItemPadding(context)
	p.Y += RoundedCornerRadius(context) + int(offsetY)
	for i := range l.abstractList.ItemCount() {
		item, _ := l.abstractList.ItemByIndex(i)
		if l.checkmarkIndexPlus1 == i+1 {
			mode := context.ColorMode()
			if l.checkmarkIndexPlus1 == hoveredItemIndex+1 {
				mode = guigui.ColorModeDark
			}
			img, err := theResourceImages.Get("check", mode)
			if err != nil {
				return err
			}
			l.checkmark.SetImage(img)

			imgSize := listItemCheckmarkSize(context)
			imgP := p
			itemH := context.Size(item.Content).Y
			imgP.Y += (itemH - imgSize) * 3 / 4
			appender.AppendChildWidgetWithBounds(&l.checkmark, image.Rectangle{
				Min: imgP,
				Max: imgP.Add(image.Pt(imgSize, imgSize)),
			})
		}

		itemP := p
		if l.checkmarkIndexPlus1 > 0 {
			itemP.X += listItemCheckmarkSize(context) + listItemTextAndImagePadding(context)
		}
		appender.AppendChildWidgetWithPosition(item.Content, itemP)
		p.Y += context.Size(item.Content).Y
	}

	appender.AppendChildWidgetWithPosition(&l.dragDropOverlay, context.Position(l))

	if l.lastHoverredItemIndexPlus1 != hoveredItemIndex+1 {
		l.lastHoverredItemIndexPlus1 = hoveredItemIndex + 1
		if l.isHoveringVisible() {
			guigui.RequestRedraw(l)
		}
	}

	return nil
}

func (l *List[T]) ItemByIndex(index int) (ListItem[T], bool) {
	return l.abstractList.ItemByIndex(index)
}

func (l *List[T]) SelectedItemIndex() int {
	return l.abstractList.SelectedItemIndex()
}

func (l *List[T]) HoveredItemIndex(context *guigui.Context) int {
	if !context.IsWidgetHitAt(l, image.Pt(ebiten.CursorPosition())) {
		return -1
	}
	_, y := ebiten.CursorPosition()
	_, offsetY := l.scrollOverlay.Offset()
	y -= RoundedCornerRadius(context)
	y -= context.Position(l).Y
	y -= int(offsetY)
	index := -1
	var cy int
	for i := range l.abstractList.ItemCount() {
		item, _ := l.abstractList.ItemByIndex(i)
		h := context.Size(item.Content).Y
		if cy <= y && y < cy+h {
			index = i
			break
		}
		cy += h
	}
	return index
}

func (l *List[T]) SetItems(items []ListItem[T]) {
	l.abstractList.SetItems(items)
	l.cachedDefaultWidth = 0
	l.cachedDefaultHeight = 0
}

func (l *List[T]) SelectItemByIndex(index int) {
	if l.abstractList.SelectItemByIndex(index) {
		guigui.RequestRedraw(l)
	}
}

func (l *List[T]) SelectItemByTag(tag T) {
	if l.abstractList.SelectItemByTag(tag) {
		guigui.RequestRedraw(l)
	}
}

func (l *List[T]) JumpToItemIndex(index int) {
	if index < 0 || index >= l.abstractList.ItemCount() {
		return
	}
	l.indexToJumpPlus1 = index + 1
}

func (l *List[T]) SetItemBorderVisible(visible bool) {
	if l.itemBorderVisible == visible {
		return
	}
	l.itemBorderVisible = visible
	guigui.RequestRedraw(l)
}

func (l *List[T]) isHoveringVisible() bool {
	return l.style == ListStyleMenu
}

func (l *List[T]) Style() ListStyle {
	return l.style
}

func (l *List[T]) SetStyle(style ListStyle) {
	if l.style == style {
		return
	}
	l.style = style
	guigui.RequestRedraw(l)
}

func (l *List[T]) calcDropDstIndex(context *guigui.Context) int {
	_, y := ebiten.CursorPosition()
	for i := range l.abstractList.ItemCount() {
		if r := l.itemRect(context, i); y < (r.Min.Y+r.Max.Y)/2 {
			return i
		}
	}
	return l.abstractList.ItemCount()
}

func (l *List[T]) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	// Process dragging.
	if l.dragDropOverlay.IsDragging() {
		_, y := ebiten.CursorPosition()
		p := context.Position(l)
		h := context.Size(l).Y
		var dy float64
		if upperY := p.Y + UnitSize(context); y < upperY {
			dy = float64(upperY-y) / 4
		}
		if lowerY := p.Y + h - UnitSize(context); y >= lowerY {
			dy = float64(lowerY-y) / 4
		}
		l.scrollOverlay.SetOffsetByDelta(context, l.contentSize(context), 0, dy)
		i := l.calcDropDstIndex(context)
		if l.dropDstIndexPlus1-1 != i {
			l.dropDstIndexPlus1 = i + 1
			guigui.RequestRedraw(l)
		}
		return guigui.HandleInputByWidget(l)
	}

	// Process dropping.
	var dropped bool
	if l.dropSrcIndexPlus1 > 0 && l.dropDstIndexPlus1 > 0 {
		dropped = true
		/*if l.callback != nil && l.callback.OnItemDropped != nil {
			l.callback.OnItemDropped(l.dropSrcIndex, l.dropDstIndex)
		}*/
	}

	l.dropSrcIndexPlus1 = 0
	if l.dropDstIndexPlus1 != 0 {
		l.dropDstIndexPlus1 = 0
		guigui.RequestRedraw(l)
	}

	if dropped {
		return guigui.HandleInputByWidget(l)
	}

	index := l.HoveredItemIndex(context)
	if index >= 0 && index < l.abstractList.ItemCount() {
		x, y := ebiten.CursorPosition()
		left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		right := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

		switch {
		case left || right:
			item, _ := l.abstractList.ItemByIndex(index)
			if !item.Selectable {
				return guigui.HandleInputByWidget(l)
			}

			wasFocused := context.IsFocusedOrHasFocusedChild(l)
			context.SetFocused(l, true)
			if l.SelectedItemIndex() != index || !wasFocused || l.style == ListStyleMenu {
				l.SelectItemByIndex(index)
				l.lastSelectingItemTime = time.Now()
			}
			l.pressStartX = x
			l.pressStartY = y
			if right {
				/*if l.callback != nil && l.callback.OnContextMenu != nil {
					x, y := ebiten.CursorPosition()
					l.callback.OnContextMenu(index, x, y)
				}*/
			}
			l.startPressingIndexPlus1 = index + 1
			l.startPressingLeft = left

		case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
			item, _ := l.abstractList.ItemByIndex(index)
			if item.Draggable && l.SelectedItemIndex() == index && l.startPressingIndexPlus1-1 == index && (l.pressStartX != x || l.pressStartY != y) {
				l.dragDropOverlay.Start(index)
			}

		case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
			if l.SelectedItemIndex() == index && l.startPressingLeft && time.Since(l.lastSelectingItemTime) > 400*time.Millisecond {
				/*if l.callback != nil && l.callback.OnItemEditStarted != nil {
					l.callback.OnItemEditStarted(index)
				}*/
			}
			l.pressStartX = 0
			l.pressStartY = 0
			l.startPressingIndexPlus1 = 0
			l.startPressingLeft = false
		}

		return guigui.HandleInputByWidget(l)
	}

	l.dropSrcIndexPlus1 = 0
	l.pressStartX = 0
	l.pressStartY = 0

	return guigui.HandleInputResult{}
}

func (l *List[T]) itemYFromIndex(context *guigui.Context, index int) int {
	y := RoundedCornerRadius(context)
	for i := range l.abstractList.ItemCount() {
		if i == index {
			break
		}
		item, _ := l.abstractList.ItemByIndex(i)
		y += context.Size(item.Content).Y
	}
	return y
}

func (l *List[T]) itemRect(context *guigui.Context, index int) image.Rectangle {
	_, offsetY := l.scrollOverlay.Offset()
	b := context.Bounds(l)
	padding := listItemPadding(context)
	b.Min.X += RoundedCornerRadius(context) + padding
	b.Max.X -= RoundedCornerRadius(context) + padding
	b.Min.Y += l.itemYFromIndex(context, index)
	b.Min.Y += int(offsetY)
	if item, ok := l.abstractList.ItemByIndex(index); ok {
		b.Max.Y = b.Min.Y + context.Size(item.Content).Y
	}
	return b
}

func (l *List[T]) selectedItemColor(context *guigui.Context) color.Color {
	if l.SelectedItemIndex() < 0 || l.SelectedItemIndex() >= l.abstractList.ItemCount() {
		return nil
	}
	if l.style == ListStyleMenu {
		return nil
	}
	if context.IsFocusedOrHasFocusedChild(l) || l.style == ListStyleSidebar {
		return draw.Color(context.ColorMode(), draw.ColorTypeAccent, 0.5)
	}
	return draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.8)
}

func (l *List[T]) Draw(context *guigui.Context, dst *ebiten.Image) {
	if l.style != ListStyleSidebar {
		clr := draw.ControlColor(context.ColorMode(), context.IsEnabled(l))
		bounds := context.Bounds(l)
		draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
	}

	// Draw item borders.
	if l.itemBorderVisible && l.abstractList.ItemCount() > 0 {
		_, offsetY := l.scrollOverlay.Offset()
		p := context.Position(l)
		w := context.Size(l).X
		y := float32(p.Y) + float32(RoundedCornerRadius(context)) + float32(offsetY)
		for i := range l.abstractList.ItemCount() {
			item, _ := l.abstractList.ItemByIndex(i)
			y += float32(context.Size(item.Content).Y)
			if i == l.SelectedItemIndex() || i+1 == l.SelectedItemIndex() {
				continue
			}
			if i == l.abstractList.ItemCount()-1 {
				continue
			}
			x0 := p.X + RoundedCornerRadius(context)
			x1 := p.X + w - RoundedCornerRadius(context)
			width := 1 * float32(context.Scale())
			clr := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.5)
			vector.StrokeLine(dst, float32(x0), y, float32(x1), y, width, clr, false)
		}
	}

	if clr := l.selectedItemColor(context); clr != nil && l.SelectedItemIndex() >= 0 && l.SelectedItemIndex() < l.abstractList.ItemCount() {
		r := l.itemRect(context, l.SelectedItemIndex())
		r.Min.X -= RoundedCornerRadius(context)
		r.Max.X += RoundedCornerRadius(context)
		if r.Overlaps(context.VisibleBounds(l)) {
			draw.DrawRoundedRect(context, dst, r, clr, RoundedCornerRadius(context))
		}
	}

	hoveredItemIndex := l.HoveredItemIndex(context)
	hoveredItem, ok := l.abstractList.ItemByIndex(hoveredItemIndex)
	if ok && l.isHoveringVisible() && hoveredItemIndex >= 0 && hoveredItemIndex < l.abstractList.ItemCount() && hoveredItem.Selectable {
		r := l.itemRect(context, hoveredItemIndex)
		r.Min.X -= RoundedCornerRadius(context)
		r.Max.X += RoundedCornerRadius(context)
		if r.Overlaps(context.VisibleBounds(l)) {
			clr := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.9)
			if l.style == ListStyleMenu {
				clr = draw.Color(context.ColorMode(), draw.ColorTypeAccent, 0.5)
			}
			draw.DrawRoundedRect(context, dst, r, clr, RoundedCornerRadius(context))
		}
	}

	// Draw a drag indicator.
	/*if l.hoveredItemIndex >= 0 && l.hoveredItemIndex < l.abstractList.ItemCount() && l.items[l.hoveredItemIndex].Draggable && !l.dragDropOverlayWidget.Behavior().(*DragDropOverlay).IsDragging() {
		img := resource.Image("dragindicator", l.settings.Theme().UIForegroundColor)
		op := &ebiten.DrawImageOptions{}
		s := float64(2*RoundedCornerRadius(context)) / float64(img.Bounds().Dy())
		op.GeoM.Scale(s, s)
		r := l.itemRect(context, widget, l.hoveredItemIndex)
		op.GeoM.Translate(float64(r.Min.X-2*RoundedCornerRadius(context)), float64(r.Min.Y)+(float64(r.Dy())-float64(img.Bounds().Dy())*s)/2)
		dst.DrawImage(img, op)
	}*/

	// Draw a dragging guideline.
	if l.dropDstIndexPlus1 > 0 {
		p := context.Position(l)
		x0 := float32(p.X) + float32(RoundedCornerRadius(context))
		x1 := float32(p.X+context.Size(l).X) - float32(RoundedCornerRadius(context))
		y := float32(p.Y)
		y += float32(l.itemYFromIndex(context, l.dropDstIndexPlus1-1))
		_, offsetY := l.scrollOverlay.Offset()
		y += float32(offsetY)
		vector.StrokeLine(dst, x0, y, x1, y, 2*float32(context.Scale()), draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.1), false)
	}
}

/*func (l *List[T]) onDrop(data any) {
	l.dropSrcIndex = data.(int)
}*/

func (l *List[T]) defaultWidth(context *guigui.Context) int {
	if l.cachedDefaultWidth > 0 {
		return l.cachedDefaultWidth
	}
	var w int
	for i := range l.abstractList.ItemCount() {
		item, _ := l.abstractList.ItemByIndex(i)
		w = max(w, context.Size(item.Content).X)
	}
	w += 2*RoundedCornerRadius(context) + 2*listItemPadding(context)
	l.cachedDefaultWidth = w
	return w
}

func (l *List[T]) defaultHeight(context *guigui.Context) int {
	if l.cachedDefaultHeight > 0 {
		return l.cachedDefaultHeight
	}

	var h int
	h += RoundedCornerRadius(context)
	for i := range l.abstractList.ItemCount() {
		item, _ := l.abstractList.ItemByIndex(i)
		h += context.Size(item.Content).Y
	}
	h += RoundedCornerRadius(context)
	l.cachedDefaultHeight = h
	return h
}

func (l *List[T]) DefaultSize(context *guigui.Context) image.Point {
	w := l.defaultWidth(context)
	if l.checkmarkIndexPlus1 > 0 {
		w += listItemCheckmarkSize(context) + listItemTextAndImagePadding(context)
	}
	h := l.defaultHeight(context)
	return image.Pt(w, h)
}

type listFrame[T comparable] struct {
	guigui.DefaultWidget

	list *List[T]
}

func (l *listFrame[T]) Draw(context *guigui.Context, dst *ebiten.Image) {
	border := draw.RoundedRectBorderTypeInset
	if l.list.style != ListStyleNormal {
		border = draw.RoundedRectBorderTypeOutset
	}
	bounds := context.Bounds(l)
	clr1, clr2 := draw.BorderColors(context.ColorMode(), border, false)
	borderWidth := float32(1 * context.Scale())
	draw.DrawRoundedRectBorder(context, dst, bounds, clr1, clr2, RoundedCornerRadius(context), borderWidth, border)
}

func (l *listFrame[T]) DefaultSize(context *guigui.Context) image.Point {
	return context.Size(l.list)
}

func moveItemInSlice[T any](slice []T, from int, count int, to int) {
	if count == 0 {
		return
	}
	if from <= to && to <= from+count {
		return
	}
	if from < to {
		to -= count
	}

	s := make([]T, count)
	copy(s, slice[from:from+count])
	slice = slices.Delete(slice, from, from+count)
	// Assume that the slice has enough capacity, then the underlying array should not change.
	_ = slices.Insert(slice, to, s...)
}

func listItemCheckmarkSize(context *guigui.Context) int {
	return int(LineHeight(context) * 3 / 4)
}

func listItemTextAndImagePadding(context *guigui.Context) int {
	return UnitSize(context) / 8
}
