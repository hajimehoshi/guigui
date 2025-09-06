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

	background        basicwidget.Background
	createButton      basicwidget.Button
	textInput         basicwidget.TextInput
	tasksPanel        basicwidget.Panel
	tasksPanelContent tasksPanelContent

	model Model

	locales           []language.Tag
	faceSourceEntries []basicwidget.FaceSourceEntry

	mainLayout layout.GridLayout
	topLayout  layout.GridLayout
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
	appender.AppendChildWidget(&r.textInput)
	appender.AppendChildWidget(&r.createButton)
	appender.AppendChildWidget(&r.tasksPanel)
}

func (r *Root) Build(context *guigui.Context) error {
	r.updateFontFaceSources(context)

	r.textInput.SetOnKeyJustPressed(func(key ebiten.Key) bool {
		if key == ebiten.KeyEnter {
			r.tryCreateTask(r.textInput.Value())
			return true
		}
		return false
	})

	r.createButton.SetText("Create")
	r.createButton.SetOnUp(func() {
		r.tryCreateTask(r.textInput.Value())
	})
	context.SetEnabled(&r.createButton, r.model.CanAddTask(r.textInput.Value()))

	r.tasksPanelContent.SetOnDeleted(func(id int) {
		r.model.DeleteTaskByID(id)
	})
	r.tasksPanel.SetContent(&r.tasksPanelContent)
	r.tasksPanel.SetAutoBorder(true)

	u := basicwidget.UnitSize(context)
	r.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(u),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	r.topLayout = layout.GridLayout{
		Bounds: r.mainLayout.CellBounds(0, 0),
		Widths: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(5 * u),
		},
		ColumnGap: u / 2,
	}

	r.tasksPanelContent.SetWidth(r.mainLayout.CellBounds(0, 1).Dx())

	return nil
}

func (r *Root) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.background:
		return context.Bounds(r)
	case &r.textInput:
		return r.topLayout.CellBounds(0, 0)
	case &r.createButton:
		return r.topLayout.CellBounds(1, 0)
	case &r.tasksPanel:
		return r.mainLayout.CellBounds(0, 1)
	}
	return image.Rectangle{}
}

func (r *Root) tryCreateTask(text string) {
	if r.model.TryAddTask(text) {
		r.textInput.ForceSetValue("")
	}
}

type taskWidget struct {
	guigui.DefaultWidget

	doneButton basicwidget.Button
	text       basicwidget.Text

	layout layout.GridLayout
}

const (
	taskWidgetEventDoneButtonPressed = "doneButtonPressed"
)

func (t *taskWidget) SetOnDoneButtonPressed(f func()) {
	guigui.RegisterEventHandler(t, taskWidgetEventDoneButtonPressed, f)
}

func (t *taskWidget) SetText(text string) {
	t.text.SetValue(text)
}

func (t *taskWidget) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&t.doneButton)
	appender.AppendChildWidget(&t.text)
}

func (t *taskWidget) Build(context *guigui.Context) error {
	t.doneButton.SetText("Done")
	t.doneButton.SetOnUp(func() {
		guigui.DispatchEventHandler(t, taskWidgetEventDoneButtonPressed)
	})

	t.text.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	u := basicwidget.UnitSize(context)
	t.layout = layout.GridLayout{
		Bounds: context.Bounds(t),
		Widths: []layout.Size{
			layout.FixedSize(3 * u),
			layout.FlexibleSize(1),
		},
		ColumnGap: u / 2,
	}

	return nil
}

func (t *taskWidget) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &t.doneButton:
		return t.layout.CellBounds(0, 0)
	case &t.text:
		return t.layout.CellBounds(1, 0)
	}
	return image.Rectangle{}
}

func (t *taskWidget) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	return image.Pt(6*int(basicwidget.UnitSize(context)), t.doneButton.Measure(context, guigui.Constraints{}).Y)
}

type tasksPanelContent struct {
	guigui.DefaultWidget

	taskWidgets []taskWidget

	width int

	layout layout.GridLayout
}

const (
	tasksPanelContentEventDeleted = "deleted"
)

func (t *tasksPanelContent) SetWidth(width int) {
	t.width = width
}

func (t *tasksPanelContent) SetOnDeleted(f func(id int)) {
	guigui.RegisterEventHandler(t, tasksPanelContentEventDeleted, f)
}

func (t *tasksPanelContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	model := context.Model(t, modelKeyModel).(*Model)
	if model.TaskCount() > len(t.taskWidgets) {
		t.taskWidgets = slices.Grow(t.taskWidgets, model.TaskCount()-len(t.taskWidgets))[:model.TaskCount()]
	} else {
		t.taskWidgets = slices.Delete(t.taskWidgets, model.TaskCount(), len(t.taskWidgets))
	}
	for i := range t.taskWidgets {
		appender.AppendChildWidget(&t.taskWidgets[i])
	}
}

func (t *tasksPanelContent) Build(context *guigui.Context) error {
	model := context.Model(t, modelKeyModel).(*Model)
	for i := range model.TaskCount() {
		task := model.TaskByIndex(i)
		t.taskWidgets[i].SetOnDoneButtonPressed(func() {
			guigui.DispatchEventHandler(t, tasksPanelContentEventDeleted, task.ID)
		})
		t.taskWidgets[i].SetText(task.Text)
	}

	u := basicwidget.UnitSize(context)

	t.layout = layout.GridLayout{
		Bounds: context.Bounds(t),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row >= len(t.taskWidgets) {
					return layout.FixedSize(0)
				}
				w := guigui.FixedWidthConstraints(context.Bounds(t).Dx())
				h := t.taskWidgets[row].Measure(context, w).Y
				return layout.FixedSize(h)
			}),
		},
		RowGap: u / 4,
	}

	return nil
}

func (t *tasksPanelContent) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	// TODO: This is not efficient as searching an index takes O(n) time.
	for i := range t.taskWidgets {
		if widget == &t.taskWidgets[i] {
			return t.layout.CellBounds(0, i)
		}
	}
	return image.Rectangle{}
}

func (t *tasksPanelContent) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	var h int
	for i := range t.taskWidgets {
		h += t.taskWidgets[i].Measure(context, constraints).Y
		h += int(u / 4)
	}
	return image.Pt(t.width, h)
}

func main() {
	op := &guigui.RunOptions{
		Title:         "TODO",
		WindowMinSize: image.Pt(320, 240),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
