package gui

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type GUI struct {
	window *app.Window
	ops    *op.Ops
	theme  *material.Theme

	header  *header
	sidebar *sidebar
	board   *board
}

// Returns a GUI struct
func NewGUI(width, height int) *GUI {
	w := new(app.Window)
	w.Option(app.Size(unit.Dp(width), unit.Dp(height)))
	w.Option(app.Title("Chess Analysis"))
	ops := new(op.Ops)
	th := material.NewTheme()
	th.Palette.Bg = color.NRGBA{48, 46, 42, 255}

	// define components
	header := newHeader()
	sidebar := newSidebar()
	board := newBoard()

	return &GUI{
		window:  w,
		ops:     ops,
		theme:   th,
		header:  header,
		sidebar: sidebar,
		board:   board,
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
			gtx := app.NewContext(g.ops, e)
			g.Layout(gtx)
			e.Frame(gtx.Ops)
		// Exit app
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}
