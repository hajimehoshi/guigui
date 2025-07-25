// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type PanelStyle int

const (
	PanelStyleDefault PanelStyle = iota
	PanelStyleSide
)

type PanelBorder struct {
	Start  bool
	Top    bool
	End    bool
	Bottom bool
}

type Panel struct {
	guigui.DefaultWidget

	content      guigui.Widget
	scollOverlay ScrollOverlay
	border       panelBorder
	style        PanelStyle

	hasNextOffset     bool
	nextOffsetX       float64
	nextOffsetY       float64
	isNextOffsetDelta bool
}

func (p *Panel) SetContent(widget guigui.Widget) {
	p.content = widget
}

func (p *Panel) SetStyle(typ PanelStyle) {
	if p.style == typ {
		return
	}
	p.style = typ
	guigui.RequestRedraw(p)
}

func (p *Panel) SetBorder(borders PanelBorder) {
	p.border.setBorders(borders)
}

func (p *Panel) SetAutoBorder(auto bool) {
	p.border.SetAutoBorder(auto)
}

func (p *Panel) SetScrollOffset(offsetX, offsetY float64) {
	p.hasNextOffset = true
	p.nextOffsetX = offsetX
	p.nextOffsetY = offsetY
	p.isNextOffsetDelta = false
}

func (p *Panel) SetScrollOffsetByDelta(offsetXDelta, offsetYDelta float64) {
	p.hasNextOffset = true
	p.nextOffsetX = offsetXDelta
	p.nextOffsetY = offsetYDelta
	p.isNextOffsetDelta = true
}

func (p *Panel) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	if p.content != nil {
		appender.AppendChildWidget(p.content)
	}
	appender.AppendChildWidget(&p.scollOverlay)
	appender.AppendChildWidget(&p.border)
}

func (p *Panel) Build(context *guigui.Context) error {
	if p.content == nil {
		return nil
	}

	if p.hasNextOffset {
		if p.isNextOffsetDelta {
			p.scollOverlay.SetOffsetByDelta(context, context.ActualSize(p.content), p.nextOffsetX, p.nextOffsetY)
		} else {
			p.scollOverlay.SetOffset(context, context.ActualSize(p.content), p.nextOffsetX, p.nextOffsetY)
		}
		p.hasNextOffset = false
		p.nextOffsetX = 0
		p.nextOffsetY = 0
	}

	offsetX, offsetY := p.scollOverlay.Offset()
	context.SetPosition(p.content, context.Position(p).Add(image.Pt(int(offsetX), int(offsetY))))

	p.scollOverlay.SetContentSize(context, context.ActualSize(p.content))
	context.SetBounds(&p.scollOverlay, context.Bounds(p), p)

	p.border.scrollOverlay = &p.scollOverlay
	context.SetBounds(&p.border, context.Bounds(p), p)

	return nil
}

func (p *Panel) Draw(context *guigui.Context, dst *ebiten.Image) {
	switch p.style {
	case PanelStyleSide:
		dst.Fill(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.9))
	}
}

type panelBorder struct {
	guigui.DefaultWidget

	scrollOverlay *ScrollOverlay
	borders       PanelBorder
	autoBorder    bool
}

func (b *panelBorder) setBorders(borders PanelBorder) {
	if b.borders == borders {
		return
	}
	b.borders = borders
	guigui.RequestRedraw(b)
}

func (b *panelBorder) SetAutoBorder(auto bool) {
	if b.autoBorder == auto {
		return
	}
	b.autoBorder = auto
	guigui.RequestRedraw(b)
}

func (p *panelBorder) Draw(context *guigui.Context, dst *ebiten.Image) {
	// Render borders.
	strokeWidth := float32(1 * context.Scale())
	bounds := context.Bounds(p)
	x0 := float32(bounds.Min.X)
	x1 := float32(bounds.Max.X)
	y0 := float32(bounds.Min.Y)
	y1 := float32(bounds.Max.Y)
	offsetX, offsetY := p.scrollOverlay.Offset()
	r := p.scrollOverlay.scrollRange(context)
	clr := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.8)
	if (p.autoBorder && offsetX < float64(r.Max.X)) || p.borders.Start {
		vector.StrokeLine(dst, x0+strokeWidth/2, y0, x0+strokeWidth/2, y1, strokeWidth, clr, false)
	}
	if (p.autoBorder && offsetY < float64(r.Max.Y)) || p.borders.Top {
		vector.StrokeLine(dst, x0, y0+strokeWidth/2, x1, y0+strokeWidth/2, strokeWidth, clr, false)
	}
	if (p.autoBorder && offsetX > float64(r.Min.X)) || p.borders.End {
		vector.StrokeLine(dst, x1-strokeWidth/2, y0, x1-strokeWidth/2, y1, strokeWidth, clr, false)
	}
	if (p.autoBorder && offsetY > float64(r.Min.Y)) || p.borders.Bottom {
		vector.StrokeLine(dst, x0, y1-strokeWidth/2, x1, y1-strokeWidth/2, strokeWidth, clr, false)
	}
}
