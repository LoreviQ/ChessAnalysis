package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()

	// Check initial positions of some pieces
	if b.squares[0][0].pType != Rook || b.squares[0][0].color != "black" {
		t.Errorf("Expected black Rook at a8, got %v", b.squares[0][0])
	}
	if b.squares[7][4].pType != King || b.squares[7][4].color != "white" {
		t.Errorf("Expected white King at e1, got %v", b.squares[7][4])
	}
	if b.squares[1][0].pType != Pawn || b.squares[1][0].color != "black" {
		t.Errorf("Expected black Pawn at a7, got %v", b.squares[1][0])
	}
	if b.squares[6][0].pType != Pawn || b.squares[6][0].color != "white" {
		t.Errorf("Expected white Pawn at a2, got %v", b.squares[6][0])
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
		{0, "a", Rook, "black"},
		{7, "e", King, "white"},
		{1, "a", Pawn, "black"},
		{6, "a", Pawn, "white"},
	}

	for _, tt := range tests {
		p := b.getPieceAtSquare(tt.rank, tt.file)
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
