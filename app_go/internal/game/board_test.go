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

	expected := `8 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 
7 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ 
6                 
5                 
4                 
3                 
2 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟ 
1 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 
  a b c d e f g h
`
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
		err      error
	}{
		{"e", 2, "e", 4, nil},
		{"e", 7, "e", 5, nil},
		{"b", 1, "c", 3, nil},
		{"g", 8, "f", 6, nil},
		{"f", 2, "f", 4, nil},
		{"a", 1, "a", 2, ErrSquareOccupied},
		{"d", 1, "e", 8, ErrSquareOccupied},
		{"g", 3, "b", 1, ErrNoPieceAtSquare},
		{"c", 4, "d", 5, ErrNoPieceAtSquare},
	}

	for _, tt := range tests {
		err := b.MovePiece(tt.fromFile, tt.fromRank, tt.toFile, tt.toRank)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err == nil && tt.err != nil {
			t.Errorf("Expected error: %v, got nil", tt.err)
		}
	}

	// Check board after moves
	board := b.PrintBoard()
	expected := `8 ♖ ♘ ♗ ♕ ♔ ♗   ♖ 
7 ♙ ♙ ♙ ♙   ♙ ♙ ♙ 
6           ♘     
5         ♙       
4         ♟ ♟     
3     ♞           
2 ♟ ♟ ♟ ♟     ♟ ♟ 
1 ♜   ♝ ♛ ♚ ♝ ♞ ♜ 
  a b c d e f g h
`
	if board != expected {
		t.Errorf("Expected board string:\n%s\nGot:\n%s", expected, board)
	}
}

func TestMovePieceCapture(t *testing.T) {
	b := NewBoard()
	tests := []struct {
		fromFile string
		fromRank int
		err      error
	}{
		{"c", 1, nil},
		{"f", 2, nil},
		{"f", 6, ErrNoPieceAtSquare},
		{"c", 3, ErrNoPieceAtSquare},
	}

	for _, tt := range tests {
		err := b.CapturePiece(tt.fromFile, tt.fromRank)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err == nil && tt.err != nil {
			t.Errorf("Expected error: %v, got nil", tt.err)
		}
	}

	// Check captured pieces
	takenByBlack := b.GetCapturedByColour("black")
	expectedByBlack := []pieceType{Bishop, Pawn}
	if len(takenByBlack) != len(expectedByBlack) {
		t.Errorf("Expected %v, got %v", expectedByBlack, takenByBlack)
	}
	for i, p := range takenByBlack {
		if p.pType != expectedByBlack[i] {
			t.Errorf("Expected %v, got %v", expectedByBlack, takenByBlack)
		}
	}
}

func TestPromotePawn(t *testing.T) {
	b := NewBoard()
	tests := []struct {
		file  string
		rank  int
		pType pieceType
		err   error
	}{
		{"e", 2, Queen, nil},
		{"a", 7, Rook, nil},
		{"e", 4, Queen, ErrNoPieceAtSquare},
		{"h", 1, Knight, ErrPieceNotPawn},
	}

	for _, tt := range tests {
		err := b.PromotePawn(tt.file, tt.rank, tt.pType)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err == nil && tt.err != nil {
			t.Errorf("Expected error: %v, got nil", tt.err)
		}
	}
}
