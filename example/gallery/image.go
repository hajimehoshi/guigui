// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"embed"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

//go:embed resource/*.png
var pngImages embed.FS

type imageCacheKey struct {
	name      string
	colorMode guigui.ColorMode
}

type imageCache struct {
	m map[imageCacheKey]*ebiten.Image
}

var theImageCache = &imageCache{}

func (i *imageCache) Get(name string, colorMode guigui.ColorMode) (*ebiten.Image, error) {
	key := imageCacheKey{
		name:      name,
		colorMode: colorMode,
	}
	if img, ok := i.m[key]; ok {
		return img, nil
	}

	f, err := pngImages.Open("resource/" + name + ".png")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	pImg, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	pImg = basicwidget.CreateMonochromeImage(colorMode, pImg)

	img := ebiten.NewImageFromImage(pImg)
	if i.m == nil {
		i.m = map[imageCacheKey]*ebiten.Image{}
	}
	i.m[key] = img
	return img, nil
}
