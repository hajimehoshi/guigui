// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package basicwidget

import (
	"slices"
)

type ider[ID comparable] interface {
	id() ID
}

type abstractList[ID comparable, Item ider[ID]] struct {
	items           []Item
	selectedIndices []int

	onItemSelected func(index int)
}

func (a *abstractList[ID, Item]) ResetEventHandlers() {
	a.onItemSelected = nil
}

func (a *abstractList[ID, Item]) SetOnItemSelected(f func(index int)) {
	a.onItemSelected = f
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

func (a *abstractList[ID, Item]) SelectItemByIndex(index int, forceFireEvents bool) bool {
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
		if a.onItemSelected != nil {
			a.onItemSelected(index)
		}
	}
	return true
}

func (a *abstractList[ID, Item]) SelectItemByID(id ID, forceFireEvents bool) bool {
	idx := slices.IndexFunc(a.items, func(item Item) bool {
		return item.id() == id
	})
	return a.SelectItemByIndex(idx, forceFireEvents)
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
