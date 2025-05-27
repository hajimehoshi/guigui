// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type IconAlign int

const (
	IconAlignStart IconAlign = iota
	IconAlignEnd
)

type Button struct {
	guigui.DefaultWidget

	button    baseButton
	content   guigui.Widget
	text      Text
	icon      Image
	iconAlign IconAlign

	textColor color.Color
}

func (b *Button) SetOnDown(f func()) {
	b.button.SetOnDown(f)
}

func (b *Button) SetOnUp(f func()) {
	b.button.SetOnUp(f)
}

func (b *Button) setOnRepeat(f func()) {
	b.button.setOnRepeat(f)
}

func (b *Button) SetContent(content guigui.Widget) {
	b.content = content
}

func (b *Button) SetText(text string) {
	b.text.SetValue(text)
}

func (b *Button) SetTextBold(bold bool) {
	b.text.SetBold(bold)
}

func (b *Button) SetIcon(icon *ebiten.Image) {
	b.icon.SetImage(icon)
}

func (b *Button) SetIconAlign(align IconAlign) {
	if b.iconAlign == align {
		return
	}
	b.iconAlign = align
	guigui.RequestRedraw(b)
}

func (b *Button) SetTextColor(clr color.Color) {
	if draw.EqualColor(b.textColor, clr) {
		return
	}
	b.textColor = clr
	guigui.RequestRedraw(b)
}

func (b *Button) setPairedButton(pair *Button) {
	b.button.setPairedButton(&pair.button)
}

func (b *Button) setKeepPressed(keep bool) {
	b.button.setKeepPressed(keep)
}

func (b *Button) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidgetWithBounds(&b.button, context.Bounds(b))

	s := context.Size(b)
	ds := b.defaultSize(context, false)

	if b.content != nil {
		contentP := context.Position(b)
		contentP.X += (s.X - ds.X) / 2
		contentP.Y += (s.Y - ds.Y) / 2

		cs := context.Size(b.content)
		contentP.X += buttonEdgeAndImagePadding(context)
		contentP.Y += (s.Y - cs.Y) / 2
		if b.button.isPressed(context) {
			contentP.Y += int(0.5 * context.Scale())
		} else {
			contentP.Y -= int(0.5 * context.Scale())
		}
		appender.AppendChildWidgetWithPosition(b.content, contentP)
	}

	imgSize := b.iconSize(context)

	tw := b.text.TextSize(context).X
	if b.textColor != nil {
		b.text.SetColor(b.textColor)
	} else {
		b.text.SetColor(draw.TextColor(context.ColorMode(), context.IsEnabled(b)))
	}
	b.text.SetHorizontalAlign(HorizontalAlignCenter)
	b.text.SetVerticalAlign(VerticalAlignMiddle)

	textP := context.Position(b)
	if b.icon.HasImage() {
		textP.X += (s.X - ds.X) / 2
		switch b.iconAlign {
		case IconAlignStart:
			textP.X += buttonEdgeAndImagePadding(context)
			textP.X += imgSize.X + buttonTextAndImagePadding(context)
		case IconAlignEnd:
			textP.X += buttonEdgeAndTextPadding(context)
		}
	} else {
		textP.X += (s.X - tw) / 2
	}
	if b.button.isPressed(context) {
		textP.Y += int(0.5 * context.Scale())
	} else {
		textP.Y -= int(0.5 * context.Scale())
	}
	appender.AppendChildWidgetWithBounds(&b.text, image.Rectangle{
		Min: textP,
		Max: textP.Add(image.Pt(tw, s.Y)),
	})

	imgP := context.Position(b)
	if b.text.Value() != "" {
		imgP.X += (s.X - ds.X) / 2
		switch b.iconAlign {
		case IconAlignStart:
			imgP.X += buttonEdgeAndImagePadding(context)
		case IconAlignEnd:
			imgP.X += buttonEdgeAndTextPadding(context)
			imgP.X += tw + buttonTextAndImagePadding(context)
		}
	} else {
		imgP.X += (s.X - imgSize.X) / 2
	}
	imgP.Y += (s.Y - imgSize.Y) / 2
	if b.button.isPressed(context) {
		imgP.Y += int(0.5 * context.Scale())
	} else {
		imgP.Y -= int(0.5 * context.Scale())
	}
	appender.AppendChildWidgetWithBounds(&b.icon, image.Rectangle{
		Min: imgP,
		Max: imgP.Add(imgSize),
	})

	return nil
}

func (b *Button) DefaultSize(context *guigui.Context) image.Point {
	return b.defaultSize(context, false)
}

func (b *Button) defaultSize(context *guigui.Context, forceBold bool) image.Point {
	h := defaultButtonSize(context).Y
	var textAndImageW int
	if b.text.Value() != "" {
		if forceBold {
			textAndImageW = b.text.boldTextSize(context).X
		} else {
			textAndImageW = b.text.TextSize(context).X
		}
		textAndImageW += buttonEdgeAndTextPadding(context)
	}
	if b.icon.HasImage() {
		if textAndImageW == 0 {
			textAndImageW += buttonEdgeAndImagePadding(context)
		}
		textAndImageW += defaultIconSize(context)
		if b.text.Value() != "" {
			textAndImageW += buttonTextAndImagePadding(context)
		}
		textAndImageW += buttonEdgeAndImagePadding(context)
	}

	var contentW int
	if b.content != nil {
		contentW = b.content.DefaultSize(context).X
		contentW += 2 * buttonEdgeAndImagePadding(context)
	}

	return image.Pt(max(textAndImageW, contentW), h)
}

func (b *Button) setSharpenCorners(sharpenCorners draw.SharpenCorners) {
	b.button.setSharpenCorners(sharpenCorners)
}

func buttonTextAndImagePadding(context *guigui.Context) int {
	return UnitSize(context) / 4
}

func buttonEdgeAndTextPadding(context *guigui.Context) int {
	return UnitSize(context) / 2
}

func buttonEdgeAndImagePadding(context *guigui.Context) int {
	return UnitSize(context) / 4
}

func (b *Button) iconSize(context *guigui.Context) image.Point {
	s := context.Size(b)
	if b.text.Value() != "" {
		s := min(defaultIconSize(context), s.X, s.Y)
		return image.Pt(s, s)
	}
	r := b.button.radius(context)
	w := max(0, s.X-2*r)
	h := max(defaultIconSize(context), s.Y-2*r)
	return image.Pt(w, h)
}

func (b *Button) setUseAccentColor(use bool) {
	b.button.setUseAccentColor(use)
}
