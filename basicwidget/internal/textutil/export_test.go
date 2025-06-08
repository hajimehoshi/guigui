// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package textutil

import (
	"iter"
)

type Line struct {
	Pos int
	Str string
}

func Lines(width int, str string, autoWrap bool, advance func(str string) float64) iter.Seq[Line] {
	return func(yield func(Line) bool) {
		for l := range lines(width, str, autoWrap, advance) {
			if !yield(Line{
				Pos: l.pos,
				Str: l.str,
			}) {
				return
			}
		}
	}
}

func NextIndentPosition(position float64, indentWidth float64) float64 {
	return nextIndentPosition(position, indentWidth)
}
