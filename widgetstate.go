// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package guigui

import (
	"errors"
	"image"
	"maps"

	"github.com/hajimehoshi/ebiten/v2"
)

type bounds3D struct {
	bounds image.Rectangle
	z      int
}

type widgetsAndVisibleBounds struct {
	bounds3Ds       map[Widget]bounds3D
	currentBounds3D map[Widget]bounds3D
}

func (w *widgetsAndVisibleBounds) reset() {
	clear(w.bounds3Ds)
}

func (w *widgetsAndVisibleBounds) append(context *Context, widget Widget) {
	if w.bounds3Ds == nil {
		w.bounds3Ds = map[Widget]bounds3D{}
	}
	bounds := context.VisibleBounds(widget)
	if bounds.Empty() {
		return
	}
	w.bounds3Ds[widget] = bounds3D{
		bounds: bounds,
		z:      z(widget),
	}
}

func (w *widgetsAndVisibleBounds) equals(context *Context, currentWidgets []Widget) bool {
	if w.currentBounds3D == nil {
		w.currentBounds3D = map[Widget]bounds3D{}
	} else {
		clear(w.currentBounds3D)
	}
	for _, widget := range currentWidgets {
		if context.VisibleBounds(widget).Empty() {
			continue
		}
		w.currentBounds3D[widget] = bounds3D{
			bounds: context.VisibleBounds(widget),
			z:      z(widget),
		}
	}
	return maps.Equal(w.bounds3Ds, w.currentBounds3D)
}

func (w *widgetsAndVisibleBounds) redrawIfDifferentParentZ(app *app) {
	for widget, bounds3D := range w.bounds3Ds {
		if widget.ZDelta() != 0 {
			app.requestRedraw(bounds3D.bounds)
			RequestRedraw(widget)
		}
	}
}

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

	offscreen *ebiten.Image

	dirty     bool
	hasZCache bool
	zCache    int

	_ noCopy
}

func (w *widgetState) isInTree() bool {
	p := w
	for ; p.parent != nil; p = p.parent.widgetState() {
	}
	return p.root
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
	widget.widgetState().dirty = true
}

func z(widget Widget) int {
	s := widget.widgetState()
	if s.hasZCache {
		return widget.widgetState().zCache
	}
	var r int
	if parent := s.parent; parent != nil {
		r = z(parent)
	}
	r += widget.ZDelta()
	s.zCache = r
	s.hasZCache = true
	return r
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
