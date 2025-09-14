// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

//go:build js

package clipboard

import "syscall/js"

func ReadAll() (string, error) {
	ch := make(chan string)
	then := js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- args[0].String()
		return nil
	})
	defer then.Release()
	js.Global().Get("navigator").Get("clipboard").Call("readText").Call("then", then)
	return <-ch, nil
}

func WriteAll(text string) error {
	ch := make(chan struct{})
	then := js.FuncOf(func(this js.Value, args []js.Value) any {
		close(ch)
		return nil
	})
	defer then.Release()
	js.Global().Get("navigator").Get("clipboard").Call("writeText", text).Call("then", then)
	<-ch
	return nil
}
