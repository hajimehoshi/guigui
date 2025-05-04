// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

//go:build !darwin

package basicwidget

import (
	"regexp"
	"syscall/js"
)

var (
	isMacintosh = regexp.MustCompile(`\bMacintosh\b`)
	isIPhone    = regexp.MustCompile(`\biPhone\b`)
	isIPad      = regexp.MustCompile(`\biPad\b`)
)

var isDarwin bool

func init() {
	ua := js.Global().Get("navigator").Get("userAgent").String()
	if isMacintosh.MatchString(ua) {
		isDarwin = true
		return
	}
	if isIPhone.MatchString(ua) {
		isDarwin = true
		return
	}
	if isIPad.MatchString(ua) {
		isDarwin = true
		return
	}
}

func useEmacsKeybind() bool {
	return isDarwin
}
