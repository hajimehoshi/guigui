// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package main

import (
	"golang.org/x/text/language"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Settings struct {
	guigui.DefaultWidget

	form                      basicwidget.Form
	colorModeText             basicwidget.Text
	colorModeSegmentedControl basicwidget.SegmentedControl[guigui.ColorMode]
	localeText                basicwidget.Text
	localeDropdownList        basicwidget.DropdownList[language.Tag]
	scaleText                 basicwidget.Text
	scaleSegmentedControl     basicwidget.SegmentedControl[float64]
}

var hongKongChinese = language.MustParse("zh-HK")

func (s *Settings) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.colorModeText.SetValue("Color Mode")
	s.colorModeSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[guigui.ColorMode]{
		{
			Text: "Light",
			Tag:  guigui.ColorModeLight,
		},
		{
			Text: "Dark",
			Tag:  guigui.ColorModeDark,
		},
	})
	s.colorModeSegmentedControl.SetOnItemSelected(func(index int) {
		item, ok := s.colorModeSegmentedControl.ItemByIndex(index)
		if !ok {
			context.SetColorMode(guigui.ColorModeLight)
			return
		}
		context.SetColorMode(item.Tag)
	})
	s.colorModeSegmentedControl.SelectItemByTag(context.ColorMode())

	s.localeText.SetValue("Locale")
	s.localeDropdownList.SetItems([]basicwidget.DropdownListItem[language.Tag]{
		{
			Text: "(Default)",
			Tag:  language.Und,
		},
		{
			Text: "English",
			Tag:  language.English,
		},
		{
			Text: "Japanese",
			Tag:  language.Japanese,
		},
		{
			Text: "Korean",
			Tag:  language.Korean,
		},
		{
			Text: "Simplified Chinese",
			Tag:  language.SimplifiedChinese,
		},
		{
			Text: "Traditional Chinese",
			Tag:  language.TraditionalChinese,
		},
		{
			Text: "Hong Kong Chinese",
			Tag:  hongKongChinese,
		},
	})
	s.localeDropdownList.SetOnItemSelected(func(index int) {
		item, ok := s.localeDropdownList.ItemByIndex(index)
		if !ok {
			context.SetAppLocales(nil)
			return
		}
		if item.Tag == language.Und {
			context.SetAppLocales(nil)
			return
		}
		context.SetAppLocales([]language.Tag{item.Tag})
	})
	if !s.localeDropdownList.IsPopupOpen() {
		if locales := context.AppendAppLocales(nil); len(locales) > 0 {
			s.localeDropdownList.SelectItemByTag(locales[0])
		} else {
			s.localeDropdownList.SelectItemByTag(language.Und)
		}
	}

	s.scaleText.SetValue("Scale")
	s.scaleSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[float64]{
		{
			Text: "80%",
			Tag:  0.8,
		},
		{
			Text: "100%",
			Tag:  1,
		},
		{
			Text: "120%",
			Tag:  1.2,
		},
	})
	s.scaleSegmentedControl.SetOnItemSelected(func(index int) {
		item, ok := s.scaleSegmentedControl.ItemByIndex(index)
		if !ok {
			context.SetAppScale(1)
			return
		}
		context.SetAppScale(item.Tag)
	})
	s.scaleSegmentedControl.SelectItemByTag(context.AppScale())

	s.form.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &s.colorModeText,
			SecondaryWidget: &s.colorModeSegmentedControl,
		},
		{
			PrimaryWidget:   &s.localeText,
			SecondaryWidget: &s.localeDropdownList,
		},
		{
			PrimaryWidget:   &s.scaleText,
			SecondaryWidget: &s.scaleSegmentedControl,
		},
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row >= 1 {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(s.form.DefaultSize(context).Y)
			}),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&s.form, gl.CellBounds(0, 0))

	return nil
}
