// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package main

import (
	"fmt"
	"image"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
)

type modelKey int

const (
	modelKeyModel modelKey = iota
)

type Root struct {
	guigui.DefaultWidget

	background   basicwidget.Background
	sidebar      Sidebar
	settings     Settings
	basic        Basic
	buttons      Buttons
	texts        Texts
	textInputs   TextInputs
	numberInputs NumberInputs
	lists        Lists
	tables       Tables
	popups       Popups

	model Model

	locales           []language.Tag
	faceSourceEntries []basicwidget.FaceSourceEntry

	layout layout.GridLayout
}

func (r *Root) updateFontFaceSources(context *guigui.Context) {
	r.locales = slices.Delete(r.locales, 0, len(r.locales))
	r.locales = context.AppendLocales(r.locales)

	r.faceSourceEntries = slices.Delete(r.faceSourceEntries, 0, len(r.faceSourceEntries))
	r.faceSourceEntries = cjkfont.AppendRecommendedFaceSourceEntries(r.faceSourceEntries, r.locales)
	basicwidget.SetFaceSources(r.faceSourceEntries)
}

func (r *Root) Model(key any) any {
	switch key {
	case modelKeyModel:
		return &r.model
	default:
		return nil
	}
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.background)
	appender.AppendChildWidget(&r.sidebar)
	switch r.model.Mode() {
	case "settings":
		appender.AppendChildWidget(&r.settings)
	case "basic":
		appender.AppendChildWidget(&r.basic)
	case "buttons":
		appender.AppendChildWidget(&r.buttons)
	case "texts":
		appender.AppendChildWidget(&r.texts)
	case "textinputs":
		appender.AppendChildWidget(&r.textInputs)
	case "numberinputs":
		appender.AppendChildWidget(&r.numberInputs)
	case "lists":
		appender.AppendChildWidget(&r.lists)
	case "tables":
		appender.AppendChildWidget(&r.tables)
	case "popups":
		appender.AppendChildWidget(&r.popups)
	}
}

func (r *Root) Update(context *guigui.Context) error {
	r.updateFontFaceSources(context)
	r.layout = layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	return nil
}

func (r *Root) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.background:
		return context.Bounds(r)
	case &r.sidebar:
		return r.layout.CellBounds(0, 0)
	case &r.settings:
		return r.layout.CellBounds(1, 0)
	case &r.basic:
		return r.layout.CellBounds(1, 0)
	case &r.buttons:
		return r.layout.CellBounds(1, 0)
	case &r.texts:
		return r.layout.CellBounds(1, 0)
	case &r.textInputs:
		return r.layout.CellBounds(1, 0)
	case &r.numberInputs:
		return r.layout.CellBounds(1, 0)
	case &r.lists:
		return r.layout.CellBounds(1, 0)
	case &r.tables:
		return r.layout.CellBounds(1, 0)
	case &r.popups:
		return r.layout.CellBounds(1, 0)
	}
	return image.Rectangle{}
}

func main() {
	op := &guigui.RunOptions{
		Title:      "Component Gallery",
		WindowSize: image.Pt(800, 600),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
