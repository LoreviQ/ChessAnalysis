package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type GUI struct {
	window *app.Window
	ops    *op.Ops
	theme  *material.Theme
}

// Returns a GUI struct
func NewGUI() *GUI {
	return &GUI{
		window: new(app.Window),
		ops:    new(op.Ops),
		theme:  material.NewTheme(),
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
