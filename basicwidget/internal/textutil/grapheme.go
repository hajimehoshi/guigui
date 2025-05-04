// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package textutil

import (
	"iter"

	"github.com/rivo/uniseg"
)

func graphemes(str string) iter.Seq[string] {
	return func(yield func(s string) bool) {
		state := -1
		for len(str) > 0 {
			var cluster string
			cluster, str, _, state = uniseg.StepString(str, state)
			if !yield(cluster) {
				return
			}
		}
	}
}

func BackspaceOnGraphemes(str string, position int) (string, int) {
	var pos int
	for c := range graphemes(str) {
		startPos := pos
		endPos := pos + len(c)
		if position > endPos {
			pos = endPos
			continue
		}
		return str[:startPos] + str[endPos:], startPos
	}
	return str, position
}

func DeleteOnGraphemes(str string, position int) (string, int) {
	var pos int
	for c := range graphemes(str) {
		startPos := pos
		endPos := pos + len(c)
		if position > startPos {
			pos = endPos
			continue
		}
		return str[:startPos] + str[endPos:], startPos
	}
	return str, position
}

func PrevPositionOnGraphemes(str string, position int) int {
	var pos int
	for c := range graphemes(str) {
		startPos := pos
		endPos := pos + len(c)
		if position > endPos {
			pos = endPos
			continue
		}
		return startPos
	}
	return position
}

func NextPositionOnGraphemes(str string, position int) int {
	var pos int
	for c := range graphemes(str) {
		startPos := pos
		endPos := pos + len(c)
		if position > startPos {
			pos = endPos
			continue
		}
		return endPos
	}
	return position
}
