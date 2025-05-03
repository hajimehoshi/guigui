// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package colormode

import (
	"os/exec"
	"strings"
)

func systemColorMode() ColorMode {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme").Output()
	if err != nil {
		return Unknown
	}
	if strings.Contains(strings.ToLower(string(out)), "dark") {
		return Dark
	}
	return Light
}
