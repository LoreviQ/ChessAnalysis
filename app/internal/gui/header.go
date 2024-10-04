package gui

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type header struct {
	gui     *GUI
	size    image.Point
	buttons []*headerButton
}

type headerButton struct {
	name            string
	widget          *widget.Clickable
	menu            *component.MenuState
	menuContextArea *component.ContextArea
}

func newHeader(g *GUI) *header {
	return &header{
		gui: g,
		buttons: []*headerButton{
			{
				name:   "Themes",
				widget: &widget.Clickable{},
				menu: &component.MenuState{
					Options: []func(gtx layout.Context) layout.Dimensions{
						material.Label(g.theme.giouiTheme, unit.Sp(16), "HotDogStand").Layout,
					},
				},
				menuContextArea: &component.ContextArea{
					Activation:       pointer.ButtonPrimary,
					AbsolutePosition: true,
				},
			},
		},
	}
}

func (h *header) Layout(gtx layout.Context) layout.Dimensions {
	// Define the fixed size for the header
	h.size = image.Point{X: gtx.Constraints.Max.X, Y: 50} // Fixed height of 50 pixels

	// Adjust the constraints to enforce the fixed size
	gtx.Constraints.Min = h.size
	gtx.Constraints.Max = h.size

	// Header bg
	rect := image.Rectangle{
		Max: h.size,
	}
	paint.FillShape(gtx.Ops, h.gui.theme.contrastFg, clip.Rect(rect).Op())
	return layout.Stack{}.Layout(gtx,
		// Header Size
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: h.size}
		}),
		// Buttons
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				h.buttonsLayout()...,
			)
		}),
	)
}

func (h *header) buttonsLayout() []layout.FlexChild {
	children := make([]layout.FlexChild, len(h.buttons))
	for i, button := range h.buttons {
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return button.Layout(gtx, h.gui.theme)
		})
	}
	return children
}

func (hb *headerButton) Layout(gtx layout.Context, th *chessAnalysisTheme) layout.Dimensions {
	button := material.Button(th.giouiTheme, hb.widget, hb.name)
	button.CornerRadius = unit.Dp(0)
	button.Inset = layout.UniformInset(unit.Dp(1))
	button.Background = color.NRGBA{0, 0, 0, 0}
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			gtx.Constraints.Min.X = gtx.Constraints.Min.X + 100
			return button.Layout(gtx)

		}),
	)
}
