// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package guigui

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/language"

	"github.com/hajimehoshi/guigui/internal/colormode"
	"github.com/hajimehoshi/guigui/internal/locale"
)

type ColorMode int

var envLocales []language.Tag

func init() {
	if locales := os.Getenv("GUIGUI_LOCALES"); locales != "" {
		for _, tag := range strings.Split(os.Getenv("GUIGUI_LOCALES"), ",") {
			l, err := language.Parse(strings.TrimSpace(tag))
			if err != nil {
				slog.Warn(fmt.Sprintf("invalid GUIGUI_LOCALES: %s", tag))
				continue
			}
			envLocales = append(envLocales, l)
		}
	}
}

var systemLocales []language.Tag

func init() {
	ls, err := locale.Locales()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	systemLocales = ls
}

const (
	ColorModeLight ColorMode = iota
	ColorModeDark
)

type Context struct {
	app     *app
	inBuild bool

	appScaleMinus1             float64
	colorMode                  ColorMode
	colorModeSet               bool
	cachedDefaultColorMode     colormode.ColorMode
	cachedDefaultColorModeTime time.Time
	defaultColorWarnOnce       sync.Once
	locales                    []language.Tag
	allLocales                 []language.Tag

	tmpWidgetStates []*widgetState
}

func (c *Context) Scale() float64 {
	return c.DeviceScale() * c.AppScale()
}

func (c *Context) DeviceScale() float64 {
	return c.app.deviceScale
}

func (c *Context) AppScale() float64 {
	return c.appScaleMinus1 + 1
}

func (c *Context) SetAppScale(scale float64) {
	if c.appScaleMinus1 == scale-1 {
		return
	}
	c.appScaleMinus1 = scale - 1
	c.app.requestRedraw(c.app.bounds())
}

func (c *Context) ColorMode() ColorMode {
	if c.colorModeSet {
		return c.colorMode
	}
	return c.autoColorMode()
}

func (c *Context) SetColorMode(mode ColorMode) {
	if c.colorModeSet && mode == c.colorMode {
		return
	}

	c.colorMode = mode
	c.colorModeSet = true
	c.app.requestRedraw(c.app.bounds())
}

func (c *Context) UseAutoColorMode() {
	if !c.colorModeSet {
		return
	}
	c.colorModeSet = false
	c.app.requestRedraw(c.app.bounds())
}

func (c *Context) IsAutoColorModeUsed() bool {
	return !c.colorModeSet
}

func (c *Context) autoColorMode() ColorMode {
	// TODO: Consider the system color mode.
	switch mode := os.Getenv("GUIGUI_COLOR_MODE"); mode {
	case "light":
		return ColorModeLight
	case "dark":
		return ColorModeDark
	case "":
		if time.Since(c.cachedDefaultColorModeTime) >= time.Second {
			m := colormode.SystemColorMode()
			if c.cachedDefaultColorMode != m {
				c.app.requestRedraw(c.app.bounds())
			}
			c.cachedDefaultColorMode = m
			c.cachedDefaultColorModeTime = time.Now()
		}
		switch c.cachedDefaultColorMode {
		case colormode.Light:
			return ColorModeLight
		case colormode.Dark:
			return ColorModeDark
		}
	default:
		c.defaultColorWarnOnce.Do(func() {
			slog.Warn(fmt.Sprintf("invalid GUIGUI_COLOR_MODE: %s", mode))
		})
	}

	return ColorModeLight
}

func (c *Context) AppendLocales(locales []language.Tag) []language.Tag {
	if len(c.allLocales) == 0 {
		// App locales
		for _, l := range c.locales {
			if slices.Contains(c.allLocales, l) {
				continue
			}
			c.allLocales = append(c.allLocales, l)
		}
		// Env locales
		for _, l := range envLocales {
			if slices.Contains(c.allLocales, l) {
				continue
			}
			c.allLocales = append(c.allLocales, l)
		}
		// System locales
		for _, l := range systemLocales {
			if slices.Contains(c.allLocales, l) {
				continue
			}
			c.allLocales = append(c.allLocales, l)
		}
	}
	return append(locales, c.allLocales...)
}

func (c *Context) AppendAppLocales(locales []language.Tag) []language.Tag {
	origLen := len(locales)
	for _, l := range c.locales {
		if slices.Contains(locales[origLen:], l) {
			continue
		}
		locales = append(locales, l)
	}
	return locales
}

func (c *Context) SetAppLocales(locales []language.Tag) {
	if slices.Equal(c.locales, locales) {
		return
	}

	c.locales = slices.Delete(c.locales, 0, len(c.locales))
	c.locales = append(c.locales, locales...)
	c.allLocales = slices.Delete(c.allLocales, 0, len(c.allLocales))

	c.app.requestRedraw(c.app.bounds())
}

func (c *Context) AppSize() image.Point {
	return c.app.bounds().Size()
}

func (c *Context) AppBounds() image.Rectangle {
	return c.app.bounds()
}

func (c *Context) Position(widget Widget) image.Point {
	return widget.widgetState().position
}

// Deprecated: use [Widget.Layout] instead.
func (c *Context) SetPosition(widget Widget, position image.Point) {
	if widget.widgetState().position == position {
		return
	}
	c.clearVisibleBoundsCacheForWidget(widget)
	widget.widgetState().position = position
	// Rerendering happens at (*.app).requestRedrawIfTreeChanged if necessary.
}

// Deprecated: use [Widget.Layout] instead.
const AutoSize = -1

// Deprecated: use [Widget.Layout] or [WidgetWithSize] instead.
func (c *Context) SetSize(widget Widget, size image.Point, specifierWidget Widget) {
	w := widget.widgetState()
	if size == image.Pt(AutoSize, AutoSize) {
		delete(w.sizes, specifierWidget.widgetState())
		return
	}
	if w.sizes[specifierWidget.widgetState()] == size {
		return
	}
	c.clearVisibleBoundsCacheForWidget(widget)
	if w.sizes == nil {
		w.sizes = map[*widgetState]image.Point{}
	}
	w.sizes[specifierWidget.widgetState()] = size
}

func (c *Context) ActualSize(widget Widget) image.Point {
	widgetState := widget.widgetState()
	w := widgetState
	for {
		c.tmpWidgetStates = append(c.tmpWidgetStates, w)
		if w.parent == nil {
			break
		}
		w = w.parent.widgetState()
	}
	s := image.Pt(AutoSize, AutoSize)
	for _, w := range c.tmpWidgetStates {
		size, ok := widgetState.sizes[w]
		if !ok {
			continue
		}
		if s.X == AutoSize {
			s.X = size.X
		}
		if s.Y == AutoSize {
			s.Y = size.Y
		}
		if s.X != AutoSize && s.Y != AutoSize {
			break
		}
	}
	c.tmpWidgetStates = slices.Delete(c.tmpWidgetStates, 0, len(c.tmpWidgetStates))
	if s.X == AutoSize || s.Y == AutoSize {
		size := widget.Measure(c, Constraints{})
		if s.X == AutoSize {
			s.X = size.X
		}
		if s.Y == AutoSize {
			s.Y = size.Y
		}
	}
	return s
}

func (c *Context) Bounds(widget Widget) image.Rectangle {
	widgetState := widget.widgetState()
	return image.Rectangle{
		Min: widgetState.position,
		Max: widgetState.position.Add(c.ActualSize(widget)),
	}
}

// Deprecated: use [Widget.Layout] instead.
func (c *Context) SetBounds(widget Widget, bounds image.Rectangle, specifierWidget Widget) {
	c.SetPosition(widget, bounds.Min)
	c.SetSize(widget, bounds.Size(), specifierWidget)
}

func (c *Context) VisibleBounds(widget Widget) image.Rectangle {
	state := widget.widgetState()
	if state.hasVisibleBoundsCache {
		return state.visibleBoundsCache
	}

	parent := widget.widgetState().parent
	if parent == nil {
		b := c.app.bounds()
		state.hasVisibleBoundsCache = true
		state.visibleBoundsCache = b
		return b
	}
	if widget.ZDelta() != 0 {
		b := c.Bounds(widget)
		state.hasVisibleBoundsCache = true
		state.visibleBoundsCache = b
		return b
	}

	var b image.Rectangle
	parentVB := c.VisibleBounds(parent)
	if !parentVB.Empty() {
		b = parentVB.Intersect(c.Bounds(widget))
	}
	state.hasVisibleBoundsCache = true
	state.visibleBoundsCache = b
	return b
}

func (c *Context) SetVisible(widget Widget, visible bool) {
	widgetState := widget.widgetState()
	if widgetState.hidden == !visible {
		return
	}
	widgetState.hidden = !visible
	if !visible {
		c.blur(widget)
	}
	RequestRedraw(widget)
}

func (c *Context) IsVisible(widget Widget) bool {
	return widget.widgetState().isVisible()
}

func (c *Context) SetEnabled(widget Widget, enabled bool) {
	widgetState := widget.widgetState()
	if widgetState.disabled == !enabled {
		return
	}
	widgetState.disabled = !enabled
	if !enabled {
		c.blur(widget)
	}
	RequestRedraw(widget)
}

func (c *Context) IsEnabled(widget Widget) bool {
	return widget.widgetState().isEnabled()
}

func (c *Context) SetFocused(widget Widget, focused bool) {
	if focused {
		c.focus(widget)
	} else {
		c.blur(widget)
	}
}

func (c *Context) focus(widget Widget) {
	ws := widget.widgetState()
	if !ws.isVisible() {
		return
	}
	if !ws.isEnabled() {
		return
	}

	if !ws.isInTree(c.app.buildCount) {
		return
	}
	if c.app.focusedWidgetState == widget.widgetState() {
		return
	}

	c.app.focusWidget(widget.widgetState())

	// Rerender everything when a focus changes.
	// A widget including a focused widget might be affected.
	c.app.requestRedraw(c.app.bounds())
}

func (c *Context) blur(widget Widget) {
	widgetState := widget.widgetState()
	if !widgetState.isInTree(c.app.buildCount) {
		return
	}
	var unfocused bool
	_ = traverseWidget(widget, func(w Widget) error {
		if c.app.focusedWidgetState == w.widgetState() {
			c.app.focusWidget(c.app.root.widgetState())
			unfocused = true
			return skipTraverse
		}
		return nil
	})
	if unfocused {
		// Rerender everything when a focus changes.
		// A widget including a focused widget might be affected.
		c.app.requestRedraw(c.app.bounds())
	}
}

func (c *Context) IsFocused(widget Widget) bool {
	widgetState := widget.widgetState()
	return widgetState.isInTree(c.app.buildCount) && widgetState.isVisible() && c.app.focusedWidgetState == widget.widgetState()
}

func (c *Context) IsFocusedOrHasFocusedChild(widget Widget) bool {
	if c.inBuild {
		panic("guigui: IsFocusedOrHasFocusedChild cannot be called in Build")
	}

	if len(widget.widgetState().children) == 0 {
		return c.app.focusedWidgetState == widget.widgetState()
	}

	w := c.app.focusedWidgetState
	if w == nil {
		return false
	}
	for {
		widgetState := widget.widgetState()
		if w == widgetState {
			return widgetState.isInTree(c.app.buildCount) && widgetState.isVisible()
		}
		if w.parent == nil {
			break
		}
		w = w.parent.widgetState()
	}
	return false
}

func (c *Context) Opacity(widget Widget) float64 {
	return widget.widgetState().opacity()
}

func (c *Context) SetOpacity(widget Widget, opacity float64) {
	opacity = min(max(opacity, 0), 1)
	widgetState := widget.widgetState()
	if widgetState.transparency == 1-opacity {
		return
	}
	widgetState.transparency = 1 - opacity
	RequestRedraw(widget)
}

func (c *Context) IsWidgetHitAtCursor(widget Widget) bool {
	return c.app.isWidgetHitAt(widget)
}

func (c *Context) SetCustomDraw(widget Widget, customDraw CustomDrawFunc) {
	widget.widgetState().customDraw = customDraw
}

func (c *Context) clearVisibleBoundsCacheForWidget(widget Widget) {
	widget.widgetState().hasVisibleBoundsCache = false
	widget.widgetState().visibleBoundsCache = image.Rectangle{}
	for _, child := range widget.widgetState().children {
		c.clearVisibleBoundsCacheForWidget(child)
	}
}

func (c *Context) Model(widget Widget, key any) any {
	for w := widget; w != nil; w = w.widgetState().parent {
		if v := w.Model(key); v != nil {
			return v
		}
	}
	return nil
}
