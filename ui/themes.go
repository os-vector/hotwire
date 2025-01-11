package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const ScalingFactor float32 = 1

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
	return theme.DefaultTheme().Size(n) * ScalingFactor
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
	return theme.DefaultTheme().Size(n) * ScalingFactor
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
	return theme.DefaultTheme().Size(n) * ScalingFactor
}
