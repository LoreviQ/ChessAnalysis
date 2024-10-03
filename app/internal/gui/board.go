package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type board struct {
}

func newBoard() *board {
	return &board{}
}

func (h *board) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// temp fill for development
	rect := image.Rectangle{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Constraints.Max.Y,
		},
	}
	paint.FillShape(gtx.Ops, color.NRGBA{0, 0, 255, 255}, clip.Rect(rect).Op())

	return layout.Dimensions{}
}
