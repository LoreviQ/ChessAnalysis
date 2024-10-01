package main

import (
	"fmt"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
)

func main() {
	b := game.NewBoard()
	board := b.PrintBoard()
	fmt.Print(board)
}
