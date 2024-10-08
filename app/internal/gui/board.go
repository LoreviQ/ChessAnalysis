package gui

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/f32"
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
	gui          *GUI
	squares      [][]layout.FlexChild
	squareSize   image.Point
	activeGameID int
	movesList    *widget.List
	gameState    *game.Game
	stateNum     int
	moves        []*MoveButton
	flipped      bool
	evaluated    bool
}

type MoveButton struct {
	move      *game.Move
	notation  string
	widget    *widget.Clickable
	gameState *game.Game
	eval      *eval.MoveEval
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
	moveStrs, err := g.db.GetMovesByID(selectedGame.ID)
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
	// evaluate the game
	done := make(chan struct{})
	go evaluateGame(g.eng, gameState.MoveHistory, moves, done)
	go func() {
		<-done
		// Draw a new frame
		g.window.Invalidate()
		g.board.evaluated = true
	}()

	// check if the board should be flipped
	flipped := true
	if selectedGame.PlayerIsWhite {
		flipped = false
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

// Draw the evaluation bar
func (b *Board) drawEvalBar(gtx layout.Context) layout.Dimensions {
	offset := 50
	height := b.squareSize.Y * 8
	scoreMult := b.getScoreMult(b.stateNum)
	rect1 := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: offset,
		},
		Max: image.Point{
			X: b.squareSize.Y / 3,
			Y: height + offset,
		},
	}
	rect2 := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: offset,
		},
		Max: image.Point{
			X: b.squareSize.Y / 3,
			Y: int(scoreMult*height/1000) + offset,
		},
	}
	colour1 := b.gui.theme.chessBoardTheme.player1
	colour2 := b.gui.theme.chessBoardTheme.player2
	if b.flipped {
		colour1, colour2 = colour2, colour1
	}
	paint.FillShape(gtx.Ops, colour1, clip.Rect(rect1).Op())
	paint.FillShape(gtx.Ops, colour2, clip.Rect(rect2).Op())
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
			if b.moves == nil {
				return layout.Dimensions{}
			}
			margins := layout.Inset{
				Top:    20,
				Bottom: 20,
				Left:   20,
				Right:  20,
			}
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
					// Eval Graph
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return b.drawEvalGraph(gtx)
					}),
					// Spacer
					layout.Rigid(layout.Spacer{Height: 20}.Layout),
					// Move list
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return b.movesList.Layout(gtx, (len(b.moves))/2, func(gtx layout.Context, i int) layout.Dimensions {
							return b.moveListElement(gtx, i)
						})
					}),
				)
			})
		}),
	)
}

// Get the colour of a square
func (b *Board) getSquareColour(i, j int) color.NRGBA {
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
	button.Color = th.text
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

func evaluateGame(engine *eval.Engine, moves []game.Move, moveButtons []*MoveButton, done chan struct{}) error {
	if engine == nil {
		return errors.New("no engine")
	}
	notations := game.ConvertMovesToUCINotation(moves)
	evals := engine.EvalGame(strings.Join(notations, " "))
	for i, eval := range evals {
		moveButtons[i].eval = eval
	}
	// signal that the evaluation is done
	done <- struct{}{}
	return nil
}

// produce a score multiplier between 100 and 900 for eval bar
func (b *Board) getScoreMult(stateNum int) int {
	if b.moves == nil {
		return 500
	}
	eval := b.moves[stateNum].eval
	if eval == nil {
		return 500 // default value
	}
	turn := 1
	if stateNum%2 == 1 {
		turn = -1
	}
	if b.flipped {
		turn *= -1
	}
	if eval.Mate {
		if b.flipped {
			return 1000
		}
		return 0
	}
	score := eval.Score * turn
	return 500 - min(400, max(-400, score))
}

func (b *Board) drawEvalGraph(gtx layout.Context) layout.Dimensions {
	height := b.squareSize.Y * 2
	player1Colour := b.gui.theme.chessBoardTheme.player1
	player2Colour := b.gui.theme.chessBoardTheme.player2
	if b.flipped {
		player1Colour, player2Colour = player2Colour, player1Colour
	}
	return layout.Stack{}.Layout(gtx,
		//loading message
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if b.evalComplete() {
				return layout.Dimensions{}
			}
			rect := image.Rectangle{
				Max: image.Point{
					X: gtx.Constraints.Max.X,
					Y: height,
				},
			}
			paint.FillShape(gtx.Ops, b.gui.theme.bg, clip.Rect(rect).Op())
			return layout.Dimensions{Size: rect.Max}
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(b.gui.theme.giouiTheme, unit.Sp(20), "Evaluating game...")
			label.Color = b.gui.theme.text
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return label.Layout(gtx)
			})
		}),
		// Fill the background
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if !b.evalComplete() {
				return layout.Dimensions{}
			}
			rect := image.Rectangle{
				Max: image.Point{
					X: gtx.Constraints.Max.X,
					Y: height,
				},
			}
			paint.FillShape(gtx.Ops, player2Colour, clip.Rect(rect).Op())
			return layout.Dimensions{Size: rect.Max}
		}),
		// Eval graph
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if !b.evalComplete() {
				return layout.Dimensions{}
			}
			// Get Scores
			scores := make([]float32, len(b.moves))
			for i, move := range b.moves {
				if move.eval == nil {
					return layout.Dimensions{}
				}
				score := b.getScoreMult(i)
				scores[i] = float32(score) / 1000
			}
			xGap := float32(gtx.Constraints.Max.X) / float32(len(scores)-1)
			// Draw path
			var path clip.Path
			centres := make([]f32.Point, len(scores)+1)
			path.Begin(gtx.Ops)
			evalPoint := f32.Pt(0, float32(b.squareSize.Y))
			path.MoveTo(evalPoint)
			centres[0] = evalPoint
			for i := 1; i < len(scores); i++ {
				evalPoint = f32.Pt(xGap*float32(i), float32(height)*(scores[i]))
				path.LineTo(evalPoint)
				centres[i] = evalPoint
			}
			path.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(height)))
			path.LineTo(f32.Pt(0, float32(height)))
			path.Close()
			outline := clip.Outline{Path: path.End()}.Op()
			paint.FillShape(gtx.Ops, player1Colour, outline)
			// Turn indicator
			centre := centres[b.stateNum]
			turnIndicator := image.Rectangle{
				Min: image.Point{
					X: int(centre.X) - 1,
					Y: 0,
				},
				Max: image.Point{
					X: int(centre.X) + 1,
					Y: height,
				},
			}
			paint.FillShape(gtx.Ops, b.gui.theme.contrastFg, clip.Rect(turnIndicator).Op())
			// halfway line
			halfwayLine := image.Rectangle{
				Min: image.Point{
					X: 0,
					Y: height/2 - 1,
				},
				Max: image.Point{
					X: gtx.Constraints.Max.X,
					Y: height/2 + 1,
				},
			}
			paint.FillShape(gtx.Ops, b.gui.theme.contrastFg, clip.Rect(halfwayLine).Op())
			return layout.Dimensions{}
		}),
	)
}

// Returns a bool indicating if the game has been evaluated
func (b *Board) evalComplete() bool {
	if b.moves == nil {
		return false
	}
	return b.moves[len(b.moves)-1].eval != nil
}
