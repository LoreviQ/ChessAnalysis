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
	expectedBoard := "8 ♖♘♗♕♔♗\u3000♖\n7 ♙♙♙♙\u3000♙♙♙\n6 \u3000\u3000\u3000\u3000\u3000♘\u3000\u3000\n5 \u3000\u3000\u3000\u3000♙\u3000\u3000\u3000\n4 \u3000\u3000\u3000\u3000♟♟\u3000\u3000\n3 \u3000\u3000♞\u3000\u3000\u3000\u3000\u3000\n2 ♟♟♟♟\u3000\u3000♟♟\n1 ♜\u3000♝♛♚♝♞♜\n  a b c d e f g h\n"
	for i := range board {
		boardRune := board[i]
		expectedRune := expectedBoard[i]
		if boardRune != expectedRune {
			fmt.Printf("Expected %c, got %c\n", expectedRune, boardRune)
		}
	}
}
