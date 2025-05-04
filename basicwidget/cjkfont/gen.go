// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

//go:build ignore

package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if err := xmain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func xmain() error {
	// Get `Variable/OTC/NotoSansCJK-VF.otf.ttc` (Full) from https://github.com/notofonts/noto-cjk/releases/tag/Sans2.004.
	const url = "https://github.com/notofonts/noto-cjk/releases/download/Sans2.004/01_NotoSansCJK-OTF-VF.zip"
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	r, err := zip.NewReader(bytes.NewReader(bs), int64(len(bs)))
	if err != nil {
		return err
	}
	in, err := r.Open("Variable/OTC/NotoSansCJK-VF.otf.ttc")
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create("NotoSansCJK-VF.otf.ttc.gz")
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return err
	}

	if _, err := io.Copy(gw, in); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
