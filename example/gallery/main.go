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

func (r *Root) Build(context *guigui.Context) error {
	r.updateFontFaceSources(context)

	context.SetBounds(&r.background, context.Bounds(r), r)

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	context.SetBounds(&r.sidebar, gl.CellBounds(0, 0), r)
	bounds := gl.CellBounds(1, 0)
	switch r.model.Mode() {
	case "settings":
		context.SetBounds(&r.settings, bounds, r)
	case "basic":
		context.SetBounds(&r.basic, bounds, r)
	case "buttons":
		context.SetBounds(&r.buttons, bounds, r)
	case "texts":
		context.SetBounds(&r.texts, bounds, r)
	case "textinputs":
		context.SetBounds(&r.textInputs, bounds, r)
	case "numberinputs":
		context.SetBounds(&r.numberInputs, bounds, r)
	case "lists":
		context.SetBounds(&r.lists, bounds, r)
	case "tables":
		context.SetBounds(&r.tables, bounds, r)
	case "popups":
		context.SetBounds(&r.popups, bounds, r)
	}

	return nil
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
