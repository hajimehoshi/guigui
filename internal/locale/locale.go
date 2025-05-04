// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package locale

import "golang.org/x/text/language"

func Locales() ([]language.Tag, error) {
	ls, err := locales()
	if err != nil {
		return nil, err
	}
	tags := make([]language.Tag, len(ls))
	for i, l := range ls {
		tags[i] = language.Make(l)
	}
	return tags, nil
}
