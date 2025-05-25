// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package guigui

import (
	"errors"
	"fmt"
	"image"
	"maps"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

type bounds3D struct {
	bounds      image.Rectangle
	zDelta      int
	visible     bool // For hit testing.
	passThrough bool // For hit testing.
}

func bounds3DFromWidget(context *Context, widget Widget) (bounds3D, bool) {
	bounds := context.VisibleBounds(widget)
	if bounds.Empty() {
		return bounds3D{}, false
	}
	return bounds3D{
		bounds:      bounds,
		zDelta:      widget.ZDelta(),
		visible:     widget.widgetState().isVisible(),
		passThrough: widget.PassThrough(),
	}, true
}

type widgetsAndVisibleBounds struct {
	bounds3Ds       map[*widgetState]bounds3D
	currentBounds3D map[*widgetState]bounds3D
}

func (w *widgetsAndVisibleBounds) reset() {
	clear(w.bounds3Ds)
}

func (w *widgetsAndVisibleBounds) append(context *Context, widget Widget) {
	if w.bounds3Ds == nil {
		w.bounds3Ds = map[*widgetState]bounds3D{}
	}
	b, ok := bounds3DFromWidget(context, widget)
	if !ok {
		return
	}
	w.bounds3Ds[widget.widgetState()] = b
}

func (w *widgetsAndVisibleBounds) equals(context *Context, currentWidgets []Widget) bool {
	if w.currentBounds3D == nil {
		w.currentBounds3D = map[*widgetState]bounds3D{}
	} else {
		clear(w.currentBounds3D)
	}
	for _, widget := range currentWidgets {
		b, ok := bounds3DFromWidget(context, widget)
		if !ok {
			continue
		}
		w.currentBounds3D[widget.widgetState()] = b
	}
	return maps.Equal(w.bounds3Ds, w.currentBounds3D)
}

func (w *widgetsAndVisibleBounds) redrawIfDifferentParentZ(app *app) {
	for widgetState, bounds3D := range w.bounds3Ds {
		if bounds3D.zDelta != 0 {
			app.requestRedraw(bounds3D.bounds)
			requestRedraw(widgetState)
		}
	}
}

type CustomDrawFunc func(dst, widgetImage *ebiten.Image, op *ebiten.DrawImageOptions)

type widgetState struct {
	root bool

	position    image.Point
	widthPlus1  int
	heightPlus1 int

	parent   Widget
	children []Widget
	prev     widgetsAndVisibleBounds

	hidden       bool
	disabled     bool
	transparency float64
	customDraw   CustomDrawFunc

	offscreen *ebiten.Image

	dirty                 bool
	dirtyAt               string
	z                     int
	hasVisibleBoundsCache bool
	visibleBoundsCache    image.Rectangle

	_ noCopy
}

func (w *widgetState) isInTree() bool {
	return w.root || w.parent != nil
}

func (w *widgetState) isVisible() bool {
	if w.parent != nil {
		if w.hidden {
			return false
		}
		return w.parent.widgetState().isVisible()
	}
	return !w.hidden
}

func (w *widgetState) isEnabled() bool {
	if w.parent != nil {
		if w.disabled {
			return false
		}
		return w.parent.widgetState().isEnabled()
	}
	return !w.disabled
}

func (w *widgetState) opacity() float64 {
	return 1 - w.transparency
}

func (w *widgetState) ensureOffscreen(bounds image.Rectangle) *ebiten.Image {
	if w.offscreen != nil {
		if !bounds.In(w.offscreen.Bounds()) {
			w.offscreen.Deallocate()
			w.offscreen = nil
		}
	}
	if w.offscreen == nil {
		w.offscreen = ebiten.NewImageWithOptions(bounds, nil)
	}
	return w.offscreen.SubImage(bounds).(*ebiten.Image)
}

var skipTraverse = errors.New("skip traverse")

func traverseWidget(widget Widget, f func(widget Widget) error) error {
	if err := f(widget); err != nil {
		return err
	}
	for _, child := range widget.widgetState().children {
		if err := traverseWidget(child, f); err != nil {
			return err
		}
	}
	return nil
}

func RequestRedraw(widget Widget) {
	requestRedraw(widget.widgetState())
}

func requestRedraw(widgetState *widgetState) {
	widgetState.dirty = true
	if theDebugMode.showRenderingRegions {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			widgetState.dirtyAt = fmt.Sprintf("%s:%d", file, line)
		}
	}
}

// noCopy is a struct to warn that the struct should not be copied.
//
// For details, see https://go.dev/issues/8005#issuecomment-190753527
type noCopy struct {
}

func (n *noCopy) Lock() {
}

func (n *noCopy) Unlock() {
}
