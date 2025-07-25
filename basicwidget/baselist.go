// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"fmt"
	"image"
	"image/color"

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

	indexToJumpPlus1        int
	dragSrcIndexPlus1       int
	dragDstIndexPlus1       int
	pressStartPlus1         image.Point
	startPressingIndexPlus1 int
	startPressingLeft       bool
	headerHeight            int
	footerHeight            int
	contentWidthPlus1       int

	cachedDefaultWidth         int
	cachedDefaultContentHeight int

	onItemsMoved func(from, count, to int)
}

func listItemPadding(context *guigui.Context) int {
	return RoundedCornerRadius(context) + UnitSize(context)/4
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

func (b *baseList[T]) SetHeaderHeight(height int) {
	if b.headerHeight == height {
		return
	}
	b.headerHeight = height
	guigui.RequestRedraw(b)
}

func (b *baseList[T]) SetFooterHeight(height int) {
	if b.footerHeight == height {
		return
	}
	b.footerHeight = height
	guigui.RequestRedraw(b)
}

func (b *baseList[T]) SetContentWidth(width int) {
	if b.contentWidthPlus1 == width+1 {
		return
	}
	b.contentWidthPlus1 = width + 1
	guigui.RequestRedraw(b)
}

func (b *baseList[T]) contentSize(context *guigui.Context) image.Point {
	w := context.ActualSize(b).X
	if b.contentWidthPlus1 > 0 {
		w = b.contentWidthPlus1 - 1
	}
	h := b.defaultHeight(context)
	h -= b.headerHeight
	h -= b.footerHeight
	return image.Pt(w, h)
}

func (b *baseList[T]) BeforeBuild(context *guigui.Context) {
	b.abstractList.ResetEventHandlers()
	b.onItemsMoved = nil
}

func (b *baseList[T]) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	for i := range b.abstractList.ItemCount() {
		item, _ := b.abstractList.ItemByIndex(i)
		if b.checkmarkIndexPlus1 == i+1 {
			appender.AppendChildWidget(&b.checkmark)
		}
		appender.AppendChildWidget(item.Content)
	}
	if b.style != ListStyleSidebar && b.style != ListStyleMenu {
		appender.AppendChildWidget(&b.listFrame)
	}
	appender.AppendChildWidget(&b.scrollOverlay)
}

func (b *baseList[T]) Build(context *guigui.Context) error {
	b.scrollOverlay.SetContentSize(context, b.contentSize(context))

	if idx := b.indexToJumpPlus1 - 1; idx >= 0 {
		y := b.itemYFromIndex(context, idx) - b.headerHeight - RoundedCornerRadius(context)
		b.scrollOverlay.SetOffset(context, b.contentSize(context), 0, float64(-y))
		b.indexToJumpPlus1 = 0
	}

	bounds := context.Bounds(b)
	bounds.Min.Y += b.headerHeight
	bounds.Max.Y -= b.footerHeight
	context.SetBounds(&b.scrollOverlay, bounds, b)

	// TODO: Do not call HoveredItemIndex in Build (#52).
	hoveredItemIndex := b.hoveredItemIndex(context)
	p := context.Position(b)
	offsetX, offsetY := b.scrollOverlay.Offset()
	p.X += listItemPadding(context) + int(offsetX)
	p.Y += RoundedCornerRadius(context) + b.headerHeight + int(offsetY)
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
			itemH := context.ActualSize(item.Content).Y
			imgP.Y += (itemH - imgSize) * 3 / 4
			imgP.Y = b.adjustItemY(context, imgP.Y)
			context.SetBounds(&b.checkmark, image.Rectangle{
				Min: imgP,
				Max: imgP.Add(image.Pt(imgSize, imgSize)),
			}, b)
		}

		itemP := p
		if b.checkmarkIndexPlus1 > 0 {
			itemP.X += listItemCheckmarkSize(context) + listItemTextAndImagePadding(context)
		}
		itemP.Y = b.adjustItemY(context, itemP.Y)

		context.SetPosition(item.Content, itemP)
		p.Y += context.ActualSize(item.Content).Y
	}

	if b.style != ListStyleSidebar && b.style != ListStyleMenu {
		b.listFrame.list = b
		context.SetBounds(&b.listFrame, context.Bounds(b), b)
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
	y -= RoundedCornerRadius(context) + b.headerHeight
	y -= context.Position(b).Y
	y -= int(offsetY)
	index := -1
	var cy int
	for i := range b.abstractList.ItemCount() {
		item, _ := b.abstractList.ItemByIndex(i)
		h := context.ActualSize(item.Content).Y
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
	b.cachedDefaultContentHeight = 0
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

func (b *baseList[T]) ScrollOffset() (float64, float64) {
	return b.scrollOverlay.Offset()
}

func (b *baseList[T]) calcDropDstIndex(context *guigui.Context) int {
	_, y := ebiten.CursorPosition()
	for i := range b.abstractList.ItemCount() {
		if b := b.itemBounds(context, i); y < (b.Min.Y+b.Max.Y)/2 {
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
			h := context.ActualSize(b).Y - (b.headerHeight + b.footerHeight)
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
	left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	right := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	if index >= 0 && index < b.abstractList.ItemCount() {
		c := image.Pt(ebiten.CursorPosition())

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
			}
			b.pressStartPlus1 = c.Add(image.Pt(1, 1))
			b.startPressingIndexPlus1 = index + 1
			b.startPressingLeft = left

		case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
			item, _ := b.abstractList.ItemByIndex(index)
			if item.Movable && b.SelectedItemIndex() == index && b.startPressingIndexPlus1-1 == index && (b.pressStartPlus1 != c.Add(image.Pt(1, 1))) {
				b.dragSrcIndexPlus1 = index + 1
			}

		case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
			b.pressStartPlus1 = image.Point{}
			b.startPressingIndexPlus1 = 0
			b.startPressingLeft = false
		}

		return guigui.HandleInputByWidget(b)
	}

	b.dragSrcIndexPlus1 = 0
	b.pressStartPlus1 = image.Point{}

	return guigui.HandleInputResult{}
}

func (b *baseList[T]) itemYFromIndex(context *guigui.Context, index int) int {
	y := RoundedCornerRadius(context) + b.headerHeight
	for i := range b.abstractList.ItemCount() {
		if i == index {
			break
		}
		item, _ := b.abstractList.ItemByIndex(i)
		y += context.ActualSize(item.Content).Y
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

func (b *baseList[T]) itemBounds(context *guigui.Context, index int) image.Rectangle {
	offsetX, offsetY := b.scrollOverlay.Offset()
	bounds := context.Bounds(b)
	bounds.Min.X += int(offsetX)
	if b.contentWidthPlus1 > 0 {
		bounds.Max.X = bounds.Min.X + b.contentWidthPlus1 - 1
	} else {
		bounds.Max.X += int(offsetX)
	}
	bounds.Min.X += listItemPadding(context)
	bounds.Max.X -= listItemPadding(context)
	bounds.Min.Y += b.itemYFromIndex(context, index)
	bounds.Min.Y += int(offsetY)
	if item, ok := b.abstractList.ItemByIndex(index); ok {
		bounds.Max.Y = bounds.Min.Y + context.ActualSize(item.Content).Y
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
	if !context.IsEnabled(b) {
		return draw.Color2(context.ColorMode(), draw.ColorTypeBase, 0.7, 0.2)
	}
	return draw.Color2(context.ColorMode(), draw.ColorTypeBase, 0.7, 0.5)
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
			bounds := b.itemBounds(context, i)
			if bounds.Min.Y > vb.Max.Y {
				break
			}
			bounds.Min.X -= RoundedCornerRadius(context)
			bounds.Max.X += RoundedCornerRadius(context)
			if !bounds.Overlaps(vb) {
				continue
			}
			clr := draw.SecondaryControlColor(context.ColorMode(), context.IsEnabled(b))
			draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
		}
	}

	// Draw the selected item background.
	if clr := b.selectedItemColor(context); clr != nil && b.SelectedItemIndex() >= 0 && b.SelectedItemIndex() < b.abstractList.ItemCount() {
		bounds := b.itemBounds(context, b.SelectedItemIndex())
		bounds.Min.X -= RoundedCornerRadius(context)
		bounds.Max.X += RoundedCornerRadius(context)
		if bounds.Overlaps(vb) {
			draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
		}
	}

	hoveredItemIndex := b.hoveredItemIndex(context)
	hoveredItem, ok := b.abstractList.ItemByIndex(hoveredItemIndex)
	if ok && b.isHoveringVisible() && hoveredItemIndex >= 0 && hoveredItemIndex < b.abstractList.ItemCount() && hoveredItem.Selectable {
		bounds := b.itemBounds(context, hoveredItemIndex)
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
			bounds := b.itemBounds(context, hoveredItemIndex)
			p := bounds.Min
			p.X = context.Position(b).X + listItemPadding(context)
			op.GeoM.Translate(float64(p.X-2*RoundedCornerRadius(context)), float64(p.Y)+(float64(bounds.Dy())-float64(img.Bounds().Dy())*s)/2)
			op.ColorScale.ScaleAlpha(0.5)
			op.Filter = ebiten.FilterLinear
			dst.DrawImage(img, op)
		}
	}

	// Draw a dragging guideline.
	if b.dragDstIndexPlus1 > 0 {
		p := context.Position(b)
		offsetX, _ := b.scrollOverlay.Offset()
		x0 := float32(p.X) + float32(RoundedCornerRadius(context))
		x0 += float32(offsetX)
		x1 := x0 + float32(b.contentSize(context).X)
		x1 -= 2 * float32(RoundedCornerRadius(context))
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
		w = max(w, context.ActualSize(item.Content).X)
	}
	w += 2 * listItemPadding(context)
	b.cachedDefaultWidth = w
	return w
}

func (b *baseList[T]) defaultHeight(context *guigui.Context) int {
	r := RoundedCornerRadius(context)
	if b.cachedDefaultContentHeight > 0 {
		return b.cachedDefaultContentHeight + 2*r + b.headerHeight + b.footerHeight
	}

	var h int
	for i := range b.abstractList.ItemCount() {
		item, _ := b.abstractList.ItemByIndex(i)
		h += context.ActualSize(item.Content).Y
	}
	b.cachedDefaultContentHeight = h
	return h + 2*r + b.headerHeight + b.footerHeight
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

func (l *listFrame[T]) headerBounds(context *guigui.Context) image.Rectangle {
	bounds := context.Bounds(l)
	bounds.Max.Y = bounds.Min.Y + l.list.headerHeight
	return bounds
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
	// Draw a header.
	if l.list.headerHeight > 0 {
		bounds := l.headerBounds(context)
		draw.DrawRoundedRectWithSharpenCorners(context, dst, bounds, draw.ControlColor(context.ColorMode(), context.IsEnabled(l)), RoundedCornerRadius(context), draw.SharpenCorners{
			UpperStart: false,
			UpperEnd:   false,
			LowerStart: true,
			LowerEnd:   true,
		})

		x0 := float32(bounds.Min.X)
		x1 := float32(bounds.Max.X)
		y0 := float32(bounds.Max.Y)
		y1 := float32(bounds.Max.Y)
		clr := draw.Color2(context.ColorMode(), draw.ColorTypeBase, 0.9, 0.4)
		if !context.IsEnabled(l) {
			clr = draw.Color2(context.ColorMode(), draw.ColorTypeBase, 0.8, 0.3)
		}
		vector.StrokeLine(dst, x0, y0, x1, y1, float32(context.Scale()), clr, false)
	}

	// Draw a footer.
	if l.list.footerHeight > 0 {
		bounds := l.footerBounds(context)
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
		if !context.IsEnabled(l) {
			clr = draw.Color2(context.ColorMode(), draw.ColorTypeBase, 0.8, 0.3)
		}
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
