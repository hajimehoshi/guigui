// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package colormode

import (
	"golang.org/x/sys/windows/registry"
)

func systemColorMode() ColorMode {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`,
		registry.QUERY_VALUE)
	defer func() {
		_ = k.Close()
	}()

	val, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return Unknown
	}

	if val == 0 {
		return Dark
	}
	return Light
}
