// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"image"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

func barMaxOpacity() int {
	return int(float64(ebiten.TPS()) / 6)
}

func barShowingTime() int {
	return ebiten.TPS() / 2
}

type ScrollOverlay struct {
	guigui.DefaultWidget

	contentSize image.Point
	offsetX     float64
	offsetY     float64

	lastSize              image.Point
	lastCursorPosition    image.Point
	lastWheelX            float64
	lastWheelY            float64
	lastOffsetX           float64
	lastOffsetY           float64
	draggingX             bool
	draggingY             bool
	draggingStartPosition image.Point
	draggingStartOffsetX  float64
	draggingStartOffsetY  float64
	onceUpdated           bool

	barOpacity     int
	barVisibleTime int

	contentSizeChanged bool

	onScroll func(offsetX, offsetY float64)
}

func (s *ScrollOverlay) SetOnScroll(f func(offsetX, offsetY float64)) {
	s.onScroll = f
}

func (s *ScrollOverlay) Reset() {
	s.offsetX = 0
	s.offsetY = 0
}

func (s *ScrollOverlay) SetContentSize(context *guigui.Context, contentSize image.Point) {
	if s.contentSize == contentSize {
		return
	}

	s.contentSize = contentSize
	s.adjustOffset(context)
	if s.onceUpdated {
		s.contentSizeChanged = true
		guigui.RequestRedraw(s)
	}
}

func (s *ScrollOverlay) SetOffsetByDelta(context *guigui.Context, contentSize image.Point, dx, dy float64) {
	s.SetOffset(context, contentSize, s.offsetX+dx, s.offsetY+dy)
}

func (s *ScrollOverlay) SetOffset(context *guigui.Context, contentSize image.Point, x, y float64) {
	s.SetContentSize(context, contentSize)

	x, y = s.doAdjustOffset(context, x, y)
	if s.offsetX == x && s.offsetY == y {
		return
	}
	s.offsetX = x
	s.offsetY = y
	if s.onceUpdated {
		guigui.RequestRedraw(s)
	}
}

func (s *ScrollOverlay) isHovered(context *guigui.Context) bool {
	return context.IsWidgetHitAt(s, image.Pt(ebiten.CursorPosition()))
}

func (s *ScrollOverlay) setDragging(draggingX, draggingY bool) {
	if s.draggingX == draggingX && s.draggingY == draggingY {
		return
	}

	s.draggingX = draggingX
	s.draggingY = draggingY
}

func adjustedWheel() (float64, float64) {
	x, y := ebiten.Wheel()
	switch runtime.GOOS {
	case "darwin":
		x *= 2
		y *= 2
	}
	return x, y
}

func (s *ScrollOverlay) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	hovered := s.isHovered(context)
	if hovered {
		x, y := ebiten.CursorPosition()
		dx, dy := adjustedWheel()
		s.lastCursorPosition = image.Pt(x, y)
		s.lastWheelX = dx
		s.lastWheelY = dy
	} else {
		s.lastCursorPosition = image.Pt(-1, -1)
		s.lastWheelX = 0
		s.lastWheelY = 0
	}

	if !s.draggingX && !s.draggingY && hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		hb, vb := s.barBounds(context)
		if image.Pt(x, y).In(hb) {
			s.setDragging(true, s.draggingY)
			s.draggingStartPosition.X = x
			s.draggingStartOffsetX = s.offsetX
		} else if image.Pt(x, y).In(vb) {
			s.setDragging(s.draggingX, true)
			s.draggingStartPosition.Y = y
			s.draggingStartOffsetY = s.offsetY
		}
		if s.draggingX || s.draggingY {
			return guigui.HandleInputByWidget(s)
		}
	}

	if dx, dy := adjustedWheel(); dx != 0 || dy != 0 {
		s.setDragging(false, false)
	}

	if (s.draggingX || s.draggingY) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		var dx, dy float64
		if s.draggingX {
			dx = float64(x - s.draggingStartPosition.X)
		}
		if s.draggingY {
			dy = float64(y - s.draggingStartPosition.Y)
		}
		if dx != 0 || dy != 0 {
			prevOffsetX := s.offsetX
			prevOffsetY := s.offsetY

			cs := context.Size(s)
			barWidth, barHeight := s.barSize(context)
			if s.draggingX && barWidth > 0 && s.contentSize.X-cs.X > 0 {
				offsetPerPixel := float64(s.contentSize.X-cs.X) / (float64(cs.X) - barWidth)
				s.offsetX = s.draggingStartOffsetX + float64(-dx)*offsetPerPixel
			}
			if s.draggingY && barHeight > 0 && s.contentSize.Y-cs.Y > 0 {
				offsetPerPixel := float64(s.contentSize.Y-cs.Y) / (float64(cs.Y) - barHeight)
				s.offsetY = s.draggingStartOffsetY + float64(-dy)*offsetPerPixel
			}
			s.adjustOffset(context)
			if prevOffsetX != s.offsetX || prevOffsetY != s.offsetY {
				if s.onScroll != nil {
					s.onScroll(s.offsetX, s.offsetY)
				}
				guigui.RequestRedraw(s)
			}
		}
		return guigui.HandleInputByWidget(s)
	}

	if (s.draggingX || s.draggingY) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.setDragging(false, false)
	}

	if dx, dy := adjustedWheel(); dx != 0 || dy != 0 {
		if !hovered {
			return guigui.HandleInputResult{}
		}
		s.setDragging(false, false)

		prevOffsetX := s.offsetX
		prevOffsetY := s.offsetY
		s.offsetX += dx * 4 * context.Scale()
		s.offsetY += dy * 4 * context.Scale()
		s.adjustOffset(context)
		if prevOffsetX != s.offsetX || prevOffsetY != s.offsetY {
			if s.onScroll != nil {
				s.onScroll(s.offsetX, s.offsetY)
			}
			guigui.RequestRedraw(s)
			return guigui.HandleInputByWidget(s)
		}
		return guigui.HandleInputResult{}
	}

	return guigui.HandleInputResult{}
}

func (s *ScrollOverlay) CursorShape(context *guigui.Context) (ebiten.CursorShapeType, bool) {
	x, y := ebiten.CursorPosition()
	hb, vb := s.barBounds(context)
	if image.Pt(x, y).In(hb) || image.Pt(x, y).In(vb) {
		return ebiten.CursorShapeDefault, true
	}
	return 0, false
}

func (s *ScrollOverlay) Offset() (float64, float64) {
	return s.offsetX, s.offsetY
}

func (s *ScrollOverlay) adjustOffset(context *guigui.Context) {
	s.offsetX, s.offsetY = s.doAdjustOffset(context, s.offsetX, s.offsetY)
}

func (s *ScrollOverlay) doAdjustOffset(context *guigui.Context, x, y float64) (float64, float64) {
	r := s.scrollRange(context)
	x = min(max(x, float64(r.Min.X)), float64(r.Max.X))
	y = min(max(y, float64(r.Min.Y)), float64(r.Max.Y))
	return x, y
}

func (s *ScrollOverlay) scrollRange(context *guigui.Context) image.Rectangle {
	bounds := context.Bounds(s)
	return image.Rectangle{
		Min: image.Pt(min(bounds.Dx()-s.contentSize.X, 0), min(bounds.Dy()-s.contentSize.Y, 0)),
		Max: image.Pt(0, 0),
	}
}

func (s *ScrollOverlay) isBarVisible(context *guigui.Context) bool {
	if s.draggingX || s.draggingY {
		return true
	}
	if s.lastWheelX != 0 || s.lastWheelY != 0 {
		return true
	}
	if s.lastOffsetX != s.offsetX || s.lastOffsetY != s.offsetY {
		return true
	}

	bounds := context.Bounds(s)
	if s.contentSize.X > bounds.Dx() && bounds.Max.Y-UnitSize(context) <= s.lastCursorPosition.Y {
		return true
	}
	if s.contentSize.Y > bounds.Dy() && bounds.Max.X-UnitSize(context) <= s.lastCursorPosition.X {
		return true
	}
	return false
}

func (s *ScrollOverlay) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	cs := context.Size(s)
	if s.lastSize != cs {
		s.adjustOffset(context)
		s.lastSize = cs
	}

	context.SetOpacity(s, float64(s.barOpacity)/float64(barMaxOpacity())*3/4)
	return nil
}

func (s *ScrollOverlay) Update(context *guigui.Context) error {
	if s.contentSizeChanged {
		s.barVisibleTime = barShowingTime()
		s.contentSizeChanged = false
	}

	if s.isBarVisible(context) || (s.barVisibleTime == barShowingTime() && s.barOpacity < barMaxOpacity()) {
		if s.barOpacity < barMaxOpacity() {
			s.barOpacity++
			guigui.RequestRedraw(s)
		}
		s.barVisibleTime = barShowingTime()
	} else {
		if s.barVisibleTime > 0 {
			s.barVisibleTime--
		}
		if s.barVisibleTime == 0 && s.barOpacity > 0 {
			s.barOpacity--
			guigui.RequestRedraw(s)
		}
	}

	s.lastOffsetX = s.offsetX
	s.lastOffsetY = s.offsetY

	s.onceUpdated = true

	return nil
}

func (s *ScrollOverlay) Draw(context *guigui.Context, dst *ebiten.Image) {
	if s.barOpacity == 0 {
		return
	}

	barColor := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.2)
	hb, vb := s.barBounds(context)

	// Show a horizontal bar.
	if !hb.Empty() {
		draw.DrawRoundedRect(context, dst, hb, barColor, RoundedCornerRadius(context))
	}

	// Show a vertical bar.
	if !vb.Empty() {
		draw.DrawRoundedRect(context, dst, vb, barColor, RoundedCornerRadius(context))
	}
}

func scrollOverlayBarStrokeWidth(context *guigui.Context) float64 {
	return 8 * context.Scale()
}

func scrollOverlayPadding(context *guigui.Context) float64 {
	return 2 * context.Scale()
}

func (s *ScrollOverlay) barSize(context *guigui.Context) (float64, float64) {
	bounds := context.Bounds(s)
	padding := scrollOverlayPadding(context)

	var w, h float64
	if s.contentSize.X > bounds.Dx() {
		w = (float64(bounds.Dx()) - 2*padding) * float64(bounds.Dx()) / float64(s.contentSize.X)
		w = max(w, scrollOverlayBarStrokeWidth(context))
	}
	if s.contentSize.Y > bounds.Dy() {
		h = (float64(bounds.Dy()) - 2*padding) * float64(bounds.Dy()) / float64(s.contentSize.Y)
		h = max(h, scrollOverlayBarStrokeWidth(context))
	}
	return w, h
}

func (s *ScrollOverlay) barBounds(context *guigui.Context) (image.Rectangle, image.Rectangle) {
	bounds := context.Bounds(s)

	offsetX, offsetY := s.Offset()
	barWidth, barHeight := s.barSize(context)

	padding := scrollOverlayPadding(context)

	var horizontalBarBounds, verticalBarBounds image.Rectangle
	if s.contentSize.X > bounds.Dx() {
		rate := -offsetX / float64(s.contentSize.X-bounds.Dx())
		x0 := float64(bounds.Min.X) + padding + rate*(float64(bounds.Dx())-2*padding-barWidth)
		x1 := x0 + float64(barWidth)
		var y0, y1 float64
		if scrollOverlayBarStrokeWidth(context) > float64(bounds.Dy())*0.3 {
			y0 = float64(bounds.Max.Y) - float64(bounds.Dy())*0.3
			y1 = float64(bounds.Max.Y)
		} else {
			y0 = float64(bounds.Max.Y) - padding - scrollOverlayBarStrokeWidth(context)
			y1 = float64(bounds.Max.Y) - padding
		}
		horizontalBarBounds = image.Rect(int(x0), int(y0), int(x1), int(y1))
	}
	if s.contentSize.Y > bounds.Dy() {
		rate := -offsetY / float64(s.contentSize.Y-bounds.Dy())
		y0 := float64(bounds.Min.Y) + padding + rate*(float64(bounds.Dy())-2*padding-barHeight)
		y1 := y0 + float64(barHeight)
		var x0, x1 float64
		if scrollOverlayBarStrokeWidth(context) > float64(bounds.Dx())*0.3 {
			x0 = float64(bounds.Max.X) - float64(bounds.Dx())*0.3
			x1 = float64(bounds.Max.X)
		} else {
			x0 = float64(bounds.Max.X) - padding - scrollOverlayBarStrokeWidth(context)
			x1 = float64(bounds.Max.X) - padding
		}
		verticalBarBounds = image.Rect(int(x0), int(y0), int(x1), int(y1))
	}
	return horizontalBarBounds, verticalBarBounds
}
