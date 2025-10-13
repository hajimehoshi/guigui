// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package textutil_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/guigui-gui/guigui/basicwidget/internal/textutil"
)

func TestNoWrapLines(t *testing.T) {
	testCases := []struct {
		str       string
		positions []int
		lines     []string
	}{
		{
			str:       "Hello, World!",
			positions: []int{0},
			lines:     []string{"Hello, World!"},
		},
		{
			str:       "Hello,\nWorld!",
			positions: []int{0, 7},
			lines:     []string{"Hello,\n", "World!"},
		},
		{
			str:       "Hello,\nWorld!\n",
			positions: []int{0, 7, 14},
			lines:     []string{"Hello,\n", "World!\n", ""},
		},
		{
			str:       "Hello,\rWorld!",
			positions: []int{0, 7},
			lines:     []string{"Hello,\r", "World!"},
		},
		{
			str:       "Hello,\u0085World!",
			positions: []int{0, 8}, // U+0085 is 2 bytes in UTF-8.
			lines:     []string{"Hello,\u0085", "World!"},
		},
		{
			str:       "Hello,\n\nWorld!",
			positions: []int{0, 7, 8},
			lines:     []string{"Hello,\n", "\n", "World!"},
		},
		{
			str:       "Hello,\r\nWorld!",
			positions: []int{0, 8},
			lines:     []string{"Hello,\r\n", "World!"},
		},
		{
			str:       "Hello,\n\rWorld!",
			positions: []int{0, 7, 8},
			lines:     []string{"Hello,\n", "\r", "World!"},
		},
		{
			str:       "",
			positions: []int{0},
			lines:     []string{""},
		},
		{
			str:       "\n",
			positions: []int{0, 1},
			lines:     []string{"\n", ""},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			var gotPositions []int
			var gotLines []string
			for l := range textutil.Lines(0, tc.str, false, nil) {
				gotPositions = append(gotPositions, l.Pos)
				gotLines = append(gotLines, l.Str)
			}
			if !slices.Equal(gotPositions, tc.positions) {
				t.Errorf("got positions %v, want %v", gotPositions, tc.positions)
			}
			if !slices.Equal(gotLines, tc.lines) {
				t.Errorf("got lines %v, want %v", gotLines, tc.lines)
			}
		})
	}
}

func TestNextIndentPosition(t *testing.T) {
	testCases := []struct {
		position    float64
		indentWidth float64
		expected    float64
	}{
		{
			position:    0,
			indentWidth: 10.5,
			expected:    10.5,
		},
		{
			position:    104,
			indentWidth: 10.5,
			expected:    105,
		},
		{
			position:    104.9995,
			indentWidth: 10.5,
			expected:    105,
		},
		{
			position:    105,
			indentWidth: 10.5,
			expected:    115.5,
		},
		{
			position:    105.0001,
			indentWidth: 10.5,
			expected:    115.5,
		},
		{
			position:    106,
			indentWidth: 10.5,
			expected:    115.5,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("position=%f indentWidth=%f", tc.position, tc.indentWidth), func(t *testing.T) {
			got := textutil.NextIndentPosition(tc.position, tc.indentWidth)
			if got != tc.expected {
				t.Errorf("got %f, want %f", got, tc.expected)
			}
		})
	}
}
