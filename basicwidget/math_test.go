// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package basicwidget_test

import (
	"math"
	"testing"

	"github.com/hajimehoshi/guigui/basicwidget"
)

func TestMaxInteger(t *testing.T) {
	if got, want := basicwidget.MaxInteger[int](), int(math.MaxInt); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[int8](), int8(math.MaxInt8); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[int16](), int16(math.MaxInt16); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[int32](), int32(math.MaxInt32); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[int64](), int64(math.MaxInt64); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[uint](), uint(math.MaxUint); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[uint8](), uint8(math.MaxUint8); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[uint16](), uint16(math.MaxUint16); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[uint32](), uint32(math.MaxUint32); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[uint64](), uint64(math.MaxUint64); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MaxInteger[uintptr](), uintptr(math.MaxUint); got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	type MyInt int
	if got, want := basicwidget.MaxInteger[MyInt](), MyInt(math.MaxInt); got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	type MyUint uint
	if got, want := basicwidget.MaxInteger[MyUint](), MyUint(math.MaxUint); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestMinInteger(t *testing.T) {
	if got, want := basicwidget.MinInteger[int](), int(math.MinInt); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[int8](), int8(math.MinInt8); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[int16](), int16(math.MinInt16); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[int32](), int32(math.MinInt32); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[int64](), int64(math.MinInt64); got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := basicwidget.MinInteger[uint](), uint(0); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[uint8](), uint8(0); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[uint16](), uint16(0); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[uint32](), uint32(0); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[uint64](), uint64(0); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	if got, want := basicwidget.MinInteger[uintptr](), uintptr(0); got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	type MyInt int
	if got, want := basicwidget.MinInteger[MyInt](), MyInt(math.MinInt); got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	type MyUint uint
	if got, want := basicwidget.MinInteger[MyUint](), MyUint(0); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
