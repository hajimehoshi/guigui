// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"math/big"

	"github.com/hajimehoshi/guigui/basicwidget"
)

type Model struct {
	mode string

	buttons      ButtonsModel
	texts        TextsModel
	textInputs   TextInputsModel
	numberInputs NumberInputsModel
}

func (m *Model) Mode() string {
	if m.mode == "" {
		return "settings"
	}
	return m.mode
}

func (m *Model) SetMode(mode string) {
	m.mode = mode
}

func (m *Model) Buttons() *ButtonsModel {
	return &m.buttons
}

func (m *Model) Texts() *TextsModel {
	return &m.texts
}

func (m *Model) TextInputs() *TextInputsModel {
	return &m.textInputs
}

func (m *Model) NumberInputs() *NumberInputsModel {
	return &m.numberInputs
}

type ButtonsModel struct {
	disabled bool
}

func (b *ButtonsModel) Enabled() bool {
	return !b.disabled
}

func (b *ButtonsModel) SetEnabled(enabled bool) {
	b.disabled = !enabled
}

type TextsModel struct {
	text    string
	textSet bool

	horizontalAlign basicwidget.HorizontalAlign
	verticalAlign   basicwidget.VerticalAlign
	noWrap          bool
	bold            bool
	selectable      bool
	editable        bool
}

func (t *TextsModel) Text() string {
	if !t.textSet {
		return `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
隴西の李徴は博学才穎、天宝の末年、若くして名を虎榜に連ね、ついで江南尉に補せられたが、性、狷介、自ら恃むところ頗る厚く、賤吏に甘んずるを潔しとしなかった。`
	}
	return t.text
}

func (t *TextsModel) SetText(text string) {
	t.text = text
	t.textSet = true
}

func (t *TextsModel) HorizontalAlign() basicwidget.HorizontalAlign {
	return t.horizontalAlign
}

func (t *TextsModel) SetHorizontalAlign(align basicwidget.HorizontalAlign) {
	t.horizontalAlign = align
}

func (t *TextsModel) VerticalAlign() basicwidget.VerticalAlign {
	return t.verticalAlign
}

func (t *TextsModel) SetVerticalAlign(align basicwidget.VerticalAlign) {
	t.verticalAlign = align
}

func (t *TextsModel) AutoWrap() bool {
	return !t.noWrap
}

func (t *TextsModel) SetAutoWrap(autoWrap bool) {
	t.noWrap = !autoWrap
}

func (t *TextsModel) Bold() bool {
	return t.bold
}

func (t *TextsModel) SetBold(bold bool) {
	t.bold = bold
}

func (t *TextsModel) Selectable() bool {
	return t.selectable
}

func (t *TextsModel) SetSelectable(selectable bool) {
	t.selectable = selectable
	if !selectable {
		t.editable = false
	}
}

func (t *TextsModel) Editable() bool {
	return t.editable
}

func (t *TextsModel) SetEditable(editable bool) {
	t.editable = editable
	if editable {
		t.selectable = true
	}
}

type TextInputsModel struct {
	singleLineText     string
	singleLinetTextSet bool
	multilineText      string
	multilineTextSet   bool

	horizontalAlign basicwidget.HorizontalAlign
	verticalAlign   basicwidget.VerticalAlign
	noWrap          bool
	uneditable      bool
	disabled        bool
}

func (t *TextInputsModel) SingleLineText() string {
	if !t.singleLinetTextSet {
		return "Hello, Guigui!"
	}
	return t.singleLineText
}

func (t *TextInputsModel) SetSingleLineText(text string) {
	t.singleLineText = text
	t.singleLinetTextSet = true
}

func (t *TextInputsModel) MultilineText() string {
	if !t.multilineTextSet {
		return "Hello, Guigui!\nThis is a multiline text field."
	}
	return t.multilineText
}

func (t *TextInputsModel) SetMultilineText(text string) {
	t.multilineText = text
	t.multilineTextSet = true
}

func (t *TextInputsModel) HorizontalAlign() basicwidget.HorizontalAlign {
	return t.horizontalAlign
}

func (t *TextInputsModel) SetHorizontalAlign(align basicwidget.HorizontalAlign) {
	t.horizontalAlign = align
}

func (t *TextInputsModel) VerticalAlign() basicwidget.VerticalAlign {
	return t.verticalAlign
}

func (t *TextInputsModel) SetVerticalAlign(align basicwidget.VerticalAlign) {
	t.verticalAlign = align
}

func (t *TextInputsModel) AutoWrap() bool {
	return !t.noWrap
}

func (t *TextInputsModel) SetAutoWrap(autoWrap bool) {
	t.noWrap = !autoWrap
}

func (t *TextInputsModel) Editable() bool {
	return !t.uneditable
}

func (t *TextInputsModel) SetEditable(editable bool) {
	t.uneditable = !editable
	if editable {
		t.disabled = false
	}
}

func (t *TextInputsModel) Enabled() bool {
	return !t.disabled
}

func (t *TextInputsModel) SetEnabled(enabled bool) {
	t.disabled = !enabled
	if !enabled {
		t.uneditable = true
	}
}

type NumberInputsModel struct {
	numberInputValue1 big.Int
	numberInputValue2 uint64
	numberInputValue3 int

	uneditable bool
	disabled   bool
}

func (n *NumberInputsModel) Editable() bool {
	return !n.uneditable
}

func (n *NumberInputsModel) SetEditable(editable bool) {
	n.uneditable = !editable
	if editable {
		n.disabled = false
	}
}

func (n *NumberInputsModel) Enabled() bool {
	return !n.disabled
}

func (n *NumberInputsModel) SetEnabled(enabled bool) {
	n.disabled = !enabled
	if !enabled {
		n.uneditable = true
	}
}

func (n *NumberInputsModel) NumberInputValue1() *big.Int {
	var v big.Int
	v.Set(&n.numberInputValue1)
	return &v
}

func (n *NumberInputsModel) SetNumberInputValue1(value *big.Int) {
	n.numberInputValue1.Set(value)
}

func (n *NumberInputsModel) NumberInputValue2() uint64 {
	return n.numberInputValue2
}

func (n *NumberInputsModel) SetNumberInputValue2(value uint64) {
	n.numberInputValue2 = value
}

func (n *NumberInputsModel) NumberInputValue3() int {
	return n.numberInputValue3
}

func (n *NumberInputsModel) SetNumberInputValue3(value int) {
	n.numberInputValue3 = value
}
