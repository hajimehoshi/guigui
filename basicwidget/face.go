// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package basicwidget

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

//go:generate go run gen.go

//go:embed InterVariable.ttf.gz
var interVariableTTFGz []byte

var theDefaultFaceSource *text.GoTextFaceSource

var (
	theFaceCache   map[faceCacheKey]text.Face
	theFaceSources []*text.GoTextFaceSource
)

func init() {
	r, err := gzip.NewReader(bytes.NewReader(interVariableTTFGz))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = r.Close()
	}()
	f, err := text.NewGoTextFaceSource(r)
	if err != nil {
		panic(err)
	}
	theDefaultFaceSource = f
	theFaceSources = []*text.GoTextFaceSource{f}
}

func fontFace(size float64, weight text.Weight, features []fontFeature, lang language.Tag) text.Face {
	key := faceCacheKey{
		size:     size,
		weight:   weight,
		features: string(serializeFontFeatures(features)),
		lang:     lang,
	}
	if f, ok := theFaceCache[key]; ok {
		return f
	}

	var fs []text.Face
	for _, faceSource := range theFaceSources {
		f := &text.GoTextFace{
			Source:   faceSource,
			Size:     size,
			Language: lang,
		}
		f.SetVariation(text.MustParseTag("wght"), float32(weight))
		for _, ff := range features {
			f.SetFeature(ff.Tag, ff.Value)
		}
		fs = append(fs, f)
	}
	mf, err := text.NewMultiFace(fs...)
	if err != nil {
		panic(err)
	}

	if theFaceCache == nil {
		theFaceCache = map[faceCacheKey]text.Face{}
	}
	theFaceCache[key] = mf

	return mf
}

func DefaultFaceSource() *text.GoTextFaceSource {
	return theDefaultFaceSource
}

func SetFaceSources(faces []*text.GoTextFaceSource) {
	if len(faces) == 0 {
		faces = []*text.GoTextFaceSource{theDefaultFaceSource}
	}
	if slices.Equal(theFaceSources, faces) {
		return
	}
	theFaceSources = slices.Clone(faces)
	clear(theFaceCache)
}

type faceCacheKey struct {
	size     float64
	weight   text.Weight
	features string
	lang     language.Tag
}

type fontFeature struct {
	Tag   text.Tag
	Value uint32
}

func serializeFontFeatures(features []fontFeature) []byte {
	if len(features) == 0 {
		return nil
	}

	b := make([]byte, 0, len(features)*8)
	for _, ff := range features {
		b = append(b, byte(ff.Tag), byte(ff.Tag>>8), byte(ff.Tag>>16), byte(ff.Tag>>24))
		b = append(b, byte(ff.Value), byte(ff.Value>>8), byte(ff.Value>>16), byte(ff.Value>>24))
	}
	return b
}
