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

	expected := "8 ♖♘♗♕♔♗♘♖\n7 ♙♙♙♙♙♙♙♙\n6 \u3000\u3000\u3000\u3000\u3000\u3000\u3000\u3000\n5 \u3000\u3000\u3000\u3000\u3000\u3000\u3000\u3000\n4 \u3000\u3000\u3000\u3000\u3000\u3000\u3000\u3000\n3 \u3000\u3000\u3000\u3000\u3000\u3000\u3000\u3000\n2 ♟♟♟♟♟♟♟♟\n1 ♜♞♝♛♚♝♞♜\n  a b c d e f g h\n"
	if boardStr != expected {
		t.Errorf("Expected board string:\n%s\nGot:\n%s", expected, boardStr)
	}
}

func TestMovePiece(t *testing.T) {
	b := NewBoard()
	tests := []struct {
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

	for _, tt := range tests {
		err := b.MovePiece(tt.fromFile, tt.fromRank, tt.toFile, tt.toRank)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	board := b.PrintBoard()
	expected := "8 ♖♘♗♕♔♗\u3000♖\n7 ♙♙♙♙\u3000♙♙♙\n6 \u3000\u3000\u3000\u3000\u3000♘\u3000\u3000\n5 \u3000\u3000\u3000\u3000♙\u3000\u3000\u3000\n4 \u3000\u3000\u3000\u3000♟♟\u3000\u3000\n3 \u3000\u3000♞\u3000\u3000\u3000\u3000\u3000\n2 ♟♟♟♟\u3000\u3000♟♟\n1 ♜\u3000♝♛♚♝♞♜\n  a b c d e f g h\n"
	if board != expected {
		t.Errorf("Expected board string:\n%s\nGot:\n%s", expected, board)
	}
}

func TestMovePieceInvalid(t *testing.T) {
	b := NewBoard()
	err := b.MovePiece("e", 3, "e", 4)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
