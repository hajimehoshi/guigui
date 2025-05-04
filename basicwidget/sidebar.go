// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type Sidebar struct {
	guigui.DefaultWidget

	scrollablePanel ScrollablePanel
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidgetWithBounds(&s.scrollablePanel, context.Bounds(s))
	return nil
}

func (s *Sidebar) SetContent(widget guigui.Widget) {
	s.scrollablePanel.SetContent(widget)
}

func (s *Sidebar) Draw(context *guigui.Context, dst *ebiten.Image) {
	dst.Fill(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.875))
	b := context.Bounds(s)
	b.Min.X = b.Max.X - int(1*context.Scale())
	dst.SubImage(b).(*ebiten.Image).Fill(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85))
}
