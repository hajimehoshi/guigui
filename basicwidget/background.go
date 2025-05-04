// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"

	"github.com/hajimehoshi/guigui"
)

type Background struct {
	guigui.DefaultWidget
}

func (b *Background) Draw(context *guigui.Context, dst *ebiten.Image) {
	dst.Fill(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.95))
}
