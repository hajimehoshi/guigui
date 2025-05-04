// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

//go:build !darwin

package locale

import "github.com/jeandeaual/go-locale"

func locales() ([]string, error) {
	return locale.GetLocales()
}
