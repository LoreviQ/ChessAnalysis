package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()

	// Check initial positions of some pieces
	if b.Squares[0][0].pType != Rook || b.Squares[0][0].color != "white" {
		t.Errorf("Expected white Rook at a8, got %v", b.Squares[0][0])
	}
	if b.Squares[7][4].pType != King || b.Squares[7][4].color != "black" {
		t.Errorf("Expected black King at e1, got %v", b.Squares[7][4])
	}
	if b.Squares[1][0].pType != Pawn || b.Squares[1][0].color != "white" {
		t.Errorf("Expected white Pawn at a7, got %v", b.Squares[1][0])
	}
	if b.Squares[6][0].pType != Pawn || b.Squares[6][0].color != "black" {
		t.Errorf("Expected black Pawn at a2, got %v", b.Squares[6][0])
	}
}

func TestGetPieceAtSquare(t *testing.T) {
	b := NewBoard()

	tests := []struct {
		rank  int
		file  string
		pType pieceType
		color string
	}{
		{1, "a", Rook, "white"},
		{8, "e", King, "black"},
		{2, "a", Pawn, "white"},
		{7, "a", Pawn, "black"},
	}

	for _, tt := range tests {
		p := b.GetPieceAtSquare(tt.file, tt.rank)
		if p.pType != tt.pType || p.color != tt.color {
			t.Errorf("Expected %s %v at %s%d, got %v", tt.color, tt.pType, tt.file, tt.rank+1, p)
		}
	}
}

func TestPrintBoard(t *testing.T) {
	b := NewBoard()
	boardStr := b.PrintBoard()

	expected := "♖♘♗♕♔♗♘♖\n♙♙♙♙♙♙♙♙\n        \n        \n        \n        \n♟♟♟♟♟♟♟♟\n♜♞♝♛♚♝♞♜\n"
	if boardStr != expected {
		t.Errorf("Expected board string:\n%s\nGot:\n%s", expected, boardStr)
	}
}

func TestMovePiece(t *testing.T) {
	b := NewBoard()
	b.MovePiece("e", 2, "e", 4)
	p := b.GetPieceAtSquare("e", 4)
	if p == nil || p.pType != Pawn || p.color != "white" {
		t.Errorf("Expected white Pawn at e4, got %v", p)
	}
}
