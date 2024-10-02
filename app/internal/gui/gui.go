package gui

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type GUI struct {
	window  *app.Window
	ops     *op.Ops
	theme   *material.Theme
	colours *colours
}

type colours struct {
	bg        color.NRGBA
	fg        color.NRGBA
	text      color.NRGBA
	highlight color.NRGBA
}

// Returns a GUI struct
func NewGUI() *GUI {
	return &GUI{
		window: new(app.Window),
		ops:    new(op.Ops),
		theme:  material.NewTheme(),
		colours: &colours{
			bg:        color.NRGBA{41, 40, 45, 255},
			fg:        color.NRGBA{53, 54, 62, 255},
			text:      color.NRGBA{255, 255, 255, 255},
			highlight: color.NRGBA{63, 81, 182, 255},
		},
	}
}

// CreateGUI creates the GUI
func (g *GUI) CreateGUI() {
	go g.draw()
	app.Main()
}

// Main event loop
func (g *GUI) draw() error {
	for {
		switch e := g.window.Event().(type) {

		// Re-render app
		case app.FrameEvent:
			g.landingPage(e)

		// Exit app
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}
