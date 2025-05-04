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

type ScrollablePanel struct {
	guigui.DefaultWidget

	content      guigui.Widget
	scollOverlay ScrollOverlay
	border       scrollablePanelBorder
}

func (s *ScrollablePanel) SetContent(widget guigui.Widget) {
	s.content = widget
}

func (s *ScrollablePanel) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if s.content == nil {
		return nil
	}

	offsetX, offsetY := s.scollOverlay.Offset()
	appender.AppendChildWidgetWithPosition(s.content, context.Position(s).Add(image.Pt(int(offsetX), int(offsetY))))

	s.scollOverlay.SetContentSize(context, context.Size(s.content))
	appender.AppendChildWidgetWithBounds(&s.scollOverlay, context.Bounds(s))

	s.border.scrollOverlay = &s.scollOverlay
	appender.AppendChildWidgetWithBounds(&s.border, context.Bounds(s))

	return nil
}

type scrollablePanelBorder struct {
	guigui.DefaultWidget

	scrollOverlay *ScrollOverlay
}

func (s *scrollablePanelBorder) Draw(context *guigui.Context, dst *ebiten.Image) {
	// Render borders.
	strokeWidth := float32(1 * context.Scale())
	bounds := context.Bounds(s)
	x0 := float32(bounds.Min.X)
	x1 := float32(bounds.Max.X)
	y0 := float32(bounds.Min.Y)
	y1 := float32(bounds.Max.Y)
	offsetX, offsetY := s.scrollOverlay.Offset()
	r := s.scrollOverlay.scrollRange(context)
	if offsetX < float64(r.Max.X) {
		vector.StrokeLine(dst, x0+strokeWidth/2, y0, x0+strokeWidth/2, y1, strokeWidth, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85), false)
	}
	if offsetY < float64(r.Max.Y) {
		vector.StrokeLine(dst, x0, y0+strokeWidth/2, x1, y0+strokeWidth/2, strokeWidth, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85), false)
	}
	if offsetX > float64(r.Min.X) {
		vector.StrokeLine(dst, x1-strokeWidth/2, y0, x1-strokeWidth/2, y1, strokeWidth, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85), false)
	}
	if offsetY > float64(r.Min.Y) {
		vector.StrokeLine(dst, x0, y1-strokeWidth/2, x1, y1-strokeWidth/2, strokeWidth, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85), false)
	}
}
