package game

import (
	"testing"
)

func TestGetSymbol(t *testing.T) {
	tests := []struct {
		piece     piece
		expected  string
		expectErr bool
	}{
		{piece{pType: King}, "K", false},
		{piece{pType: Queen}, "Q", false},
		{piece{pType: Rook}, "R", false},
		{piece{pType: Bishop}, "B", false},
		{piece{pType: Knight}, "N", false},
		{piece{pType: Pawn}, "", false},
	}

	for _, test := range tests {
		result, err := test.piece.getSymbol()
		if (err != nil) != test.expectErr {
			t.Errorf("Expected error: %v, got: %v", test.expectErr, err)
		}
		if result != test.expected {
			t.Errorf("Expected: %s, got: %s", test.expected, result)
		}
	}
}

func TestGetPrintable(t *testing.T) {
	tests := []struct {
		piece     piece
		expected  string
		expectErr bool
	}{
		{piece{pType: Pawn, color: "white"}, "♟", false},
		{piece{pType: Pawn, color: "black"}, "♙", false},
		{piece{pType: King, color: "white"}, "♚", false},
		{piece{pType: King, color: "black"}, "♔", false},
		{piece{pType: Queen, color: "white"}, "♛", false},
		{piece{pType: Queen, color: "black"}, "♕", false},
		{piece{pType: Rook, color: "white"}, "♜", false},
		{piece{pType: Rook, color: "black"}, "♖", false},
		{piece{pType: Bishop, color: "white"}, "♝", false},
		{piece{pType: Bishop, color: "black"}, "♗", false},
		{piece{pType: Knight, color: "white"}, "♞", false},
		{piece{pType: Knight, color: "black"}, "♘", false},
	}

	for _, test := range tests {
		result, err := test.piece.getPrintable()
		if (err != nil) != test.expectErr {
			t.Errorf("Expected error: %v, got: %v", test.expectErr, err)
		}
		if result != test.expected {
			t.Errorf("Expected: %s, got: %s", test.expected, result)
		}
	}
}
