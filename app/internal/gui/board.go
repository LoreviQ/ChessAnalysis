package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type board struct {
	squares [8][8]*image.Rectangle
	g       *GUI
}

func newBoard(g *GUI) *board {
	return &board{
		g: g,
	}
}

func (b *board) Layout(gtx layout.Context) layout.Dimensions {
	margins := layout.Inset{
		Top:    100,
		Bottom: 100,
	}

	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
			layout.Flexed(1, layout.Spacer{}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
					layout.Flexed(1, layout.Spacer{}.Layout),
					layout.Rigid(b.drawBoard),
					layout.Flexed(1, layout.Spacer{}.Layout),
				)
			}),
			layout.Flexed(1, layout.Spacer{}.Layout),
		)
	})
}

func (b *board) drawBoard(gtx layout.Context) layout.Dimensions {
	// Calculate square size
	largest := gtx.Constraints.Max.X
	if (gtx.Constraints.Max.Y) < largest {
		largest = gtx.Constraints.Max.Y
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
				paint.FillShape(gtx.Ops, b.g.theme.chessBoardTheme.square1Colour, clip.Rect(square).Op())
			} else {
				paint.FillShape(gtx.Ops, b.g.theme.chessBoardTheme.square2Colour, clip.Rect(square).Op())
			}
		}
	}
	return layout.Dimensions{Size: boardSize}
}
