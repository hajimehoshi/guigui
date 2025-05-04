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

type Button struct {
	guigui.DefaultWidget

	pressed         bool
	keepPressed     bool
	useAccentColor  bool
	borderInvisible bool
	prevHovered     bool
	sharpenCorners  draw.SharpenCorners
	pairedButton    *Button

	onDown   func()
	onUp     func()
	onRepeat func()
}

func (b *Button) SetOnDown(f func()) {
	b.onDown = f
}

func (b *Button) SetOnUp(f func()) {
	b.onUp = f
}

func (b *Button) setOnRepeat(f func()) {
	b.onRepeat = f
}

func (b *Button) setPairedButton(pair *Button) {
	b.pairedButton = pair
}

func (b *Button) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	hovered := b.isHovered(context)
	if b.prevHovered != hovered {
		b.prevHovered = hovered
		guigui.RequestRedraw(b)
	}
	return nil
}

func (b *Button) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if context.IsEnabled(b) && b.isHovered(context) && !b.keepPressed {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			context.SetFocused(b, true)
			b.pressed = true
			if b.onDown != nil {
				b.onDown()
			}
			if isMouseButtonRepeating(ebiten.MouseButtonLeft) {
				if b.onRepeat != nil {
					b.onRepeat()
				}
			}
			return guigui.HandleInputByWidget(b)
		}
		if (b.pressed || b.pairedButton != nil && b.pairedButton.pressed) && isMouseButtonRepeating(ebiten.MouseButtonLeft) {
			if b.onRepeat != nil {
				b.onRepeat()
			}
			return guigui.HandleInputByWidget(b)
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && b.pressed {
			b.pressed = false
			if b.onUp != nil {
				b.onUp()
			}
			return guigui.HandleInputByWidget(b)
		}
	}
	if !context.IsEnabled(b) || !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.pressed = false
	}
	return guigui.HandleInputResult{}
}

func (b *Button) CursorShape(context *guigui.Context) (ebiten.CursorShapeType, bool) {
	if (b.canPress(context) || b.pressed || b.pairedButton != nil && b.pairedButton.pressed) && !b.keepPressed {
		return ebiten.CursorShapePointer, true
	}
	return 0, true
}

func (b *Button) Draw(context *guigui.Context, dst *ebiten.Image) {
	cm := context.ColorMode()
	backgroundColor := draw.ControlColor(context.ColorMode(), context.IsEnabled(b))
	if context.IsEnabled(b) {
		if b.isPressed(context) {
			if b.useAccentColor {
				backgroundColor = draw.Color2(cm, draw.ColorTypeAccent, 0.875, 0.5)
			} else {
				backgroundColor = draw.Color2(cm, draw.ColorTypeBase, 0.95, 0.25)
			}
		} else if b.canPress(context) {
			backgroundColor = draw.Color2(cm, draw.ColorTypeBase, 0.975, 0.275)
		}
	}

	bounds := context.Bounds(b)
	r := min(RoundedCornerRadius(context), bounds.Dx()/4, bounds.Dy()/4)
	border := !b.borderInvisible
	if context.IsEnabled(b) && (b.isHovered(context) || b.keepPressed) {
		border = true
	}
	if border || b.isPressed(context) {
		draw.DrawRoundedRectWithSharpenCorners(context, dst, bounds, backgroundColor, r, b.sharpenCorners)
	}

	if border {
		borderType := draw.RoundedRectBorderTypeRegular
		if context.IsEnabled(b) {
			if b.isPressed(context) {
				borderType = draw.RoundedRectBorderTypeInset
			} else {
				borderType = draw.RoundedRectBorderTypeOutset
			}
		}
		clr1, clr2 := draw.BorderColors(context.ColorMode(), borderType, b.useAccentColor && b.isPressed(context))
		draw.DrawRoundedRectBorderWithSharpenCorners(context, dst, bounds, clr1, clr2, r, float32(1*context.Scale()), borderType, b.sharpenCorners)
	}
}

func (b *Button) canPress(context *guigui.Context) bool {
	return context.IsEnabled(b) && b.isHovered(context) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !b.keepPressed
}

func (b *Button) isHovered(context *guigui.Context) bool {
	return context.IsWidgetHitAt(b, image.Pt(ebiten.CursorPosition()))
}

func (b *Button) isActive(context *guigui.Context) bool {
	return context.IsEnabled(b) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && b.isHovered(context) && (b.pressed || b.pairedButton != nil && b.pairedButton.pressed)
}

func (b *Button) isPressed(context *guigui.Context) bool {
	return context.IsEnabled(b) && (b.isActive(context) || b.keepPressed)
}

func (b *Button) setKeepPressed(keep bool) {
	if b.keepPressed == keep {
		return
	}
	b.keepPressed = keep
	guigui.RequestRedraw(b)
}

func defaultButtonSize(context *guigui.Context) image.Point {
	return image.Pt(6*UnitSize(context), UnitSize(context))
}

func (b *Button) DefaultSize(context *guigui.Context) image.Point {
	return defaultButtonSize(context)
}

func (b *Button) setSharpenCorners(sharpenCorners draw.SharpenCorners) {
	if b.sharpenCorners == sharpenCorners {
		return
	}
	b.sharpenCorners = sharpenCorners
	guigui.RequestRedraw(b)
}

func (b *Button) setUseAccentColor(use bool) {
	if b.useAccentColor == use {
		return
	}
	b.useAccentColor = use
	guigui.RequestRedraw(b)
}
