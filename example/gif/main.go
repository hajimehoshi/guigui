// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.DefaultWidget

	background basicwidget.Background
	pic        basicwidget.Image
	gif1       basicwidget.Gif
	gif2       basicwidget.Gif

	thumbnail *ebiten.Image
	gifFrames []*ebiten.Image
	gifDelay  []int

	sync sync.Once
	err  error
}

func (r *Root) loadGifThumbnail(path string) (*ebiten.Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	gifImage, err := gif.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(gifImage), nil
}

func (r *Root) loadGif(path string) ([]*ebiten.Image, []int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	gifImage, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	var frames []*ebiten.Image
	var delay []int

	bounds := gifImage.Image[0].Bounds()
	canvas := image.NewRGBA(bounds)

	for i, frame := range gifImage.Image {
		if i > 0 {
			switch gifImage.Disposal[i-1] {
			case gif.DisposalNone:
				//
			case gif.DisposalBackground:
				prevBounds := gifImage.Image[i-1].Bounds()
				for y := prevBounds.Min.Y; y < prevBounds.Max.Y; y++ {
					for x := prevBounds.Min.X; x < prevBounds.Max.X; x++ {
						canvas.Set(x, y, image.Transparent)
					}
				}
			case gif.DisposalPrevious:
				canvas = image.NewRGBA(bounds)
			default:
				canvas = image.NewRGBA(bounds)
			}
		}

		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		frameImg := ebiten.NewImageFromImage(canvas)

		frames = append(frames, frameImg)
		delay = append(delay, gifImage.Delay[i])
	}

	return frames, delay, nil
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.sync.Do(func() {
		filename := filepath.Join("example", "gif", "dancing-gopher.gif")

		thumbnail, err := r.loadGifThumbnail(filename)
		if err != nil {
			r.err = err
			return
		}
		r.thumbnail = thumbnail

		frames, delays, err := r.loadGif(filename)
		if err != nil {
			r.err = err
			return
		}
		r.gifFrames = frames
		r.gifDelay = delays
	})

	if r.err != nil {
		return r.err
	}

	r.gif1.SetPlay(basicwidget.GifPlayOnce)

	r.pic.SetImage(r.thumbnail)
	r.gif1.SetGif(r.thumbnail, r.gifFrames, r.gifDelay)
	r.gif2.SetGif(r.thumbnail, r.gifFrames, r.gifDelay)

	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
		},
	}

	appender.AppendChildWidgetWithBounds(&r.pic, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&r.gif1, gl.CellBounds(0, 1))
	appender.AppendChildWidgetWithBounds(&r.gif2, gl.CellBounds(0, 2))

	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title: "Gif",
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
