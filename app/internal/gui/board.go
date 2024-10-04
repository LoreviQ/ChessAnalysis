package gui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/LoreviQ/ChessAnalysis/app/internal/game"
)

type Board struct {
	gui          *GUI
	squares      [][]layout.FlexChild
	squareSize   image.Point
	activeGameID int
	movesList    *widget.List
	gameState    *game.Game
	stateNum     int
	moves        []*MoveButton
	flipped      bool
}

type MoveButton struct {
	move      *game.Move
	notation  string
	widget    *widget.Clickable
	gameState *game.Game
}

func newBoard(g *GUI, activeGameID int) *Board {
	// Get the moves for the active game
	errBoard := &Board{
		gui:          g,
		activeGameID: activeGameID,
		movesList:    &widget.List{},
		gameState:    game.NewGame(),
		stateNum:     0,
		moves:        nil,
		flipped:      false,
	}
	moveStrs, err := g.db.GetMovesByID(activeGameID)
	if err != nil {
		return errBoard
	}
	// Get the board state for the active game
	gameState := game.NewGame()
	// check if the move string is valid
	err = gameState.Moves(moveStrs)
	if err != nil {
		return errBoard
	}
	// turn the move strings into moves
	gameState.NewGame()
	moves := make([]*MoveButton, len(moveStrs)+1)
	moves[0] = &MoveButton{
		move:      nil,
		notation:  "",
		widget:    &widget.Clickable{},
		gameState: gameState.Clone(),
	}
	for i, moveStr := range moveStrs {
		move, err := gameState.Move(moveStr)
		if err != nil {
			break
		}
		moves[i+1] = &MoveButton{
			move:      &move,
			notation:  moveStr,
			widget:    &widget.Clickable{},
			gameState: gameState.Clone(),
		}
	}
	return &Board{
		gui:          g,
		activeGameID: activeGameID,
		movesList: &widget.List{
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: true,
			},
		},
		gameState: gameState,
		stateNum:  len(moves) - 1,
		moves:     moves,
		flipped:   false,
	}
}

func (b *Board) Layout(gtx layout.Context) layout.Dimensions {
	b.updateState(gtx)
	margins := layout.Inset{
		Top:    50,
		Bottom: 50,
		Left:   50,
		Right:  50,
	}

	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		squareHeight := (gtx.Constraints.Max.Y - 100) / 8
		squareWidth := (gtx.Constraints.Max.X * 15) / 191
		smallest := squareHeight
		if squareWidth < squareHeight {
			smallest = squareWidth
		}
		b.squareSize = image.Point{X: smallest, Y: smallest}

		// board width = squareSize * 8
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

// Draw the chess board
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

// Draw the rows of the chess board
func (b *Board) drawRows() []layout.FlexChild {
	rows := make([]layout.FlexChild, len(b.squares))
	for i, row := range b.squares {
		rows[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx, row...)
		})
	}
	return rows
}

// Draw a square on the chess board
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
				row := 7 - i
				if b.flipped {
					row = i
				}
				// Draw the piece
				piece := b.gameState.Board.Squares[row][j]
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

// Draw all the squares on the chess board
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

// Draw an image
func (b *Board) drawImage(image image.Image) func(layout.Context) layout.Dimensions {
	scale := float32(b.squareSize.X) / float32(image.Bounds().Dx())
	return func(gtx layout.Context) layout.Dimensions {
		return widget.Image{
			Src:   paint.NewImageOp(image),
			Scale: scale,
		}.Layout(gtx)
	}
}

// Draw the evaluation bar
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

// Draw the analysis pane
func (b *Board) drawAnalysis(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Fill the background
			rect := image.Rectangle{
				Max: gtx.Constraints.Max,
			}
			paint.FillShape(gtx.Ops, b.gui.theme.contrastBg, clip.Rect(rect).Op())
			return layout.Dimensions{Size: rect.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			margins := layout.Inset{
				Top:    20,
				Bottom: 20,
				Left:   20,
				Right:  20,
			}
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return b.movesList.Layout(gtx, (len(b.moves))/2, func(gtx layout.Context, i int) layout.Dimensions {
					return b.moveListElement(gtx, i)
				})
			})
		}),
	)
}

// Get the colour of a square
func (b *Board) getSquareColour(i, j int) color.NRGBA {
	if b.flipped {
		i++
	}
	if (i+j)%2 == 0 {
		return b.gui.theme.chessBoardTheme.square1
	}
	return b.gui.theme.chessBoardTheme.square2
}

// Draw a move list element
func (b *Board) moveListElement(gtx layout.Context, i int) layout.Dimensions {
	buttonWidth := b.squareSize.X*2 - 40
	return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
		// move number
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return button(gtx, b.gui.theme, fmt.Sprintf("%d.", i+1), i, 40, &widget.Clickable{})
		}),
		// player 1 move
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if i*2+1 < len(b.moves) {
				return b.moves[i*2+1].Layout(gtx, b.gui.theme, i, buttonWidth)
			}
			return layout.Dimensions{}
		}),
		// player 2 move
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if i*2+2 < len(b.moves) {
				return b.moves[i*2+2].Layout(gtx, b.gui.theme, i, buttonWidth)
			}
			return layout.Dimensions{}
		}),
		// spacer
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return button(gtx, b.gui.theme, "", i, -1, &widget.Clickable{})
		}),
	)

}

// Layout the move button
func (m *MoveButton) Layout(gtx layout.Context, th *chessAnalysisTheme, i, width int) layout.Dimensions {
	return button(gtx, th, m.notation, i, width, m.widget)
}

// Layout the button
func button(gtx layout.Context, th *chessAnalysisTheme, text string, i, width int, widget *widget.Clickable) layout.Dimensions {
	button := material.Button(th.giouiTheme, widget, text)
	button.CornerRadius = unit.Dp(0)
	if i%2 == 0 {
		button.Background = th.bg
	} else {
		button.Background = color.NRGBA{0, 0, 0, 0}
	}
	button.Inset = layout.UniformInset(unit.Dp(1))
	height := 40
	if width == -1 {
		width = gtx.Constraints.Max.X
	}
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = height
			gtx.Constraints.Max.Y = height
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width
			return button.Layout(gtx)
		}),
	)
}

// Checks if a button has been clicked
func (b *Board) updateState(gtx layout.Context) {
	// keypress
	for {
		keyEvent, ok := gtx.Event(
			key.Filter{},
		)
		if !ok {
			break
		}
		if ev, ok := keyEvent.(key.Event); ok {
			if (ev.Name == key.NameLeftArrow ||
				ev.Name == key.NameRightArrow) &&
				ev.State == key.Press {
				b.arrowKeys(ev)
			}
		}
	}

	// Buttons
	if b.moves == nil || len(b.moves) == 0 {
		return
	}
	for i, move := range b.moves {
		if move.widget.Clicked(gtx) {
			b.gameState = move.gameState
			b.stateNum = i
		}
	}
}

func (b *Board) arrowKeys(e key.Event) {
	if b.moves == nil {
		return
	}
	var newState *game.Game
	switch e.Name {
	case key.NameLeftArrow:
		if b.stateNum == 0 {
			return
		}
		newState = b.moves[b.stateNum-1].gameState
		b.stateNum--
	case key.NameRightArrow:
		if b.stateNum == len(b.moves)-1 {
			return
		}
		newState = b.moves[b.stateNum+1].gameState
		b.stateNum++
	default:
		return
	}
	b.gameState = newState
}
