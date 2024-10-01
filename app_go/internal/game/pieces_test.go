package game

import (
	"testing"
)

func TestGetSymbol(t *testing.T) {
	tests := []struct {
		piece     Piece
		expected  rune
		expectErr bool
	}{
		{Piece{PieceType: King}, 'K', false},
		{Piece{PieceType: Queen}, 'Q', false},
		{Piece{PieceType: Rook}, 'R', false},
		{Piece{PieceType: Bishop}, 'B', false},
		{Piece{PieceType: Knight}, 'N', false},
		{Piece{PieceType: Pawn}, 0, false},
	}

	for _, test := range tests {
		result, err := test.piece.getSymbol()
		if (err != nil) != test.expectErr {
			t.Errorf("Expected error: %v, got: %v", test.expectErr, err)
		}
		if result != test.expected {
			t.Errorf("Expected: %c, got: %c", test.expected, result)
		}
	}
}

func TestGetPrintable(t *testing.T) {
	tests := []struct {
		piece     Piece
		expected  rune
		expectErr bool
	}{
		{Piece{PieceType: Pawn, Color: "white"}, '♟', false},
		{Piece{PieceType: Pawn, Color: "black"}, '♙', false},
		{Piece{PieceType: King, Color: "white"}, '♚', false},
		{Piece{PieceType: King, Color: "black"}, '♔', false},
		{Piece{PieceType: Queen, Color: "white"}, '♛', false},
		{Piece{PieceType: Queen, Color: "black"}, '♕', false},
		{Piece{PieceType: Rook, Color: "white"}, '♜', false},
		{Piece{PieceType: Rook, Color: "black"}, '♖', false},
		{Piece{PieceType: Bishop, Color: "white"}, '♝', false},
		{Piece{PieceType: Bishop, Color: "black"}, '♗', false},
		{Piece{PieceType: Knight, Color: "white"}, '♞', false},
		{Piece{PieceType: Knight, Color: "black"}, '♘', false},
	}

	for _, test := range tests {
		result, err := test.piece.getPrintable()
		if (err != nil) != test.expectErr {
			t.Errorf("Expected error: %v, got: %v", test.expectErr, err)
		}
		if result != test.expected {
			t.Errorf("Expected: %c, got: %c", test.expected, result)
		}
	}
}

func TestGetPawnMoves(t *testing.T) {
	g := NewGame()
	g.Board = CustomBoard([8][8]*Piece{
		{
			&Piece{PieceType: Rook, Color: "white", Active: true},
			&Piece{PieceType: Knight, Color: "white", Active: true},
			&Piece{PieceType: Bishop, Color: "white", Active: true},
			&Piece{PieceType: Queen, Color: "white", Active: true},
			&Piece{PieceType: King, Color: "white", Active: true},
			&Piece{PieceType: Bishop, Color: "white", Active: true},
			&Piece{PieceType: Knight, Color: "white", Active: true},
			&Piece{PieceType: Rook, Color: "white", Active: true},
		},
		{
			&Piece{PieceType: Pawn, Color: "white", Active: true},
			nil,
			nil,
			&Piece{PieceType: Pawn, Color: "white", Active: true},
			&Piece{PieceType: Pawn, Color: "white", Active: true},
			&Piece{PieceType: Pawn, Color: "white", Active: true},
			nil,
			nil,
		},
		{nil, nil, &Piece{PieceType: Pawn, Color: "white", Active: true}, nil, nil, nil, nil, nil},
		{nil, nil, nil, &Piece{PieceType: Pawn, Color: "black", Active: true}, nil, nil, nil, nil},
		{nil, &Piece{PieceType: Pawn, Color: "white", Active: true}, &Piece{PieceType: Pawn, Color: "black", Active: true}, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{
			&Piece{PieceType: Pawn, Color: "black", Active: true},
			&Piece{PieceType: Pawn, Color: "black", Active: true},
			nil,
			nil,
			&Piece{PieceType: Pawn, Color: "black", Active: true},
			&Piece{PieceType: Pawn, Color: "black", Active: true},
			&Piece{PieceType: Pawn, Color: "white", Active: true},
			&Piece{PieceType: Pawn, Color: "white", Active: true},
		},
		{
			&Piece{PieceType: Rook, Color: "black", Active: true},
			&Piece{PieceType: Knight, Color: "black", Active: true},
			&Piece{PieceType: Bishop, Color: "black", Active: true},
			&Piece{PieceType: Queen, Color: "black", Active: true},
			&Piece{PieceType: King, Color: "black", Active: true},
			&Piece{PieceType: Bishop, Color: "black", Active: true},
			nil,
			nil,
		},
	})
	b := g.Board
	g.MoveHistory = append(g.MoveHistory, Move{FromFile: 'c', FromRank: 7, ToFile: 'c', ToRank: 5})
	// Board looks like this:
	// 8 ♖ ♘ ♗ ♕ ♔ ♗
	// 7 ♙ ♙       ♙ ♙ ♟ ♟
	// 6
	// 5    ♟ ♙
	// 4         ♙
	// 3       ♟
	// 2 ♟       ♟ ♟ ♟
	// 1 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜
	//   a  b  c  d  e  f  g  h
	pa2, _ := b.GetPieceAtSquare('a', 2) // Testing forward 1 and forward 2
	pb5, _ := b.GetPieceAtSquare('b', 5) // Testing En passant
	pc3, _ := b.GetPieceAtSquare('c', 3) // Testing capture
	pg7, _ := b.GetPieceAtSquare('g', 7) // Testing capture + promotion
	ph7, _ := b.GetPieceAtSquare('h', 7) // Testing promotion

	tests := []struct {
		piece    *Piece
		fromFile rune
		fromRank int
		expected []Move
	}{
		{pa2, 'a', 2, []Move{
			{FromFile: 'a', FromRank: 2, ToFile: 'a', ToRank: 3},
			{FromFile: 'a', FromRank: 2, ToFile: 'a', ToRank: 4},
		}},
		{pb5, 'b', 5, []Move{
			{FromFile: 'b', FromRank: 5, ToFile: 'b', ToRank: 6},
			{FromFile: 'b', FromRank: 5, Capture: 'x', ToFile: 'c', ToRank: 6},
		}},
		{pc3, 'c', 3, []Move{
			{FromFile: 'c', FromRank: 3, ToFile: 'c', ToRank: 4},
			{FromFile: 'c', FromRank: 3, Capture: 'x', ToFile: 'd', ToRank: 4},
		}},
		{pg7, 'g', 7, []Move{
			{FromFile: 'g', FromRank: 7, ToFile: 'g', ToRank: 8, Promotion: 'Q'},
			{FromFile: 'g', FromRank: 7, ToFile: 'g', ToRank: 8, Promotion: 'R'},
			{FromFile: 'g', FromRank: 7, ToFile: 'g', ToRank: 8, Promotion: 'N'},
			{FromFile: 'g', FromRank: 7, ToFile: 'g', ToRank: 8, Promotion: 'B'},
			{FromFile: 'g', FromRank: 7, Capture: 'x', ToFile: 'f', ToRank: 8, Promotion: 'Q'},
			{FromFile: 'g', FromRank: 7, Capture: 'x', ToFile: 'f', ToRank: 8, Promotion: 'R'},
			{FromFile: 'g', FromRank: 7, Capture: 'x', ToFile: 'f', ToRank: 8, Promotion: 'N'},
			{FromFile: 'g', FromRank: 7, Capture: 'x', ToFile: 'f', ToRank: 8, Promotion: 'B'},
		}},
		{ph7, 'h', 7, []Move{
			{FromFile: 'h', FromRank: 7, ToFile: 'h', ToRank: 8, Promotion: 'Q'},
			{FromFile: 'h', FromRank: 7, ToFile: 'h', ToRank: 8, Promotion: 'R'},
			{FromFile: 'h', FromRank: 7, ToFile: 'h', ToRank: 8, Promotion: 'N'},
			{FromFile: 'h', FromRank: 7, ToFile: 'h', ToRank: 8, Promotion: 'B'},
		}},
	}

	for _, test := range tests {
		moves := test.piece.GetPossibleMoves(&g, test.fromFile, test.fromRank)
		if len(moves) != len(test.expected) {
			t.Errorf("Expected %d moves, got %d", len(test.expected), len(moves))
		}
		// check if move is in expected mvoes
		for _, expectedMove := range test.expected {
			found := false
			for _, move := range moves {
				if expectedMove == move {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected move %v not found", expectedMove)
			}
		}
	}
}

func TestGetRookMoves(t *testing.T) {
	g := NewGame()
	b := g.Board
	// Put rook at d4 so we can test all directions
	b.MovePiece(Move{FromFile: 'a', FromRank: 1, ToFile: 'd', ToRank: 4})
	rook, _ := b.GetPieceAtSquare('d', 4)
	fromFile := 'd'
	fromRank := 4
	expected := []Move{
		// Forward
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 5},
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 6},
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 7, Capture: 'x'},
		// Backward
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 3},
		// Left
		{FromFile: 'd', FromRank: 4, ToFile: 'c', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'b', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'a', ToRank: 4},
		// Right
		{FromFile: 'd', FromRank: 4, ToFile: 'e', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'f', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'g', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'h', ToRank: 4},
	}
	moves := rook.GetPossibleMoves(&g, fromFile, fromRank)
	if len(moves) != len(expected) {
		t.Errorf("Expected %d moves, got %d", len(expected), len(moves))
	}
	// check if move is in expected mvoes
	for _, expectedMove := range expected {
		found := false
		for _, move := range moves {
			if expectedMove == move {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %v not found", expectedMove)
		}
	}
}

func TestGetBishopMoves(t *testing.T) {
	g := NewGame()
	b := g.Board
	// Put bishop at d4 so we can test all directions
	b.MovePiece(Move{FromFile: 'c', FromRank: 1, ToFile: 'd', ToRank: 4})
	bishop, _ := b.GetPieceAtSquare('d', 4)
	fromFile := 'd'
	fromRank := 4
	expected := []Move{
		// Forward left
		{FromFile: 'd', FromRank: 4, ToFile: 'c', ToRank: 5},
		{FromFile: 'd', FromRank: 4, ToFile: 'b', ToRank: 6},
		{FromFile: 'd', FromRank: 4, ToFile: 'a', ToRank: 7, Capture: 'x'},
		// Forward right
		{FromFile: 'd', FromRank: 4, ToFile: 'e', ToRank: 5},
		{FromFile: 'd', FromRank: 4, ToFile: 'f', ToRank: 6},
		{FromFile: 'd', FromRank: 4, ToFile: 'g', ToRank: 7, Capture: 'x'},
		// Backward left
		{FromFile: 'd', FromRank: 4, ToFile: 'c', ToRank: 3},
		// Backward right
		{FromFile: 'd', FromRank: 4, ToFile: 'e', ToRank: 3},
	}
	moves := bishop.GetPossibleMoves(&g, fromFile, fromRank)
	if len(moves) != len(expected) {
		t.Errorf("Expected %d moves, got %d", len(expected), len(moves))
	}
	// check if move is in expected mvoes
	for _, expectedMove := range expected {
		found := false
		for _, move := range moves {
			if expectedMove == move {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %v not found", expectedMove)
		}
	}
}

func TestGetQueenMoves(t *testing.T) {
	g := NewGame()
	b := g.Board
	// Put queen at d4 so we can test all directions
	b.MovePiece(Move{FromFile: 'd', FromRank: 1, ToFile: 'd', ToRank: 4})
	queen, _ := b.GetPieceAtSquare('d', 4)
	fromFile := 'd'
	fromRank := 4
	expected := []Move{
		// Forward
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 5},
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 6},
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 7, Capture: 'x'},
		// Backward
		{FromFile: 'd', FromRank: 4, ToFile: 'd', ToRank: 3},
		// Left
		{FromFile: 'd', FromRank: 4, ToFile: 'c', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'b', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'a', ToRank: 4},
		// Right
		{FromFile: 'd', FromRank: 4, ToFile: 'e', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'f', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'g', ToRank: 4},
		{FromFile: 'd', FromRank: 4, ToFile: 'h', ToRank: 4},
		// Forward left
		{FromFile: 'd', FromRank: 4, ToFile: 'c', ToRank: 5},
		{FromFile: 'd', FromRank: 4, ToFile: 'b', ToRank: 6},
		{FromFile: 'd', FromRank: 4, ToFile: 'a', ToRank: 7, Capture: 'x'},
		// Forward right
		{FromFile: 'd', FromRank: 4, ToFile: 'e', ToRank: 5},
		{FromFile: 'd', FromRank: 4, ToFile: 'f', ToRank: 6},
		{FromFile: 'd', FromRank: 4, ToFile: 'g', ToRank: 7, Capture: 'x'},
		// Backward left
		{FromFile: 'd', FromRank: 4, ToFile: 'c', ToRank: 3},
		// Backward right
		{FromFile: 'd', FromRank: 4, ToFile: 'e', ToRank: 3},
	}
	moves := queen.GetPossibleMoves(&g, fromFile, fromRank)
	if len(moves) != len(expected) {
		t.Errorf("Expected %d moves, got %d", len(expected), len(moves))
	}
	// check if move is in expected mvoes
	for _, expectedMove := range expected {
		found := false
		for _, move := range moves {
			if expectedMove == move {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %v not found", expectedMove)
		}
	}
}
