// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package main

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
)

func init() {
}

type Root struct {
	guigui.RootWidget

	background   basicwidget.Background
	sidebar      Sidebar
	settings     Settings
	basic        Basic
	buttons      Buttons
	texts        Texts
	textInputs   TextInputs
	numberInputs NumberInputs
	lists        Lists
	popups       Popups

	model Model
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	faceSources := []*text.GoTextFaceSource{
		basicwidget.DefaultFaceSource(),
	}
	for _, locale := range context.AppendLocales(nil) {
		fs := cjkfont.FaceSourceFromLocale(locale)
		if fs != nil {
			faceSources = append(faceSources, fs)
			break
		}
	}
	if len(faceSources) == 1 {
		// Set a Japanese font as a fallback. You can use any font you like here.
		faceSources = append(faceSources, cjkfont.FaceSourceJP())
	}
	basicwidget.SetFaceSources(faceSources)

	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	r.buttons.SetModel(&r.model)
	r.texts.SetModel(&r.model)
	r.textInputs.SetModel(&r.model)
	r.numberInputs.SetModel(&r.model)
	r.sidebar.SetModel(&r.model)

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
	case "popups":
		appender.AppendChildWidgetWithBounds(&r.popups, bounds)
	}

	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title:      "Component Gallery",
		WindowSize: image.Pt(800, 600),
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
