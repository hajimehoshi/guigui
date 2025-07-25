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

func (b *Button) BeforeBuild(context *guigui.Context) {
	b.button.ResetEventHandlers()
}

func (b *Button) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&b.button)
	if b.content != nil {
		appender.AppendChildWidget(b.content)
	}
	appender.AppendChildWidget(&b.text)
	appender.AppendChildWidget(&b.icon)
}

func (b *Button) Build(context *guigui.Context) error {
	context.SetBounds(&b.button, context.Bounds(b), b)

	s := context.ActualSize(b)

	if b.content != nil {
		contentP := context.Position(b)

		if b.button.isPressed(context) {
			contentP.Y += int(0.5 * context.Scale())
		} else {
			contentP.Y -= int(0.5 * context.Scale())
		}
		context.SetBounds(b.content, image.Rectangle{
			Min: contentP,
			Max: contentP.Add(s),
		}, b)
	}

	imgSize := b.iconSize(context)

	tw := b.text.DefaultSizeInContainer(context, context.ActualSize(b).X).X
	if b.textColor != nil {
		b.text.SetColor(b.textColor)
	} else {
		b.text.SetColor(draw.TextColor(context.ColorMode(), context.IsEnabled(b)))
	}
	b.text.SetHorizontalAlign(HorizontalAlignCenter)
	b.text.SetVerticalAlign(VerticalAlignMiddle)

	ds := b.defaultSize(context, false)

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
	context.SetBounds(&b.text, image.Rectangle{
		Min: textP,
		Max: textP.Add(image.Pt(tw, s.Y)),
	}, b)

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
	context.SetBounds(&b.icon, image.Rectangle{
		Min: imgP,
		Max: imgP.Add(imgSize),
	}, b)

	return nil
}

func (b *Button) DefaultSize(context *guigui.Context) image.Point {
	return b.defaultSize(context, false)
}

func (b *Button) defaultSize(context *guigui.Context, forceBold bool) image.Point {
	h := defaultButtonSize(context).Y
	var textAndImageW int
	if b.text.Value() != "" {
		textAndImageW += buttonEdgeAndTextPadding(context)
		if forceBold {
			textAndImageW += b.text.boldTextSize(context, 0).X
		} else {
			textAndImageW += b.text.DefaultSizeInContainer(context, 0).X
		}
	}
	if b.icon.HasImage() {
		if textAndImageW == 0 {
			textAndImageW += buttonEdgeAndImagePadding(context)
		}
		if b.text.Value() != "" {
			textAndImageW += buttonTextAndImagePadding(context)
		}
		textAndImageW += defaultIconSize(context)
		textAndImageW += buttonEdgeAndImagePadding(context)
	} else {
		textAndImageW += buttonEdgeAndTextPadding(context)
	}

	var contentW int
	if b.content != nil {
		contentW = b.content.DefaultSize(context).X
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
	s := context.ActualSize(b)
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
