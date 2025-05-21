// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
)

type GifPlay int

const (
	GifPlayLoop GifPlay = iota
	GifPlayOnce
)

type Gif struct {
	guigui.DefaultWidget

	thumbnail *ebiten.Image
	frame     []*ebiten.Image
	delay     []int

	loop     int
	lastTick int64
	gifPlay  GifPlay

	canvas Image
}

func (g *Gif) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if g.thumbnail == nil {
		return nil
	}

	if len(g.frame) == 0 {
		g.canvas.SetImage(g.thumbnail)
	} else {
		g.canvas.SetImage(g.frame[g.loop])
	}

	appender.AppendChildWidgetWithBounds(&g.canvas, context.Bounds(g))
	return nil
}

func (g *Gif) Tick(context *guigui.Context) error {
	if len(g.frame) == 0 || len(g.delay) == 0 {
		return nil
	}

	currentTick := ebiten.Tick()
	tps := ebiten.TPS()

	requiredDelayTicks := int64((g.delay[g.loop] * tps) / 100)
	if requiredDelayTicks < 1 {
		requiredDelayTicks = 1
	}

	if currentTick-g.lastTick >= requiredDelayTicks {
		newFrame := g.loop + 1

		switch g.gifPlay {
		case GifPlayOnce:
			if newFrame >= len(g.frame) {
				newFrame = len(g.frame) - 1
			} else {
				newFrame %= len(g.frame)
			}
		case GifPlayLoop:
			newFrame %= len(g.frame)
		}

		if newFrame != g.loop {
			g.loop = newFrame
			g.lastTick = currentTick
		}
	}

	return nil
}

func (g *Gif) Play() GifPlay {
	return g.gifPlay
}

func (g *Gif) SetGif(thumbnail *ebiten.Image, frame []*ebiten.Image, delay []int) {
	if g.thumbnail == thumbnail {
		return
	}
	g.thumbnail = thumbnail
	g.frame = frame
	g.delay = delay
	g.loop = 0
	g.lastTick = ebiten.Tick()
	guigui.RequestRedraw(g)
}

func (g *Gif) SetPlay(play GifPlay) {
	if g.gifPlay == play {
		return
	}
	g.gifPlay = play
	g.loop = 0
	g.lastTick = ebiten.Tick()
}
