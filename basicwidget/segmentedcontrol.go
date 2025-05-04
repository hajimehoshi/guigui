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
	Text     string
	Image    *ebiten.Image
	Disabled bool
	Tag      T
}

func (s SegmentedControlItem[T]) tag() T {
	return s.Tag
}

type SegmentedControl[T comparable] struct {
	guigui.DefaultWidget

	abstractList abstractList[T, SegmentedControlItem[T]]
	textButtons  []TextButton

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
	if s.abstractList.SelectItemByIndex(index) {
		guigui.RequestRedraw(s)
	}
}

func (s *SegmentedControl[T]) SelectItemByTag(tag T) {
	if s.abstractList.SelectItemByTag(tag) {
		guigui.RequestRedraw(s)
	}
}

func (s *SegmentedControl[T]) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.textButtons = adjustSliceSize(s.textButtons, s.abstractList.ItemCount())

	sizes := make([]layout.Size, s.abstractList.ItemCount())
	for i := range s.abstractList.ItemCount() {
		item, _ := s.abstractList.ItemByIndex(i)
		s.textButtons[i].SetText(item.Text)
		s.textButtons[i].SetImage(item.Image)
		s.textButtons[i].SetTextBold(s.abstractList.SelectedItemIndex() == i && context.IsEnabled(s))
		s.textButtons[i].setUseAccentColor(true)
		if s.abstractList.ItemCount() > 1 {
			switch i {
			case 0:
				switch s.direction {
				case SegmentedControlDirectionHorizontal:
					s.textButtons[i].setSharpenCorners(draw.SharpenCorners{
						UpperRight: true,
						LowerRight: true,
					})
				case SegmentedControlDirectionVertical:
					s.textButtons[i].setSharpenCorners(draw.SharpenCorners{
						LowerLeft:  true,
						LowerRight: true,
					})
				}
			case s.abstractList.ItemCount() - 1:
				switch s.direction {
				case SegmentedControlDirectionHorizontal:
					s.textButtons[i].setSharpenCorners(draw.SharpenCorners{
						UpperLeft: true,
						LowerLeft: true,
					})
				case SegmentedControlDirectionVertical:
					s.textButtons[i].setSharpenCorners(draw.SharpenCorners{
						UpperRight: true,
						UpperLeft:  true,
					})
				}
			default:
				s.textButtons[i].setSharpenCorners(draw.SharpenCorners{
					UpperLeft:  true,
					LowerLeft:  true,
					UpperRight: true,
					LowerRight: true,
				})
			}
		}
		context.SetEnabled(&s.textButtons[i], !item.Disabled)
		s.textButtons[i].setKeepPressed(s.abstractList.SelectedItemIndex() == i)
		s.textButtons[i].SetOnDown(func() {
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

	for i := range s.textButtons {
		switch s.direction {
		case SegmentedControlDirectionHorizontal:
			appender.AppendChildWidgetWithBounds(&s.textButtons[i], g.CellBounds(i, 0))
		case SegmentedControlDirectionVertical:
			appender.AppendChildWidgetWithBounds(&s.textButtons[i], g.CellBounds(0, i))
		}
	}

	return nil
}

func (s *SegmentedControl[T]) DefaultSize(context *guigui.Context) image.Point {
	var w, h int
	for i := range s.textButtons {
		size := s.textButtons[i].defaultSize(context, true)
		w = max(w, size.X)
		h = max(h, size.Y)
	}
	switch s.direction {
	case SegmentedControlDirectionHorizontal:
		return image.Pt(w*len(s.textButtons), h)
	case SegmentedControlDirectionVertical:
		return image.Pt(w, h*len(s.textButtons))
	default:
		panic(fmt.Sprintf("basicwidget: unknown direction %d", s.direction))
	}
}
