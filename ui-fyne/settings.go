package main

import (
	"hotwire/pkg/log"
	"sort"
	"time"

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

	loadingBar := widget.NewProgressBar()
	loadBarCard := widget.NewCard("Loading...", "", loadingBar)
	switchSection(content, loadBarCard)

	for i := 0; i <= 100; i++ {
		loadingBar.SetValue(float64(i) / 100)
		time.Sleep(time.Millisecond * 20)
	}

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

	// go func() {
	// 	for {
	// 		volumeSelect.SetSelected("High")
	// 		time.Sleep(time.Second)
	// 		volumeSelect.SetSelected("Low")
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	tempSelect := widget.NewSelect(stringsFromMapStringInt(tempOptions), func(temp string) {})

	return container.New(
		NewFlowLayout(
			5,
		),
		disconnectButton,
		card("Volume", "", volumeSelect),
		card("Button Action", "", buttonSelect),
		card("Location", "ex. 'Brooklyn, New York, United States", locationFinal),
		card("Time Zone", "", timeZoneSelect),
		card("Temp Units", "", tempSelect),
	)
}

type FlowLayout struct {
	Spacing float32
}

func NewFlowLayout(spacing float32) *FlowLayout {
	return &FlowLayout{Spacing: spacing}
}

func (f *FlowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	x := float32(0)
	y := float32(0)
	rowHeight := float32(0)

	for _, o := range objects {
		if !o.Visible() {
			continue
		}

		min := o.MinSize()
		// If adding this object exceeds the container width, wrap to next line
		if x+min.Width > size.Width {
			x = 0
			y += rowHeight + f.Spacing
			rowHeight = 0
		}

		o.Move(fyne.NewPos(x, y))
		o.Resize(min)

		x += min.Width + f.Spacing

		if min.Height > rowHeight {
			rowHeight = min.Height
		}
	}
}

func (f *FlowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var maxWidth float32
	var totalHeight float32

	visibleCount := 0
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		visibleCount++
		min := o.MinSize()
		if min.Width > maxWidth {
			maxWidth = min.Width
		}
		totalHeight += min.Height
	}
	if visibleCount > 1 {
		totalHeight += f.Spacing * float32(visibleCount-1)
	}

	return fyne.NewSize(maxWidth, totalHeight)
}
