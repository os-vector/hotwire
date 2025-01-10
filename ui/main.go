package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type darkPurpleSharpTheme struct{}

func (d *darkPurpleSharpTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x2E, G: 0x00, B: 0x3E, A: 0xFF}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x4A, G: 0x14, B: 0x8C, A: 0xFF}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x7C, G: 0x43, B: 0xBD, A: 0xFF}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xA0, G: 0x59, B: 0xFF, A: 0xFF}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xAA}
	}
	return theme.DefaultTheme().Color(n, v)
}

func (d *darkPurpleSharpTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}
func (d *darkPurpleSharpTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (d *darkPurpleSharpTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameInputRadius {
		return 0
	}
	return theme.DefaultTheme().Size(n) * 1.1
}

func switchSection(content *fyne.Container, newObj fyne.CanvasObject) {
	oldObj := content.Objects[0]
	w := content.Size().Width
	newObj.Move(fyne.NewPos(-w, 0))
	content.Add(newObj)

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

func main() {
	a := app.New()
	a.Settings().SetTheme(&darkPurpleSharpTheme{})
	w := a.NewWindow("Hotwire")

	// top bar stuff
	title := widget.NewLabelWithStyle("Hotwire", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	dropdownLabel := widget.NewLabel("Advanced:")
	dropdown := widget.NewSelect([]string{"test1", "test2", "test3"}, func(val string) {})
	topBar := container.NewHBox(title, layout.NewSpacer(), dropdownLabel, dropdown)

	// thick top separator
	topSep := canvas.NewRectangle(color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF})
	topSep.SetMinSize(fyne.NewSize(0, 3))
	topBarArea := container.NewVBox(topBar, topSep)

	gen := container.NewVBox(
		widget.NewLabel("Dashboard"),
		widget.NewRichText(&widget.TextSegment{
			Style: widget.RichTextStyleCodeBlock,
			Text:  "Rich Text Test",
		}),
		widget.NewButton("Restart", func() { println("restart") }),
	)

	adv := container.NewVBox(
		widget.NewLabel("Bot Management"),
		widget.NewSelect([]string{"00e20145", "00601b50"}, func(string) {}),
		widget.NewCard("00e20145", "Status: doin his thing", widget.NewButton("connect", func() {})),
	)

	net := container.NewVBox(
		widget.NewLabel("Configuration"),
		widget.NewEntry(),
		widget.NewButton("Save Port", func() { println("network saved") }),
	)

	content := container.NewStack(gen)

	menu := container.NewVBox(
		widget.NewButton("Dashboard", func() { switchSection(content, gen) }),
		widget.NewButton("Bot Management", func() { switchSection(content, adv) }),
		widget.NewButton("Configuration", func() { switchSection(content, net) }),
	)

	// thick vertical separator
	leftSep := canvas.NewRectangle(color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF})
	leftSep.SetMinSize(fyne.NewSize(3, 0))
	mainArea := container.NewHBox(menu, leftSep, content)

	// overall layout
	ui := container.NewVBox(topBarArea, mainArea)

	w.SetContent(ui)
	w.Resize(fyne.NewSize(700, 400))
	w.ShowAndRun()
}
