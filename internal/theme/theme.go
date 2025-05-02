package theme

type ColorMode int

const (
	ThemeLight ColorMode = iota
	ThemeDark
	ThemeUnknown
)

func DetectSystemTheme() ColorMode {
	return detectSystemTheme() // implemented in per-OS files
}
