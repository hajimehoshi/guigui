//go:build !darwin && !linux && !windows

// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package colormode

func systemColorMode() ColorMode {
	// Unknown default to Light mode or what should I put in this function?
	return Light
}
