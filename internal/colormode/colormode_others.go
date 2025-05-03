// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

//go:build !darwin && !windows

package colormode

func systemColorMode() ColorMode {
	return Light
}
