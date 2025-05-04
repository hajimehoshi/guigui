// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors

package main

import (
	"slices"
	"strings"
)

type Model struct {
	tasks []Task
}

func (m *Model) CanAddTask(text string) bool {
	return strings.TrimSpace(text) != ""
}

func (m *Model) TryAddTask(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	m.tasks = slices.Insert(m.tasks, 0, NewTask(text))
	return true
}

func (m *Model) DeleteTaskByID(id int) {
	m.tasks = slices.DeleteFunc(m.tasks, func(t Task) bool {
		return t.ID == id
	})
}

func (m *Model) TaskCount() int {
	return len(m.tasks)
}

func (m *Model) TaskByIndex(i int) Task {
	return m.tasks[i]
}

var theCurrentID int

func nextTaskID() int {
	theCurrentID++
	return theCurrentID
}

type Task struct {
	ID   int
	Text string
}

func NewTask(text string) Task {
	return Task{
		ID:   nextTaskID(),
		Text: text,
	}
}
