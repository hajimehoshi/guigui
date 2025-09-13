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
	Value     T
}

func (s SegmentedControlItem[T]) value() T {
	return s.Value
}

type SegmentedControl[T comparable] struct {
	guigui.DefaultWidget

	abstractList abstractList[T, SegmentedControlItem[T]]
	buttons      []Button

	direction   SegmentedControlDirection
	layoutSizes []layout.Size
}

func (s *SegmentedControl[T]) SetDirection(direction SegmentedControlDirection) {
	if s.direction == direction {
		return
	}
	s.direction = direction
	guigui.RequestRedraw(s)
}

func (s *SegmentedControl[T]) SetOnItemSelected(f func(index int)) {
	s.abstractList.SetOnItemSelected(s, f)
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
	if s.abstractList.SelectItemByIndex(s, index, false) {
		guigui.RequestRedraw(s)
	}
}

func (s *SegmentedControl[T]) SelectItemByValue(value T) {
	if s.abstractList.SelectItemByValue(s, value, false) {
		guigui.RequestRedraw(s)
	}
}

func (s *SegmentedControl[T]) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	for i := range s.buttons {
		adder.AddChild(&s.buttons[i])
	}
}

func (s *SegmentedControl[T]) Update(context *guigui.Context) error {
	s.buttons = adjustSliceSize(s.buttons, s.abstractList.ItemCount())

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
	}

	return nil
}

func (s *SegmentedControl[T]) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	s.layoutSizes = adjustSliceSize(s.layoutSizes, s.abstractList.ItemCount())
	for i := range s.abstractList.ItemCount() {
		s.layoutSizes[i] = layout.FlexibleSize(1)
	}

	var g layout.GridLayout
	switch s.direction {
	case SegmentedControlDirectionHorizontal:
		g = layout.GridLayout{
			Bounds: context.Bounds(s),
			Widths: s.layoutSizes,
		}
	case SegmentedControlDirectionVertical:
		g = layout.GridLayout{
			Bounds:  context.Bounds(s),
			Heights: s.layoutSizes,
		}
	}

	idx := -1
	for i := range s.buttons {
		if &s.buttons[i] == widget {
			idx = i
			break
		}

	}
	if idx >= 0 {
		switch s.direction {
		case SegmentedControlDirectionHorizontal:
			return g.CellBounds(idx, 0)
		case SegmentedControlDirectionVertical:
			return g.CellBounds(0, idx)
		}
	}

	return image.Rectangle{}
}

func (s *SegmentedControl[T]) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	var w, h int
	for i := range s.buttons {
		size := s.buttons[i].defaultSize(context, constraints, true)
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
