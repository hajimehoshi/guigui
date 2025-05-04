// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package locale

import (
	"errors"
	"fmt"
	"os/exec"

	"howett.net/plist"
)

func locales() ([]string, error) {
	// `[[NSBundle mainBundle] preferredLocalizations]` might be available,
	// but this could return only one language.
	cmd := exec.Command("defaults", "read", "-g", "AppleLanguages")
	out, err := cmd.Output()
	if err != nil {
		if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
			out = append(out, ee.Stderr...)
		}
		return nil, fmt.Errorf("locale: `defaults read -g AppleLanguages` failed: %w\n%s", err, out)
	}

	var locales []string
	if _, err := plist.Unmarshal(out, &locales); err != nil {
		return nil, fmt.Errorf("locale: plist.Unmarshal failed: %w", err)
	}
	return locales, nil
}
