// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package colormode

import (
	"os/exec"
	"strings"
)

func systemColorMode() ColorMode {
	out, err := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle").Output()
	if err != nil {
		return Light
	}
	if strings.TrimSpace(string(out)) == "Dark" {
		return Dark
	}
	return Light
}
