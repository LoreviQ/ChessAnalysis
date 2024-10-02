package main

import (
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
)

func main() {
	g := game.NewGame()
	b := g.Board
	b.MovePiece(game.Move{FromFile: 'b', FromRank: 1, ToFile: 'b', ToRank: 3})
	b.MovePiece(game.Move{FromFile: 'c', FromRank: 1, ToFile: 'c', ToRank: 3})
	b.MovePiece(game.Move{FromFile: 'd', FromRank: 1, ToFile: 'd', ToRank: 3})
	b.MovePiece(game.Move{FromFile: 'f', FromRank: 1, ToFile: 'f', ToRank: 3})
	b.MovePiece(game.Move{FromFile: 'g', FromRank: 1, ToFile: 'g', ToRank: 3})
	moves := g.GetPossibleMoves()
	notations, err := game.ConvertMovesToShortAlgebraicNotation(moves)
	if err != nil {
		panic(err)
	}
	for notation := range notations {
		println(notation)
	}

}
