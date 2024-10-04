package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/LoreviQ/ChessAnalysis/app/internal/game"
)

type Board struct {
	gui          *GUI
	squares      [][]layout.FlexChild
	squareSize   image.Point
	activeGameID int
	movesList    *widget.List
	gameState    *game.Game
	moves        []*MoveButton
}

type MoveButton struct {
	move          game.Move
	shortNotation string
	widget        *widget.Clickable
}

func newBoard(g *GUI) *Board {
	return &Board{
		gui: g,
		movesList: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}
}

func (b *Board) Layout(gtx layout.Context) layout.Dimensions {
	margins := layout.Inset{
		Top:    50,
		Bottom: 50,
		Left:   50,
		Right:  50,
	}

	// Get the moves for the active game
	moves, err := b.gui.db.GetMovesByID(b.activeGameID)
	if err != nil {
		panic(err)
	}

	// Get the board state for the active game
	b.gameState = game.NewGame()
	b.gameState.Moves(moves)
	b.moves = make([]*MoveButton, len(moves))
	moveHistory := b.gameState.MoveHistory
	shortNotation, err := game.ConvertMovesToShortAlgebraicNotation(moveHistory)
	if err != nil {
		panic(err)
	}
	notationList := make([]string, 0, len(shortNotation))
	for k := range shortNotation {
		notationList = append(notationList, k)
	}
	for i, move := range b.gameState.MoveHistory {
		b.moves = append(b.moves, &MoveButton{
			move:          move,
			shortNotation: notationList[i],
			widget:        &widget.Clickable{},
		})
	}

	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		squareHeight := (gtx.Constraints.Max.Y - 100) / 8
		squareWidth := (gtx.Constraints.Max.X * 15) / 191
		smallest := squareHeight
		if squareWidth < squareHeight {
			smallest = squareWidth
		}
		b.squareSize = image.Point{X: smallest, Y: smallest}

		// Board width = squareSize * 8
		// eval width = squareSize * 5/15
		// analysis width = squareSize * 4
		// spacer width = squareSize * 3/15
		// total width = squareSize * 191/15

		return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
			layout.Flexed(1, layout.Spacer{}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
					layout.Rigid(b.drawEvalBar),
					layout.Rigid(layout.Spacer{Width: unit.Dp(b.squareSize.Y / 5)}.Layout),
					layout.Rigid(b.drawBoard),
					layout.Rigid(layout.Spacer{Width: unit.Dp(b.squareSize.Y / 5)}.Layout),
					layout.Flexed(1, b.drawAnalysis),
				)
			}),
			layout.Flexed(1, layout.Spacer{}.Layout),
		)
	})
}

// Alternative drawBoard function using flex layout
func (b *Board) drawBoard(gtx layout.Context) layout.Dimensions {
	b.squares = b.drawSquares(8, 8)
	return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 50}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
				b.drawRows()...,
			)
		}),
		layout.Rigid(layout.Spacer{Height: 50}.Layout),
	)
}

func (b *Board) drawRows() []layout.FlexChild {
	rows := make([]layout.FlexChild, len(b.squares))
	for i, row := range b.squares {
		rows[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx, row...)
		})
	}
	return rows
}

func (b *Board) drawSquare(i, j int) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				// Draw the square
				square := image.Rectangle{
					Max: b.squareSize,
				}
				paint.FillShape(gtx.Ops, b.getSquareColour(i, j), clip.Rect(square).Op())
				return layout.Dimensions{Size: square.Max}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				// Draw the piece
				piece := b.gameState.Board.Squares[i][j]
				if piece == nil {
					return layout.Dimensions{}
				}
				img := b.gui.theme.chessBoardTheme.pieces[piece.GetImageName()]
				return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
							layout.Rigid(b.drawImage(*img)),
						)
					}),
				)
			}),
		)
	})
}

func (b *Board) drawSquares(maxRow, maxCol int) [][]layout.FlexChild {
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

func (b *Board) drawImage(image image.Image) func(layout.Context) layout.Dimensions {
	scale := float32(b.squareSize.X) / float32(image.Bounds().Dx())
	return func(gtx layout.Context) layout.Dimensions {
		return widget.Image{
			Src:   paint.NewImageOp(image),
			Scale: scale,
		}.Layout(gtx)
	}
}

func (b *Board) drawEvalBar(gtx layout.Context) layout.Dimensions {
	rect1 := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 50,
		},
		Max: image.Point{
			X: b.squareSize.Y / 3,
			Y: b.squareSize.Y*8 + 50,
		},
	}
	rect2 := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 50,
		},
		Max: image.Point{
			X: b.squareSize.Y / 3,
			Y: b.squareSize.Y*4 + 50,
		},
	}
	paint.FillShape(gtx.Ops, b.gui.theme.chessBoardTheme.player1, clip.Rect(rect1).Op())
	paint.FillShape(gtx.Ops, b.gui.theme.chessBoardTheme.player2, clip.Rect(rect2).Op())
	return layout.Dimensions{Size: rect1.Max}
}

func (b *Board) drawAnalysis(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Fill the background
			rect := image.Rectangle{
				Max: gtx.Constraints.Max,
			}
			paint.FillShape(gtx.Ops, b.gui.theme.chessBoardTheme.contrastBg, clip.Rect(rect).Op())
			return layout.Dimensions{Size: rect.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			/*
				margins := layout.Inset{
					Top:    20,
					Bottom: 20,
					Left:   20,
					Right:  20,
				}
			*/
			return layout.Dimensions{}
		}),
	)
}

func (b *Board) getSquareColour(i, j int) color.NRGBA {
	if (i+j)%2 == 0 {
		return b.gui.theme.chessBoardTheme.square1
	}
	return b.gui.theme.chessBoardTheme.square2
}
