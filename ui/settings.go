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

var tempOptions map[string]int = map[string]int{
	"Celsius":    0,
	"Fahrenheit": 1,
}

var timeZones []string = []string{
	"Pacific/Honolulu",
	"America/Juneau",
	"America/Los_Angeles",
	"America/Phoenix",
	"America/Denver",
	"America/Lima",
	"America/Chicago",
	"America/Bogota",
	"America/New_York",
	"America/Argentina/Buenos_Aires",
	"America/Santiago",
	"America/Sao_Paulo",
	"America/Halifax",
	"America/St_Johns",
	"GMT",
	"Europe/Lisbon",
	"Europe/London",
	"Europe/Paris",
	"Europe/Athens",
	"Europe/Istanbul",
	"Europe/Moscow",
	"Africa/Lagos",
	"Africa/Harare",
	"Africa/Addis_Ababa",
	"Asia/Dubai",
	"Asia/Tehran",
	"Asia/Karachi",
	"Asia/Kolkata",
	"Asia/Dhaka",
	"Asia/Bangkok",
	"Asia/Jakarta",
	"Asia/Hong_Kong",
	"Asia/Singapore",
	"Asia/Manila",
	"Asia/Seoul",
	"Asia/Tokyo",
	"Australia/Perth",
	"Australia/Darwin",
	"Australia/Adelaide",
	"Australia/Brisbane",
	"Australia/Sydney",
	"Australia/Auckland",
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
	// disconnect button

	go func() {
		// watchdog, ensure we have connection to bot. maybe this could be event stream?
	}()

	disconnectButton := widget.NewButton("Disconnect", func() {
		switchSection(content, makeBotMan(content))
	})

	// volume
	volumeSelect := widget.NewSelect(stringsFromMapStringInt(volumeOptions), func(s string) {
		log.Debug("change volume of ", esn, " to ", s, ", aka. ", volumeOptions[s])
	})

	// button
	buttonSelect := widget.NewSelect(stringsFromMapStringInt(buttonOptions), func(s string) {
		log.Debug("change button of ", esn, " to ", s, ", aka. ", buttonOptions[s])
	})

	// location
	// TODO: hook up to some auto-completions API
	locationEntry := widget.NewEntry()
	locationEntryButton := widget.NewButton("Set", func() {})
	locationFinal := container.NewVBox(locationEntry, locationEntryButton)

	// time zone
	timeZoneSelect := widget.NewSelect(timeZones, func(zone string) {
		//TODO
	})

	tempSelect := widget.NewSelect(stringsFromMapStringInt(tempOptions), func(temp string) {})

	return container.NewVBox(disconnectButton,
		container.NewHBox(
			card("Volume", "", volumeSelect),
			card("Button Action", "", buttonSelect),
		),
		container.NewHBox(
			card("Location", "ex. 'Brooklyn, New York, United States", locationFinal),
		),
		container.NewHBox(
			card("Time Zone", "", timeZoneSelect),
			card("Temp Units", "", tempSelect),
		),
	)
}
