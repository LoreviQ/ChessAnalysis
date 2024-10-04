package gui

import (
	"errors"
	"fmt"
	"image"

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
	if err == nil {
		gameId = games[0].ID
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
	paint.FillShape(gtx.Ops, s.gui.theme.chessBoardTheme.contrastBg, clip.Rect(rect).Op())

	err := s.updateState()
	if err != nil {
		return layout.Dimensions{Size: sidebarSize}
	}

	return s.list.Layout(gtx, len(s.games), func(gtx layout.Context, i int) layout.Dimensions {
		return s.sidebarListelement(gtx, i)
	})
}

func (s *sidebar) sidebarListelement(gtx layout.Context, i int) layout.Dimensions {
	game := s.games[i].game
	return layout.Inset{Top: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				gameid := material.Label(s.gui.theme.giouiTheme,
					unit.Sp(16),
					fmt.Sprintf("%d:%s", game.ID, game.ChessdotcomID))
				gameid.Alignment = text.Start
				gameid.Color = s.gui.theme.chessBoardTheme.text
				return gameid.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				date := material.Label(s.gui.theme.giouiTheme,
					unit.Sp(16),
					fmt.Sprintf(game.CreatedAt))
				date.Alignment = text.End
				date.Color = s.gui.theme.chessBoardTheme.text
				return date.Layout(gtx)
			}),
		)
	})
}

func (s *sidebar) updateState() error {
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
	// Change board if selected game is different
	if s.gui.board != nil && s.gui.board.activeGameID != s.selectedGameID {
		s.gui.board = newBoard(s.gui, s.selectedGameID)
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
