// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package textutil

import (
	"iter"
)

func Lines(width int, str string, autoWrap bool, advance func(str string) float64) iter.Seq2[int, string] {
	return lines(width, str, autoWrap, advance)
}

func NextIndentPosition(position float64, indentWidth float64) float64 {
	return nextIndentPosition(position, indentWidth)
}
