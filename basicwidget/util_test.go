// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget_test

import (
	"slices"
	"testing"

	"github.com/guigui-gui/guigui/basicwidget"
)

func TestMoveItemsInSlice(t *testing.T) {
	type testCase struct {
		slice    []int
		from     int
		count    int
		to       int
		outSlice []int
		outIndex int
	}

	testCases := []testCase{
		{
			slice:    nil,
			from:     0,
			count:    0,
			to:       0,
			outSlice: nil,
			outIndex: 0,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     0,
			count:    1,
			to:       1,
			outSlice: []int{1, 2, 3, 4, 5},
			outIndex: 0,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     0,
			count:    2,
			to:       2,
			outSlice: []int{1, 2, 3, 4, 5},
			outIndex: 0,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     0,
			count:    2,
			to:       3,
			outSlice: []int{3, 1, 2, 4, 5},
			outIndex: 1,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     2,
			count:    1,
			to:       3,
			outSlice: []int{1, 2, 3, 4, 5},
			outIndex: 2,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     2,
			count:    1,
			to:       0,
			outSlice: []int{3, 1, 2, 4, 5},
			outIndex: 0,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     2,
			count:    2,
			to:       0,
			outSlice: []int{3, 4, 1, 2, 5},
			outIndex: 0,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     2,
			count:    2,
			to:       5,
			outSlice: []int{1, 2, 5, 3, 4},
			outIndex: 3,
		},
		{
			slice:    []int{1, 2, 3, 4, 5},
			from:     0,
			count:    5,
			to:       5,
			outSlice: []int{1, 2, 3, 4, 5},
			outIndex: 0,
		},
	}

	for _, tc := range testCases {
		origSlice := make([]int, len(tc.slice))
		copy(origSlice, tc.slice)
		idx := basicwidget.MoveItemsInSlice(tc.slice, tc.from, tc.count, tc.to)
		if got, want := tc.slice, tc.outSlice; !slices.Equal(got, want) {
			t.Errorf("MoveItemsInSlice(%v, %d, %d, %d); got %v, want %v", origSlice, tc.from, tc.count, tc.to, got, want)
		}
		if got, want := idx, tc.outIndex; got != want {
			t.Errorf("MoveItemsInSlice(%v, %d, %d, %d) = %d; want %d", origSlice, tc.from, tc.count, tc.to, got, want)
		}
	}
}
