package gui

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type sidebar struct {
	gui  *GUI
	list *widget.List
}

func newSidebar(g *GUI) *sidebar {
	return &sidebar{
		gui: g,
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
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

	//Make games components
	games, err := s.gui.db.GetGames()
	if err != nil {
		return layout.Dimensions{}
	}
	s.gui.board.activeGameID = games[0].ID
	return s.list.Layout(gtx, len(games), func(gtx layout.Context, i int) layout.Dimensions {
		game := games[i]
		return layout.Inset{Top: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					gameid := material.Label(s.gui.theme.giouiTheme,
						unit.Sp(16),
						fmt.Sprintf("%d:%s", game.ID, game.ChessdotcomID))
					gameid.Alignment = text.Start
					gameid.Color = s.gui.theme.chessBoardTheme.fg
					return gameid.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					date := material.Label(s.gui.theme.giouiTheme,
						unit.Sp(16),
						fmt.Sprintf(game.CreatedAt))
					date.Alignment = text.End
					date.Color = s.gui.theme.chessBoardTheme.fg
					return date.Layout(gtx)
				}),
			)
		})
	})
}
