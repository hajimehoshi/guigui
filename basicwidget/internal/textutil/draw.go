// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package textutil

import (
	"image"
	"image/color"
	"math"
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

var theAdvanceCache = map[string]float64{}

func Draw(bounds image.Rectangle, dst *ebiten.Image, str string, options *DrawOptions) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))
	op.ColorScale.ScaleWithColor(options.TextColor)
	if dst.Bounds() != bounds {
		dst = dst.SubImage(bounds).(*ebiten.Image)
	}

	op.LineSpacing = options.LineHeight
	// Do not use op.PrimaryAlign due to tab width.

	yOffset := textPositionYOffset(bounds.Size(), str, &options.Options)
	op.GeoM.Translate(0, yOffset)

	// Cache is valid only for this call, as the options can be changed.
	clear(theAdvanceCache)
	cachedAdvance := func(str string, face text.Face, tabWidth float64, keepTailingSpace bool) float64 {
		if v, ok := theAdvanceCache[str]; ok {
			return v
		}
		a := advance(str, face, tabWidth, keepTailingSpace)
		theAdvanceCache[str] = a
		return a
	}

	for pos, line := range lines(bounds.Dx(), str, options.AutoWrap, func(str string) float64 {
		return cachedAdvance(str, options.Face, options.TabWidth, options.KeepTailingSpace)
	}) {
		y := op.GeoM.Element(1, 2)
		if int(math.Ceil(y+options.LineHeight)) < bounds.Min.Y {
			continue
		}
		if int(math.Floor(y)) >= bounds.Max.Y {
			break
		}

		start := pos
		end := pos + len(line) - tailingLineBreakLen(line)

		if options.DrawSelection {
			if start <= options.SelectionEnd && end >= options.SelectionStart {
				start := max(start, options.SelectionStart)
				end := min(end, options.SelectionEnd)
				if start != end {
					// TextPositionFromIndex is too slow here.
					posStart0, posStart1, countStart := textPositionFromIndex(bounds.Dx(), str, start, &options.Options, cachedAdvance)
					posEnd0, _, countEnd := textPositionFromIndex(bounds.Dx(), str, end, &options.Options, cachedAdvance)
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
					posStart0, posStart1, countStart := textPositionFromIndex(bounds.Dx(), str, start, &options.Options, cachedAdvance)
					posEnd0, _, countEnd := textPositionFromIndex(bounds.Dx(), str, end, &options.Options, cachedAdvance)
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
					posStart0, posStart1, countStart := textPositionFromIndex(bounds.Dx(), str, start, &options.Options, cachedAdvance)
					posEnd0, _, countEnd := textPositionFromIndex(bounds.Dx(), str, end, &options.Options, cachedAdvance)
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
		origGeoM := op.GeoM
		if !options.KeepTailingSpace {
			line = strings.TrimRightFunc(line, unicode.IsSpace)
		}
		x := oneLineLeft(bounds.Dx(), line, options.Face, options.HorizontalAlign, options.TabWidth, options.KeepTailingSpace)
		op.GeoM.Translate(x, 0)
		if options.TabWidth == 0 {
			text.Draw(dst, line, options.Face, op)
		} else {
			var origX float64
			for {
				head, tail, ok := strings.Cut(line, "\t")
				text.Draw(dst, head, options.Face, op)
				if !ok {
					break
				}
				x := origX + text.Advance(head, options.Face)
				nextX := nextIndentPosition(x, options.TabWidth)
				op.GeoM.Translate(nextX-origX, 0)
				origX = nextX
				line = tail
			}
		}
		op.GeoM = origGeoM
		op.GeoM.Translate(0, options.LineHeight)
	}
}
