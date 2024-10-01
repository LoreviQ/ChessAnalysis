package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()

	// Check initial positions of some pieces
	if b.Squares[0][0].PieceType != Rook || b.Squares[0][0].Color != "white" {
		t.Errorf("Expected white Rook at a8, got %v", b.Squares[0][0])
	}
	if b.Squares[7][4].PieceType != King || b.Squares[7][4].Color != "black" {
		t.Errorf("Expected black King at e1, got %v", b.Squares[7][4])
	}
	if b.Squares[1][0].PieceType != Pawn || b.Squares[1][0].Color != "white" {
		t.Errorf("Expected white Pawn at a7, got %v", b.Squares[1][0])
	}
	if b.Squares[6][0].PieceType != Pawn || b.Squares[6][0].Color != "black" {
		t.Errorf("Expected black Pawn at a2, got %v", b.Squares[6][0])
	}
}

func TestGetPieceAtSquare(t *testing.T) {
	b := NewBoard()

	tests := []struct {
		file  rune
		rank  int
		pType PieceType
		color string
		err   error
	}{
		{'a', 1, Rook, "white", nil},
		{'e', 8, King, "black", nil},
		{'e', 8, King, "black", nil},
		{'a', 2, Pawn, "white", nil},
		{'a', 7, Pawn, "black", nil},
		{'a', 12, 0, "", ErrInvalidRank},
		{'i', 1, 0, "", ErrInvalidFile},
	}

	for _, tt := range tests {
		p, err := b.GetPieceAtSquare(tt.file, tt.rank)
		if err != tt.err {
			t.Errorf("Expected error %v, got %v", tt.err, err)
		}
		if p != nil && (p.PieceType != tt.pType || p.Color != tt.color) {
			t.Errorf("Expected piece %v, got %v", tt, p)
		}
	}
}

func TestGetPieceAtSquareAfterMove(t *testing.T) {
	g := NewGame()
	b := g.Board
	piece, _ := b.GetPieceAtSquare('e', 2)
	move := Move{FromFile: 'e', FromRank: 2, ToFile: 'e', ToRank: 4}
	b.MovePiece(move)

	tests := []struct {
		file  rune
		rank  int
		piece *Piece
	}{
		{'e', 4, piece},
		{'e', 2, nil},
	}

	for _, tt := range tests {
		p, _ := b.GetPieceAtSquare(tt.file, tt.rank)
		if p != tt.piece {
			t.Errorf("Expected piece %v, got %v", tt.piece, p)
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
		move Move
		err  error
	}{
		{Move{FromFile: 'e', FromRank: 2, ToFile: 'e', ToRank: 4}, nil},
		{Move{FromFile: 'e', FromRank: 7, ToFile: 'e', ToRank: 5}, nil},
		{Move{FromFile: 'b', FromRank: 1, ToFile: 'c', ToRank: 3}, nil},
		{Move{FromFile: 'g', FromRank: 8, ToFile: 'f', ToRank: 6}, nil},
		{Move{FromFile: 'f', FromRank: 2, ToFile: 'f', ToRank: 4}, nil},
		{Move{FromFile: 'a', FromRank: 1, ToFile: 'a', ToRank: 2}, ErrSquareOccupied},
		{Move{FromFile: 'd', FromRank: 1, ToFile: 'e', ToRank: 8}, ErrSquareOccupied},
		{Move{FromFile: 'g', FromRank: 3, ToFile: 'b', ToRank: 1}, ErrNoPieceAtSquare},
		{Move{FromFile: 'c', FromRank: 4, ToFile: 'd', ToRank: 5}, ErrNoPieceAtSquare},
	}

	for _, tt := range tests {
		err := b.MovePiece(tt.move)
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
		move Move
		err  error
	}{
		{Move{FromFile: 'a', FromRank: 2, ToFile: 'c', ToRank: 7, Capture: 'x'}, nil},
		{Move{FromFile: 'b', FromRank: 2, ToFile: 'f', ToRank: 8, Capture: 'x'}, nil},
		{Move{FromFile: 'c', FromRank: 2, ToFile: 'f', ToRank: 6, Capture: 'x'}, ErrNoPieceAtSquare},
		{Move{FromFile: 'd', FromRank: 2, ToFile: 'c', ToRank: 3, Capture: 'x'}, ErrNoPieceAtSquare},
		{Move{FromFile: 'e', FromRank: 2, ToFile: 'e', ToRank: 1, Capture: 'x'}, ErrSquareOccupied},
	}

	for _, tt := range tests {
		err := b.MovePiece(tt.move)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err == nil && tt.err != nil {
			t.Errorf("Expected error: %v, got nil", tt.err)
		}
	}

	// Check captured pieces
	takenByWhite := b.GetCapturedByColour("white")
	expectedByWhite := []PieceType{Pawn, Bishop}
	if len(takenByWhite) != len(expectedByWhite) {
		t.Errorf("Expected %v, got %v", expectedByWhite, takenByWhite)
	}
	for i, p := range takenByWhite {
		if p.PieceType != expectedByWhite[i] {
			t.Errorf("Expected %v, got %v", expectedByWhite, takenByWhite)
		}
	}
}

func TestPromotePawn(t *testing.T) {
	b := NewBoard()
	tests := []struct {
		file  rune
		rank  int
		pType PieceType
		err   error
	}{
		{'e', 2, Queen, nil},
		{'a', 7, Rook, nil},
		{'e', 4, Queen, ErrNoPieceAtSquare},
		{'h', 1, Knight, ErrPieceNotPawn},
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
