// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package guigui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MouseOverlay struct {
	DefaultWidget

	hovering      bool
	prevPos       image.Point
	pressingLeft  bool
	pressingRight bool

	sizeSet bool
	width   int
	height  int

	onDown  func(mouseButton ebiten.MouseButton, cursorPosition image.Point)
	onUp    func(mouseButton ebiten.MouseButton, cursorPosition image.Point)
	onEnter func(cursorPosition image.Point)
	onLeave func(cursorPosition image.Point)
	onMove  func(cursorPosition image.Point)
}

func (m *MouseOverlay) SetOnDown(f func(mouseButton ebiten.MouseButton, cursorPosition image.Point)) {
	m.onDown = f
}

func (m *MouseOverlay) SetOnUp(f func(mouseButton ebiten.MouseButton, cursorPosition image.Point)) {
	m.onUp = f
}

func (m *MouseOverlay) SetOnEnter(f func(cursorPosition image.Point)) {
	m.onEnter = f
}

func (m *MouseOverlay) SetOnLeave(f func(cursorPosition image.Point)) {
	m.onLeave = f
}

func (m *MouseOverlay) SetOnMove(f func(cursorPosition image.Point)) {
	m.onMove = f
}

func (m *MouseOverlay) HandleInput(context *Context) HandleInputResult {
	pos := image.Pt(ebiten.CursorPosition())
	defer func() {
		m.prevPos = pos
	}()
	m.setHovering(pos.In(VisibleBounds(m)) && IsVisible(m))
	if m.hovering && m.prevPos != pos && m.onMove != nil {
		m.onMove(pos)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if !image.Pt(ebiten.CursorPosition()).In(VisibleBounds(m)) {
			return HandleInputResult{}
		}
		if IsEnabled(m) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				m.setPressing(true, ebiten.MouseButtonLeft)
			}
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				m.setPressing(true, ebiten.MouseButtonRight)
			}
		}
		Focus(m)
		return HandleInputByWidget(m)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && m.pressingLeft ||
		inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) && m.pressingRight {
		if m.pressingLeft {
			m.setPressing(false, ebiten.MouseButtonLeft)
		}
		if m.pressingRight {
			m.setPressing(false, ebiten.MouseButtonRight)
		}
		if !image.Pt(ebiten.CursorPosition()).In(VisibleBounds(m)) {
			return HandleInputResult{}
		}
		if IsEnabled(m) {
			return HandleInputByWidget(m)
		}
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		m.setPressing(false, ebiten.MouseButtonLeft)
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		m.setPressing(false, ebiten.MouseButtonRight)
	}

	return HandleInputResult{}
}

func (m *MouseOverlay) Update(context *Context) error {
	if !IsVisible(m) {
		m.setHovering(false)
	}
	return nil
}

func (m *MouseOverlay) Size(context *Context) (int, int) {
	if m.sizeSet {
		return m.width, m.height
	}
	return Parent(m).Size(context)
}

func (m *MouseOverlay) SetSize(width, height int) {
	m.sizeSet = true
	m.width = width
	m.height = height
}

func (m *MouseOverlay) ResetSize() {
	m.sizeSet = false
	m.width = 0
	m.height = 0
}

func (m *MouseOverlay) setPressing(pressing bool, mouseButton ebiten.MouseButton) {
	switch mouseButton {
	case ebiten.MouseButtonLeft:
		if m.pressingLeft == pressing {
			return
		}
	case ebiten.MouseButtonRight:
		if m.pressingRight == pressing {
			return
		}
	}

	if IsEnabled(m) {
		if p := image.Pt(ebiten.CursorPosition()); p.In(VisibleBounds(Parent(m))) {
			if pressing {
				if m.onDown != nil {
					m.onDown(mouseButton, p)
				}
			} else {
				if m.onUp != nil {
					m.onUp(mouseButton, p)
				}
			}
		}
	}

	switch mouseButton {
	case ebiten.MouseButtonLeft:
		m.pressingLeft = pressing
	case ebiten.MouseButtonRight:
		m.pressingRight = pressing
	}
	RequestRedraw(m)
}

func (m *MouseOverlay) setHovering(hovering bool) {
	if m.hovering == hovering {
		return
	}

	if IsEnabled(m) {
		if p := image.Pt(ebiten.CursorPosition()); p.In(VisibleBounds(Parent(m))) {
			if hovering {
				if m.onEnter != nil {
					m.onEnter(p)
				}
			} else {
				if m.onLeave != nil {
					m.onLeave(p)
				}
			}
		}
	}

	m.hovering = hovering
	RequestRedraw(m)
}

func (m *MouseOverlay) IsPressing() bool {
	return m.pressingLeft || m.pressingRight
}

func (m *MouseOverlay) IsHovering() bool {
	return m.hovering
}
