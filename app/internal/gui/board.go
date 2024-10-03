package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type board struct {
	gui          *GUI
	squares      [][]layout.FlexChild
	squareSize   image.Point
	activeGameID int
}

func newBoard(g *GUI) *board {
	return &board{
		gui: g,
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
		smallest := gtx.Constraints.Max.X
		if gtx.Constraints.Max.Y < smallest {
			smallest = gtx.Constraints.Max.Y
		}
		b.squareSize = image.Point{X: smallest / 8, Y: smallest / 8}
		return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
			layout.Flexed(1, layout.Spacer{}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if smallest == gtx.Constraints.Max.Y {
							return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
								layout.Rigid(b.drawEvalBar),
								layout.Rigid(layout.Spacer{Width: unit.Dp(b.squareSize.Y / 5)}.Layout),
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

// Alternative drawBoard function using flex layout
func (b *board) drawBoard(gtx layout.Context) layout.Dimensions {
	b.squares = b.drawSquares(8, 8)
	return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
		b.drawRows()...,
	)
}

func (b *board) drawRows() []layout.FlexChild {
	rows := make([]layout.FlexChild, len(b.squares))
	for i, row := range b.squares {
		rows[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx, row...)
		})
	}
	return rows
}

func (b *board) drawSquare(i, j int) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		square := image.Rectangle{
			Max: b.squareSize,
		}
		paint.FillShape(gtx.Ops, b.getSquareColour(i, j), clip.Rect(square).Op())
		return layout.Dimensions{Size: square.Max}
	})
}

func (b *board) drawSquares(maxRow, maxCol int) [][]layout.FlexChild {
	children := make([][]layout.FlexChild, maxRow)
	for i := 0; i < maxRow; i++ {
		children[i] = make([]layout.FlexChild, maxCol)
		for j := 0; j < maxCol; j++ {
			child := b.drawSquare(i, j)
			children[i][j] = child
		}
	}
	return children
}

func drawImage()

func (b *board) drawEvalBar(gtx layout.Context) layout.Dimensions {
	rect1 := image.Rectangle{
		Max: image.Point{
			X: b.squareSize.Y / 3,
			Y: b.squareSize.Y * 8,
		},
	}
	rect2 := image.Rectangle{
		Max: image.Point{
			X: b.squareSize.Y / 3,
			Y: b.squareSize.Y * 4,
		},
	}
	paint.FillShape(gtx.Ops, b.gui.theme.chessBoardTheme.player1Colour, clip.Rect(rect1).Op())
	paint.FillShape(gtx.Ops, b.gui.theme.chessBoardTheme.player2Colour, clip.Rect(rect2).Op())
	return layout.Dimensions{Size: rect1.Max}
}

func (b *board) getSquareColour(i, j int) color.NRGBA {
	if (i+j)%2 == 0 {
		return b.gui.theme.chessBoardTheme.square1Colour
	}
	return b.gui.theme.chessBoardTheme.square2Colour
}
