// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget_test

import (
	"fmt"
	"testing"

	"github.com/guigui-gui/guigui/basicwidget"
)

func TestReplaceNewLineWithSpace(t *testing.T) {
	testCases := []struct {
		text          string
		start         int
		end           int
		shiftIndex    int
		outText       string
		outStart      int
		outEnd        int
		outShiftIndex int
	}{
		{
			text:          "",
			start:         0,
			end:           0,
			shiftIndex:    -1,
			outText:       "",
			outStart:      0,
			outEnd:        0,
			outShiftIndex: -1,
		},
		{
			text:          "Hello,\nWorld!",
			start:         7,
			end:           13,
			shiftIndex:    7,
			outText:       "Hello, World!",
			outStart:      7,
			outEnd:        13,
			outShiftIndex: 7,
		},
		{
			text:          "Hello,\nWorld!",
			start:         7,
			end:           13,
			shiftIndex:    13,
			outText:       "Hello, World!",
			outStart:      7,
			outEnd:        13,
			outShiftIndex: 13,
		},
		{
			text:          "Hello,\r\nWorld!",
			start:         6,
			end:           6,
			shiftIndex:    6,
			outText:       "Hello, World!",
			outStart:      6,
			outEnd:        6,
			outShiftIndex: 6,
		},
		{
			text:          "Hello,\r\nWorld!",
			start:         8,
			end:           14,
			shiftIndex:    14,
			outText:       "Hello, World!",
			outStart:      7,
			outEnd:        13,
			outShiftIndex: 13,
		},
		{
			text:          "Hello,\u2028World!",
			start:         9,
			end:           15,
			shiftIndex:    15,
			outText:       "Hello, World!",
			outStart:      7,
			outEnd:        13,
			outShiftIndex: 13,
		},
		{
			text:          "Hello,\r\nWorld!",
			start:         6,
			end:           7, // In between \r and \n
			shiftIndex:    7,
			outText:       "Hello, World!",
			outStart:      6,
			outEnd:        7,
			outShiftIndex: 7,
		},
		{
			text:          "a\r\u2028\nb",
			start:         5,
			end:           7,
			shiftIndex:    5,
			outText:       "a   b",
			outStart:      3,
			outEnd:        5,
			outShiftIndex: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.text), func(t *testing.T) {
			gotText, gotStart, gotEnd, gotShiftIndex := basicwidget.ReplaceNewLinesWithSpace(tc.text, tc.start, tc.end, tc.shiftIndex)
			if gotText != tc.outText || gotStart != tc.outStart || gotEnd != tc.outEnd || gotShiftIndex != tc.outShiftIndex {
				t.Errorf("got (%q, %d, %d, %d), want (%q, %d, %d, %d)", gotText, gotStart, gotEnd, gotShiftIndex, tc.outText, tc.outStart, tc.outEnd, tc.outShiftIndex)
			}
		})
	}
}
