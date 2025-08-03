// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"slices"

	"github.com/hajimehoshi/guigui"
)

type ider[ID comparable] interface {
	id() ID
}

const (
	abstractListEventItemSelected = "itemSelected"
)

type abstractList[ID comparable, Item ider[ID]] struct {
	items           []Item
	selectedIndices []int
}

func (a *abstractList[ID, Item]) SetOnItemSelected(widget guigui.Widget, f func(context *guigui.Context, index int)) {
	guigui.RegisterEventHandler(widget, abstractListEventItemSelected, f)
}

func (a *abstractList[ID, Item]) SetItems(items []Item) {
	a.items = adjustSliceSize(items, len(items))
	copy(a.items, items)
}

func (a *abstractList[ID, Item]) ItemCount() int {
	return len(a.items)
}

func (a *abstractList[ID, Item]) ItemByIndex(index int) (Item, bool) {
	if index < 0 || index >= len(a.items) {
		var item Item
		return item, false
	}
	return a.items[index], true
}

func (a *abstractList[ID, Item]) SelectItemByIndex(widget guigui.Widget, index int, forceFireEvents bool) bool {
	if index < 0 || index >= len(a.items) {
		if len(a.selectedIndices) == 0 {
			return false
		}
		a.selectedIndices = a.selectedIndices[:0]
		return true
	}

	if len(a.selectedIndices) == 1 && a.selectedIndices[0] == index && !forceFireEvents {
		return false
	}

	selected := slices.Contains(a.selectedIndices, index)
	a.selectedIndices = adjustSliceSize(a.selectedIndices, 1)
	a.selectedIndices[0] = index
	if !selected || forceFireEvents {
		guigui.InvokeEventHandler(widget, abstractListEventItemSelected, index)
	}
	return true
}

func (a *abstractList[ID, Item]) SelectItemByID(widget guigui.Widget, id ID, forceFireEvents bool) bool {
	idx := slices.IndexFunc(a.items, func(item Item) bool {
		return item.id() == id
	})
	return a.SelectItemByIndex(widget, idx, forceFireEvents)
}

func (a *abstractList[ID, Item]) SelectedItem() (Item, bool) {
	if len(a.selectedIndices) == 0 {
		var item Item
		return item, false
	}
	return a.items[a.selectedIndices[0]], true
}

func (a *abstractList[ID, Item]) SelectedItemIndex() int {
	if len(a.selectedIndices) == 0 {
		return -1
	}
	return a.selectedIndices[0]
}
