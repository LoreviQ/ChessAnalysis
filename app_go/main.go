package main

import (
	"fmt"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
)

func main() {
	g := game.NewGame()
	moves := g.GetPossibleMoves()
	for _, move := range moves {
		notation, _ := move.ShortAlgebraicNotation(false, false)
		fmt.Println(notation)
	}

}
