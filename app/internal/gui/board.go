package gui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type board struct {
}

func newBoard() *board {
	return &board{}
}

func (h *board) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Dimensions{}
}
