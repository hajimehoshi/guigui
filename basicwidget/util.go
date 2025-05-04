// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import "slices"

func adjustSliceSize[T any](slice []T, size int) []T {
	if len(slice) == size {
		return slice
	}
	if len(slice) < size {
		return slices.Grow(slice, size-len(slice))[:size]
	}
	return slices.Delete(slice, size, len(slice))
}
