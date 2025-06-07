// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package textutil

import (
	"image"
	"image/color"
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DrawOptions struct {
	Options

	TextColor color.Color

	DrawSelection  bool
	SelectionStart int
	SelectionEnd   int
	SelectionColor color.Color

	DrawComposition          bool
	CompositionStart         int
	CompositionEnd           int
	CompositionActiveStart   int
	CompositionActiveEnd     int
	InactiveCompositionColor color.Color
	ActiveCompositionColor   color.Color
	CompositionBorderWidth   float32
}

func Draw(bounds image.Rectangle, dst *ebiten.Image, str string, options *DrawOptions) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))
	op.ColorScale.ScaleWithColor(options.TextColor)
	if dst.Bounds() != bounds {
		dst = dst.SubImage(bounds).(*ebiten.Image)
	}

	op.LineSpacing = options.LineHeight

	switch options.HorizontalAlign {
	case HorizontalAlignStart, HorizontalAlignLeft:
		// TODO: For RTL languages, HorizontalAlignStart should be the same as HorizontalAlignRight.
		op.PrimaryAlign = text.AlignStart
	case HorizontalAlignCenter:
		op.GeoM.Translate(float64(bounds.Dx())/2, 0)
		op.PrimaryAlign = text.AlignCenter
	case HorizontalAlignEnd, HorizontalAlignRight:
		// TODO: For RTL languages, HorizontalAlignEnd should be the same as HorizontalAlignLeft.
		op.GeoM.Translate(float64(bounds.Dx()), 0)
		op.PrimaryAlign = text.AlignEnd
	}

	yOffset := textPositionYOffset(bounds.Size(), str, &options.Options)
	op.GeoM.Translate(0, yOffset)

	for pos, line := range lines(bounds.Dx(), str, options.AutoWrap, func(str string) float64 {
		return advance(str, options.Face, options.KeepTailingSpace)
	}) {
		start := pos
		end := pos + len(line) - tailingLineBreakLen(line)

		if options.DrawSelection {
			if start <= options.SelectionEnd && end >= options.SelectionStart {
				start := max(start, options.SelectionStart)
				end := min(end, options.SelectionEnd)
				if start != end {
					posStart0, posStart1, countStart := TextPositionFromIndex(bounds.Dx(), str, start, &options.Options)
					posEnd0, _, countEnd := TextPositionFromIndex(bounds.Dx(), str, end, &options.Options)
					if countStart > 0 && countEnd > 0 {
						posStart := posStart0
						if countStart == 2 {
							posStart = posStart1
						}
						posEnd := posEnd0
						x := float32(posStart.X) + float32(bounds.Min.X)
						y := float32(posStart.Top) + float32(bounds.Min.Y)
						width := float32(posEnd.X - posStart.X)
						height := float32(posStart.Bottom - posStart.Top)
						vector.DrawFilledRect(dst, x, y, width, height, options.SelectionColor, false)
					}
				}
			}
		}

		if options.DrawComposition {
			if start <= options.CompositionEnd && end >= options.CompositionStart {
				start := max(start, options.CompositionStart)
				end := min(end, options.CompositionEnd)
				if start != end {
					posStart0, posStart1, countStart := TextPositionFromIndex(bounds.Dx(), str, start, &options.Options)
					posEnd0, _, countEnd := TextPositionFromIndex(bounds.Dx(), str, end, &options.Options)
					if countStart > 0 && countEnd > 0 {
						posStart := posStart0
						if countStart == 2 {
							posStart = posStart1
						}
						posEnd := posEnd0
						x := float32(posStart.X) + float32(bounds.Min.X)
						y := float32(posStart.Bottom) + float32(bounds.Min.Y) - options.CompositionBorderWidth
						w := float32(posEnd.X - posStart.X)
						h := options.CompositionBorderWidth
						vector.DrawFilledRect(dst, x, y, w, h, options.InactiveCompositionColor, false)
					}
				}
			}
			if start <= options.CompositionActiveEnd && end >= options.CompositionActiveStart {
				start := max(start, options.CompositionActiveStart)
				end := min(end, options.CompositionActiveEnd)
				if start != end {
					posStart0, posStart1, countStart := TextPositionFromIndex(bounds.Dx(), str, start, &options.Options)
					posEnd0, _, countEnd := TextPositionFromIndex(bounds.Dx(), str, end, &options.Options)
					if countStart > 0 && countEnd > 0 {
						posStart := posStart0
						if countStart == 2 {
							posStart = posStart1
						}
						posEnd := posEnd0
						x := float32(posStart.X) + float32(bounds.Min.X)
						y := float32(posStart.Bottom) + float32(bounds.Min.Y) - options.CompositionBorderWidth
						w := float32(posEnd.X - posStart.X)
						h := options.CompositionBorderWidth
						vector.DrawFilledRect(dst, x, y, w, h, options.ActiveCompositionColor, false)
					}
				}
			}
		}

		// Draw the text.
		if options.KeepTailingSpace {
			text.Draw(dst, line, options.Face, op)
		} else {
			text.Draw(dst, strings.TrimRightFunc(line, unicode.IsSpace), options.Face, op)
		}
		op.GeoM.Translate(0, options.LineHeight)
	}
}
