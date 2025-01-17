package main

import (
	"fmt"
	"image/color"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func switchSection(content *fyne.Container, newObj fyne.CanvasObject) {
	// Clear out any existing children.
	content.Objects = nil

	// Add the new child, so that Fyne can lay it out.
	content.Add(newObj)
	content.Refresh()
}

func formatDuration(sec int) string {
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	var parts []string
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%d minutes", m))
	}
	if s > 0 {
		parts = append(parts, fmt.Sprintf("%d seconds", s))
	}
	return strings.Join(parts, " ")
}

type ManagementBot struct {
	Name                  string
	ESN                   string
	Status                string
	SecondsSinceLastComms int
}

var bots = []ManagementBot{
	{"Starbutt", "00601b50", "Starbutt is looking around", 1},
	{"DVT3", "00200010", "DVT3 wants a fistbump", 3},
	{"Vector", "0060059b", "Vector is inactive", 2},
}

func containerFromRobots(content *fyne.Container, robots []ManagementBot) *fyne.Container {
	var f int
	mainContainer := container.NewVBox()
	hbox := container.NewHBox()
	for i, r := range robots {
		c := widget.NewCard(r.Name, r.Status,
			container.NewVBox(
				widget.NewRichTextWithText("ESN: "+r.ESN),
				widget.NewRichTextWithText("Last conncheck: "+formatDuration(r.SecondsSinceLastComms)),
				widget.NewButton("Manage", func() {
					switchSection(content, settingsMenu(content, r.ESN))
				}),
			),
		)
		if f == 1 {
			hbox.Add(c)
			mainContainer.Add(hbox)
			hbox = container.NewHBox()
		} else if i == len(robots)-1 {
			hbox.Add(c)
			mainContainer.Add(hbox)
		} else {
			hbox.Add(c)
			f = 0
		}
		f++
	}
	return mainContainer
}

func makeBotMan(content *fyne.Container) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle("Bot Management", 0, fyne.TextStyle{Bold: true}),
		containerFromRobots(content, bots),
	)
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&darkTheme{})
	w := a.NewWindow("Hotwire")

	title := widget.NewLabelWithStyle("Hotwire", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	dropLabel := widget.NewLabel("Theme:")
	drop := widget.NewSelect([]string{"Dark", "Purple", "Green"}, func(val string) {
		switch val {
		case "Dark":
			a.Settings().SetTheme(&darkTheme{})
		case "Purple":
			a.Settings().SetTheme(&purpleTheme{})
		case "Green":
			a.Settings().SetTheme(&greenTheme{})
		}
	})

	topSep := canvas.NewRectangle(color.NRGBA{0xAA, 0xAA, 0xAA, 0xFF})
	topSep.SetMinSize(fyne.NewSize(0, 3))
	topBar := container.NewHBox(title, layout.NewSpacer(), dropLabel, drop)
	topBarArea := container.NewVBox(topBar, topSep)

	gen := container.NewVBox(
		widget.NewLabelWithStyle("Dashboard", 0, fyne.TextStyle{Bold: true}),
		widget.NewRichText(&widget.TextSegment{
			Text:  "Rich Text Test",
			Style: widget.RichTextStyleCodeBlock,
		}),
		widget.NewButton("Restart", func() { println("restart") }),
	)

	content := container.NewStack(gen)

	botMan := makeBotMan(content)

	net := container.NewVBox(
		widget.NewLabelWithStyle("Configuration", 0, fyne.TextStyle{Bold: true}),
		widget.NewEntry(),
		widget.NewButton("Save Port", func() {
			http.Post("http://localhost:8000", "application/json", nil)
		}),
	)

	var dashBtn, botBtn, configBtn *widget.Button
	var buttons []*widget.Button

	refreshHighlight := func(active *widget.Button) {
		for _, b := range buttons {
			if b == active {
				b.Importance = widget.HighImportance
			} else {
				b.Importance = widget.MediumImportance
			}
			b.Refresh()
		}
	}

	dashBtn = widget.NewButton("Dashboard", func() {
		switchSection(content, gen)
		refreshHighlight(dashBtn)
	})
	botBtn = widget.NewButton("Bot Management", func() {
		switchSection(content, botMan)
		refreshHighlight(botBtn)
	})
	configBtn = widget.NewButton("Configuration", func() {
		switchSection(content, net)
		refreshHighlight(configBtn)
	})

	buttons = []*widget.Button{dashBtn, botBtn, configBtn}
	dashBtn.Importance = widget.HighImportance

	leftSep := canvas.NewRectangle(color.NRGBA{0xAA, 0xAA, 0xAA, 0xFF})
	leftSep.SetMinSize(fyne.NewSize(3, 0))
	//menu := container.NewVBox(dashBtn, botBtn, configBtn)
	//mainArea := container.NewHBox(menu, leftSep, content)
	ui := container.NewBorder(
		topBarArea, // top
		nil,        // bottom
		container.NewVBox(dashBtn, botBtn, configBtn), // left
		nil,     // right
		content, // center
	)

	//ui := container.NewVBox(topBarArea, mainArea)

	w.SetContent(ui)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
