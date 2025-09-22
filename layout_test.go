// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package guigui_test

import (
	"image"
	"testing"

	"github.com/hajimehoshi/guigui"
)

type dummyWidget struct {
	guigui.DefaultWidget

	size image.Point
}

func (d *dummyWidget) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	return d.size
}

func TestLinearLayoutMeasure(t *testing.T) {
	l := &guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &dummyWidget{
					size: image.Pt(100, 200),
				},
			},
		},
	}
	var context guigui.Context
	if got, want := l.Measure(&context, guigui.Constraints{}), image.Pt(100, 200); got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}

	for _, dir := range []guigui.LayoutDirection{guigui.LayoutDirectionHorizontal, guigui.LayoutDirectionVertical} {
		l2 := &guigui.LinearLayout{
			Direction: dir,
			Items: []guigui.LinearLayoutItem{
				{
					Layout: l,
				},
			},
		}
		if got, want := l2.Measure(&context, guigui.Constraints{}), image.Pt(100, 200); got != want {
			t.Errorf("dir: %v, got: %v, want: %v", dir, got, want)
		}
	}
}
