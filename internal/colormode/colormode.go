// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package colormode

type ColorMode int

const (
	Unknown ColorMode = iota
	Light
	Dark
)

func SystemColorMode() ColorMode {
	return systemColorMode()
}
