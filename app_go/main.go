package main

import (
	"gioui.org/app"
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/gui"
)

func main() {
	go gui.CreateGUI()
	app.Main()
}
