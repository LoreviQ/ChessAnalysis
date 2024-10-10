package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type settingsMenu struct {
	gui      *GUI
	settings []*setting
}

type setting struct {
	name        string
	settingType string
	editor      *widget.Editor
	button      *widget.Clickable
	data        string
}

func newSettingsMenu(g *GUI) *settingsMenu {
	settings := []*setting{
		{
			name:        "Engine Path",
			settingType: "button",
			editor:      nil,
			button:      &widget.Clickable{},
			data: func() string {
				if g.eng == nil {
					return ""
				}
				return g.eng.Path
			}(),
		},
	}
	return &settingsMenu{
		gui:      g,
		settings: settings,
	}
}

func (sm *settingsMenu) Layout(gtx layout.Context) layout.Dimensions {
	if !sm.gui.header.buttons[1].show {
		return layout.Dimensions{}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Background{}.Layout(gtx,
						sm.BG,
						sm.Menu,
					)
				}),
				layout.Flexed(1, layout.Spacer{}.Layout),
			)
		}),
		layout.Flexed(1, layout.Spacer{}.Layout),
	)
}

func (sm *settingsMenu) BG(gtx layout.Context) layout.Dimensions {
	defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, sm.gui.theme.bg)
	return layout.Dimensions{Size: gtx.Constraints.Min}
}

func (sm *settingsMenu) Menu(gtx layout.Context) layout.Dimensions {
	// dimensions
	height := sm.gui.board.squareSize.Y * 6
	width := sm.gui.board.squareSize.X * 8
	bounds := image.Point{
		X: width,
		Y: height,
	}
	margins := layout.UniformInset(unit.Dp(10))

	// labels
	title := material.Label(sm.gui.theme.giouiTheme, unit.Sp(32), "Settings")
	title.Color = sm.gui.theme.text
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Maintain the existing dimensions
			return layout.Dimensions{Size: bounds}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Layout the labels
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(title.Layout),
				)
			})
		}),
	)
}
