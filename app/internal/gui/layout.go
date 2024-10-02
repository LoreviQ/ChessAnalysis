package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func (g *GUI) layout(gtx layout.Context) layout.Dimensions {
	// Set Background
	rect := image.Rectangle{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Constraints.Max.Y,
		},
	}
	paint.FillShape(gtx.Ops, g.theme.Palette.Bg, clip.Rect(rect).Op())
	return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return g.header.Layout(gtx, g.theme)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return g.sidebar.Layout(gtx, g.theme)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return g.board.Layout(gtx, g.theme)
				}),
			)
		}),
	)

}
