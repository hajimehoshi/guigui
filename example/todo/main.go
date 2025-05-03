// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package main

import (
	"fmt"
	"image"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.RootWidget

	background        basicwidget.Background
	createButton      basicwidget.TextButton
	textInput         basicwidget.TextInput
	tasksPanel        basicwidget.ScrollablePanel
	tasksPanelContent tasksPanelContent

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

	r.textInput.SetOnEnterPressed(func(text string) {
		r.tryCreateTask(text)
	})

	r.createButton.SetText("Create")
	r.createButton.SetOnUp(func() {
		r.tryCreateTask(r.textInput.Value())
	})
	context.SetEnabled(&r.createButton, r.model.CanAddTask(r.textInput.Value()))

	r.tasksPanelContent.SetModel(&r.model)
	r.tasksPanelContent.SetOnDeleted(func(id int) {
		r.model.DeleteTaskByID(id)
	})
	r.tasksPanel.SetContent(&r.tasksPanelContent)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(u),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	{
		gl := layout.GridLayout{
			Bounds: gl.CellBounds(0, 0),
			Widths: []layout.Size{
				layout.FlexibleSize(1),
				layout.FixedSize(5 * u),
			},
			ColumnGap: u / 2,
		}
		appender.AppendChildWidgetWithBounds(&r.textInput, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&r.createButton, gl.CellBounds(1, 0))
	}
	{
		bounds := gl.CellBounds(0, 1)
		context.SetSize(&r.tasksPanelContent, image.Pt(bounds.Dx(), guigui.DefaultSize))
		appender.AppendChildWidgetWithBounds(&r.tasksPanel, bounds)
	}

	return nil
}

func (r *Root) tryCreateTask(text string) {
	if r.model.TryAddTask(text) {
		r.textInput.ForceSetValue("")
	}
}

type taskWidget struct {
	guigui.DefaultWidget

	doneButton basicwidget.TextButton
	text       basicwidget.Text

	onDoneButtonPressed func()
}

func (t *taskWidget) SetOnDoneButtonPressed(f func()) {
	t.onDoneButtonPressed = f
}

func (t *taskWidget) SetText(text string) {
	t.text.SetValue(text)
}

func (t *taskWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	t.doneButton.SetText("Done")
	t.doneButton.SetOnUp(func() {
		if t.onDoneButtonPressed != nil {
			t.onDoneButtonPressed()
		}
	})

	t.text.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(t),
		Widths: []layout.Size{
			layout.FixedSize(3 * u),
			layout.FlexibleSize(1),
		},
		ColumnGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&t.doneButton, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&t.text, gl.CellBounds(1, 0))

	return nil
}

func (t *taskWidget) DefaultSize(context *guigui.Context) image.Point {
	return image.Pt(6*int(basicwidget.UnitSize(context)), context.Size(&t.doneButton).Y)
}

type tasksPanelContent struct {
	guigui.DefaultWidget

	taskWidgets []taskWidget

	onDeleted func(id int)

	model *Model
}

func (t *tasksPanelContent) SetOnDeleted(f func(id int)) {
	t.onDeleted = f
}

func (t *tasksPanelContent) SetModel(model *Model) {
	t.model = model
}

func (t *tasksPanelContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if t.model.TaskCount() > len(t.taskWidgets) {
		t.taskWidgets = slices.Grow(t.taskWidgets, t.model.TaskCount()-len(t.taskWidgets))[:t.model.TaskCount()]
	} else {
		t.taskWidgets = slices.Delete(t.taskWidgets, t.model.TaskCount(), len(t.taskWidgets))
	}
	for i := range t.model.TaskCount() {
		task := t.model.TaskByIndex(i)
		t.taskWidgets[i].SetOnDoneButtonPressed(func() {
			if t.onDeleted != nil {
				t.onDeleted(task.ID)
			}
		})
		t.taskWidgets[i].SetText(task.Text)
	}

	u := basicwidget.UnitSize(context)

	gl := layout.GridLayout{
		Bounds: context.Bounds(t),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row >= len(t.taskWidgets) {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(context.Size(&t.taskWidgets[row]).Y)
			}),
		},
		RowGap: u / 4,
	}
	for i := range t.taskWidgets {
		bounds := gl.CellBounds(0, i)
		appender.AppendChildWidgetWithBounds(&t.taskWidgets[i], bounds)
	}

	return nil
}

func (t *tasksPanelContent) DefaultSize(context *guigui.Context) image.Point {
	u := basicwidget.UnitSize(context)
	var h int
	for i := range t.taskWidgets {
		h += context.Size(&t.taskWidgets[i]).Y
		h += int(u / 4)
	}
	return image.Pt(6*int(u), h)
}

func main() {
	op := &guigui.RunOptions{
		Title:         "TODO",
		WindowMinSize: image.Pt(320, 240),
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
