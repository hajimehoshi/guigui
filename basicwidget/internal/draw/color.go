// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package draw

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/oklab"

	"github.com/hajimehoshi/guigui"
)

func EqualColor(c0, c1 color.Color) bool {
	if c0 == c1 {
		return true
	}
	if c0 == nil || c1 == nil {
		return false
	}
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1 && a0 == a1
}

var (
	blue   = oklab.OklchModel.Convert(color.RGBA{R: 0x00, G: 0x5a, B: 0xff, A: 0xff}).(oklab.Oklch)
	green  = oklab.OklchModel.Convert(color.RGBA{R: 0x03, G: 0xaf, B: 0x7a, A: 0xff}).(oklab.Oklch)
	yellow = oklab.OklchModel.Convert(color.RGBA{R: 0xff, G: 0xf1, B: 0x00, A: 0xff}).(oklab.Oklch)
	red    = oklab.OklchModel.Convert(color.RGBA{R: 0xff, G: 0x4b, B: 0x00, A: 0xff}).(oklab.Oklch)
)

var (
	white = oklab.OklchModel.Convert(color.White).(oklab.Oklch)
	black = oklab.OklchModel.Convert(oklab.Oklab{L: 0.2, A: 0, B: 0, Alpha: 1}).(oklab.Oklch)
	gray  = oklab.OklchModel.Convert(oklab.Oklab{L: 0.6, A: 0, B: 0, Alpha: 1}).(oklab.Oklch)
)

type ColorType int

const (
	ColorTypeBase ColorType = iota
	ColorTypeAccent
	ColorTypeInfo
	ColorTypeSuccess
	ColorTypeWarning
	ColorTypeDanger
)

func Color(colorMode guigui.ColorMode, typ ColorType, lightnessInLightMode float64) color.Color {
	return Color2(colorMode, typ, lightnessInLightMode, 1-lightnessInLightMode)
}

func Color2(colorMode guigui.ColorMode, typ ColorType, lightnessInLightMode, lightnessInDarkMode float64) color.Color {
	var base color.Color
	switch typ {
	case ColorTypeBase:
		base = gray
	case ColorTypeAccent:
		base = blue
	case ColorTypeInfo:
		base = blue
	case ColorTypeSuccess:
		base = green
	case ColorTypeWarning:
		base = yellow
	case ColorTypeDanger:
		base = red
	default:
		panic(fmt.Sprintf("basicwidget: invalid color type: %d", typ))
	}
	switch colorMode {
	case guigui.ColorModeLight:
		return getColor(base, lightnessInLightMode, black, white)
	case guigui.ColorModeDark:
		return getColor(base, lightnessInDarkMode, black, white)
	default:
		panic(fmt.Sprintf("basicwidget: invalid color mode: %d", colorMode))
	}
}

func getColor(base color.Color, lightness float64, back, front color.Color) color.Color {
	c0 := oklab.OklchModel.Convert(back).(oklab.Oklch)
	c1 := oklab.OklchModel.Convert(front).(oklab.Oklch)
	l := oklab.OklchModel.Convert(base).(oklab.Oklch).L
	l = max(min(l, c1.L), c0.L)
	l2 := c0.L*(1-lightness) + c1.L*lightness
	if l2 < l {
		rate := (l2 - c0.L) / (l - c0.L)
		return MixColor(c0, base, rate)
	}
	rate := (l2 - l) / (c1.L - l)
	return MixColor(base, c1, rate)
}

func MixColor(clr0, clr1 color.Color, rate float64) color.Color {
	if rate == 0 {
		return clr0
	}
	if rate == 1 {
		return clr1
	}
	okClr0 := oklab.OklabModel.Convert(clr0).(oklab.Oklab)
	okClr1 := oklab.OklabModel.Convert(clr1).(oklab.Oklab)
	return oklab.Oklab{
		L:     okClr0.L*(1-rate) + okClr1.L*rate,
		A:     okClr0.A*(1-rate) + okClr1.A*rate,
		B:     okClr0.B*(1-rate) + okClr1.B*rate,
		Alpha: okClr0.Alpha*(1-rate) + okClr1.Alpha*rate,
	}
}

func ScaleAlpha(clr color.Color, alpha float64) color.Color {
	r, g, b, a := clr.RGBA()
	r = uint32(float64(r) * alpha)
	g = uint32(float64(g) * alpha)
	b = uint32(float64(b) * alpha)
	a = uint32(float64(a) * alpha)
	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

func BorderColors(colorMode guigui.ColorMode, borderType RoundedRectBorderType, accent bool) (color.Color, color.Color) {
	typ1 := ColorTypeBase
	typ2 := ColorTypeBase
	if accent {
		typ1 = ColorTypeAccent
	}
	switch borderType {
	case RoundedRectBorderTypeRegular:
		return Color2(colorMode, typ1, 0.8, 0.1), Color2(colorMode, typ2, 0.8, 0.1)
	case RoundedRectBorderTypeInset:
		return Color2(colorMode, typ1, 0.7, 0), Color2(colorMode, typ2, 0.85, 0.15)
	case RoundedRectBorderTypeOutset:
		return Color2(colorMode, typ1, 0.85, 0.5), Color2(colorMode, typ2, 0.7, 0.2)
	}
	panic(fmt.Sprintf("draw: invalid border type: %d", borderType))
}

func TextColor(colorMode guigui.ColorMode, enabled bool) color.Color {
	if !enabled {
		return Color(colorMode, ColorTypeBase, 0.5)
	}
	return Color(colorMode, ColorTypeBase, 0.1)
}

func ControlColor(colorMode guigui.ColorMode, enabled bool) color.Color {
	if !enabled {
		return Color2(colorMode, ColorTypeBase, 0.9, 0.1)
	}
	return Color2(colorMode, ColorTypeBase, 1, 0.3)
}

func ThumbColor(colorMode guigui.ColorMode, enabled bool) color.Color {
	if !enabled {
		return Color2(colorMode, ColorTypeBase, 0.9, 0.55)
	}
	return Color2(colorMode, ColorTypeBase, 1, 0.6)
}
