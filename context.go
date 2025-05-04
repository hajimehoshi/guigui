// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package guigui

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/hajimehoshi/guigui/internal/colormode"
	"github.com/hajimehoshi/guigui/internal/locale"

	"golang.org/x/text/language"

	"github.com/hajimehoshi/guigui/internal/locale"
)

type ColorMode int

var defaultColorMode ColorMode

func init() {
	// TODO: Consider the system color mode.
	switch mode := os.Getenv("GUIGUI_COLOR_MODE"); mode {
	case "light":
		defaultColorMode = ColorModeLight
	case "dark":
		defaultColorMode = ColorModeDark
	case "":
		switch colormode.SystemColorMode() {
		case colormode.Light:
			defaultColorMode = ColorModeLight
		case colormode.Dark:
			defaultColorMode = ColorModeDark
		default:
			defaultColorMode = ColorModeLight
		}
	default:
		slog.Warn(fmt.Sprintf("invalid GUIGUI_COLOR_MODE: %s", mode))
	}
}

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
	app *app

	appScaleMinus1 float64
	colorMode      ColorMode
	hasColorMode   bool
	locales        []language.Tag
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
	if c.hasColorMode {
		return c.colorMode
	}
	return defaultColorMode
}

func (c *Context) SetColorMode(mode ColorMode) {
	if c.hasColorMode && c.colorMode == mode {
		return
	}
	c.colorMode = mode
	c.hasColorMode = true
	c.app.requestRedraw(c.app.bounds())
}

func (c *Context) ResetColorMode() {
	c.hasColorMode = false
}

func (c *Context) AppendLocales(locales []language.Tag) []language.Tag {
	origLen := len(locales)
	// App locales
	for _, l := range c.locales {
		if slices.Contains(locales[origLen:], l) {
			continue
		}
		locales = append(locales, l)
	}
	// Env locales
	for _, l := range envLocales {
		if slices.Contains(locales[origLen:], l) {
			continue
		}
		locales = append(locales, l)
	}
	// System locales
	for _, l := range systemLocales {
		if slices.Contains(locales[origLen:], l) {
			continue
		}
		locales = append(locales, l)
	}
	return locales
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

	c.locales = append([]language.Tag(nil), locales...)
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

func (c *Context) SetPosition(widget Widget, position image.Point) {
	widget.widgetState().position = position
	// Rerendering happens at (*.app).requestRedrawIfTreeChanged if necessary.
}

const DefaultSize = -1

func (c *Context) SetSize(widget Widget, size image.Point) {
	widget.widgetState().widthPlus1 = size.X + 1
	widget.widgetState().heightPlus1 = size.Y + 1
}

func (c *Context) Size(widget Widget) image.Point {
	widgetState := widget.widgetState()
	ds := widget.DefaultSize(c)
	var s image.Point
	if widgetState.widthPlus1 == 0 {
		s.X = ds.X
	} else {
		s.X = widgetState.widthPlus1 - 1
	}
	if widgetState.heightPlus1 == 0 {
		s.Y = ds.Y
	} else {
		s.Y = widgetState.heightPlus1 - 1
	}
	return s
}

func (c *Context) Bounds(widget Widget) image.Rectangle {
	widgetState := widget.widgetState()
	return image.Rectangle{
		Min: widgetState.position,
		Max: widgetState.position.Add(c.Size(widget)),
	}
}

func (c *Context) VisibleBounds(widget Widget) image.Rectangle {
	parent := widget.widgetState().parent
	if parent == nil {
		return c.app.bounds()
	}
	if widget.ZDelta() != 0 {
		return c.Bounds(widget)
	}
	return c.VisibleBounds(parent).Intersect(c.Bounds(widget))
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
	widgetState := widget.widgetState()
	if !widgetState.isVisible() {
		return
	}
	if !widgetState.isEnabled() {
		return
	}

	if !widgetState.isInTree() {
		return
	}
	if c.app.focusedWidget == widget {
		return
	}

	var oldWidget Widget
	if c.app.focusedWidget != nil {
		oldWidget = c.app.focusedWidget
	}

	c.app.focusedWidget = widget
	RequestRedraw(c.app.focusedWidget)
	if oldWidget != nil {
		RequestRedraw(oldWidget)
	}
}

func (c *Context) blur(widget Widget) {
	widgetState := widget.widgetState()
	if !widgetState.isInTree() {
		return
	}
	_ = traverseWidget(widget, func(w Widget) error {
		if c.app.focusedWidget == w {
			c.app.focusedWidget = nil
			return skipTraverse
		}
		return nil
	})
	RequestRedraw(widget)
}

func (c *Context) IsFocused(widget Widget) bool {
	// Check this first to avoid unnecessary evaluation.
	if c.app.focusedWidget != widget {
		return false
	}
	widgetState := widget.widgetState()
	return widgetState.isInTree() && widgetState.isVisible()
}

func (c *Context) HasFocusedChildWidget(widget Widget) bool {
	if c.IsFocused(widget) {
		return true
	}

	widgetState := widget.widgetState()
	for _, child := range widgetState.children {
		if c.HasFocusedChildWidget(child) {
			return true
		}
	}
	return false
}

func (c *Context) Opacity(widget Widget) float64 {
	return widget.widgetState().opacity()
}

func (c *Context) SetOpacity(widget Widget, opacity float64) {
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 1 {
		opacity = 1
	}
	widgetState := widget.widgetState()
	if widgetState.transparency == 1-opacity {
		return
	}
	widgetState.transparency = 1 - opacity
	RequestRedraw(widget)
}

func (c *Context) IsWidgetHitAt(widget Widget, point image.Point) bool {
	return c.app.isWidgetHitAt(widget, point)
}
