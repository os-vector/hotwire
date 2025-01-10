package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type darkTheme struct{}

func (d *darkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{0x22, 0x22, 0x22, 0xFF}
	case theme.ColorNameButton:
		return color.NRGBA{0x33, 0x33, 0x33, 0xFF}
	case theme.ColorNameForeground:
		return color.White
	}
	return theme.DefaultTheme().Color(n, v)
}
func (d *darkTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (d *darkTheme) Font(s fyne.TextStyle) fyne.Resource     { return theme.DefaultTheme().Font(s) }
func (d *darkTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameInputRadius {
		return 0
	}
	return theme.DefaultTheme().Size(n) * 1.3
}

type purpleTheme struct{}

func (p *purpleTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{0x2E, 0x00, 0x3E, 0xFF}
	case theme.ColorNameButton:
		return color.NRGBA{0x4A, 0x14, 0x8C, 0xFF}
	case theme.ColorNameForeground:
		return color.White
	}
	return theme.DefaultTheme().Color(n, v)
}
func (p *purpleTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (p *purpleTheme) Font(s fyne.TextStyle) fyne.Resource     { return theme.DefaultTheme().Font(s) }
func (d *purpleTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameInputRadius {
		return 0
	}
	return theme.DefaultTheme().Size(n) * 1.3
}

type greenTheme struct{}

func (g *greenTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{0x00, 0x20, 0x00, 0xFF}
	case theme.ColorNameButton:
		return color.NRGBA{0x00, 0x55, 0x00, 0xFF}
	case theme.ColorNameForeground:
		return color.White
	}
	return theme.DefaultTheme().Color(n, v)
}
func (g *greenTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (g *greenTheme) Font(s fyne.TextStyle) fyne.Resource     { return theme.DefaultTheme().Font(s) }
func (d *greenTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameInputRadius {
		return 0
	}
	return theme.DefaultTheme().Size(n) * 1.3
}

func switchSection(content *fyne.Container, newObj fyne.CanvasObject) {
	oldObj := content.Objects[0]
	w := content.Size().Width
	content.Add(newObj)
	newObj.Move(fyne.NewPos(w, 0))

	canvas.NewPositionAnimation(oldObj.Position(), fyne.NewPos(w, 0), time.Millisecond*200, func(pos fyne.Position) {
		oldObj.Move(pos)
		oldObj.Refresh()
	}).Start()

	canvas.NewPositionAnimation(newObj.Position(), fyne.NewPos(0, 0), time.Millisecond*200, func(pos fyne.Position) {
		newObj.Move(pos)
		newObj.Refresh()
	}).Start()

	go func() {
		time.Sleep(time.Millisecond * 200)
		content.Remove(oldObj)
	}()
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

func containerFromRobots(robots []ManagementBot) *fyne.Container {
	var f int
	mainContainer := container.NewVBox()
	hbox := container.NewHBox()
	for i, r := range robots {
		c := widget.NewCard(r.Name, r.Status,
			container.NewVBox(
				widget.NewRichTextWithText("ESN: "+r.ESN),
				widget.NewRichTextWithText("Last conncheck: "+formatDuration(r.SecondsSinceLastComms)),
				widget.NewButton("Manage", func() {}),
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
	adv := container.NewVBox(
		widget.NewLabelWithStyle("Bot Management", 0, fyne.TextStyle{Bold: true}),
		containerFromRobots(bots),
	)
	net := container.NewVBox(
		widget.NewLabelWithStyle("Configuration", 0, fyne.TextStyle{Bold: true}),
		widget.NewEntry(),
		widget.NewButton("Save Port", func() { println("network saved") }),
	)

	content := container.NewStack(gen)

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
		switchSection(content, adv)
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
	menu := container.NewVBox(dashBtn, botBtn, configBtn)
	mainArea := container.NewHBox(menu, leftSep, content)
	ui := container.NewVBox(topBarArea, mainArea)

	w.SetContent(ui)
	w.Resize(fyne.NewSize(800, 800))
	w.ShowAndRun()
}
