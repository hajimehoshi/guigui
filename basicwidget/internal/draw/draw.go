// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package draw

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
)

type RoundedRectBorderType int

const (
	RoundedRectBorderTypeRegular RoundedRectBorderType = iota
	RoundedRectBorderTypeInset
	RoundedRectBorderTypeOutset
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	// This is hacky, but WritePixels is better than Fill in term of automatic texture packing.
	whiteImage.WritePixels(pix)
}

func appendRectVectorPath(path *vector.Path, x0, y0, x1, y1 float32, radius float32) {
	path.MoveTo(x0, y0)
	path.LineTo(x1, y0)
	path.LineTo(x1, y1)
	path.LineTo(x0, y1)
	path.LineTo(x0, y0)
}

func appendRoundedRectVectorPath(path *vector.Path, rx0, ry0, rx1, ry1 float32, radius float32) {
	x0 := rx0
	x1 := rx0 + radius
	x2 := rx1 - radius
	x3 := rx1
	y0 := ry0
	y1 := ry0 + radius
	y2 := ry1 - radius
	y3 := ry1

	path.MoveTo(x1, y0)
	path.LineTo(x2, y0)
	path.ArcTo(x3, y0, x3, y1, radius)
	path.LineTo(x3, y2)
	path.ArcTo(x3, y3, x2, y3, radius)
	path.LineTo(x1, y3)
	path.ArcTo(x0, y3, x0, y2, radius)
	path.LineTo(x0, y1)
	path.ArcTo(x0, y0, x1, y0, radius)
}

type imageKey struct {
	radius      int
	borderWidth float32
	borderType  RoundedRectBorderType
	colorMode   guigui.ColorMode
}

var (
	whiteRoundedRects       = map[imageKey]*ebiten.Image{}
	whiteRoundedShadowRects = map[imageKey]*ebiten.Image{}
	whiteRectBorders        = map[imageKey]*ebiten.Image{}
	whiteRoundedRectBorders = map[imageKey]*ebiten.Image{}
)

func ensureWhiteRoundedRect(radius int) *ebiten.Image {
	key := imageKey{
		radius: radius,
	}
	if img, ok := whiteRoundedRects[key]; ok {
		return img
	}

	s := radius * 3
	img := ebiten.NewImage(s, s)

	var path vector.Path
	appendRoundedRectVectorPath(&path, 0, 0, float32(s), float32(s), float32(radius))
	path.Close()

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	img.DrawTriangles(vs, is, whiteSubImage, op)

	whiteRoundedRects[key] = img

	return img
}

func ensureWhiteRoundedShadowRect(radius int) *ebiten.Image {
	key := imageKey{
		radius: radius,
	}
	if img, ok := whiteRoundedShadowRects[key]; ok {
		return img
	}

	s := radius * 3
	img := ebiten.NewImage(s, s)

	pix := make([]byte, 4*s*s)

	easeInQuad := func(x float64) float64 {
		return x * x
	}

	for j := 0; j < radius; j++ {
		for i := 0; i < radius; i++ {
			x := float64(radius - i)
			y := float64(radius - j)
			d := max(0, min(1, math.Hypot(x, y)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
		for i := radius; i < 2*radius; i++ {
			d := max(0, min(1, float64(radius-j)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
		for i := 2 * radius; i < 3*radius; i++ {
			x := float64(i - 2*radius)
			y := float64(radius - j)
			d := max(0, min(1, math.Hypot(x, y)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
	}
	for j := radius; j < 2*radius; j++ {
		for i := 0; i < radius; i++ {
			d := max(0, min(1, float64(radius-i)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
		for i := radius; i < 2*radius; i++ {
			pix[4*(j*s+i)] = 0xff
			pix[4*(j*s+i)+1] = 0xff
			pix[4*(j*s+i)+2] = 0xff
			pix[4*(j*s+i)+3] = 0xff
		}
		for i := 2 * radius; i < 3*radius; i++ {
			d := max(0, min(1, float64(i-2*radius)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
	}
	for j := 2 * radius; j < 3*radius; j++ {
		for i := 0; i < radius; i++ {
			x := float64(radius - i)
			y := float64(j - 2*radius)
			d := max(0, min(1, math.Hypot(x, y)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
		for i := radius; i < 2*radius; i++ {
			d := max(0, min(1, float64(j-2*radius)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
		for i := 2 * radius; i < 3*radius; i++ {
			x := float64(i - 2*radius)
			y := float64(j - 2*radius)
			d := max(0, min(1, math.Hypot(x, y)/float64(radius)))
			a := byte(0xff * easeInQuad(1-d))
			pix[4*(j*s+i)] = a
			pix[4*(j*s+i)+1] = a
			pix[4*(j*s+i)+2] = a
			pix[4*(j*s+i)+3] = a
		}
	}

	img.WritePixels(pix)

	whiteRoundedShadowRects[key] = img

	return img
}

func ensureWhiteRectBorder(partSize int, borderWidth float32, borderType RoundedRectBorderType, colorMode guigui.ColorMode) *ebiten.Image {
	key := imageKey{
		radius:      partSize,
		borderWidth: borderWidth,
		borderType:  borderType,
		colorMode:   colorMode,
	}
	if img, ok := whiteRectBorders[key]; ok {
		return img
	}

	img := whiteRoundedRectBorder(partSize, borderWidth, borderType, colorMode, appendRectVectorPath)
	whiteRectBorders[key] = img
	return img
}

func ensureWhiteRoundedRectBorder(radius int, borderWidth float32, borderType RoundedRectBorderType, colorMode guigui.ColorMode) *ebiten.Image {
	key := imageKey{
		radius:      radius,
		borderWidth: borderWidth,
		borderType:  borderType,
		colorMode:   colorMode,
	}
	if img, ok := whiteRoundedRectBorders[key]; ok {
		return img
	}

	img := whiteRoundedRectBorder(radius, borderWidth, borderType, colorMode, appendRoundedRectVectorPath)
	whiteRoundedRectBorders[key] = img
	return img
}

func whiteRoundedRectBorder(radius int, borderWidth float32, borderType RoundedRectBorderType, colorMode guigui.ColorMode, appendPathFunc func(path *vector.Path, rx0, ry0, rx1, ry1 float32, radius float32)) *ebiten.Image {
	// Use it's own anti-aliasing instead of Ebitengine's anti-aliasing for higher quality result.
	// Ebitengine's anti-aliasing just scales vertice and doesn't create finer paths for anti-aliasing scale.
	const aaScale = 2
	radius *= aaScale
	s := radius * 3
	inset := borderWidth * aaScale

	var path vector.Path
	appendPathFunc(&path, 0, 0, float32(s), float32(s), float32(radius))
	switch borderType {
	case RoundedRectBorderTypeRegular:
		appendPathFunc(&path, inset, inset, float32(s)-inset, float32(s)-inset, float32(radius)-inset)
	case RoundedRectBorderTypeInset:
		// Use a thicker border for the dark mode, as colors tend to be contracting colors.
		if colorMode == guigui.ColorModeDark {
			appendPathFunc(&path, inset, inset*2, float32(s)-inset, float32(s)-inset/2, float32(radius)-inset/2)
		} else {
			appendPathFunc(&path, inset, inset*3/2, float32(s)-inset, float32(s)-inset/2, float32(radius)-inset/2)
		}
	case RoundedRectBorderTypeOutset:
		// Use a thicker border for the dark mode, as colors tend to be contracting colors.
		if colorMode == guigui.ColorModeDark {
			appendPathFunc(&path, inset, inset/2, float32(s)-inset, float32(s)-inset*2, float32(radius)-inset/2)
		} else {
			appendPathFunc(&path, inset, inset/2, float32(s)-inset, float32(s)-inset*3/2, float32(radius)-inset/2)
		}
	}
	path.Close()

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.FillRuleEvenOdd

	offscreen := ebiten.NewImage(s, s)
	offscreen.DrawTriangles(vs, is, whiteSubImage, op)
	defer offscreen.Deallocate()

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.0/aaScale, 1.0/aaScale)
	op2.Filter = ebiten.FilterLinear
	img := ebiten.NewImage(s/aaScale, s/aaScale)
	img.DrawImage(offscreen, op2)

	return img
}

type SharpenCorners struct {
	UpperLeft  bool
	UpperRight bool
	LowerLeft  bool
	LowerRight bool
}

func (s *SharpenCorners) bools() [3][3]bool {
	return [3][3]bool{
		{!s.UpperLeft, true, !s.UpperRight},
		{true, true, true},
		{!s.LowerLeft, true, !s.LowerRight},
	}
}

func (s *SharpenCorners) invertedBools() [3][3]bool {
	bs := s.bools()
	for j := range 3 {
		for i := range 3 {
			bs[j][i] = !bs[j][i]
		}
	}
	return bs
}

func DrawRoundedRect(context *guigui.Context, dst *ebiten.Image, bounds image.Rectangle, clr color.Color, radius int) {
	DrawRoundedRectWithSharpenCorners(context, dst, bounds, clr, radius, SharpenCorners{})
}

func DrawRoundedRectWithSharpenCorners(context *guigui.Context, dst *ebiten.Image, bounds image.Rectangle, clr color.Color, radius int, sharpenCorners SharpenCorners) {
	if !dst.Bounds().Overlaps(bounds) {
		return
	}
	if bounds.Dx()/2-1 < radius {
		radius = bounds.Dx()/2 - 1
	}
	if bounds.Dy()/2-1 < radius {
		radius = bounds.Dy()/2 - 1
	}
	if radius == 0 {
		return
	}

	if sharpenCorners == (SharpenCorners{}) {
		drawNinePatch(dst, bounds, ensureWhiteRoundedRect(radius), clr, clr)
		return
	}

	drawNinePatchParts(dst, bounds, ensureWhiteRoundedRect(radius), clr, clr, sharpenCorners.bools())
	if !dst.Bounds().Intersect(bounds).Empty() {
		vs, is, op := ninePatchVertices(bounds, whiteImage.Bounds(), radius, radius, clr, clr, sharpenCorners.invertedBools())
		dst.DrawTriangles32(vs, is, whiteImage, op)
	}
}

func DrawRoundedShadowRect(context *guigui.Context, dst *ebiten.Image, bounds image.Rectangle, clr color.Color, radius int) {
	if !dst.Bounds().Overlaps(bounds) {
		return
	}
	if bounds.Dx()/2-1 < radius {
		radius = bounds.Dx()/2 - 1
	}
	if bounds.Dy()/2-1 < radius {
		radius = bounds.Dy()/2 - 1
	}
	if radius == 0 {
		return
	}
	drawNinePatch(dst, bounds, ensureWhiteRoundedShadowRect(radius), clr, clr)
}

func DrawRoundedRectBorder(context *guigui.Context, dst *ebiten.Image, bounds image.Rectangle, clr1, clr2 color.Color, radius int, borderWidth float32, borderType RoundedRectBorderType) {
	DrawRoundedRectBorderWithSharpenCorners(context, dst, bounds, clr1, clr2, radius, borderWidth, borderType, SharpenCorners{})
}

func DrawRoundedRectBorderWithSharpenCorners(context *guigui.Context, dst *ebiten.Image, bounds image.Rectangle, clr1, clr2 color.Color, radius int, borderWidth float32, borderType RoundedRectBorderType, sharpenCorners SharpenCorners) {
	if !dst.Bounds().Overlaps(bounds) {
		return
	}
	if bounds.Dx()/2-1 < radius {
		radius = bounds.Dx()/2 - 1
	}
	if bounds.Dy()/2-1 < radius {
		radius = bounds.Dy()/2 - 1
	}
	if radius == 0 {
		return
	}

	if sharpenCorners == (SharpenCorners{}) {
		drawNinePatch(dst, bounds, ensureWhiteRoundedRectBorder(radius, borderWidth, borderType, context.ColorMode()), clr1, clr2)
		return
	}

	drawNinePatchParts(dst, bounds, ensureWhiteRoundedRectBorder(radius, borderWidth, borderType, context.ColorMode()), clr1, clr2, sharpenCorners.bools())
	drawNinePatchParts(dst, bounds, ensureWhiteRectBorder(radius, borderWidth, borderType, context.ColorMode()), clr1, clr2, sharpenCorners.invertedBools())
}

func drawNinePatch(dst *ebiten.Image, bounds image.Rectangle, src *ebiten.Image, clr1, clr2 color.Color) {
	drawNinePatchParts(dst, bounds, src, clr1, clr2, [3][3]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	})
}

func drawNinePatchParts(dst *ebiten.Image, bounds image.Rectangle, src *ebiten.Image, clr1, clr2 color.Color, renderingParts [3][3]bool) {
	if dst.Bounds().Intersect(bounds).Empty() {
		return
	}
	cornerW, cornerH := src.Bounds().Dx()/3, src.Bounds().Dy()/3
	vs, is, op := ninePatchVertices(bounds, src.Bounds(), cornerW, cornerH, clr1, clr2, renderingParts)
	dst.DrawTriangles32(vs, is, src, op)
}

func ninePatchVertices(dstBounds, srcBounds image.Rectangle, cornerW, cornerH int, clr1, clr2 color.Color, renderingParts [3][3]bool) ([]ebiten.Vertex, []uint32, *ebiten.DrawTrianglesOptions) {
	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha

	var c1 [4]uint32
	var c2 [4]uint32
	c1[0], c1[1], c1[2], c1[3] = clr1.RGBA()
	c2[0], c2[1], c2[2], c2[3] = clr2.RGBA()

	mix := func(a, b uint32, rate float32) float32 {
		return (1-rate)*float32(a)/0xffff + rate*float32(b)/0xffff
	}
	rates := [...]float32{
		0,
		float32(cornerH) / float32(dstBounds.Dy()),
		float32(dstBounds.Dy()-cornerH) / float32(dstBounds.Dy()),
		1,
	}
	var clrs [4][4]float32
	for j, rate := range rates {
		for i := range clrs[j] {
			clrs[j][i] = mix(c1[i], c2[i], rate)
		}
	}

	var vs []ebiten.Vertex
	var is []uint32
	for j := range 3 {
		for i := range 3 {
			if !renderingParts[j][i] {
				continue
			}

			var scaleX float32 = 1.0
			var scaleY float32 = 1.0
			var tx, ty int

			switch i {
			case 0:
				tx = dstBounds.Min.X
			case 1:
				scaleX = float32(dstBounds.Dx()-2*cornerW) / float32(cornerW)
				tx = dstBounds.Min.X + cornerW
			case 2:
				tx = dstBounds.Max.X - cornerW
			}
			switch j {
			case 0:
				ty = dstBounds.Min.Y
			case 1:
				scaleY = float32(dstBounds.Dy()-2*cornerH) / float32(cornerH)
				ty = dstBounds.Min.Y + cornerH
			case 2:
				ty = dstBounds.Max.Y - cornerH
			}

			tx0 := float32(tx)
			tx1 := float32(tx) + scaleX*float32(cornerW)
			ty0 := float32(ty)
			ty1 := float32(ty) + scaleY*float32(cornerH)
			sx0 := float32(i * srcBounds.Dx() / 3)
			sy0 := float32(j * srcBounds.Dy() / 3)
			sx1 := float32(i+1) * float32(srcBounds.Dx()/3)
			sy1 := float32(j+1) * float32(srcBounds.Dy()/3)

			base := uint32(len(vs))
			vs = append(vs,
				ebiten.Vertex{
					DstX:   tx0,
					DstY:   ty0,
					SrcX:   sx0,
					SrcY:   sy0,
					ColorR: clrs[j][0],
					ColorG: clrs[j][1],
					ColorB: clrs[j][2],
					ColorA: clrs[j][3],
				},
				ebiten.Vertex{
					DstX:   tx1,
					DstY:   ty0,
					SrcX:   sx1,
					SrcY:   sy0,
					ColorR: clrs[j][0],
					ColorG: clrs[j][1],
					ColorB: clrs[j][2],
					ColorA: clrs[j][3],
				},
				ebiten.Vertex{
					DstX:   tx0,
					DstY:   ty1,
					SrcX:   sx0,
					SrcY:   sy1,
					ColorR: clrs[j+1][0],
					ColorG: clrs[j+1][1],
					ColorB: clrs[j+1][2],
					ColorA: clrs[j+1][3],
				},
				ebiten.Vertex{
					DstX:   tx1,
					DstY:   ty1,
					SrcX:   sx1,
					SrcY:   sy1,
					ColorR: clrs[j+1][0],
					ColorG: clrs[j+1][1],
					ColorB: clrs[j+1][2],
					ColorA: clrs[j+1][3],
				},
			)
			is = append(is, base+0, base+1, base+2, base+1, base+2, base+3)
		}
	}

	return vs, is, op
}
