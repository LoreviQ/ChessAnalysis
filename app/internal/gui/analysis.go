package gui

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/LoreviQ/ChessAnalysis/app/internal/eval"
)

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

// Draw the analysis panel
func (b *Board) drawAnalysis(gtx layout.Context) layout.Dimensions {
	margins := layout.Inset{
		Top:    20,
		Bottom: 20,
		Left:   20,
		Right:  20,
	}
	return layout.Stack{}.Layout(gtx,
		// Fill the background
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			rect := image.Rectangle{
				Max: gtx.Constraints.Max,
			}
			paint.FillShape(gtx.Ops, b.gui.theme.contrastBg, clip.Rect(rect).Op())
			return layout.Dimensions{Size: rect.Max}
		}),
		// Panel contents
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if b.moves == nil {
				return layout.Dimensions{}
			}
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
					// Eval Graph
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return b.drawEvalGraph(gtx)
					}),
					// Spacer
					layout.Rigid(layout.Spacer{Height: 20}.Layout),
					// Best lines
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return b.BestLineLists[0].Layout(gtx, len(b.moves[b.stateNum].evals[0].BestLine), func(gtx layout.Context, i int) layout.Dimensions {
							return b.drawBestLineSegment(gtx, i)
						})
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
			if b.evalComplete() {
				return layout.Dimensions{}
			}
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
				if move.evals[0] == nil {
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
					Y: int(centre.Y) - 1,
				},
				Max: image.Point{
					X: gtx.Constraints.Max.X,
					Y: int(centre.Y) + 1,
				},
			}
			paint.FillShape(gtx.Ops, b.gui.theme.contrastFg, clip.Rect(halfwayLine).Op())
			return layout.Dimensions{}
		}),
	)
}

// Draw a segment of the best line
func (b *Board) drawBestLineSegment(gtx layout.Context, i int) layout.Dimensions {
	th := b.gui.theme
	eval := b.moves[b.stateNum].evals[0]
	var label string
	bestLine := eval.BestLine
	if i == 0 {
		label = b.getScoreStr(b.stateNum)
	} else {
		label = bestLine[i-1]
	}
	button := material.Button(th.giouiTheme, &widget.Clickable{}, label)
	button.CornerRadius = unit.Dp(5)
	button.Background = th.bg
	button.Inset = layout.UniformInset(unit.Dp(10))
	button.Color = th.text
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			margins := layout.Inset{
				Right: 5,
			}
			gtx.Constraints.Min.Y = 40
			gtx.Constraints.Max.Y = 40
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return button.Layout(gtx)
			})
		}),
	)
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

// Layout a generic button
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

// produce a score multiplier between 100 and 900 for eval bar
func (b *Board) getScoreMult(stateNum int) int {
	if b.moves == nil {
		return 500
	}
	e := eval.GetEvalNum(b.moves[stateNum].evals, 1)
	if e == nil {
		return 500 // default value
	}
	turn := 1
	if stateNum%2 == 1 {
		turn = -1
	}
	if b.flipped {
		turn *= -1
	}
	if e.Mate {
		if b.flipped {
			return 1000
		}
		return 0
	}
	score := e.Score * turn
	return 500 - min(400, max(-400, score))
}

// Produce a formatted string representing the score
func (b *Board) getScoreStr(stateNum int) string {
	if b.moves == nil {
		return ""
	}
	e := eval.GetEvalNum(b.moves[stateNum].evals, 1)
	if e == nil {
		return "" // default value
	}
	turn := 1
	if stateNum%2 == 1 {
		turn = -1
	}
	if b.flipped {
		turn *= -1
	}
	if e.Mate {
		return fmt.Sprintf("M%d", int(math.Abs(float64(e.MateIn))))
	}
	score := e.Score * turn
	return fmt.Sprintf("%.1f", float32(score)/100)
}

// Returns a bool indicating if the game has been evaluated
func (b *Board) evalComplete() bool {
	if b.moves == nil {
		return false
	}
	return b.moves[len(b.moves)-1].evals[0] != nil
}
