package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type board struct {
	g          *GUI
	squares    [8][8]*image.Rectangle
	squareSize int
	boardSize  image.Point
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
		Left:   50,
		Right:  50,
	}

	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		largest := gtx.Constraints.Max.X
		if (gtx.Constraints.Max.Y) < largest {
			largest = gtx.Constraints.Max.Y
		}
		b.squareSize = (largest) / 8
		b.boardSize = image.Point{X: b.squareSize * 8, Y: b.squareSize * 8}
		return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
			layout.Flexed(1, layout.Spacer{}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if largest == gtx.Constraints.Max.Y {
							return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
								layout.Rigid(b.drawEvalBar),
								layout.Rigid(layout.Spacer{Width: unit.Dp(b.boardSize.Y / 40)}.Layout),
							)
						}
						return layout.Dimensions{}
					}),
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

	// Layout board squares
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			square := image.Rect(
				i*b.squareSize,
				j*b.squareSize,
				(i+1)*b.squareSize,
				(j+1)*b.squareSize,
			)
			b.squares[i][j] = &square
			if (i+j)%2 == 0 {
				paint.FillShape(gtx.Ops, b.g.theme.chessBoardTheme.square1Colour, clip.Rect(square).Op())
			} else {
				paint.FillShape(gtx.Ops, b.g.theme.chessBoardTheme.square2Colour, clip.Rect(square).Op())
			}
		}
	}
	return layout.Dimensions{Size: b.boardSize}
}

func (b *board) drawEvalBar(gtx layout.Context) layout.Dimensions {
	rect1 := image.Rectangle{
		Max: image.Point{
			X: b.boardSize.Y / 20,
			Y: b.boardSize.Y,
		},
	}
	rect2 := image.Rectangle{
		Max: image.Point{
			X: b.boardSize.Y / 20,
			Y: b.boardSize.Y / 2,
		},
	}
	paint.FillShape(gtx.Ops, b.g.theme.chessBoardTheme.player1Colour, clip.Rect(rect1).Op())
	paint.FillShape(gtx.Ops, b.g.theme.chessBoardTheme.player2Colour, clip.Rect(rect2).Op())
	return layout.Dimensions{Size: rect1.Max}
}
