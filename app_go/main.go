package main

import (
	"fmt"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
)

func main() {
	b := game.CustomBoard([8][8]*game.Piece{
		{
			&game.Piece{PieceType: game.Rook, Color: "white", Active: true},
			&game.Piece{PieceType: game.Knight, Color: "white", Active: true},
			&game.Piece{PieceType: game.Bishop, Color: "white", Active: true},
			&game.Piece{PieceType: game.Queen, Color: "white", Active: true},
			&game.Piece{PieceType: game.King, Color: "white", Active: true},
			&game.Piece{PieceType: game.Bishop, Color: "white", Active: true},
			&game.Piece{PieceType: game.Knight, Color: "white", Active: true},
			&game.Piece{PieceType: game.Rook, Color: "white", Active: true},
		},
		{
			&game.Piece{PieceType: game.Pawn, Color: "white", Active: true},
			nil,
			nil,
			&game.Piece{PieceType: game.Pawn, Color: "white", Active: true},
			&game.Piece{PieceType: game.Pawn, Color: "white", Active: true},
			&game.Piece{PieceType: game.Pawn, Color: "white", Active: true},
			nil,
			nil,
		},
		{nil, nil, &game.Piece{PieceType: game.Pawn, Color: "white", Active: true}, nil, nil, nil, nil, nil},
		{nil, nil, nil, &game.Piece{PieceType: game.Pawn, Color: "black", Active: true}, nil, nil, nil, nil},
		{nil, &game.Piece{PieceType: game.Pawn, Color: "white", Active: true}, &game.Piece{PieceType: game.Pawn, Color: "black", Active: true}, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{
			&game.Piece{PieceType: game.Pawn, Color: "black", Active: true},
			&game.Piece{PieceType: game.Pawn, Color: "black", Active: true},
			nil,
			nil,
			&game.Piece{PieceType: game.Pawn, Color: "black", Active: true},
			&game.Piece{PieceType: game.Pawn, Color: "black", Active: true},
			&game.Piece{PieceType: game.Pawn, Color: "white", Active: true},
			&game.Piece{PieceType: game.Pawn, Color: "white", Active: true},
		},
		{
			&game.Piece{PieceType: game.Rook, Color: "black", Active: true},
			&game.Piece{PieceType: game.Knight, Color: "black", Active: true},
			&game.Piece{PieceType: game.Bishop, Color: "black", Active: true},
			&game.Piece{PieceType: game.Queen, Color: "black", Active: true},
			&game.Piece{PieceType: game.King, Color: "black", Active: true},
			&game.Piece{PieceType: game.Bishop, Color: "black", Active: true},
			nil,
			nil,
		},
	})
	board := b.PrintBoard()
	fmt.Println(board)
}
