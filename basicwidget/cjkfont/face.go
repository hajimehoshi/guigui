// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package cjkfont

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

//go:generate go run gen.go

//go:embed NotoSansCJK-VF.otf.ttc.gz
var notoSansCJKVFOTFTTCGz []byte

var (
	theFaceSourceSC *text.GoTextFaceSource
	theFaceSourceTC *text.GoTextFaceSource
	theFaceSourceHK *text.GoTextFaceSource
	theFaceSourceJP *text.GoTextFaceSource
	theFaceSourceKR *text.GoTextFaceSource
)

func init() {
	r, err := gzip.NewReader(bytes.NewReader(notoSansCJKVFOTFTTCGz))
	if err != nil {
		panic(err)
	}
	fs, err := text.NewGoTextFaceSourcesFromCollection(r)
	if err != nil {
		panic(err)
	}
	var (
		faceSC *text.GoTextFaceSource
		faceTC *text.GoTextFaceSource
		faceHK *text.GoTextFaceSource
		faceJP *text.GoTextFaceSource
		faceKR *text.GoTextFaceSource
	)
	for _, f := range fs {
		switch f.Metadata().Family {
		case "Noto Sans CJK SC":
			faceSC = f
		case "Noto Sans CJK TC":
			faceTC = f
		case "Noto Sans CJK HK":
			faceHK = f
		case "Noto Sans CJK JP":
			faceJP = f
		case "Noto Sans CJK KR":
			faceKR = f
		default:
			panic(fmt.Sprintf("cjkfont: unknown family: %s", f.Metadata().Family))
		}
	}

	theFaceSourceSC = faceSC
	theFaceSourceTC = faceTC
	theFaceSourceHK = faceHK
	theFaceSourceJP = faceJP
	theFaceSourceKR = faceKR
}

func FaceSourceSC() *text.GoTextFaceSource {
	return theFaceSourceSC
}

func FaceSourceTC() *text.GoTextFaceSource {
	return theFaceSourceTC
}

func FaceSourceHK() *text.GoTextFaceSource {
	return theFaceSourceHK
}

func FaceSourceJP() *text.GoTextFaceSource {
	return theFaceSourceJP
}

func FaceSourceKR() *text.GoTextFaceSource {
	return theFaceSourceKR
}

var (
	ja = language.MustParseBase("ja")
	ko = language.MustParseBase("ko")
	zh = language.MustParseBase("zh")
	cn = language.MustParseRegion("CN")
	tw = language.MustParseRegion("TW")
	mo = language.MustParseRegion("MO")
	hk = language.MustParseRegion("HK")
)

func FaceSourceFromLocale(locale language.Tag) *text.GoTextFaceSource {
	if locale == language.Und {
		return nil
	}

	base, conf := locale.Base()
	if conf == language.No {
		return nil
	}
	switch base {
	case ja:
		return theFaceSourceJP
	case ko:
		return theFaceSourceKR
	case zh:
		region, conf := locale.Region()
		if conf == language.No {
			return nil
		}
		switch region {
		case cn:
			return theFaceSourceSC
		case tw, mo:
			return theFaceSourceTC
		case hk:
			return theFaceSourceHK
		}
	}
	return nil
}
