package gui

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/LoreviQ/ChessAnalysis/app/internal/database"
)

type sidebar struct {
	gui            *GUI
	list           *widget.List
	selectedGameID int
	games          []*gameButton
}

type gameButton struct {
	game   database.Game
	widget *widget.Clickable
}

func newSidebar(g *GUI) *sidebar {
	var gameId int
	games, err := g.db.GetGames()
	if err == nil && len(games) > 0 {
		// default to last game
		gameId = games[len(games)-1].ID
	}
	gameButtons := make([]*gameButton, len(games))
	for i, game := range games {
		gameButtons[i] = &gameButton{
			game:   game,
			widget: &widget.Clickable{},
		}
	}
	return &sidebar{
		gui: g,
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		selectedGameID: gameId,
		games:          gameButtons,
	}
}

func (s *sidebar) Layout(gtx layout.Context) layout.Dimensions {
	// Define the fixed size for the sidebar
	sidebarSize := image.Point{X: 400, Y: gtx.Constraints.Max.Y}

	// Adjust the constraints to enforce the fixed size
	gtx.Constraints.Min = sidebarSize
	gtx.Constraints.Max = sidebarSize

	// Sidebar bg
	rect := image.Rectangle{
		Max: sidebarSize,
	}
	paint.FillShape(gtx.Ops, s.gui.theme.contrastBg, clip.Rect(rect).Op())

	if len(s.games) == 0 {
		label := material.Label(s.gui.theme.giouiTheme, unit.Sp(16), "No games found")
		label.Color = s.gui.theme.text
		return layout.Center.Layout(gtx, label.Layout)
	}
	err := s.updateState(gtx)
	if err != nil {
		return layout.Dimensions{Size: sidebarSize}
	}
	return s.list.Layout(gtx, len(s.games), func(gtx layout.Context, i int) layout.Dimensions {
		return s.games[i].Layout(gtx, s.gui.theme, i, -1)
	})
}

func (s *sidebar) updateState(gtx layout.Context) error {
	// Update games from db
	var err error
	games, err := s.gui.db.GetGames()
	if len(games) > len(s.games) {
		// new game added to db, need to create a new gameButton
		for _, game := range games {
			if !s.gameExists(game.ID) {
				s.games = append(s.games, &gameButton{
					game:   game,
					widget: &widget.Clickable{},
				})
			}
		}
	}
	if err != nil || len(s.games) == 0 {
		return errors.New("failed to get games")
	}

	// Buttons
	if s.games == nil || len(s.games) == 0 {
		return nil
	}

	var selectedGame database.Game
	for _, gameButton := range s.games {
		if gameButton.widget.Clicked(gtx) {
			s.selectedGameID = gameButton.game.ID
			selectedGame = gameButton.game
		}
	}

	// Change board if selected game is different
	if s.gui.board != nil && s.gui.board.activeGameID != s.selectedGameID {
		if selectedGame.ID == 0 {
			for _, gameButton := range s.games {
				if gameButton.game.ID == s.selectedGameID {
					selectedGame = gameButton.game
					break
				}
			}
		}
		s.gui.board = newBoard(s.gui, &selectedGame)
	}
	return nil
}

func (s *sidebar) gameExists(id int) bool {
	for _, gameButton := range s.games {
		if gameButton.game.ID == id {
			return true
		}
	}
	return false
}

func (gb *gameButton) Layout(gtx layout.Context, th *chessAnalysisTheme, i, width int) layout.Dimensions {
	button := material.Button(th.giouiTheme, gb.widget, "")
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
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			margins := layout.Inset{Left: unit.Dp(8)}
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				labelStr := fmt.Sprintf("%d:%s", gb.game.ID, gb.game.ChessdotcomID)
				gameLabel := material.Label(th.giouiTheme, unit.Sp(16), labelStr)
				gameLabel.Color = th.text
				gameLabel.Alignment = text.Start
				return layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, gameLabel.Layout)
				})
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			labelStr := gb.game.CreatedAt
			dateLabel := material.Label(th.giouiTheme, unit.Sp(12), labelStr)
			dateLabel.Color = th.textMuted
			dateLabel.Alignment = text.End
			return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				margins := layout.Inset{Right: unit.Dp(6), Bottom: unit.Dp(4)}
				return margins.Layout(gtx, dateLabel.Layout)
			})
		}),
	)
}
