// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package basicwidget

func MaxInteger[T Integer]() T {
	return maxInteger[T]()
}

func MinInteger[T Integer]() T {
	return minInteger[T]()
}
