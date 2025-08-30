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
	simplePopupContent guigui.WidgetWithSize[*simplePopupContent]

	contextMenuPopup basicwidget.PopupMenu[int]

	layout                   layout.GridLayout
	contextMenuPopupPosition image.Point
}

func (p *Popups) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	for i := range p.forms {
		appender.AppendChildWidget(&p.forms[i])
	}
	appender.AppendChildWidget(&p.simplePopup)
	appender.AppendChildWidget(&p.contextMenuPopup)
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
	p.layout = layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row >= len(p.forms) {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(p.forms[row].Measure(context, guigui.FixedWidthConstraints(context.Bounds(p).Dx()-u)).Y)
			}),
		},
		RowGap: u / 2,
	}

	p.simplePopupContent.Widget().SetPopup(&p.simplePopup)
	p.simplePopup.SetContent(&p.simplePopupContent)
	p.simplePopup.SetBackgroundBlurred(p.blurBackgroundToggle.Value())
	p.simplePopup.SetCloseByClickingOutside(p.closeByClickingOutsideToggle.Value())
	p.simplePopup.SetAnimationDuringFade(true)

	p.simplePopupContent.SetFixedSize(p.contentSize(context))

	p.contextMenuPopup.SetItemsByStrings([]string{"Item 1", "Item 2", "Item 3"})
	// A context menu's position is updated at HandlePointingInput.

	return nil
}

func (p *Popups) contentSize(context *guigui.Context) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(int(12*u), int(6*u))
}

func (p *Popups) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &p.forms[0]:
		return p.layout.CellBounds(0, 0)
	case &p.forms[1]:
		return p.layout.CellBounds(0, 1)
	case &p.simplePopup:
		appBounds := context.AppBounds()
		contentSize := p.contentSize(context)
		p := image.Point{
			X: appBounds.Min.X + (appBounds.Dx()-contentSize.X)/2,
			Y: appBounds.Min.Y + (appBounds.Dy()-contentSize.Y)/2,
		}
		return image.Rectangle{
			Min: p,
			Max: p.Add(contentSize),
		}
	case &p.contextMenuPopup:
		return image.Rectangle{
			Min: p.contextMenuPopupPosition,
			Max: p.contextMenuPopupPosition.Add(p.contextMenuPopup.Measure(context, guigui.Constraints{})),
		}
	}
	return image.Rectangle{}
}

func (p *Popups) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		// Use IsWidgetOrBackgroundHitAtCursor. context.IsWidgetHitAtCursor doesn't work when a popup's transparent background exists.
		if p.contextMenuPopup.IsWidgetOrBackgroundHitAtCursor(context, &p.contextMenuPopupClickHereText) {
			p.contextMenuPopupPosition = image.Pt(ebiten.CursorPosition())
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

	mainLayout   layout.GridLayout
	footerLayout layout.GridLayout
}

func (s *simplePopupContent) SetPopup(popup *basicwidget.Popup) {
	s.popup = popup
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

	s.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.LazySize(func(row int) layout.Size {
				if row != 1 {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(s.closeButton.Measure(context, guigui.Constraints{}).Y)
			}),
		},
	}
	s.footerLayout = layout.GridLayout{
		Bounds: s.mainLayout.CellBounds(0, 1),
		Widths: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(s.closeButton.Measure(context, guigui.Constraints{}).X),
		},
	}

	return nil
}

func (s *simplePopupContent) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &s.titleText:
		return s.mainLayout.CellBounds(0, 0)
	case &s.closeButton:
		return s.footerLayout.CellBounds(1, 0)
	}
	return image.Rectangle{}
}
