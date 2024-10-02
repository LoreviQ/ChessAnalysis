package main

import (
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/server"
)

func main() {
	g := game.NewGame()
	s := server.NewServer()
	go s.ListenAndServe()
	g.Play()
}
