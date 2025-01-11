package main

import (
	"hotwire/pkg/log"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var volumeOptions map[string]int = map[string]int{
	"Mute":        0,
	"Low":         1,
	"Medium Low":  2,
	"Medium":      3,
	"Medium High": 4,
	"High":        5,
}

var buttonOptions map[string]int = map[string]int{
	"Alexa":      0,
	"Hey Vector": 1,
}

func stringsFromMapStringInt(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func card(title string, subtitle string, obj fyne.CanvasObject) *widget.Card {
	return widget.NewCard(title, subtitle, obj)
}

func settingsMenu(content *fyne.Container, esn string) fyne.CanvasObject {
	// back button
	// backButton := widget.NewButton("Back", func() {
	// 	switchSection(content, )
	// })

	// volume
	volumeSelect := widget.NewSelect(stringsFromMapStringInt(volumeOptions), func(s string) {
		log.Debug("change volume of ", esn, " to ", s, ", aka. ", volumeOptions[s])
	})

	// button
	buttonSelect := widget.NewSelect(stringsFromMapStringInt(buttonOptions), func(s string) {
		log.Debug("change button of ", esn, " to ", s, ", aka. ", buttonOptions[s])
	})

	final := container.NewVScroll(

		container.NewHBox(
			card("Volume", "", volumeSelect),
			card("Button Action", "", buttonSelect),
		),
	)
	return final
}
