// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

const popupZ = 16

const (
	popupEventClosed = "closed"
)

func easeOutQuad(t float64) float64 {
	// https://greweb.me/2012/02/bezier-curve-based-easing-functions-from-concept-to-implementation
	// easeOutQuad
	return t * (2 - t)
}

func popupMaxOpeningCount() int {
	return ebiten.TPS() / 5
}

type PopupClosedReason int

const (
	PopupClosedReasonNone PopupClosedReason = iota
	PopupClosedReasonFuncCall
	PopupClosedReasonClickOutside
	PopupClosedReasonReopen
)

type Popup struct {
	guigui.DefaultWidget

	background popupBackground
	shadow     popupShadow
	content    popupContent
	frame      popupFrame

	openingCount           int
	showing                bool
	hiding                 bool
	closedReason           PopupClosedReason
	backgroundBlurred      bool
	closeByClickingOutside bool
	animateOnFading        bool
	contentPosition        image.Point
	nextContentPosition    image.Point
	hasNextContentPosition bool
	openAfterClose         bool
}

func (p *Popup) IsOpen() bool {
	return p.showing || p.hiding || p.openingCount > 0
}

func (p *Popup) SetContent(widget guigui.Widget) {
	p.content.setContent(widget)
}

func (p *Popup) openingRate() float64 {
	return easeOutQuad(float64(p.openingCount) / float64(popupMaxOpeningCount()))
}

func (p *Popup) ContentBounds(context *guigui.Context) image.Rectangle {
	if p.content.content == nil {
		return image.Rectangle{}
	}
	pt := p.contentPosition
	if p.animateOnFading {
		rate := p.openingRate()
		dy := int(-float64(UnitSize(context)) * (1 - rate))
		pt = pt.Add(image.Pt(0, dy))
	}
	return image.Rectangle{
		Min: pt,
		Max: pt.Add(context.ActualSize(p)),
	}
}

func (p *Popup) SetBackgroundBlurred(blurBackground bool) {
	p.backgroundBlurred = blurBackground
}

func (p *Popup) SetCloseByClickingOutside(closeByClickingOutside bool) {
	p.closeByClickingOutside = closeByClickingOutside
}

func (p *Popup) SetAnimationDuringFade(animateOnFading bool) {
	// TODO: Rename Popup to basePopup and create Popup with animateOnFading true.
	p.animateOnFading = animateOnFading
}

func (p *Popup) SetOnClosed(f func(reason PopupClosedReason)) {
	guigui.RegisterEventHandler(p, popupEventClosed, f)
}

func (p *Popup) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	// Append children even if openingRate is 0, so that the Build methods of the children are invoked.
	// The children's state might affect their size and position.
	appender.AppendChildWidget(&p.background)
	appender.AppendChildWidget(&p.shadow)
	appender.AppendChildWidget(&p.content)
	appender.AppendChildWidget(&p.frame)
}

func (p *Popup) Build(context *guigui.Context) error {
	if (p.showing || p.hiding) && p.openingCount > 0 {
		p.nextContentPosition = context.Position(p)
		p.hasNextContentPosition = true
	} else {
		p.contentPosition = context.Position(p)
		p.nextContentPosition = image.Point{}
		p.hasNextContentPosition = false
	}

	p.background.popup = p
	p.shadow.popup = p
	p.content.popup = p
	p.frame.popup = p

	// SetOpacity cannot be called for p.background so far.
	// If opacity is less than 1, the dst argument of Draw will an empty image in the current implementation.
	// TODO: This is too tricky. Refactor this.
	context.SetOpacity(&p.shadow, p.openingRate())
	context.SetOpacity(&p.content, p.openingRate())
	context.SetOpacity(&p.frame, p.openingRate())

	context.SetBounds(&p.background, context.AppBounds(), p)
	context.SetBounds(&p.shadow, context.AppBounds(), p)
	context.SetBounds(&p.content, p.ContentBounds(context), p)
	context.SetBounds(&p.frame, context.AppBounds(), p)

	return nil
}

func (p *Popup) Open(context *guigui.Context) {
	if p.showing {
		return
	}
	if p.openingCount > 0 {
		p.close(PopupClosedReasonReopen)
		p.openAfterClose = true
		return
	}
	p.showing = true
	p.hiding = false
	context.SetFocused(p, true)
}

func (p *Popup) Close() {
	p.close(PopupClosedReasonFuncCall)
}

func (p *Popup) setClosedReason(reason PopupClosedReason) {
	if p.closedReason == PopupClosedReasonNone {
		p.closedReason = reason
		return
	}
	if reason != PopupClosedReasonReopen {
		return
	}
	// Overwrite the closed reason if it is PopupClosedReasonReopen.
	// A popup might already be closed by clicking outside.
	p.closedReason = reason
}

func (p *Popup) close(reason PopupClosedReason) {
	if p.hiding {
		p.setClosedReason(reason)
		return
	}
	if p.openingCount == 0 {
		return
	}

	p.setClosedReason(reason)
	p.showing = false
	p.hiding = true
	p.openAfterClose = false
}

func (p *Popup) IsWidgetOrBackgroundHitAtCursor(context *guigui.Context, target guigui.Widget) bool {
	if context.IsWidgetHitAtCursor(target) {
		return true
	}
	if context.IsWidgetHitAtCursor(&p.background) && image.Pt(ebiten.CursorPosition()).In(context.VisibleBounds(target)) {
		return true
	}
	return false
}

func (p *Popup) Tick(context *guigui.Context) error {
	if p.showing {
		if p.openingCount < popupMaxOpeningCount() {
			p.openingCount += 3
			p.openingCount = min(p.openingCount, popupMaxOpeningCount())
		}
		if p.openingCount == popupMaxOpeningCount() {
			p.showing = false
			if p.hasNextContentPosition {
				p.contentPosition = p.nextContentPosition
				p.hasNextContentPosition = false
			}
		}
	}
	if p.hiding {
		if 0 < p.openingCount {
			if p.closedReason == PopupClosedReasonReopen {
				p.openingCount -= 3
			} else {
				p.openingCount--
			}
			p.openingCount = max(p.openingCount, 0)
		}
		if p.openingCount == 0 {
			context.SetFocused(p, false)
			p.hiding = false
			guigui.InvokeEventHandler(p, popupEventClosed, p.closedReason)
			p.closedReason = PopupClosedReasonNone
			if p.openAfterClose {
				if p.hasNextContentPosition {
					p.contentPosition = p.nextContentPosition
					p.hasNextContentPosition = false
				}
				p.Open(context)
				p.openAfterClose = false
			}
		}
	}
	return nil
}

func (p *Popup) ZDelta() int {
	return popupZ
}

func (p *Popup) PassThrough() bool {
	return !p.IsOpen()
}

type popupContent struct {
	guigui.DefaultWidget

	popup *Popup

	content guigui.Widget
}

func (p *popupContent) setContent(widget guigui.Widget) {
	p.content = widget
}

func (p *popupContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	if p.content != nil {
		appender.AppendChildWidget(p.content)
	}
}

func (p *popupContent) Build(context *guigui.Context) error {
	if p.content != nil {
		context.SetPosition(p.content, context.Position(p))
	}
	return nil
}

func (p *popupContent) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if context.IsWidgetHitAtCursor(p) {
		return guigui.AbortHandlingInputByWidget(p)
	}
	return guigui.HandleInputResult{}
}

func (p *popupContent) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := context.Bounds(p)
	clr := draw.Color(context.ColorMode(), draw.ColorTypeBase, 1)
	draw.DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
}

func (p *popupContent) ZDelta() int {
	return 1
}

type popupFrame struct {
	guigui.DefaultWidget

	popup *Popup
}

func (p *popupFrame) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := p.popup.ContentBounds(context)
	clr1, clr2 := draw.BorderColors(context.ColorMode(), draw.RoundedRectBorderTypeOutset, false)
	draw.DrawRoundedRectBorder(context, dst, bounds, clr1, clr2, RoundedCornerRadius(context), float32(1*context.Scale()), draw.RoundedRectBorderTypeOutset)
}

func (p *popupFrame) DefaultSize(context *guigui.Context) image.Point {
	return context.ActualSize(p.popup)
}

func (p *popupFrame) ZDelta() int {
	return 1
}

type popupBackground struct {
	guigui.DefaultWidget

	popup *Popup

	backgroundCache *ebiten.Image
}

func (p *popupBackground) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if p.popup.showing || p.popup.hiding {
		return guigui.AbortHandlingInputByWidget(p)
	}

	if context.IsWidgetHitAtCursor(p) {
		if p.popup.closeByClickingOutside {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				p.popup.close(PopupClosedReasonClickOutside)
				// Continue handling inputs so that clicking a right button can be handled by other widgets.
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
					return guigui.HandleInputResult{}
				}
			}
		}
	}

	return guigui.AbortHandlingInputByWidget(p)
}

func (p *popupBackground) Draw(context *guigui.Context, dst *ebiten.Image) {
	if !p.popup.backgroundBlurred {
		return
	}

	bounds := context.Bounds(p)
	if p.backgroundCache != nil && !bounds.In(p.backgroundCache.Bounds()) {
		p.backgroundCache.Deallocate()
		p.backgroundCache = nil
	}
	if p.backgroundCache == nil {
		p.backgroundCache = ebiten.NewImageWithOptions(bounds, nil)
	}

	rate := p.popup.openingRate()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dst.Bounds().Min.X), float64(dst.Bounds().Min.Y))
	op.Blend = ebiten.BlendCopy
	p.backgroundCache.DrawImage(dst, op)

	draw.DrawBlurredImage(context, dst, p.backgroundCache, rate)
}

func (p *popupBackground) DefaultSize(context *guigui.Context) image.Point {
	return context.ActualSize(p.popup)
}

func (p *popupBackground) ZDelta() int {
	return 1
}

type popupShadow struct {
	guigui.DefaultWidget

	popup *Popup
}

func (p *popupShadow) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := p.popup.ContentBounds(context)
	bounds.Min.X -= int(16 * context.Scale())
	bounds.Max.X += int(16 * context.Scale())
	bounds.Min.Y -= int(8 * context.Scale())
	bounds.Max.Y += int(16 * context.Scale())
	clr := draw.ScaleAlpha(color.Black, 0.2)
	draw.DrawRoundedShadowRect(context, dst, bounds, clr, int(16*context.Scale())+RoundedCornerRadius(context))
}

func (p *popupShadow) DefaultSize(context *guigui.Context) image.Point {
	return context.ActualSize(p.popup)
}

func (p *popupShadow) ZDelta() int {
	return 1
}
