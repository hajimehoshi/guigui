// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
)

type Image struct {
	guigui.DefaultWidget

	image *ebiten.Image
}

func (i *Image) Draw(context *guigui.Context, dst *ebiten.Image) {
	if i.image == nil {
		return
	}

	p := context.Position(i)
	s := context.Size(i)
	imgScale := min(float64(s.X)/float64(i.image.Bounds().Dx()), float64(s.Y)/float64(i.image.Bounds().Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(imgScale, imgScale)
	op.GeoM.Translate(float64(p.X), float64(p.Y))
	if !context.IsEnabled(i) {
		// TODO: Reduce the saturation?
		op.ColorScale.ScaleAlpha(0.25)
	}
	// TODO: Use a better filter.
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(i.image, op)
}

func (i *Image) HasImage() bool {
	return i.image != nil
}

func (i *Image) SetImage(image *ebiten.Image) {
	if i.image == image {
		return
	}
	i.image = image
	guigui.RequestRedraw(i)
}
