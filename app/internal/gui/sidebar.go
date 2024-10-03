package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type sidebar struct {
	g *GUI
}

func newSidebar(g *GUI) *sidebar {
	return &sidebar{
		g: g,
	}
}

func (s *sidebar) Layout(gtx layout.Context) layout.Dimensions {
	// Define the fixed size for the header
	sidebarSize := image.Point{X: 500, Y: gtx.Constraints.Max.Y} // Fixed height of 50 pixels

	// Adjust the constraints to enforce the fixed size
	gtx.Constraints.Min = sidebarSize
	gtx.Constraints.Max = sidebarSize

	// temp fill for development
	rect := image.Rectangle{
		Max: sidebarSize,
	}
	paint.FillShape(gtx.Ops, color.NRGBA{0, 255, 0, 255}, clip.Rect(rect).Op())

	return layout.Dimensions{Size: sidebarSize}
}
