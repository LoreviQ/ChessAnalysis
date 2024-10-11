package gui

import (
	"fmt"
	"image"
	"strconv"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/LoreviQ/ChessAnalysis/app/internal/eval"
	"github.com/ncruces/zenity"
)

type settingsMenu struct {
	gui          *GUI
	settings     []*setting
	submitButton *widget.Clickable
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
		{
			name:        "Threads",
			settingType: "editor",
			editor:      &widget.Editor{},
			button:      nil,
			data: func() string {
				if g.eng == nil {
					return ""
				}
				return fmt.Sprintf("%d", g.eng.Threads)
			}(),
		},
	}
	return &settingsMenu{
		gui:          g,
		settings:     settings,
		submitButton: &widget.Clickable{},
	}
}

func (sm *settingsMenu) Layout(gtx layout.Context) layout.Dimensions {
	sm.updateState(gtx)
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
	margins := layout.UniformInset(unit.Dp(20))

	// labels
	title := material.Label(sm.gui.theme.giouiTheme, unit.Sp(32), "Settings")
	title.Color = sm.gui.theme.text
	children := make([]layout.FlexChild, len(sm.settings)+1)
	children[0] = layout.Rigid(title.Layout)
	for i, setting := range sm.settings {
		children[i+1] = setting.Layout(gtx, sm.gui.theme)
	}
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Maintain the existing dimensions
			return layout.Dimensions{Size: bounds}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Layout the settings
					children...,
				)
			})
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			offset := layout.Inset{
				Top:  unit.Dp(bounds.Y - 100),
				Left: unit.Dp(bounds.X - 200),
			}
			submit := material.Button(sm.gui.theme.giouiTheme, sm.submitButton, "Submit")
			submit.Background = sm.gui.theme.bg
			submit.Color = sm.gui.theme.text
			submit.TextSize = unit.Sp(32)
			return offset.Layout(gtx, submit.Layout)
		}),
	)
}

func (s *setting) Layout(gtx layout.Context, th *chessAnalysisTheme) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		name := material.Label(th.giouiTheme, unit.Sp(16), s.name)
		name.Color = th.text
		margin := layout.Inset{
			Top: unit.Dp(30),
		}
		return margin.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					offset := layout.Inset{
						Top: unit.Dp(7),
					}
					if s.settingType == "editor" {
						offset.Top = unit.Dp(0)
					}
					return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return name.Layout(gtx)
					})
				}),
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					offset := layout.Inset{
						Left: unit.Dp(200),
					}
					switch s.settingType {
					case "button":
						button := material.Button(th.giouiTheme, s.button, s.data)
						button.Background = th.bg
						return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return button.Layout(gtx)
						})
					case "editor":
						editor := material.Editor(th.giouiTheme, s.editor, s.data)
						editor.Color = th.text
						editor.HintColor = th.text
						return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return editor.Layout(gtx)
						})
					}
					return layout.Dimensions{}
				}),
			)
		})
	})
}

func (sm *settingsMenu) updateState(gtx layout.Context) {
	// change settings
	for _, setting := range sm.settings {
		switch setting.settingType {
		case "button":
			if setting.button.Clicked(gtx) {
				switch setting.name {
				case "Engine Path":
					// Get the file path
					filePath, _ := zenity.SelectFile()
					// Save the file path
					if filePath != "" {
						setting.data = filePath
					}
				}
			}
		}
	}
	if sm.submitButton.Clicked(gtx) {
		sm.submitSettings()
	}

}

func (sm *settingsMenu) submitSettings() error {
	// submit settings
	settings := map[string]string{}
	for _, setting := range sm.settings {
		settings[setting.name] = setting.data
	}

	// check if new engine needs to be loaded
	if settings["Engine Path"] != sm.gui.eng.Path {
		// load engine
		threads, err := strconv.Atoi(settings["Threads"])
		if err != nil {
			return err
		}
		newEngine, err := eval.InitializeStockfish(
			settings["Engine Path"],
			60,
			threads,
			1,
		)
		if err != nil {
			return err
		}
		sm.gui.eng = newEngine
	} else {
		// change the settings of current engine (TODO)
	}

	sm.gui.header.buttons[1].show = false
	return nil
}
