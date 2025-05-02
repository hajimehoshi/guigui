//go:build darwin
// +build darwin

package theme

import (
	"fmt"
	"os/exec"
	"strings"
)

func detectSystemTheme() ColorMode {
	out, err := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle").Output()
	if err != nil {
		return ThemeLight
	}
	if strings.Contains(strings.ToLower(string(out)), "dark") {
		return ThemeDark
	}
	fmt.Println(string(out))
	return ThemeLight
}
