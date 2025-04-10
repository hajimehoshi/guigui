// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package basicwidget

import (
	"image"

	"github.com/hajimehoshi/guigui"
)

type PopupMenu struct {
	guigui.DefaultWidget

	textList TextList
	popup    Popup

	onClosed func(index int)
}

func (p *PopupMenu) SetOnClosed(f func(index int)) {
	p.onClosed = f
}

func (p *PopupMenu) SetCheckmarkIndex(index int) {
	p.textList.SetCheckmarkIndex(index)
}

func (p *PopupMenu) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	p.popup.SetContent(func(context *guigui.Context, childAppender *ContainerChildWidgetAppender) {
		p.textList.SetStyle(ListStyleMenu)
		p.textList.list.SetOnItemSelected(func(index int) {
			p.popup.Close()
			if p.onClosed != nil {
				p.onClosed(index)
			}
		})
		guigui.SetPosition(&p.textList, p.popup.ContentBounds(context).Min)
		childAppender.AppendChildWidget(&p.textList)
	})
	p.popup.SetCloseByClickingOutside(true)
	p.popup.SetOnClosed(func(reason PopupClosedReason) {
		if p.onClosed != nil {
			p.onClosed(p.textList.SelectedItemIndex())
		}
	})
	p.popup.SetContentBounds(p.contentBounds(context))
	appender.AppendChildWidget(&p.popup)

	// Sync the visibility with the popup.
	// TODO: This is tricky. Refactor this. Perhaps, introducing Widget.PassThrough might be a good idea.
	if guigui.IsVisible(&p.popup) {
		guigui.Show(p)
	} else {
		guigui.Hide(p)
	}

	return nil
}

func (p *PopupMenu) contentBounds(context *guigui.Context) image.Rectangle {
	pos := guigui.Position(p)
	w, h := p.textList.Size(context)
	if h > 24*UnitSize(context) {
		h = 24 * UnitSize(context)
		p.textList.SetHeight(h)
	}
	r := image.Rectangle{
		Min: pos,
		Max: pos.Add(image.Pt(w, h)),
	}
	aw, ah := context.AppSize()
	if r.Max.X > aw {
		r.Min.X = aw - w
		r.Max.X = aw
	}
	if r.Min.X < 0 {
		r.Min.X = 0
		r.Max.X = w
	}
	if r.Max.Y > ah {
		r.Min.Y = ah - h
		r.Max.Y = ah
	}
	if r.Min.Y < 0 {
		r.Min.Y = 0
		r.Max.Y = h
	}
	return r
}

func (p *PopupMenu) Z() int {
	return guigui.Parent(p).Z() + popupZ
}

func (p *PopupMenu) Open() {
	guigui.Show(p)
	p.popup.Open()
}

func (p *PopupMenu) Close() {
	p.popup.Close()
}

func (p *PopupMenu) SetItemsByStrings(items []string) {
	p.textList.SetItemsByStrings(items)
}

func (p *PopupMenu) SelectedItem() (TextListItem, bool) {
	return p.textList.SelectedItem()
}

func (p *PopupMenu) SelectedItemIndex() int {
	return p.textList.SelectedItemIndex()
}

func (p *PopupMenu) SetSelectedItemIndex(index int) {
	p.textList.SetSelectedItemIndex(index)
}
