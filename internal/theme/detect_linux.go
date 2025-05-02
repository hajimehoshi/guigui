//go:build linux
// +build linux

package theme

import (
	"os/exec"
	"strings"
)

func detectSystemTheme() SystemTheme {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme").Output()
	if err != nil {
		return ThemeUnknown
	}
	if strings.Contains(strings.ToLower(string(out)), "dark") {
		return ThemeDark
	}
	return ThemeLight
}
