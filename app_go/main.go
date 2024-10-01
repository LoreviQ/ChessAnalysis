package main

import (
	"fmt"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
)

func main() {
	b := game.NewBoard()
	p1 := b.GetPieceAtSquare("e", 2)
	p2 := b.Squares[1][4]
	fmt.Print(p1)
	fmt.Print(p2)
}
