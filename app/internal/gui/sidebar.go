package gui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type sidebar struct {
}

func newSidebar() *sidebar {
	return &sidebar{}
}

func (h *sidebar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Dimensions{}
}
