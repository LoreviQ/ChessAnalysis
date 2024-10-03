package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type header struct {
}

func newHeader() *header {
	return &header{}
}

func (h *header) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Define the fixed size for the header
	headerSize := image.Point{X: gtx.Constraints.Max.X, Y: 50} // Fixed height of 50 pixels

	// Adjust the constraints to enforce the fixed size
	gtx.Constraints.Min = headerSize
	gtx.Constraints.Max = headerSize

	// temp fill for development
	rect := image.Rectangle{
		Max: headerSize,
	}
	paint.FillShape(gtx.Ops, color.NRGBA{255, 0, 0, 255}, clip.Rect(rect).Op())

	return layout.Dimensions{Size: headerSize}
}
