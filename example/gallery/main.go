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

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.updateFontFaceSources(context)

	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&r.sidebar, gl.CellBounds(0, 0))
	bounds := gl.CellBounds(1, 0)
	switch r.model.Mode() {
	case "settings":
		appender.AppendChildWidgetWithBounds(&r.settings, bounds)
	case "basic":
		appender.AppendChildWidgetWithBounds(&r.basic, bounds)
	case "buttons":
		appender.AppendChildWidgetWithBounds(&r.buttons, bounds)
	case "texts":
		appender.AppendChildWidgetWithBounds(&r.texts, bounds)
	case "textinputs":
		appender.AppendChildWidgetWithBounds(&r.textInputs, bounds)
	case "numberinputs":
		appender.AppendChildWidgetWithBounds(&r.numberInputs, bounds)
	case "lists":
		appender.AppendChildWidgetWithBounds(&r.lists, bounds)
	case "tables":
		appender.AppendChildWidgetWithBounds(&r.tables, bounds)
	case "popups":
		appender.AppendChildWidgetWithBounds(&r.popups, bounds)
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
