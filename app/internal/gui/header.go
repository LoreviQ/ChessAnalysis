package gui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type header struct {
}

func newHeader() *header {
	return &header{}
}

func (h *header) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Dimensions{}
}
