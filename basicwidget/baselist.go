// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"fmt"
	"image"
	"image/color"
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

type baseListItem[T comparable] struct {
	Content    guigui.Widget
	Selectable bool
	Movable    bool
	ID         T
}

func (b baseListItem[T]) id() T {
	return b.ID
}

func DefaultActiveListItemTextColor(context *guigui.Context) color.Color {
	return draw.Color2(context.ColorMode(), draw.ColorTypeBase, 1, 1)
}

func DefaultDisabledListItemTextColor(context *guigui.Context) color.Color {
	return draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.5)
}

type baseList[T comparable] struct {
	guigui.DefaultWidget

	checkmark     Image
	listFrame     listFrame[T]
	scrollOverlay ScrollOverlay

	abstractList               abstractList[T, baseListItem[T]]
	stripeVisible              bool
	style                      ListStyle
	checkmarkIndexPlus1        int
	lastHoverredItemIndexPlus1 int
	lastSelectingItemTime      time.Time // TODO: Use ebiten.Tick.

	indexToJumpPlus1        int
	dragSrcIndexPlus1       int
	dragDstIndexPlus1       int
	pressStartX             int
	pressStartY             int
	startPressingIndexPlus1 int
	startPressingLeft       bool
	footerHeight            int

	cachedDefaultWidth  int
	cachedDefaultHeight int

	onItemsMoved func(from, count, to int)
}

func listItemPadding(context *guigui.Context) int {
	return UnitSize(context) / 4
}

func (b *baseList[T]) SetOnItemSelected(f func(index int)) {
	b.abstractList.SetOnItemSelected(f)
}

func (b *baseList[T]) SetOnItemsMoved(f func(from, count, to int)) {
	b.onItemsMoved = f
}

func (b *baseList[T]) SetCheckmarkIndex(index int) {
	if index < 0 {
		index = -1
	}
	if b.checkmarkIndexPlus1 == index+1 {
		return
	}
	b.checkmarkIndexPlus1 = index + 1
	guigui.RequestRedraw(b)
}

func (b *baseList[T]) SetFooterHeight(height int) {
	if b.footerHeight == height {
		return
	}
	b.footerHeight = height
	guigui.RequestRedraw(b)
}

func (b *baseList[T]) contentSize(context *guigui.Context) image.Point {
	return image.Pt(context.Size(b).X, b.defaultHeight(context)-b.footerHeight)
}

func (b *baseList[T]) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	b.scrollOverlay.SetContentSize(context, b.contentSize(context))

	if idx := b.indexToJumpPlus1 - 1; idx >= 0 {
		y := b.itemYFromIndex(context, idx) - RoundedCornerRadius(context) - b.footerHeight
		b.scrollOverlay.SetOffset(context, b.contentSize(context), 0, float64(-y))
		b.indexToJumpPlus1 = 0
	}

	bounds := context.Bounds(b)
	bounds.Max.Y -= b.footerHeight
	appender.AppendChildWidgetWithBounds(&b.scrollOverlay, bounds)

	// TODO: Do not call HoveredItemIndex in Build (#52).
	hoveredItemIndex := b.hoveredItemIndex(context)
	p := context.Position(b)
	_, offsetY := b.scrollOverlay.Offset()
	p.X += RoundedCornerRadius(context) + listItemPadding(context)
	p.Y += RoundedCornerRadius(context) + int(offsetY)
	for i := range b.abstractList.ItemCount() {
		item, _ := b.abstractList.ItemByIndex(i)
		if b.checkmarkIndexPlus1 == i+1 {
			mode := context.ColorMode()
			if b.checkmarkIndexPlus1 == hoveredItemIndex+1 {
				mode = guigui.ColorModeDark
			}
			img, err := theResourceImages.Get("check", mode)
			if err != nil {
				return err
			}
			b.checkmark.SetImage(img)

			imgSize := listItemCheckmarkSize(context)
			imgP := p
			itemH := context.Size(item.Content).Y
			imgP.Y += (itemH - imgSize) * 3 / 4
			imgP.Y = b.adjustItemY(context, imgP.Y)
			appender.AppendChildWidgetWithBounds(&b.checkmark, image.Rectangle{
				Min: imgP,
				Max: imgP.Add(image.Pt(imgSize, imgSize)),
			})
		}

		itemP := p
		if b.checkmarkIndexPlus1 > 0 {
			itemP.X += listItemCheckmarkSize(context) + listItemTextAndImagePadding(context)
		}
		itemP.Y = b.adjustItemY(context, itemP.Y)

		appender.AppendChildWidgetWithPosition(item.Content, itemP)
		p.Y += context.Size(item.Content).Y
	}

	if b.style != ListStyleSidebar && b.style != ListStyleMenu {
		b.listFrame.list = b
		appender.AppendChildWidgetWithBounds(&b.listFrame, context.Bounds(b))
	}

	return nil
}

func (b *baseList[T]) hasMovableItems() bool {
	for i := range b.abstractList.ItemCount() {
		item, ok := b.abstractList.ItemByIndex(i)
		if !ok {
			continue
		}
		if item.Movable {
			return true
		}
	}
	return false
}

func (b *baseList[T]) ItemByIndex(index int) (baseListItem[T], bool) {
	return b.abstractList.ItemByIndex(index)
}

func (b *baseList[T]) SelectedItemIndex() int {
	return b.abstractList.SelectedItemIndex()
}

func (b *baseList[T]) hoveredItemIndex(context *guigui.Context) int {
	if !context.IsWidgetHitAtCursor(b) {
		return -1
	}
	_, y := ebiten.CursorPosition()
	_, offsetY := b.scrollOverlay.Offset()
	y -= RoundedCornerRadius(context)
	y -= context.Position(b).Y
	y -= int(offsetY)
	index := -1
	var cy int
	for i := range b.abstractList.ItemCount() {
		item, _ := b.abstractList.ItemByIndex(i)
		h := context.Size(item.Content).Y
		if cy <= y && y < cy+h {
			index = i
			break
		}
		cy += h
	}
	return index
}

func (b *baseList[T]) SetItems(items []baseListItem[T]) {
	b.abstractList.SetItems(items)
	b.cachedDefaultWidth = 0
	b.cachedDefaultHeight = 0
}

func (b *baseList[T]) SelectItemByIndex(index int) {
	b.selectItemByIndex(index, false)
}

func (b *baseList[T]) selectItemByIndex(index int, forceFireEvents bool) {
	if b.abstractList.SelectItemByIndex(index, forceFireEvents) {
		guigui.RequestRedraw(b)
	}
}

func (b *baseList[T]) SelectItemByID(id T) {
	if b.abstractList.SelectItemByID(id, false) {
		guigui.RequestRedraw(b)
	}
}

func (b *baseList[T]) JumpToItemIndex(index int) {
	if index < 0 || index >= b.abstractList.ItemCount() {
		return
	}
	b.indexToJumpPlus1 = index + 1
}

func (b *baseList[T]) SetStripeVisible(visible bool) {
	if b.stripeVisible == visible {
		return
	}
	b.stripeVisible = visible
	guigui.RequestRedraw(b)
}

func (b *baseList[T]) isHoveringVisible() bool {
	return b.style == ListStyleMenu
}

func (b *baseList[T]) Style() ListStyle {
	return b.style
}

func (b *baseList[T]) SetStyle(style ListStyle) {
	if b.style == style {
		return
	}
	b.style = style
	guigui.RequestRedraw(b)
}

func (b *baseList[T]) calcDropDstIndex(context *guigui.Context) int {
	_, y := ebiten.CursorPosition()
	for i := range b.abstractList.ItemCount() {
		if b := b.itemBounds(context, i, true); y < (b.Min.Y+b.Max.Y)/2 {
			return i
		}
	}
	return b.abstractList.ItemCount()
}

func (b *baseList[T]) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if b.isHoveringVisible() || b.hasMovableItems() {
		if hoveredItemIndex := b.hoveredItemIndex(context); b.lastHoverredItemIndexPlus1 != hoveredItemIndex+1 {
			b.lastHoverredItemIndexPlus1 = hoveredItemIndex + 1
			guigui.RequestRedraw(b)
		}
	}

	// Process dragging.
	if b.dragSrcIndexPlus1 > 0 {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			_, y := ebiten.CursorPosition()
			p := context.Position(b)
			h := context.Size(b).Y - b.footerHeight
			var dy float64
			if upperY := p.Y + UnitSize(context); y < upperY {
				dy = float64(upperY-y) / 4
			}
			if lowerY := p.Y + h - UnitSize(context); y >= lowerY {
				dy = float64(lowerY-y) / 4
			}
			b.scrollOverlay.SetOffsetByDelta(context, b.contentSize(context), 0, dy)
			if i := b.calcDropDstIndex(context); b.dragDstIndexPlus1-1 != i {
				b.dragDstIndexPlus1 = i + 1
				guigui.RequestRedraw(b)
			}
			return guigui.HandleInputByWidget(b)
		}
		if b.dragDstIndexPlus1 > 0 {
			if b.onItemsMoved != nil {
				// TODO: Implement multiple items drop.
				b.onItemsMoved(b.dragSrcIndexPlus1-1, 1, b.dragDstIndexPlus1-1)
			}
			b.dragDstIndexPlus1 = 0
		}
		b.dragSrcIndexPlus1 = 0
		guigui.RequestRedraw(b)
		return guigui.HandleInputByWidget(b)
	}

	index := b.hoveredItemIndex(context)
	if index >= 0 && index < b.abstractList.ItemCount() {
		x, y := ebiten.CursorPosition()
		left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		right := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

		switch {
		case left || right:
			item, _ := b.abstractList.ItemByIndex(index)
			if !item.Selectable {
				return guigui.HandleInputByWidget(b)
			}

			wasFocused := context.IsFocusedOrHasFocusedChild(b)
			if item, ok := b.abstractList.ItemByIndex(index); ok {
				context.SetFocused(item.Content, true)
			} else {
				context.SetFocused(b, true)
			}
			if b.SelectedItemIndex() != index || !wasFocused || b.style == ListStyleMenu {
				b.selectItemByIndex(index, true)
				b.lastSelectingItemTime = time.Now()
			}
			b.pressStartX = x
			b.pressStartY = y
			if right {
				/*if l.callback != nil && l.callback.OnContextMenu != nil {
					x, y := ebiten.CursorPosition()
					l.callback.OnContextMenu(index, x, y)
				}*/
			}
			b.startPressingIndexPlus1 = index + 1
			b.startPressingLeft = left

		case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
			item, _ := b.abstractList.ItemByIndex(index)
			if item.Movable && b.SelectedItemIndex() == index && b.startPressingIndexPlus1-1 == index && (b.pressStartX != x || b.pressStartY != y) {
				b.dragSrcIndexPlus1 = index + 1
			}

		case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
			if b.SelectedItemIndex() == index && b.startPressingLeft && time.Since(b.lastSelectingItemTime) > 400*time.Millisecond {
				/*if l.callback != nil && l.callback.OnItemEditStarted != nil {
					l.callback.OnItemEditStarted(index)
				}*/
			}
			b.pressStartX = 0
			b.pressStartY = 0
			b.startPressingIndexPlus1 = 0
			b.startPressingLeft = false
		}

		return guigui.HandleInputByWidget(b)
	}

	b.dragSrcIndexPlus1 = 0
	b.pressStartX = 0
	b.pressStartY = 0

	return guigui.HandleInputResult{}
}

func (b *baseList[T]) itemYFromIndex(context *guigui.Context, index int) int {
	y := RoundedCornerRadius(context)
	for i := range b.abstractList.ItemCount() {
		if i == index {
			break
		}
		item, _ := b.abstractList.ItemByIndex(i)
		y += context.Size(item.Content).Y
	}
	y = b.adjustItemY(context, y)
	return y
}

func (b *baseList[T]) adjustItemY(context *guigui.Context, y int) int {
	// Adjust the bounds based on the list style (inset or outset).
	switch b.style {
	case ListStyleNormal:
		y += int(0.5 * context.Scale())
	case ListStyleMenu:
		y += int(-0.5 * context.Scale())
	}
	return y
}

func (b *baseList[T]) itemBounds(context *guigui.Context, index int, fullWidth bool) image.Rectangle {
	_, offsetY := b.scrollOverlay.Offset()
	bounds := context.Bounds(b)
	if !fullWidth {
		padding := listItemPadding(context)
		bounds.Min.X += RoundedCornerRadius(context) + padding
		bounds.Max.X -= RoundedCornerRadius(context) + padding
	}
	bounds.Min.Y += b.itemYFromIndex(context, index)
	bounds.Min.Y += int(offsetY)
	if item, ok := b.abstractList.ItemByIndex(index); ok {
		bounds.Max.Y = bounds.Min.Y + context.Size(item.Content).Y
	}
	return bounds
}

func (b *baseList[T]) selectedItemColor(context *guigui.Context) color.Color {
	if b.SelectedItemIndex() < 0 || b.SelectedItemIndex() >= b.abstractList.ItemCount() {
		return nil
	}
	if b.style == ListStyleMenu {
		return nil
	}
	if context.IsFocusedOrHasFocusedChild(b) || b.style == ListStyleSidebar {
		return draw.Color(context.ColorMode(), draw.ColorTypeAccent, 0.5)
	}
	return draw.Color2(context.ColorMode(), draw.ColorTypeBase, 0.7, 0.5)
}

func (b *baseList[T]) drawStripe(context *guigui.Context, dst *ebiten.Image, bounds image.Rectangle) {
	r := RoundedCornerRadius(context)
	if b.style != ListStyleNormal {
		r = 0
	}
	clr := draw.SecondaryControlColor(context.ColorMode(), context.IsEnabled(b))
	if r == 0 || !draw.OverlapsWithRoundedCorner(context.Bounds(b), r, bounds) {
		dst.SubImage(bounds).(*ebiten.Image).Fill(clr)
	} else {
		draw.FillInRoundedCornerRect(context, dst, context.Bounds(b), r, bounds, clr)
	}
}

func (b *baseList[T]) Draw(context *guigui.Context, dst *ebiten.Image) {
	var clr color.Color
	switch b.style {
	case ListStyleSidebar:
	case ListStyleNormal:
		clr = draw.ControlColor(context.ColorMode(), context.IsEnabled(b))
	case ListStyleMenu:
		clr = draw.SecondaryControlColor(context.ColorMode(), context.IsEnabled(b))
	}
	if clr != nil {
		bounds := context.Bounds(b)
		draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
	}

	vb := context.VisibleBounds(b)

	if b.stripeVisible && b.abstractList.ItemCount() > 0 {
		// Draw item stripes.
		// TODO: Get indices of items that are visible.
		for i := range b.abstractList.ItemCount() {
			if i%2 == 0 {
				continue
			}
			bounds := b.itemBounds(context, i, true)
			if bounds.Min.Y > vb.Max.Y {
				break
			}
			if !bounds.Overlaps(vb) {
				continue
			}
			b.drawStripe(context, dst, bounds)
		}

		// Draw the top stripe.
		{
			bounds := b.itemBounds(context, 0, true)
			bounds.Min.Y, bounds.Max.Y = bounds.Min.Y-RoundedCornerRadius(context), bounds.Min.Y
			if bounds.Overlaps(vb) {
				b.drawStripe(context, dst, bounds)
			}
		}

		// Draw the bottom stripe.
		if b.abstractList.ItemCount()%2 == 1 {
			bounds := b.itemBounds(context, b.abstractList.ItemCount()-1, true)
			bounds.Max.Y, bounds.Min.Y = bounds.Max.Y+RoundedCornerRadius(context), bounds.Max.Y
			if bounds.Overlaps(vb) {
				b.drawStripe(context, dst, bounds)
			}
		}
	}

	// Draw the selected item background.
	if clr := b.selectedItemColor(context); clr != nil && b.SelectedItemIndex() >= 0 && b.SelectedItemIndex() < b.abstractList.ItemCount() {
		bounds := b.itemBounds(context, b.SelectedItemIndex(), b.stripeVisible)
		if bounds.Overlaps(vb) {
			if b.stripeVisible {
				r := RoundedCornerRadius(context)
				if !draw.OverlapsWithRoundedCorner(context.Bounds(b), r, bounds) {
					dst.SubImage(bounds).(*ebiten.Image).Fill(clr)
				} else {
					draw.FillInRoundedCornerRect(context, dst, context.Bounds(b), r, bounds, clr)
				}
			} else {
				bounds.Min.X -= RoundedCornerRadius(context)
				bounds.Max.X += RoundedCornerRadius(context)
				draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
			}
		}
	}

	hoveredItemIndex := b.hoveredItemIndex(context)
	hoveredItem, ok := b.abstractList.ItemByIndex(hoveredItemIndex)
	if ok && b.isHoveringVisible() && hoveredItemIndex >= 0 && hoveredItemIndex < b.abstractList.ItemCount() && hoveredItem.Selectable {
		bounds := b.itemBounds(context, hoveredItemIndex, false)
		bounds.Min.X -= RoundedCornerRadius(context)
		bounds.Max.X += RoundedCornerRadius(context)
		if bounds.Overlaps(vb) {
			clr := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.9)
			if b.style == ListStyleMenu {
				clr = draw.Color(context.ColorMode(), draw.ColorTypeAccent, 0.5)
			}
			draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
		}
	}

	// Draw a drag indicator.
	if context.IsEnabled(b) && b.dragSrcIndexPlus1 == 0 {
		if item, ok := b.abstractList.ItemByIndex(hoveredItemIndex); ok && item.Movable {
			img, err := theResourceImages.Get("drag_indicator", context.ColorMode())
			if err != nil {
				panic(fmt.Sprintf("basicwidget: failed to get drag indicator image: %v", err))
			}
			op := &ebiten.DrawImageOptions{}
			s := float64(2*RoundedCornerRadius(context)) / float64(img.Bounds().Dy())
			op.GeoM.Scale(s, s)
			bounds := b.itemBounds(context, hoveredItemIndex, false)
			op.GeoM.Translate(float64(bounds.Min.X-2*RoundedCornerRadius(context)), float64(bounds.Min.Y)+(float64(bounds.Dy())-float64(img.Bounds().Dy())*s)/2)
			op.ColorScale.ScaleAlpha(0.5)
			dst.DrawImage(img, op)
		}
	}

	// Draw a dragging guideline.
	if b.dragDstIndexPlus1 > 0 {
		p := context.Position(b)
		x0 := float32(p.X)
		x1 := float32(p.X + context.Size(b).X)
		if !b.stripeVisible {
			x0 += float32(listItemPadding(context))
			x1 -= float32(listItemPadding(context))
		}
		y := float32(p.Y)
		y += float32(b.itemYFromIndex(context, b.dragDstIndexPlus1-1))
		_, offsetY := b.scrollOverlay.Offset()
		y += float32(offsetY)
		vector.StrokeLine(dst, x0, y, x1, y, 2*float32(context.Scale()), draw.Color(context.ColorMode(), draw.ColorTypeAccent, 0.5), false)
	}
}

func (b *baseList[T]) defaultWidth(context *guigui.Context) int {
	if b.cachedDefaultWidth > 0 {
		return b.cachedDefaultWidth
	}
	var w int
	for i := range b.abstractList.ItemCount() {
		item, _ := b.abstractList.ItemByIndex(i)
		w = max(w, context.Size(item.Content).X)
	}
	w += 2*RoundedCornerRadius(context) + 2*listItemPadding(context)
	b.cachedDefaultWidth = w
	return w
}

func (b *baseList[T]) defaultHeight(context *guigui.Context) int {
	if b.cachedDefaultHeight > 0 {
		return b.cachedDefaultHeight + b.footerHeight
	}

	var h int
	h += RoundedCornerRadius(context)
	for i := range b.abstractList.ItemCount() {
		item, _ := b.abstractList.ItemByIndex(i)
		h += context.Size(item.Content).Y
	}
	h += RoundedCornerRadius(context)
	b.cachedDefaultHeight = h
	h += b.footerHeight
	return h
}

func (b *baseList[T]) DefaultSize(context *guigui.Context) image.Point {
	w := b.defaultWidth(context)
	if b.checkmarkIndexPlus1 > 0 {
		w += listItemCheckmarkSize(context) + listItemTextAndImagePadding(context)
	}
	h := b.defaultHeight(context)
	return image.Pt(w, h)
}

type listFrame[T comparable] struct {
	guigui.DefaultWidget

	list *baseList[T]
}

func (l *listFrame[T]) footerBounds(context *guigui.Context) image.Rectangle {
	bounds := context.Bounds(l)
	bounds.Min.Y = bounds.Max.Y - l.list.footerHeight
	return bounds
}

func (l *listFrame[T]) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if context.IsWidgetHitAtCursor(l) {
		if image.Pt(ebiten.CursorPosition()).In(l.footerBounds(context)) {
			return guigui.HandleInputByWidget(l)
		}
	}
	return guigui.HandleInputResult{}
}

func (l *listFrame[T]) Draw(context *guigui.Context, dst *ebiten.Image) {

	// Draw a footer border.
	if l.list.footerHeight > 0 {
		bounds := l.footerBounds(context)
		bounds.Min.Y = bounds.Max.Y - l.list.footerHeight
		draw.DrawRoundedRectWithSharpenCorners(context, dst, bounds, draw.ControlColor(context.ColorMode(), context.IsEnabled(l)), RoundedCornerRadius(context), draw.SharpenCorners{
			UpperStart: true,
			UpperEnd:   true,
			LowerStart: false,
			LowerEnd:   false,
		})

		x0 := float32(bounds.Min.X)
		x1 := float32(bounds.Max.X)
		y0 := float32(bounds.Min.Y)
		y1 := float32(bounds.Min.Y)
		clr := draw.Color2(context.ColorMode(), draw.ColorTypeBase, 0.9, 0.4)
		vector.StrokeLine(dst, x0, y0, x1, y1, float32(context.Scale()), clr, false)
	}

	bounds := context.Bounds(l)
	border := draw.RoundedRectBorderTypeInset
	if l.list.style != ListStyleNormal {
		border = draw.RoundedRectBorderTypeOutset
	}
	clr1, clr2 := draw.BorderColors(context.ColorMode(), border, false)
	borderWidth := float32(1 * context.Scale())
	draw.DrawRoundedRectBorder(context, dst, bounds, clr1, clr2, RoundedCornerRadius(context), borderWidth, border)
}

func listItemCheckmarkSize(context *guigui.Context) int {
	return int(LineHeight(context) * 3 / 4)
}

func listItemTextAndImagePadding(context *guigui.Context) int {
	return UnitSize(context) / 8
}
