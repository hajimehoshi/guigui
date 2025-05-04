// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
)

type DragDropOverlay struct {
	guigui.DefaultWidget

	object any

	onDropped func(object any)
}

func (d *DragDropOverlay) SetOnDropped(f func(object any)) {
	d.onDropped = f
}

func (d *DragDropOverlay) IsDragging() bool {
	return d.object != nil
}

func (d *DragDropOverlay) Start(object any) {
	d.object = object
}

func (d *DragDropOverlay) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if d.object != nil {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if image.Pt(ebiten.CursorPosition()).In(context.VisibleBounds(d)) {
				if d.onDropped != nil {
					d.onDropped(d.object)
				}
			}
			d.object = nil
			return guigui.HandleInputResult{}
		}
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			d.object = nil
		}
		return guigui.HandleInputResult{}
	}

	return guigui.HandleInputResult{}
}
