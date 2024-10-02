package gui

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func (g *GUI) landingPage(e app.FrameEvent) {
	// This graphics context is used for managing the rendering state.
	gtx := app.NewContext(g.ops, e)

	// Define an large label with an appropriate text:
	title := material.H1(g.theme, "Hello, Gio")

	// Change the color of the label.
	maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	title.Color = maroon

	// Change the position of the label.
	title.Alignment = text.Middle

	// Draw the label to the graphics context.
	title.Layout(gtx)

	// Pass the drawing operations to the GPU.
	e.Frame(gtx.Ops)
}
