// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package basicwidget

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
)

const popupZ = 16

func easeOutQuad(t float64) float64 {
	// https://greweb.me/2012/02/bezier-curve-based-easing-functions-from-concept-to-implementation
	// easeOutQuad
	return t * (2 - t)
}

func popupMaxOpacityCount() int {
	return ebiten.TPS() / 10
}

type PopupClosedReason int

const (
	PopupClosedReasonFuncCall PopupClosedReason = iota
	PopupClosedReasonClickOutside
	PopupClosedReasonReopen
)

type Popup struct {
	guigui.DefaultWidget

	background popupBackground
	content    popupContent
	frame      popupFrame

	opacityCount           int
	showing                bool
	hiding                 bool
	closedReason           PopupClosedReason
	backgroundBlurred      bool
	closeByClickingOutside bool
	animateOnFading        bool
	contentBounds          image.Rectangle
	nextContentBounds      image.Rectangle
	openAfterClose         bool

	initOnce sync.Once

	onClosed func(reason PopupClosedReason)
}

func (p *Popup) SetContent(f func(context *guigui.Context, childAppender *ContainerChildWidgetAppender)) {
	p.content.setContent(f)
}

func (p *Popup) opacity() float64 {
	return easeOutQuad(float64(p.opacityCount) / float64(popupMaxOpacityCount()))
}

func (p *Popup) ContentBounds(context *guigui.Context) image.Rectangle {
	if !p.animateOnFading {
		return p.contentBounds
	}
	rate := p.opacity()
	bounds := p.contentBounds
	dy := int(-float64(UnitSize(context)) * (1 - rate))
	return bounds.Add(image.Pt(0, dy))
}

func (p *Popup) SetContentBounds(bounds image.Rectangle) {
	if (p.showing || p.hiding) && p.opacityCount > 0 {
		p.nextContentBounds = bounds
		return
	}
	p.contentBounds = bounds
	p.nextContentBounds = image.Rectangle{}
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
	p.onClosed = f
}

func (p *Popup) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	p.initOnce.Do(func() {
		guigui.Hide(p)
	})

	if p.backgroundBlurred {
		appender.AppendChildWidget(&p.background)
	}

	bounds := p.ContentBounds(context)
	guigui.SetPosition(&p.content, bounds.Min)
	p.content.setSize(bounds.Dx(), bounds.Dy())
	appender.AppendChildWidget(&p.content)

	appender.AppendChildWidget(&p.frame)

	return nil
}

func (p *Popup) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if p.showing || p.hiding {
		return guigui.AbortHandlingInputByWidget(p)
	}

	// As this editor is a modal dialog, do not let other widgets to handle inputs.
	if image.Pt(ebiten.CursorPosition()).In(guigui.VisibleBounds(p)) {
		if p.closeByClickingOutside {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				p.close(PopupClosedReasonClickOutside)
				// Continue handling inputs so that clicking a right button can be handled by other widgets.
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
					return guigui.HandleInputResult{}
				}
			}
		}
	}
	return guigui.AbortHandlingInputByWidget(p)
}

func (p *Popup) Open() {
	if p.showing {
		return
	}
	if p.opacityCount > 0 {
		p.close(PopupClosedReasonReopen)
		p.openAfterClose = true
		return
	}
	guigui.Show(p)
	p.showing = true
	p.hiding = false
}

func (p *Popup) Close() {
	p.close(PopupClosedReasonFuncCall)
}

func (p *Popup) close(reason PopupClosedReason) {
	if p.hiding {
		return
	}

	p.closedReason = reason
	p.showing = false
	p.hiding = true
	p.openAfterClose = false
}

func (p *Popup) Update(context *guigui.Context) error {
	if p.showing {
		if p.opacityCount < popupMaxOpacityCount() {
			p.opacityCount++
		}
		guigui.SetOpacity(&p.content, p.opacity())
		guigui.RequestRedraw(&p.background)
		if p.opacityCount == popupMaxOpacityCount() {
			p.showing = false
			if !p.nextContentBounds.Empty() {
				p.contentBounds = p.nextContentBounds
				p.nextContentBounds = image.Rectangle{}
			}
		}
	}
	if p.hiding {
		if 0 < p.opacityCount {
			p.opacityCount--
		}
		guigui.SetOpacity(&p.content, p.opacity())
		guigui.RequestRedraw(&p.background)
		if p.opacityCount == 0 {
			p.hiding = false
			if p.onClosed != nil {
				p.onClosed(p.closedReason)
			}
			if p.openAfterClose {
				if !p.nextContentBounds.Empty() {
					p.contentBounds = p.nextContentBounds
					p.nextContentBounds = image.Rectangle{}
				}
				p.Open()
				p.openAfterClose = false
			} else {
				guigui.Hide(p)
			}
		}
	}
	return nil
}

func (p *Popup) CursorShape(context *guigui.Context) (ebiten.CursorShapeType, bool) {
	return ebiten.CursorShapeDefault, true
}

func (p *Popup) Z() int {
	return guigui.Parent(p).Z() + popupZ
}

func (p *Popup) Size(context *guigui.Context) (int, int) {
	return context.AppSize()
}

type popupContent struct {
	guigui.DefaultWidget

	setContentFunc func(context *guigui.Context, childAppender *ContainerChildWidgetAppender)
	childWidgets   ContainerChildWidgetAppender

	width  int
	height int
}

func (p *popupContent) setContent(f func(context *guigui.Context, childAppender *ContainerChildWidgetAppender)) {
	p.setContentFunc = f
}

func (p *popupContent) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	p.childWidgets.reset()
	if p.setContentFunc != nil {
		p.setContentFunc(context, &p.childWidgets)
	}
	for _, childWidget := range p.childWidgets.iter() {
		appender.AppendChildWidget(childWidget)
	}

	return nil
}

func (p *popupContent) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if image.Pt(ebiten.CursorPosition()).In(guigui.VisibleBounds(p)) {
		return guigui.AbortHandlingInputByWidget(p)
	}
	return guigui.HandleInputResult{}
}

func (p *popupContent) Draw(context *guigui.Context, dst *ebiten.Image) {
	popup := guigui.Parent(p).(*Popup)
	bounds := popup.ContentBounds(context)
	clr := ScaleAlpha(Color(context.ColorMode(), ColorTypeBase, 1), popup.opacity())
	DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))
}

func (p *popupContent) setSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *popupContent) Size(context *guigui.Context) (int, int) {
	return p.width, p.height
}

type popupFrame struct {
	guigui.DefaultWidget
}

func (p *popupFrame) Draw(context *guigui.Context, dst *ebiten.Image) {
	popup := guigui.Parent(p).(*Popup)
	bounds := popup.ContentBounds(context)
	clr := ScaleAlpha(Color(context.ColorMode(), ColorTypeBase, 0.7), popup.opacity())
	DrawRoundedRectBorder(context, dst, bounds, clr, RoundedCornerRadius(context), float32(1*context.Scale()), RoundedRectBorderTypeOutset)
}

type popupBackground struct {
	guigui.DefaultWidget

	backgroundCache *ebiten.Image
}

func (p *popupBackground) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := guigui.Bounds(p)
	if p.backgroundCache != nil && !bounds.In(p.backgroundCache.Bounds()) {
		p.backgroundCache.Deallocate()
		p.backgroundCache = nil
	}
	if p.backgroundCache == nil {
		p.backgroundCache = ebiten.NewImageWithOptions(bounds, nil)
	}

	popup := guigui.Parent(p).(*Popup)
	rate := popup.opacity()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dst.Bounds().Min.X), float64(dst.Bounds().Min.Y))
	p.backgroundCache.DrawImage(dst, op)

	DrawBlurredImage(dst, p.backgroundCache, rate)
}
