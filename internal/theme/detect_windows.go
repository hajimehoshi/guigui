//go:build windows
// +build windows

package theme

import (
	"golang.org/x/sys/windows/registry"
)

func detectSystemTheme() ColorMode {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`,
		registry.QUERY_VALUE)
	if err != nil {
		return ThemeUnknown
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return ThemeUnknown
	}

	if val == 0 {
		return ThemeDark
	}
	return ThemeLight
}
