package main

import (
	"fmt"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
)

func main() {
	b := game.NewBoard()
	moves := []struct {
		fromFile string
		fromRank int
		toFile   string
		toRank   int
	}{
		{"e", 2, "e", 4},
		{"e", 7, "e", 5},
		{"b", 1, "c", 3},
		{"g", 8, "f", 6},
		{"f", 2, "f", 4},
	}
	for _, move := range moves {
		b.MovePiece(move.fromFile, move.fromRank, move.toFile, move.toRank)
	}
	board := b.PrintBoard()
	fmt.Println(board)
}
