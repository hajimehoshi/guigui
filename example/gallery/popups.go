// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Popups struct {
	guigui.DefaultWidget

	forms                        [2]basicwidget.Form
	blurBackgroundText           basicwidget.Text
	blurBackgroundToggle         basicwidget.Toggle
	closeByClickingOutsideText   basicwidget.Text
	closeByClickingOutsideToggle basicwidget.Toggle
	showButton                   basicwidget.Button

	contextMenuPopupText          basicwidget.Text
	contextMenuPopupClickHereText basicwidget.Text

	simplePopup        basicwidget.Popup
	simplePopupContent simplePopupContent

	contextMenuPopup basicwidget.PopupMenu[int]
}

func (p *Popups) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	for i := range p.forms {
		appender.AppendChildWidget(&p.forms[i])
	}
	appender.AppendChildWidget(&p.simplePopup)
	appender.AppendChildWidget(&p.contextMenuPopup)
	p.simplePopup.SetContent(&p.simplePopupContent)
}

func (p *Popups) Build(context *guigui.Context) error {
	p.blurBackgroundText.SetValue("Blur background")
	p.closeByClickingOutsideText.SetValue("Close by clicking outside")
	p.showButton.SetText("Show")
	p.showButton.SetOnUp(func() {
		p.simplePopup.Open(context)
	})

	p.forms[0].SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &p.blurBackgroundText,
			SecondaryWidget: &p.blurBackgroundToggle,
		},
		{
			PrimaryWidget:   &p.closeByClickingOutsideText,
			SecondaryWidget: &p.closeByClickingOutsideToggle,
		},
		{
			SecondaryWidget: &p.showButton,
		},
	})

	p.contextMenuPopupText.SetValue("Context menu")
	p.contextMenuPopupClickHereText.SetValue("Click here by the right button")

	p.forms[1].SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &p.contextMenuPopupText,
			SecondaryWidget: &p.contextMenuPopupClickHereText,
		},
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row >= len(p.forms) {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(p.forms[row].DefaultSizeInContainer(context, context.Bounds(p).Dx()-u).Y)
			}),
		},
		RowGap: u / 2,
	}
	for i := range p.forms {
		context.SetBounds(&p.forms[i], gl.CellBounds(0, i), p)
	}

	p.simplePopupContent.popup = &p.simplePopup
	p.simplePopup.SetBackgroundBlurred(p.blurBackgroundToggle.Value())
	p.simplePopup.SetCloseByClickingOutside(p.closeByClickingOutsideToggle.Value())
	p.simplePopup.SetAnimationDuringFade(true)

	appBounds := context.AppBounds()
	contentSize := image.Pt(int(12*u), int(6*u))
	simplePopupPosition := image.Point{
		X: appBounds.Min.X + (appBounds.Dx()-contentSize.X)/2,
		Y: appBounds.Min.Y + (appBounds.Dy()-contentSize.Y)/2,
	}
	simplePopupBounds := image.Rectangle{
		Min: simplePopupPosition,
		Max: simplePopupPosition.Add(contentSize),
	}
	context.SetSize(&p.simplePopupContent, simplePopupBounds.Size(), p)
	context.SetBounds(&p.simplePopup, simplePopupBounds, p)

	p.contextMenuPopup.SetItemsByStrings([]string{"Item 1", "Item 2", "Item 3"})
	// A context menu's position is updated at HandlePointingInput.

	return nil
}

func (p *Popups) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		// Use IsWidgetOrBackgroundHitAtCursor. context.IsWidgetHitAtCursor doesn't work when a popup's transparent background exists.
		if p.contextMenuPopup.IsWidgetOrBackgroundHitAtCursor(context, &p.contextMenuPopupClickHereText) {
			pt := image.Pt(ebiten.CursorPosition())
			context.SetPosition(&p.contextMenuPopup, pt)
			p.contextMenuPopup.Open(context)
		}
	}
	return guigui.HandleInputResult{}
}

type simplePopupContent struct {
	guigui.DefaultWidget

	popup *basicwidget.Popup

	titleText   basicwidget.Text
	closeButton basicwidget.Button
}

func (s *simplePopupContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.titleText)
	appender.AppendChildWidget(&s.closeButton)
}

func (s *simplePopupContent) Build(context *guigui.Context) error {
	u := basicwidget.UnitSize(context)

	s.titleText.SetValue("Hello!")
	s.titleText.SetBold(true)

	s.closeButton.SetText("Close")
	s.closeButton.SetOnUp(func() {
		s.popup.Close()
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.LazySize(func(row int) layout.Size {
				if row != 1 {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(s.closeButton.DefaultSize(context).Y)
			}),
		},
	}
	context.SetBounds(&s.titleText, gl.CellBounds(0, 0), s)
	{
		gl := layout.GridLayout{
			Bounds: gl.CellBounds(0, 1),
			Widths: []layout.Size{
				layout.FlexibleSize(1),
				layout.FixedSize(s.closeButton.DefaultSize(context).X),
			},
		}
		context.SetBounds(&s.closeButton, gl.CellBounds(1, 0), s)
	}

	return nil
}
