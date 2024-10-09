package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type settingsMenu struct {
	gui *GUI
}

func newSettingsMenu(g *GUI) *settingsMenu {
	return &settingsMenu{
		gui: g,
	}
}

func (sm *settingsMenu) Layout(gtx layout.Context) layout.Dimensions {
	if !sm.gui.header.buttons[1].show {
		return layout.Dimensions{}
	}
	height := sm.gui.board.squareSize.Y * 6
	width := sm.gui.board.squareSize.X * 8
	rect := image.Rectangle{
		Max: image.Point{
			X: width,
			Y: height,
		},
	}

	return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					paint.FillShape(gtx.Ops, sm.gui.theme.bg, clip.Rect(rect).Op())
					paint.FillShape(gtx.Ops, sm.gui.theme.fg, clip.Stroke{Path: clip.Rect(rect).Path(), Width: 2}.Op())
					return layout.Dimensions{Size: rect.Max}
				}),
				layout.Flexed(1, layout.Spacer{}.Layout),
			)
		}),
		layout.Flexed(1, layout.Spacer{}.Layout),
	)
}
