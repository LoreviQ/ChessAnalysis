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
	squares [8][8]*image.Rectangle
}

// temp fill for development will be incorporated into theme
var (
	square1_colour = color.NRGBA{114, 148, 82, 255}
	square2_colour = color.NRGBA{234, 236, 206, 255}
)

func newBoard() *board {
	return &board{}
}

func (b *board) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// temp fill for development
	rect := image.Rectangle{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Constraints.Max.Y,
		},
	}
	paint.FillShape(gtx.Ops, color.NRGBA{0, 0, 255, 255}, clip.Rect(rect).Op())

	return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(b.drawBoard),
				layout.Flexed(1, layout.Spacer{}.Layout),
			)
		}),
		layout.Flexed(1, layout.Spacer{}.Layout),
	)
}

func (b *board) drawBoard(gtx layout.Context) layout.Dimensions {
	// Calculate square size
	largest := gtx.Constraints.Max.X
	if (gtx.Constraints.Max.Y - 200) < largest {
		largest = gtx.Constraints.Max.Y - 200
	}
	squareSize := (largest) / 8
	boardSize := image.Point{X: squareSize * 8, Y: squareSize * 8}

	// Layout board squares
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			square := image.Rect(
				i*squareSize,
				j*squareSize,
				(i+1)*squareSize,
				(j+1)*squareSize,
			)
			b.squares[i][j] = &square
			if (i+j)%2 != 0 {
				paint.FillShape(gtx.Ops, square1_colour, clip.Rect(square).Op())
			} else {
				paint.FillShape(gtx.Ops, square2_colour, clip.Rect(square).Op())
			}
		}
	}
	return layout.Dimensions{Size: boardSize}
}
