// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type TextInputStyle int

const (
	TextInputStyleNormal TextInputStyle = iota
	TextInputStyleInline
)

type TextInput struct {
	guigui.DefaultWidget

	background     textInputBackground
	text           Text
	iconBackground textInputIconBackground
	icon           Image
	frame          textInputFrame
	scrollOverlay  ScrollOverlay
	focus          textInputFocus

	style        TextInputStyle
	readonly     bool
	paddingStart int
	paddingEnd   int

	prevFocused bool
	prevStart   int
	prevEnd     int

}

func (t *TextInput) SetOnValueChanged(f func(context *guigui.Context, text string, committed bool)) {
	t.text.SetOnValueChanged(f)
}

func (t *TextInput) SetOnKeyJustPressed(f func(context *guigui.Context, key ebiten.Key) (handled bool)) {
	t.text.SetOnKeyJustPressed(f)
}

const (
	textInputEventTextAndSelectionChanged = "textAndSelectionChanged"
)

func (t *TextInput) SetOnTextAndSelectionChanged(f func(context *guigui.Context, text string, start, end int)) {
	guigui.RegisterEventHandler(t, textInputEventTextAndSelectionChanged, f)
}

func (t *TextInput) Value() string {
	return t.text.Value()
}

func (t *TextInput) SetValue(text string) {
	t.text.SetValue(text)
}

func (t *TextInput) ForceSetValue(text string) {
	t.text.ForceSetValue(text)
}

func (t *TextInput) ReplaceValueAtSelection(text string) {
	t.text.ReplaceValueAtSelection(text)
}

func (t *TextInput) CommitWithCurrentInputValue() {
	t.text.CommitWithCurrentInputValue()
}

func (t *TextInput) SetMultiline(multiline bool) {
	t.text.SetMultiline(multiline)
}

func (t *TextInput) SetHorizontalAlign(halign HorizontalAlign) {
	t.text.SetHorizontalAlign(halign)
}

func (t *TextInput) SetVerticalAlign(valign VerticalAlign) {
	t.text.SetVerticalAlign(valign)
}

func (t *TextInput) SetAutoWrap(autoWrap bool) {
	t.text.SetAutoWrap(autoWrap)
}

func (t *TextInput) SelectAll() {
	t.text.selectAll()
}

func (t *TextInput) SetTabular(tabular bool) {
	t.text.SetTabular(tabular)
}

func (t *TextInput) IsEditable() bool {
	return !t.readonly
}

func (t *TextInput) SetStyle(style TextInputStyle) {
	if t.style == style {
		return
	}
	t.style = style
	guigui.RequestRedraw(t)
}

func (t *TextInput) SetEditable(editable bool) {
	if t.readonly == !editable {
		return
	}
	t.readonly = !editable
	t.text.SetEditable(editable)
	guigui.RequestRedraw(t)
}

func (t *TextInput) setPaddingStart(padding int) {
	if t.paddingStart == padding {
		return
	}
	t.paddingStart = padding
	guigui.RequestRedraw(t)
}

func (t *TextInput) setPaddingEnd(padding int) {
	if t.paddingEnd == padding {
		return
	}
	t.paddingEnd = padding
	guigui.RequestRedraw(t)
}

func (t *TextInput) SetIcon(icon *ebiten.Image) {
	t.icon.SetImage(icon)
}

func (t *TextInput) textInputPaddingInScrollableContent(context *guigui.Context) (start, top, end, bottom int) {
	var x, y int
	switch t.style {
	case TextInputStyleNormal:
		x = UnitSize(context) / 2
		y = int(float64(min(context.ActualSize(t).Y, UnitSize(context)))-LineHeight(context)*(t.text.scaleMinus1+1)) / 2
	case TextInputStyleInline:
		x = UnitSize(context) / 4
	}
	start = x + t.paddingStart
	if t.icon.HasImage() {
		start += defaultIconSize(context)
	}
	top = y
	end = x + t.paddingEnd
	bottom = y
	return
}

func (t *TextInput) scrollContentSize(context *guigui.Context) image.Point {
	start, top, end, bottom := t.textInputPaddingInScrollableContent(context)
	return t.text.DefaultSizeInContainer(context, context.ActualSize(t).X-start-end).Add(image.Pt(start+end, top+bottom))
}

func (t *TextInput) isFocused(context *guigui.Context) bool {
	return context.IsFocused(t) || context.IsFocused(&t.text)
}


func (t *TextInput) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&t.background)
	appender.AppendChildWidget(&t.text)
	if t.icon.HasImage() {
		appender.AppendChildWidget(&t.iconBackground)
		appender.AppendChildWidget(&t.icon)
	}
	appender.AppendChildWidget(&t.frame)
	appender.AppendChildWidget(&t.scrollOverlay)
	if t.style != TextInputStyleInline && (context.IsFocused(t) || context.IsFocused(&t.text)) {
		appender.AppendChildWidget(&t.focus)
	}
}

func (t *TextInput) Build(context *guigui.Context) error {
	if t.prevFocused != (context.IsFocused(t) || context.IsFocused(&t.text)) {
		t.prevFocused = (context.IsFocused(t) || context.IsFocused(&t.text))
		guigui.RequestRedraw(t)
	}

	paddingStart, paddingTop, paddingEnd, paddingBottom := t.textInputPaddingInScrollableContent(context)

	t.scrollOverlay.SetContentSize(context, t.scrollContentSize(context))

	t.background.textInput = t
	context.SetBounds(&t.background, context.Bounds(t), t)

	t.text.SetEditable(!t.readonly)
	t.text.SetSelectable(true)
	t.text.SetColor(draw.TextColor(context.ColorMode(), context.IsEnabled(t)))
	t.text.setKeepTailingSpace(!t.text.autoWrap)

	pt := context.Position(t)
	s := t.text.DefaultSizeInContainer(context, context.ActualSize(t).X-paddingStart-paddingEnd)
	s.X = max(s.X, context.ActualSize(t).X-paddingStart-paddingEnd)
	s.Y = max(s.Y, context.ActualSize(t).Y-paddingTop-paddingBottom)
	textBounds := image.Rectangle{
		Min: pt,
		Max: pt.Add(s),
	}
	textBounds = textBounds.Add(image.Pt(paddingStart, paddingTop))

	// As the text is rendered in an inset box, shift the text bounds down by 0.5 pixel.
	textBounds = textBounds.Add(image.Pt(0, int(0.5*context.Scale())))

	// Set the content size before adjustScrollOffset, as the size affects the adjustment.
	context.SetSize(&t.text, textBounds.Size(), t)
	t.adjustScrollOffsetIfNeeded(context)
	offsetX, offsetY := t.scrollOverlay.Offset()
	textBounds.Min = textBounds.Min.Add(image.Pt(int(offsetX), int(offsetY)))
	context.SetPosition(&t.text, textBounds.Min)
	if draw.OverlapsWithRoundedCorner(context.Bounds(t), RoundedCornerRadius(context), textBounds) {
		// CustomDraw might be too generic and overkill for this case.
		context.SetCustomDraw(&t.text, func(dst, widgetImage *ebiten.Image, op *ebiten.DrawImageOptions) {
			draw.DrawInRoundedCornerRect(context, dst, context.Bounds(t), RoundedCornerRadius(context), widgetImage, op)
		})
	} else {
		context.SetCustomDraw(&t.text, nil)
	}
	// Focusing the text widget works only after appending it.
	if context.IsFocused(t) {
		context.SetFocused(&t.text, true)
	}

	if t.icon.HasImage() {
		t.iconBackground.textInput = t

		b := context.Bounds(t)
		iconSize := defaultIconSize(context)
		var imgBounds image.Rectangle
		imgBounds.Min = b.Min.Add(image.Point{
			X: UnitSize(context)/4 + int(0.5*context.Scale()),
			Y: (b.Dy() - iconSize) / 2,
		})
		imgBounds.Max = imgBounds.Min.Add(image.Pt(iconSize, iconSize))

		imgBgBounds := b
		imgBgBounds.Max.X = imgBounds.Max.X + UnitSize(context)/4

		context.SetBounds(&t.iconBackground, imgBgBounds, t)
		context.SetBounds(&t.icon, imgBounds, t)
	}

	context.SetBounds(&t.frame, context.Bounds(t), t)

	context.SetVisible(&t.scrollOverlay, t.text.IsMultiline())
	context.SetBounds(&t.scrollOverlay, context.Bounds(t), t)

	if t.style != TextInputStyleInline && (context.IsFocused(t) || context.IsFocused(&t.text)) {
		t.focus.textInput = t
		w := textInputFocusBorderWidth(context)
		p := context.Position(t).Add(image.Pt(-w, -w))
		context.SetPosition(&t.focus, p)
	}

	return nil
}

func (t *TextInput) adjustScrollOffsetIfNeeded(context *guigui.Context) {
	start, end, ok := t.text.selectionToDraw(context)
	if !ok {
		return
	}
	if t.prevStart == start && t.prevEnd == end {
		return
	}
	t.prevStart = start
	t.prevEnd = end
	bounds := context.Bounds(t)
	paddingStart, paddingTop, paddingEnd, paddingBottom := t.textInputPaddingInScrollableContent(context)
	if pos, ok := t.text.textPosition(context, end, true); ok {
		dx := min(float64(bounds.Max.X-paddingEnd)-pos.X, 0)
		dy := min(float64(bounds.Max.Y-paddingBottom)-pos.Bottom, 0)
		t.scrollOverlay.SetOffsetByDelta(context, t.scrollContentSize(context), dx, dy)
	}
	if pos, ok := t.text.textPosition(context, start, true); ok {
		dx := max(float64(bounds.Min.X+paddingStart)-pos.X, 0)
		dy := max(float64(bounds.Min.Y+paddingTop)-pos.Top, 0)
		t.scrollOverlay.SetOffsetByDelta(context, t.scrollContentSize(context), dx, dy)
	}
}

func (t *TextInput) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	cp := image.Pt(ebiten.CursorPosition())
	if context.IsWidgetHitAtCursor(t) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			t.text.handleClick(context, cp)
			return guigui.HandleInputByWidget(t)
		}
	}
	return guigui.HandleInputResult{}
}

func (t *TextInput) CursorShape(context *guigui.Context) (ebiten.CursorShapeType, bool) {
	return t.text.CursorShape(context)
}

func (t *TextInput) DefaultSize(context *guigui.Context) image.Point {
	if t.style == TextInputStyleInline {
		start, _, end, _ := t.textInputPaddingInScrollableContent(context)
		w := max(t.text.DefaultSize(context).X+start+end, UnitSize(context))
		h := t.text.DefaultSize(context).Y
		return image.Pt(w, h)
	}
	if t.text.IsMultiline() {
		return image.Pt(6*UnitSize(context), 4*UnitSize(context))
	}
	return image.Pt(6*UnitSize(context), UnitSize(context))
}

type textInputBackground struct {
	guigui.DefaultWidget

	textInput *TextInput
}

func (t *textInputBackground) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := context.Bounds(t)
	clr := draw.ControlColor(context.ColorMode(), context.IsEnabled(t) && t.textInput.IsEditable())
	draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
}

type textInputIconBackground struct {
	guigui.DefaultWidget

	textInput *TextInput
}

func (t *textInputIconBackground) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := context.Bounds(t)
	clr := draw.ControlColor(context.ColorMode(), context.IsEnabled(t) && t.textInput.IsEditable())
	draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
}

type textInputFrame struct {
	guigui.DefaultWidget
}

func (t *textInputFrame) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := context.Bounds(t)
	clr1, clr2 := draw.BorderColors(context.ColorMode(), draw.RoundedRectBorderTypeInset, false)
	draw.DrawRoundedRectBorder(context, dst, bounds, clr1, clr2, RoundedCornerRadius(context), float32(1*context.Scale()), draw.RoundedRectBorderTypeInset)
}

func (t *textInputFrame) PassThrough() bool {
	return true
}

func textInputFocusBorderWidth(context *guigui.Context) int {
	return int(4 * context.Scale())
}

type textInputFocus struct {
	guigui.DefaultWidget

	textInput *TextInput
}

func (t *textInputFocus) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := context.Bounds(t.textInput)
	w := textInputFocusBorderWidth(context)
	clr := draw.Color(context.ColorMode(), draw.ColorTypeAccent, 0.8)
	bounds = bounds.Inset(-w)
	draw.DrawRoundedRectBorder(context, dst, bounds, clr, clr, w+RoundedCornerRadius(context), float32(w), draw.RoundedRectBorderTypeRegular)
}

func (t *textInputFocus) ZDelta() int {
	return 1
}

func (t *textInputFocus) DefaultSize(context *guigui.Context) image.Point {
	return context.ActualSize(t.textInput).Add(image.Pt(2*textInputFocusBorderWidth(context), 2*textInputFocusBorderWidth(context)))
}

func (t *textInputFocus) PassThrough() bool {
	return true
}
