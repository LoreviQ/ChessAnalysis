package gui

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/LoreviQ/ChessAnalysis/app/internal/database"
	"github.com/LoreviQ/ChessAnalysis/app/internal/eval"
	"github.com/LoreviQ/ChessAnalysis/app/internal/game"
)

type Board struct {
	gui           *GUI
	squares       [][]layout.FlexChild
	squareSize    image.Point
	activeGameID  int
	movesList     *widget.List
	gameState     *game.Game
	stateNum      int
	moves         []*MoveButton
	flipped       bool
	evaluated     bool
	bestLines     *widget.List
	BestLineLists []*widget.List
}

type MoveButton struct {
	move      *game.Move
	notation  string
	widget    *widget.Clickable
	gameState *game.Game
	evals     []*eval.MoveEval
}

func newBoard(g *GUI, selectedGame *database.Game) *Board {
	// default board if no game is selected
	errBoard := &Board{
		gui:       g,
		gameState: game.NewGame(),
	}
	if selectedGame == nil {
		return errBoard
	}
	// Get the moves for the active game
	movesFromDB, err := g.db.GetMovesByID(selectedGame.ID)
	if err != nil {
		return errBoard
	}
	// Get the board state for the active game
	gameState := game.NewGame()
	// check if the move string is valid
	err = gameState.Moves(movesFromDB.Moves)
	if err != nil {
		return errBoard
	}
	// turn the move strings into moves
	gameState.NewGame()
	moves := make([]*MoveButton, len(movesFromDB.Moves)+1)
	moves[0] = &MoveButton{
		move:      nil,
		notation:  "",
		widget:    &widget.Clickable{},
		gameState: gameState.Clone(),
	}
	for i, moveStr := range movesFromDB.Moves {
		move, err := gameState.Move(moveStr)
		if err != nil {
			break
		}
		moves[i+1] = &MoveButton{
			move:      &move,
			notation:  moveStr,
			widget:    &widget.Clickable{},
			gameState: gameState.Clone(),
			evals:     []*eval.MoveEval{},
		}
	}
	if movesFromDB.Depth > 0 {
		// provisionally use scores from the database
		for i, scoreStr := range movesFromDB.Scores {
			if i >= len(moves) {
				break
			}
			eval := eval.ParseScoreStr(scoreStr)
			moves[i].evals = append(moves[i].evals, eval)
		}
	}
	// evaluate the game
	done := make(chan struct{})
	go evaluateGame(g.eng, gameState.MoveHistory, moves, done)
	go func() {
		<-done
		// Draw a new frame
		g.board.evaluated = true
		g.window.Invalidate()
		// Update the database with the new evals
		evals := make([][]*eval.MoveEval, len(moves))
		for i, move := range moves {
			evals[i] = move.evals
		}
		g.db.UpdateEval(movesFromDB.ID, evals)
	}()

	// check if the board should be flipped
	flipped := true
	if selectedGame.PlayerIsWhite {
		flipped = false
	}
	// create lists for best lines
	BestLineLists := make([]*widget.List, g.eng.MultiPV)
	for i := range BestLineLists {
		BestLineLists[i] = &widget.List{
			List: layout.List{
				Axis: layout.Horizontal,
			},
		}
	}
	return &Board{
		gui:          g,
		activeGameID: selectedGame.ID,
		movesList: &widget.List{
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: true,
			},
		},
		gameState: gameState,
		stateNum:  len(moves) - 1,
		moves:     moves,
		flipped:   flipped,
		bestLines: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		BestLineLists: BestLineLists,
	}
}

// Draw the board
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
		row := 7 - i
		if b.flipped {
			row = i
		}
		col := j
		if b.flipped {
			col = 7 - j
		}
		return layout.Stack{}.Layout(gtx,
			// Draw the square
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				square := image.Rectangle{
					Max: b.squareSize,
				}
				paint.FillShape(gtx.Ops, b.getSquareColour(i, j), clip.Rect(square).Op())
				return layout.Dimensions{Size: square.Max}
			}),
			// Draw the file labels
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if i != 7 {
					return layout.Dimensions{}
				}

				label := material.Label(b.gui.theme.giouiTheme, unit.Sp(30), fmt.Sprintf("%c", 'a'+col))
				label.Color = b.getSquareColour(i+1, j)
				return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					margins := layout.Inset{
						Right:  6,
						Bottom: 1,
					}
					return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return label.Layout(gtx)
					})
				})
			}),
			// Draw the rank labels
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if j != 0 {
					return layout.Dimensions{}
				}
				label := material.Label(b.gui.theme.giouiTheme, unit.Sp(30), fmt.Sprintf("%d", row+1))
				label.Color = b.getSquareColour(i, j-1)
				return layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					margins := layout.Inset{
						Left: 6,
						Top:  4,
					}
					return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return label.Layout(gtx)
					})
				})
			}),
			// Draw the piece
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				piece := b.gameState.Board.Squares[row][col]
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

// Get the colour of a square
func (b *Board) getSquareColour(i, j int) color.NRGBA {
	if (i+j)%2 == 0 {
		return b.gui.theme.chessBoardTheme.square1
	}
	return b.gui.theme.chessBoardTheme.square2
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

// Act upon key events
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

// Get the engine to evaluate the game
func evaluateGame(engine *eval.Engine, moves []game.Move, moveButtons []*MoveButton, done chan struct{}) error {
	if engine == nil {
		return errors.New("no engine")
	}
	notations := game.ConvertMovesToUCINotation(moves)
	evalss := engine.EvalGame(strings.Join(notations, " "))
	for i, evals := range evalss {
		moveButtons[i].evals = evals
	}
	// signal that the evaluation is done
	done <- struct{}{}
	return nil
}
