package gui

import (
	"image"
	"image/color"

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
	name       string
	widget     *widget.Clickable
	menu       *component.MenuState
	subButtons []*headerDropDownButton
	show       bool
}

type headerDropDownButton struct {
	name     string
	widget   *widget.Clickable
	callback func(*headerButton)
}

func newHeader(g *GUI) *header {
	// Themes header button
	buttons := make([]*headerButton, 1)
	themes := []string{"chess.com", "HotDogStand"}
	subButtons := make([]*headerDropDownButton, len(themes))
	for i, theme := range themes {
		subButtons[i] = &headerDropDownButton{
			name:   theme,
			widget: &widget.Clickable{},
			callback: func(hb *headerButton) {
				hb.show = false
				g.theme = NewTheme(theme)
			},
		}
	}
	buttons[0] = &headerButton{
		name:       "Themes",
		widget:     &widget.Clickable{},
		menu:       &component.MenuState{},
		subButtons: subButtons,
		show:       false,
	}
	// Add more buttons here
	return &header{
		gui:     g,
		buttons: buttons,
		size:    image.Point{X: 0, Y: 0},
	}
}

func (h *header) Layout(gtx layout.Context) layout.Dimensions {
	h.updateState(gtx)
	// Header size
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
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
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

func (hddb *headerDropDownButton) Layout(gtx layout.Context, th *chessAnalysisTheme) layout.Dimensions {
	button := material.Button(th.giouiTheme, hddb.widget, hddb.name)
	button.CornerRadius = unit.Dp(0)
	button.Inset = layout.UniformInset(unit.Dp(1))
	button.Color = th.text
	button.Background = color.NRGBA{0, 0, 0, 0}
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = 50
			gtx.Constraints.Min.X = gtx.Constraints.Min.X + 100
			return button.Layout(gtx)
		}),
	)
}

func (hb *headerButton) getSubButtonsLayout(th *chessAnalysisTheme) []func(gtx layout.Context) layout.Dimensions {
	children := make([]func(gtx layout.Context) layout.Dimensions, len(hb.subButtons))
	for i, button := range hb.subButtons {
		children[i] = func(gtx layout.Context) layout.Dimensions {
			return button.Layout(gtx, th)
		}
	}
	return children
}

func (h *header) updateState(gtx layout.Context) {
	h.size = image.Point{X: gtx.Constraints.Max.X, Y: 50}
	// Execute subbutton callbacks
	for _, button := range h.buttons {
		button.menu.Options = button.getSubButtonsLayout(h.gui.theme)
		for _, subButton := range button.subButtons {
			for subButton.widget.Clicked(gtx) {
				subButton.callback(button)
			}
		}
	}
	// Header button click
	for _, headerButton := range h.buttons {
		if headerButton.widget.Clicked(gtx) {
			headerButton.show = !headerButton.show
		}
	}
}

func (hb *headerButton) layoutDropDown(gtx layout.Context, w layout.Widget) layout.Dimensions {
	if hb.show {
		return w(gtx)
	}
	return layout.Dimensions{}
}
