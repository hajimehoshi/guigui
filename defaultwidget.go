// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package guigui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type DefaultWidget struct {
	s widgetState
}

var _ Widget = (*DefaultWidget)(nil)

func (*DefaultWidget) Build(context *Context, appender *ChildWidgetAppender) error {
	return nil
}

func (*DefaultWidget) HandlePointingInput(context *Context) HandleInputResult {
	return HandleInputResult{}
}

func (*DefaultWidget) HandleButtonInput(context *Context) HandleInputResult {
	return HandleInputResult{}
}

func (*DefaultWidget) Update(context *Context) error {
	return nil
}

func (*DefaultWidget) CursorShape(context *Context) (ebiten.CursorShapeType, bool) {
	return 0, false
}

func (*DefaultWidget) Draw(context *Context, dst *ebiten.Image) {
}

func (d *DefaultWidget) ZDelta() int {
	return 0
}

func (d *DefaultWidget) DefaultSize(context *Context) image.Point {
	return image.Pt(int(144*context.Scale()), int(144*context.Scale()))
}

func (d *DefaultWidget) PassThrough() bool {
	return false
}

func (d *DefaultWidget) widgetState() *widgetState {
	return &d.s
}

type RootWidget struct {
	DefaultWidget
}

func (d *RootWidget) DefaultSize(context *Context) image.Point {
	return context.app.bounds().Size()
}
