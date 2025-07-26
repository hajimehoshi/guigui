// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
	"github.com/hajimehoshi/guigui/layout"
)

type SegmentedControlDirection int

const (
	SegmentedControlDirectionHorizontal SegmentedControlDirection = iota
	SegmentedControlDirectionVertical
)

type SegmentedControlItem[T comparable] struct {
	Text      string
	Icon      *ebiten.Image
	IconAlign IconAlign
	Disabled  bool
	ID        T
}

func (s SegmentedControlItem[T]) id() T {
	return s.ID
}

type SegmentedControl[T comparable] struct {
	guigui.DefaultWidget

	abstractList abstractList[T, SegmentedControlItem[T]]
	buttons      []Button

	direction SegmentedControlDirection
}

func (s *SegmentedControl[T]) SetDirection(direction SegmentedControlDirection) {
	if s.direction == direction {
		return
	}
	s.direction = direction
	guigui.RequestRedraw(s)
}

func (s *SegmentedControl[T]) SetOnItemSelected(f func(index int)) {
	s.abstractList.SetOnItemSelected(f)
}

func (s *SegmentedControl[T]) SetItems(items []SegmentedControlItem[T]) {
	s.abstractList.SetItems(items)
}

func (s *SegmentedControl[T]) SelectedItem() (SegmentedControlItem[T], bool) {
	return s.abstractList.SelectedItem()
}

func (s *SegmentedControl[T]) SelectedItemIndex() int {
	return s.abstractList.SelectedItemIndex()
}

func (s *SegmentedControl[T]) ItemByIndex(index int) (SegmentedControlItem[T], bool) {
	return s.abstractList.ItemByIndex(index)
}

func (s *SegmentedControl[T]) SelectItemByIndex(index int) {
	if s.abstractList.SelectItemByIndex(index, false) {
		guigui.RequestRedraw(s)
	}
}

func (s *SegmentedControl[T]) SelectItemByID(id T) {
	if s.abstractList.SelectItemByID(id, false) {
		guigui.RequestRedraw(s)
	}
}

func (s *SegmentedControl[T]) ResetEventHandlers() {
	s.abstractList.ResetEventHandlers()
}

func (s *SegmentedControl[T]) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	for i := range s.buttons {
		appender.AppendChildWidget(&s.buttons[i])
	}
}

func (s *SegmentedControl[T]) Build(context *guigui.Context) error {
	s.buttons = adjustSliceSize(s.buttons, s.abstractList.ItemCount())

	sizes := make([]layout.Size, s.abstractList.ItemCount())
	for i := range s.abstractList.ItemCount() {
		item, _ := s.abstractList.ItemByIndex(i)
		s.buttons[i].SetText(item.Text)
		s.buttons[i].SetIcon(item.Icon)
		s.buttons[i].SetIconAlign(item.IconAlign)
		s.buttons[i].SetTextBold(s.abstractList.SelectedItemIndex() == i)
		s.buttons[i].setUseAccentColor(true)
		if s.abstractList.ItemCount() > 1 {
			switch i {
			case 0:
				switch s.direction {
				case SegmentedControlDirectionHorizontal:
					s.buttons[i].setSharpenCorners(draw.SharpenCorners{
						UpperEnd: true,
						LowerEnd: true,
					})
				case SegmentedControlDirectionVertical:
					s.buttons[i].setSharpenCorners(draw.SharpenCorners{
						LowerStart: true,
						LowerEnd:   true,
					})
				}
			case s.abstractList.ItemCount() - 1:
				switch s.direction {
				case SegmentedControlDirectionHorizontal:
					s.buttons[i].setSharpenCorners(draw.SharpenCorners{
						UpperStart: true,
						LowerStart: true,
					})
				case SegmentedControlDirectionVertical:
					s.buttons[i].setSharpenCorners(draw.SharpenCorners{
						UpperEnd:   true,
						UpperStart: true,
					})
				}
			default:
				s.buttons[i].setSharpenCorners(draw.SharpenCorners{
					UpperStart: true,
					LowerStart: true,
					UpperEnd:   true,
					LowerEnd:   true,
				})
			}
		}
		context.SetEnabled(&s.buttons[i], !item.Disabled)
		s.buttons[i].setKeepPressed(s.abstractList.SelectedItemIndex() == i)
		s.buttons[i].SetOnDown(func() {
			s.SelectItemByIndex(i)
		})
		sizes[i] = layout.FlexibleSize(1)
	}

	var g layout.GridLayout
	switch s.direction {
	case SegmentedControlDirectionHorizontal:
		g = layout.GridLayout{
			Bounds: context.Bounds(s),
			Widths: sizes,
		}
	case SegmentedControlDirectionVertical:
		g = layout.GridLayout{
			Bounds:  context.Bounds(s),
			Heights: sizes,
		}
	}

	for i := range s.buttons {
		switch s.direction {
		case SegmentedControlDirectionHorizontal:
			context.SetBounds(&s.buttons[i], g.CellBounds(i, 0), s)
		case SegmentedControlDirectionVertical:
			context.SetBounds(&s.buttons[i], g.CellBounds(0, i), s)
		}
	}

	return nil
}

func (s *SegmentedControl[T]) DefaultSize(context *guigui.Context) image.Point {
	var w, h int
	for i := range s.buttons {
		size := s.buttons[i].defaultSize(context, true)
		w = max(w, size.X)
		h = max(h, size.Y)
	}
	switch s.direction {
	case SegmentedControlDirectionHorizontal:
		return image.Pt(w*len(s.buttons), h)
	case SegmentedControlDirectionVertical:
		return image.Pt(w, h*len(s.buttons))
	default:
		panic(fmt.Sprintf("basicwidget: unknown direction %d", s.direction))
	}
}
